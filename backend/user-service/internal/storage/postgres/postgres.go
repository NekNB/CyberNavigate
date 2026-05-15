package postgres

import (
	"database/sql"
	"errors"

	"github.com/NekNB/CyberNavigate/backend/user-service/internal/domain/models"
	sessionService "github.com/NekNB/CyberNavigate/backend/user-service/internal/services/session"
	userService "github.com/NekNB/CyberNavigate/backend/user-service/internal/services/user"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/storage"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/sirupsen/logrus"
)

type PostgresStorage struct {
	db *sql.DB
}

var _ userService.UserDataProvider = (*PostgresStorage)(nil)

var _ sessionService.SessionProvider = (*PostgresStorage)(nil)

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

// ============ USERS =====================
func (p *PostgresStorage) Users() ([]*models.UserDTO, error) {
	rows, err := p.db.Query(
		`
			SELECT uuid, username, is_admin, created_at
			FROM users;
		`,
	)
	if err != nil {
		return nil, err
	}

	var users []*models.UserDTO

	for rows.Next() {
		user := models.UserDTO{}
		if err := rows.Scan(
			&user.UserId,
			&user.Username,
			&user.IsAdmin,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (p *PostgresStorage) UserByUserId(userId string) (*models.UserDTO, error) {
	stmt, err := p.db.Prepare(
		`
			SELECT uuid, username, is_admin, created_at
			FROM users
			WHERE uuid = $1;
		`,
	)
	if err != nil {
		return nil, err
	}

	var user models.UserDTO

	if err = stmt.QueryRow(userId).Scan(
		&user.UserId,
		&user.Username,
		&user.IsAdmin,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
func (p *PostgresStorage) UserByUsername(username string) (*models.UserDTO, error) {
	stmt, err := p.db.Prepare(
		`
			SELECT uuid, username, is_admin, created_at
			FROM users
			WHERE username = $1;
		`,
	)
	if err != nil {
		return nil, err
	}

	var user models.UserDTO

	if err = stmt.QueryRow(username).Scan(
		&user.UserId,
		&user.Username,
		&user.IsAdmin,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
func (p *PostgresStorage) PasswordSaltByUsername(username string) (passwordHash, salt string, err error) {
	if err = p.db.QueryRow(
		`
			SELECT password_hash, salt
			FROM users
			WHERE username = $1;
		`, username).
		Scan(
			&passwordHash,
			&salt,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = storage.ErrUserNotFound
			return
		}
		return
	}

	return
}

func (p *PostgresStorage) NewUser(username, passwordHash, salt string) error {
	if _, err := p.db.Exec(`
		INSERT INTO users (username, password_hash, salt)
		VALUES ($1, $2, $3);
	`, username, passwordHash, salt); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.UniqueViolation {
			return storage.ErrUserExists
		}
		return err
	}

	return nil
}

// ============ SESSIONS =====================
func (p *PostgresStorage) UserInfoByRefreshToken(refreshToken string) (*models.UserDTO, error) {

	var user models.UserDTO
	if err := p.db.QueryRow(
		`
			SELECT 
					u.uuid,
					u.is_admin
			FROM sessions s
			JOIN users u ON s.user_id = u.uuid
			WHERE 
					s.refresh_token = $1;
		`,
		refreshToken).
		Scan(
			&user.UserId,
			&user.IsAdmin,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrRefreshNotValid
		}
		return nil, err
	}

	return &user, nil
}

func (p *PostgresStorage) ExtendSession(refreshToken string, expires int) (sessionId string, err error) {

	if err = p.db.QueryRow(`
		UPDATE sessions 
		SET 
				expires_at = CURRENT_TIMESTAMP + $2::INTERVAL,
				updated_at = CURRENT_TIMESTAMP
		WHERE 
				refresh_token = $1 
				AND expires_at > CURRENT_TIMESTAMP 
				AND revoked = false
		RETURNING uuid;
	`, refreshToken, expires).
		Scan(
			&sessionId,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = storage.ErrRefreshNotValid
			return
		}
		return
	}

	return
}

func (p *PostgresStorage) NewSession(userId string, refreshTokenExpiration int) (*models.SessionDTO, error) {
	var session models.SessionDTO
	if err := p.db.QueryRow(`
		INSERT INTO sessions (user_id, expires_at)
		VALUES ($1, CURRENT_TIMESTAMP + $2::INTERVAL)
		RETURNING uuid, refresh_token;
	`, userId, refreshTokenExpiration).Scan(
		&session.SessionId,
		&session.RefreshToken,
	); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.ForeignKeyViolation {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (p *PostgresStorage) RevokeSession(sessionId string) error {
	if _, err := p.db.Exec(`
		UPDATE sessions
		SET revoked = True
		WHERE session_id = $1
	`, sessionId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrSessionNotFound
		}
		return err
	}
	return nil
}
