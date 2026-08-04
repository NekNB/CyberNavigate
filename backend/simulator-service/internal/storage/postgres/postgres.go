package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	cErr "github.com/NekNB/CyberNavigate/backend/simulator-service/internal/lib/errors"
	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/models"
	simulatorService "github.com/NekNB/CyberNavigate/backend/simulator-service/internal/services/simulator"
	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/storage"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/sirupsen/logrus"
)

type PostgresStorage struct {
	db  *sql.DB
	log *logrus.Logger
}

var _ simulatorService.SimulatorDataProvider = (*PostgresStorage)(nil)

func New(log *logrus.Logger, uri string) (*PostgresStorage, error) {

	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Errorf("%s", err.Error())
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("%s", err.Error())
	}
	return &PostgresStorage{db: db, log: log}, nil
}

func (p *PostgresStorage) CreateSession(userId, stepId string) (string, error) {
	var sessionId string

	// Блокируем все старые сессии
	if _, err := p.db.Exec(`
    UPDATE sessions 
    SET finished_at = CURRENT_TIMESTAMP
    WHERE user_id = $1 AND finished_at IS NULL;
	`, userId); err != nil {
		p.log.Error(err)
		return "", err
	}

	// Создаем новую сессию
	if err := p.db.QueryRow(`
		INSERT INTO sessions (user_id, current_step)
		VALUES ($1, $2)
		RETURNING uuid;
	`, userId, stepId).Scan(&sessionId); err != nil {
		p.log.Error(err)
		return "", err
	}
	return sessionId, nil
}
func (p *PostgresStorage) GetCurrentSession(userId string) (*models.SessionData, error) {
	var sessionData models.SessionData
	var finishedAt sql.NullTime
	if err := p.db.QueryRow(`
	SELECT uuid, created_at, current_step, current_trust, finished_at, is_finished
	FROM sessions
	WHERE user_id = $1 AND finished_at IS NULL;
	`, userId).Scan(
		&sessionData.UUID,
		&sessionData.CreatedAt,
		&sessionData.CurrentStepId,
		&sessionData.CurrentTrust,
		&finishedAt,
		&sessionData.IsFinished,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, cErr.NewTypedError(storage.ErrResourseNotFound, userId)
		}
		p.log.Error(err)
		return nil, err
	}

	if finishedAt.Valid {
		sessionData.FinishedAt = &finishedAt.Time
	}
	return &sessionData, nil
}

func (p *PostgresStorage) GetSessionBySessionId(sessionId string) (*models.SessionData, error) {
	var sessionData models.SessionData
	var finishedAt sql.NullTime
	if err := p.db.QueryRow(`
	SELECT uuid, created_at, current_step, current_trust, finished_at
	FROM sessions
	WHERE uuid = $1;
	`, sessionId).Scan(
		&sessionData.UUID,
		&sessionData.CreatedAt,
		&sessionData.CurrentStepId,
		&sessionData.CurrentTrust,
		&finishedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return nil, err
	}

	if finishedAt.Valid {
		sessionData.FinishedAt = &finishedAt.Time
	}
	return &sessionData, nil
}
func (p *PostgresStorage) SetCurrentStep(sessionId string, stepId *string) error {
	if _, err := p.db.Exec(`
		UPDATE sessions
		SET
			current_step = $1
		WHERE uuid = $2
	`, stepId, sessionId); err != nil {
		p.log.Error(err)
		return err
	}
	return nil
}

func (p *PostgresStorage) SetSessionError(sessionId, userError string) error {
	if _, err := p.db.Exec(`
	INSERT INTO errors_to_session (error, session_id) 
	VALUES ($1, $2);
	`, userError, sessionId); err != nil {
		p.log.Error(err)
		return err
	}
	return nil
}
func (p *PostgresStorage) MarkSessionAsFinished(sessionId string) error {
	if _, err := p.db.Exec(`
		UPDATE sessions
		SET
			is_finished = true
		WHERE uuid = $1;
	`, sessionId); err != nil {
		p.log.Error(err)
		return err
	}
	return nil
}

