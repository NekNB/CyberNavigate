package postgres

import (
	"database/sql"
	"errors"

	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/storage"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/sirupsen/logrus"
)

type PostgresStorage struct {
	db *sql.DB
}

var _ simulatorService.simulatorMetaProvider = (*PostgresStorage)(nil)

func New(log *logrus.Logger, uri string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Errorf("%s", err.Error())
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("%s", err.Error())
	}
	return &PostgresStorage{db: db}, nil
}

// Выборка все simulator metadata
func (p *PostgresStorage) simulators() (*[]simulator.simulatorMetaData, error) {
	rows, err := p.db.Query(
		`
			SELECT uuid, title, status 
			FROM metadata;
		`,
	)
	if err != nil {
		return nil, err
	}

	var simulatorSlice []simulator.simulatorMetaData

	for rows.Next() {
		metadata := simulator.simulatorMetaData{}
		if err := rows.Scan(
			&metadata.Id,
			&metadata.Title,
			&metadata.Status,
		); err != nil {
			return nil, err
		}

		simulatorSlice = append(simulatorSlice, metadata)
	}

	return &simulatorSlice, nil
}

// Получение конкретного simulator metadata по uuid
func (p *PostgresStorage) simulatorByUUID(simulatorUUID string) (*simulator.simulatorMetaData, error) {
	stmt, err := p.db.Prepare(
		`
			SELECT uuid, title, status 
			FROM metadata
			WHERE uuid = $1;
		`,
	)
	if err != nil {
		return nil, err
	}

	var simulatorMetadata simulator.simulatorMetaData

	if err = stmt.QueryRow(simulatorUUID).Scan(
		&simulatorMetadata.Id,
		&simulatorMetadata.Title,
		&simulatorMetadata.Status,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrsimulatorNotFound
		}
		return nil, err
	}

	return &simulatorMetadata, nil
}

// Получение конкретного simulator text по uuid
func (p *PostgresStorage) simulatorTextIDByUUID(simulatorUUID string) (string, error) {

	var textIdNull sql.NullString
	if err := p.db.QueryRow(
		`
			SELECT text_id
			FROM metadata
			WHERE uuid = $1;
		`,
		simulatorUUID,
	).Scan(&textIdNull); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrsimulatorNotFound
		}
		return "", err
	}
	if textIdNull.Valid {
		return textIdNull.String, nil
	}
	return "", storage.ErrsimulatorTextNotCreatedYet
}

// Создание сущности simulator
func (p *PostgresStorage) Createsimulator(simulatorTitle *string) (*simulator.simulatorMetaData, error) {
	var metadata simulator.simulatorMetaData

	if err := p.db.QueryRow(`
		INSERT INTO metadata (title)
		VALUES ($1)
		RETURNING uuid, title, status;
	`, simulatorTitle).
		Scan(
			&metadata.Id,
			&metadata.Title,
			&metadata.Status,
		); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.UniqueViolation {
			return nil, storage.ErrsimulatorExists
		}
		return nil, err
	}

	return &metadata, nil
}

// Обновление сущности simulator по UUID
func (p *PostgresStorage) UpdatesimulatorByUUID(uuid, title, textID, status, videoUrl string) (*simulator.simulatorMetaData, error) {
	var metadata simulator.simulatorMetaData

	if err := p.db.QueryRow(`
		UPDATE metadata
		SET
			title = COALESCE(NULLIF($2, ''), title),
			text_id = COALESCE(NULLIF($3, ''), text_id),
			status = COALESCE(NULLIF($4, '')::simulator_status, status),
			video_url = COALESCE(NULLIF($5, ''), video_url)
		WHERE uuid = $1
		RETURNING uuid, title, status;
	`, uuid, title, textID, status, videoUrl).
		Scan(
			&metadata.Id,
			&metadata.Title,
			&metadata.Status,
		); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.UniqueViolation {
			return nil, storage.ErrsimulatorNotFound
		}
		return nil, err
	}

	return &metadata, nil
}