func (p *PostgresStorage) CloseSession(sessionId string) error {
	if _, err := p.db.Exec(`
		UPDATE sessions
		SET
			finished_at = CURRENT_TIMESTAMP
		WHERE uuid = $1;
	`, sessionId); err != nil {
		p.log.Error(err)
		return err
	}
	return nil
}

func (p *PostgresStorage) CreateScenario(title, description, difficulty string) (string, error) {
	var scenarioId string

	if err := p.db.QueryRow(`
		INSERT INTO scenarios (title, description, difficulty)
		VALUES ($1, $2, $3)
		RETURNING uuid;
	`, title, description, difficulty).
		Scan(
			&scenarioId,
		); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.UniqueViolation {
			return "", cErr.NewTypedError(storage.ErrAlreadyExists, fmt.Sprintf("Scenario with title: %s already exists", title))
		} else if errors.Is(err, sql.ErrNoRows) {
			return "", cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return "", err
	}

	return scenarioId, nil
}
func (p *PostgresStorage) EditScenario(scenarioId string, title, description, difficulty *string) error {
	if _, err := p.db.Exec(`
		UPDATE scenarios
		SET
			title = COALESCE($1, title),
			description = COALESCE($2, description),
			difficulty = COALESCE($3, difficulty)
		WHERE uuid = $4;
	`, title, description, difficulty, scenarioId); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.UniqueViolation {
			return cErr.NewTypedError(storage.ErrAlreadyExists, fmt.Sprintf("Scenario with title: %s already exists", *title))
		} else if errors.Is(err, sql.ErrNoRows) {
			return cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return err
	}
	return nil
}
func (p *PostgresStorage) GetAllScenarios() (*[]models.ScenarioData, error) {
	rows, err := p.db.Query(`
		SELECT 
			s.uuid,
			s.title,
			s.description,
			s.difficulty,
			s.first_step,
			s.created_at,
			s.updated_at,
		COALESCE(
				json_agg(a.article_id) FILTER (WHERE a.article_id IS NOT NULL),
				'[]'
		) AS article_ids
		FROM scenarios s
		LEFT JOIN article_ids_to_scenario_id a ON s.uuid = a.scenario_id 
		GROUP BY s.uuid
		ORDER BY s.created_at DESC;
	`)
	if err != nil {
		p.log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var scenarios []models.ScenarioData
	for rows.Next() {
		scenarioData := models.ScenarioData{}
		var artilceIdsJSON []uint8
		var firstStep sql.NullString
		if err := rows.Scan(
			&scenarioData.UUID,
			&scenarioData.Title,
			&scenarioData.Description,
			&scenarioData.Difficulty,
			&firstStep,
			&scenarioData.CreatedAt,
			&scenarioData.UpdatedAt,
			&artilceIdsJSON,
		); err != nil {
			p.log.Error(err)
			return nil, err
		}

		if err := json.Unmarshal(artilceIdsJSON, &scenarioData.ArticleIds); err != nil {
			p.log.Error(err)
			return nil, err
		}

		if firstStep.Valid {
			scenarioData.FirstStep = firstStep.String
		}
		scenarios = append(scenarios, scenarioData)
	}

	return &scenarios, nil
}

func (p *PostgresStorage) GetScenario(scenarioId string) (*models.ScenarioData, error) {
	var firstStep sql.NullString
	var scenarioData models.ScenarioData
	var artilceIdsJSON []uint8
	if err := p.db.QueryRow(`
		SELECT 
			s.uuid,
			s.title,
			s.description,
			s.difficulty,
			s.first_step,
			s.created_at,
			s.updated_at,
		COALESCE(
				json_agg(a.article_id) FILTER (WHERE a.article_id IS NOT NULL),
				'[]'
		) AS article_ids
		FROM scenarios s
		LEFT JOIN article_ids_to_scenario_id a ON s.uuid = a.scenario_id 
		WHERE s.uuid = $1
		GROUP BY s.uuid;
	`, scenarioId).Scan(
		&scenarioData.UUID,
		&scenarioData.Title,
		&scenarioData.Description,
		&scenarioData.Difficulty,
		&firstStep,
		&scenarioData.CreatedAt,
		&scenarioData.UpdatedAt,
		&artilceIdsJSON,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return nil, err
	}

	if err := json.Unmarshal(artilceIdsJSON, &scenarioData.ArticleIds); err != nil {
		p.log.Error(err)
		return nil, err
	}

	if firstStep.Valid {
		scenarioData.FirstStep = firstStep.String
	}
	return &scenarioData, nil
}

func (p *PostgresStorage) MatchArticleToScenario(articleIds *[]string, scenarioId string) error {
	if _, err := p.db.Exec(`
		INSERT INTO article_ids_to_scenario_id (scenario_id, article_id)
		SELECT 
				$2::uuid,
				unnest($1::uuid[]);
	`, articleIds, scenarioId); err != nil {
		p.log.Error(err)
		return err
	}
	return nil
}
func (p *PostgresStorage) DeleteAllMatchArtilceToScenario(scenarioId string) error {
	if _, err := p.db.Exec(`
		DELETE FROM article_ids_to_scenario_id 
		WHERE scenario_id = $1;
	`, scenarioId); err != nil {
		p.log.Error(err)
		return err
	}
	return nil
}
func (p *PostgresStorage) SetFirstStep(scenarioId, stepId string) error {
	if _, err := p.db.Exec(`
		UPDATE scenarios
		SET
			first_step = $1
		WHERE uuid = $2;
	`, stepId, scenarioId); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.ForeignKeyViolation {
			return cErr.NewTypedError(storage.ErrNotFound, fmt.Sprintf("Step With Id: %s not found", stepId))
		}
		p.log.Error(err)
		return err
	}
	return nil
}

func (p *PostgresStorage) GetErrors(sessionId string) (*[]string, error) {
	rows, err := p.db.Query(`
		SELECT error
		FROM errors_to_session
		WHERE session_id = $1;
	`, sessionId)
	if err != nil {
		p.log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var errors []string
	for rows.Next() {
		var userError string
		if err := rows.Scan(
			&userError,
		); err != nil {
			p.log.Error(err)
			return nil, err
		}

		errors = append(errors, userError)
	}
	return &errors, nil
}

func (p *PostgresStorage) GetTrusts(sessionId string) (*[]int, error) {
	rows, err := p.db.Query(`
		SELECT trust
		FROM trusts
		WHERE session_id = $1;
	`, sessionId)
	if err != nil {
		p.log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var trusts []int
	for rows.Next() {
		var trust int
		if err := rows.Scan(
			&trust,
		); err != nil {
			p.log.Error(err)
			return nil, err
		}

		trusts = append(trusts, trust)
	}
	return &trusts, nil
}

func (p *PostgresStorage) CreateStep(scenariodId string, previousAnswer, previosStep *string, maxTrust, minTrust *int) (string, error) {
	var stepId string

	if err := p.db.QueryRow(`
		INSERT INTO steps (previous_answer, previous_step, max_trust, min_trust, scenario_id)
		VALUES (COALESCE($1::uuid, NULL), COALESCE($2::uuid, NULL), COALESCE($3::int, 100), COALESCE($4::int, -100), $5)
		RETURNING uuid;
	`, previousAnswer, previosStep, maxTrust, minTrust, scenariodId).
		Scan(
			&stepId,
		); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.CheckViolation {
			return "", cErr.NewTypedError(storage.ErrDataConfllict, "")
		} else if errors.Is(err, sql.ErrNoRows) {
			return "", cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return "", err
	}

	return stepId, nil
}
func (p *PostgresStorage) EditStep(stepID string, maxTrust, minTrust *int) error {
	if _, err := p.db.Exec(`
		UPDATE steps 
		SET 
			max_trust = COALESCE($1::int, max_trust),
			min_trust = COALESCE($2::int, min_trust)
		WHERE uuid = $3::uuid;
	`, maxTrust, minTrust, stepID); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.CheckViolation {
			return cErr.NewTypedError(storage.ErrDataConfllict, "")
		} else if errors.Is(err, sql.ErrNoRows) {
			return cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return err
	}

	return nil
}
func (p *PostgresStorage) GetStep(stepId string) (*models.StepData, error) {
	// var previousStep, previosAnswer sql.NullString
	// var MinTrust, MaxTrust sql.NullInt64
	var step models.StepData
	var actionIds []uint8

	var previousStep, previousAnswer sql.NullString

	if err := p.db.QueryRow(`
		SELECT 
			s.uuid,
			s.previous_step,
			s.previous_answer,
			s.min_trust,
			s.max_trust,
			s.scenario_id,
			s.created_at,
			s.updated_at,
			COALESCE(
				json_agg(a.uuid) FILTER (WHERE a.uuid IS NOT NULL),
				'[]'
			) AS action_ids
		FROM steps s
		LEFT JOIN actions a ON s.uuid = a.step_id
		WHERE s.uuid = $1::uuid
		GROUP BY 
			s.uuid,
			s.previous_step,
			s.previous_answer,
			s.min_trust,
			s.max_trust,
			s.scenario_id,
			s.created_at,
			s.updated_at;
	`, stepId).Scan(
		&step.UUID,
		&previousStep,
		&previousAnswer,
		&step.MinTrust,
		&step.MaxTrust,
		&step.ScenarioId,
		&step.CreatedAt,
		&step.UpdatedAt,
		&actionIds,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return nil, err
	}

	if err := json.Unmarshal(actionIds, &step.ActionIds); err != nil {
		p.log.Error(err)
		return nil, err
	}

	if previousStep.Valid {
		step.PreviousStep = &previousStep.String
	}
	if previousAnswer.Valid {
		step.PreviousAnswer = &previousAnswer.String
	}

	return &step, nil
}

func (p *PostgresStorage) MatchActionToStep(actionIds *[]string, stepId string) error {
	if _, err := p.db.Exec(`
		UPDATE actions 
		SET 
			step_id = $2::uuid
		WHERE uuid = ANY($1::uuid[]);
	`, actionIds, stepId); err != nil {
		p.log.Error(err)
		return err
	}
	return nil
}
func (p *PostgresStorage) DeleteAllMatchActionToStep(stepId string) error {
	if _, err := p.db.Exec(`
		UPDATE actions 
		SET 
			step_id = NULL
		WHERE step_id = $1;
	`, stepId); err != nil {
		p.log.Error(err)
		return err
	}
	return nil
}
func (p *PostgresStorage) GetNextStepByStepId(currentStepId string, currentTrust int) (*string, error) {
	var nextStepId string
	if err := p.db.QueryRow(`
	SELECT uuid
	FROM steps 
	WHERE previous_step = $1 AND uuid != $1 AND $2 BETWEEN min_trust AND max_trust;
	`, currentStepId, currentTrust).Scan(&nextStepId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		p.log.Error(err)
		return nil, err
	}

	return &nextStepId, nil
}
func (p *PostgresStorage) GetNextStepByAnswerId(answerId string, currentTrust int) (*string, error) {
	var nextStepId string
	if err := p.db.QueryRow(`
		SELECT uuid
		FROM steps 
		WHERE previous_answer = $1 AND $2 BETWEEN min_trust AND max_trust;
	`, answerId, currentTrust).Scan(&nextStepId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		p.log.Error(err)
		return nil, err
	}

	return &nextStepId, nil
}

func (p *PostgresStorage) CreateAction(actionType, messageId string, delay int) (string, error) {
	var actionId string

	if err := p.db.QueryRow(`
		INSERT INTO actions (type, message_id, delay)
		VALUES ($1, $2, $3)
		RETURNING uuid;
	`, actionType, messageId, delay).Scan(&actionId); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.ForeignKeyViolation {
			return "", cErr.NewTypedError(storage.ErrDataConfllict, "")
		}
		p.log.Error(err)
		return "", err
	}

	return actionId, nil
}
func (p *PostgresStorage) EditAction(actionId, actionType, messageId string, delay *int) error {
	if _, err := p.db.Exec(`
		UPDATE actions 
		SET 
				type = $1,
				message_id = $2,
				delay = COALESCE($3, delay)
		WHERE uuid = $4;
    `, actionType, messageId, delay, actionId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return err
	}

	return nil
}
func (p *PostgresStorage) GetActionById(actionId string) (*models.ActionData, error) {
	var action models.ActionData
	if err := p.db.QueryRow(`
		SELECT 
				uuid,
				type,
				message_id,
				delay
		FROM actions
		WHERE uuid = $1;
    `, actionId).Scan(
		&action.UUID,
		&action.Type,
		&action.MessageId,
		&action.Delay,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return nil, err
	}

	return &action, nil
}

func (p *PostgresStorage) CreateMessage(senderId, text *string, senderName string) (string, error) {
	var messageId string

	err := p.db.QueryRow(`
		INSERT INTO messages (sender_id, sender_name, text)
		VALUES (COALESCE($1, uuid_generate_v4()), $2, COALESCE($3, NULL))
		RETURNING uuid;
	`, senderId, senderName, text).Scan(&messageId)

	if err != nil {
		p.log.Error(err)
		return "", err
	}

	return messageId, nil
}
func (p *PostgresStorage) EditMessage(messageId string, senderId, senderName, text *string) error {
	if _, err := p.db.Exec(`
		UPDATE messages
		SET 
			sender_id = COALESCE($1, sender_id),
			sender_name = COALESCE($2, sender_name),
			text = COALESCE($3, text)
		WHERE uuid = $4::uuid;
	`, senderId, senderName, text, messageId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return err
	}

	return nil
}
func (p *PostgresStorage) MatchAnswerToMessage(answerIds *[]string, messageId string) error {
	if _, err := p.db.Exec(`
		UPDATE answers 
		SET 
			message_id = $2::uuid
		WHERE uuid = ANY($1::uuid[]);
	`, answerIds, messageId); err != nil {
		p.log.Error(err)
		return err
	}
	return nil
}
func (p *PostgresStorage) DeleteAllMatchAnswerToMessage(messageId string) error {
	if _, err := p.db.Exec(`
		UPDATE answers 
		SET 
			message_id = NULL
		WHERE message_id = $1;
	`, messageId); err != nil {
		p.log.Error(err)
		return err
	}
	return nil
}
func (p *PostgresStorage) DeleteAllMatchFileToMessage(messageId string) error {
	if _, err := p.db.Exec(`
		UPDATE files 
		SET 
			message_id = NULL,
		WHERE message_id = $1;
	`, messageId); err != nil {
		p.log.Error(err)
		return err
	}
	return nil
}
func (p *PostgresStorage) MatchFileToMessage(filesIds *[]string, messageId string) error {
	if _, err := p.db.Exec(`
		UPDATE files 
		SET 
			message_id = $2::uuid
		WHERE uuid = ANY($1::uuid[]);
	`, filesIds, messageId); err != nil {
		p.log.Error(err)
		return err
	}
	return nil
}
func (p *PostgresStorage) GetMessageById(messageId string) (*models.MessageData, error) {
	var message models.MessageData
	var text sql.NullString

	if err := p.db.QueryRow(`
		SELECT 
			uuid,
			sender_id,
			sender_name,
			text
		FROM messages
		WHERE uuid = $1::uuid;
	`, messageId).Scan(
		&message.UUID,
		&message.SenderId,
		&message.SenderName,
		&text,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return nil, err
	}

	if text.Valid {
		message.Text = &text.String
	} else {
		message.Text = nil
	}

	return &message, nil
}

func (p *PostgresStorage) CreateAnswer(addTrust int, errorText, text string) (string, error) {
	var answerId string

	if err := p.db.QueryRow(`
		INSERT INTO answers (add_trust, error, text)
		VALUES ($1, $2, $3)
		RETURNING uuid;
	`, addTrust, errorText, text).Scan(&answerId); err != nil {
		p.log.Error(err)
		return "", err
	}

	return answerId, nil
}
func (p *PostgresStorage) EditAnswer(answerId string, errorText, text *string, addTrust *int) error {
	if _, err := p.db.Exec(`
		UPDATE answers 
		SET 
				error_text = COALESCE($1, error_text),
				text = COALESCE($2, text),
				add_trust = COALESCE($3, add_trust),
				updated_at = CURRENT_TIMESTAMP
		WHERE uuid = $4;
    `, errorText, text, addTrust, answerId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return err
	}

	return nil
}
func (p *PostgresStorage) GetAnswerById(answerId string) (*models.AnswerData, error) {
	var answer models.AnswerData
	if err := p.db.QueryRow(`
		SELECT 
				uuid,
				text,
				add_trust,
				error
		FROM answers
		WHERE uuid = $1;
    `, answerId).Scan(
		&answer.UUID,
		&answer.Text,
		&answer.AddTrust,
		&answer.Error,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return nil, err
	}

	return &answer, nil
}
func (p *PostgresStorage) GetAnswersByMessageId(messageId string) (*[]models.AnswerData, error) {
	rows, err := p.db.Query(`
		SELECT 
				uuid,
				text,
				add_trust,
				error
		FROM answers
		WHERE message_id = $1;
  `, messageId)
	if err != nil {
		p.log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var answers []models.AnswerData
	for rows.Next() {
		var answer models.AnswerData
		if err := rows.Scan(
			&answer.UUID,
			&answer.Text,
			&answer.AddTrust,
			&answer.Error,
		); err != nil {
			p.log.Error(err)
			return nil, err
		}
		answers = append(answers, answer)
	}

	if err := rows.Err(); err != nil {
		p.log.Error(err)
		return nil, err
	}

	return &answers, nil
}

func (p *PostgresStorage) SaveBeginTrustLevel(sessionId string) error {
	if _, err := p.db.Exec(`
		INSERT INTO trusts (session_id, trust)
		VALUES ($1, 0);
	`, sessionId); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.ForeignKeyViolation {
			return cErr.NewTypedError(storage.ErrNotFound, fmt.Sprintf("Session With Id: %s not found", sessionId))
		}
		p.log.Error(err)
		return err
	}
	return nil
}

func (p *PostgresStorage) RegisterTrustLevel(sessionId string, addTrust int) error {
	if _, err := p.db.Exec(`
		WITH updated AS (
			UPDATE sessions 
			SET current_trust = GREATEST(-100, LEAST(100, current_trust + $1::int))
			WHERE uuid = $2
			RETURNING uuid, current_trust
		)
		INSERT INTO trusts (session_id, trust)
		SELECT $2, current_trust
		FROM updated;
	`, addTrust, sessionId); err != nil {

	}
	return nil
}
func (p *PostgresStorage) CreateFile(filename string, isSafe bool, size int, fileError *string) (string, error) {
	var id string
	if err := p.db.QueryRow(`
		INSERT INTO files (
			filename,
			is_safe,
			size,
			error
		) VALUES ($1, $2, $3, $4)
		RETURNING uuid
    `, filename, isSafe, size, fileError).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", cErr.NewTypedError(storage.ErrResourseNotFound, "")
		}
		p.log.Error(err)
		return "", err
	}

	return id, nil
}
func (p *PostgresStorage) EditFile(fileId string, filename, fileError *string, isSafe *bool, size *int) error {
	if _, err := p.db.Exec(`
        UPDATE files 
        SET 
            filename = COALESCE($1, filename),
            error = COALESCE($2, error),
            is_safe = COALESCE($3, is_safe),
            size = COALESCE($4, size),
            updated_at = CURRENT_TIMESTAMP
        WHERE uuid = $5
    `, filename, fileError, isSafe, size, fileId); err != nil {
		p.log.Error(err)
		return err
	}

	return nil
}
func (p *PostgresStorage) GetFilesByMessageId(messageId string) (*[]models.FileData, error) {
	rows, err := p.db.Query(`
		SELECT uuid, filename, is_safe, size, error
		FROM files
		WHERE message_id = $1;
    `, messageId)
	if err != nil {
		p.log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var files []models.FileData
	for rows.Next() {
		var file models.FileData
		if err := rows.Scan(
			&file.UUID,
			&file.Filename,
			&file.IsSafe,
			&file.Size,
			&file.Error,
		); err != nil {
			p.log.Error(err)
			return nil, err
		}
		files = append(files, file)
	}

	return &files, nil
}
func (p *PostgresStorage) GetFileById(fileId string) (*models.FileData, error) {
	var file models.FileData
	if err := p.db.QueryRow(`
		SELECT uuid, filename, is_safe,  size, error
		FROM files
		WHERE uuid = $1;
    `, fileId).Scan(
		&file.UUID,
		&file.Filename,
		&file.IsSafe,
		&file.Size,
		&file.Error,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, cErr.NewTypedError(storage.ErrResourseNotFound, fileId)
		}
		p.log.Error(err)
		return nil, err
	}

	return &file, nil
}
