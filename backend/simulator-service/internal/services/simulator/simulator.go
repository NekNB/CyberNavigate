package simulator

import (
	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/http"
	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/lib/convert"
	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/models"
	"github.com/NekNB/CyberNavigate/swagger/gen/simulator"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"

	"github.com/sirupsen/logrus"
)

var _ http.SimulatorServiceInterface = (*simulatorService)(nil)

type SimulatorDataProvider interface {
	CreateSession(userId, stepId string) (string, error)
	GetCurrentSession(userId string) (*models.SessionData, error)
	GetSessionBySessionId(sessionId string) (*models.SessionData, error)
	SetCurrentStep(sessionId string, stepId *string) error
	SetSessionError(sessionId, error string) error
	SaveBeginTrustLevel(sessionId string) error
	MarkSessionAsFinished(sessionId string) error

	CreateScenario(title, description, difficulty string) (string, error)
	EditScenario(scenarioId string, title, description, difficulty *string) error
	GetAllScenarios() (*[]models.ScenarioData, error)
	GetScenario(scenarioId string) (*models.ScenarioData, error)
	MatchArticleToScenario(articleIds *[]string, scenarioId string) error
	DeleteAllMatchArtilceToScenario(scenarioId string) error
	SetFirstStep(scenarioId, stepId string) error

	GetErrors(sessionId string) (*[]string, error)
	GetTrusts(sessionId string) (*[]int, error)

	CreateStep(scenariodId string, previousAnswer, previosStep *string, maxTrust, minTrust *int) (string, error)
	EditStep(stepId string, maxTrust, minTrust *int) error
	GetStep(stepId string) (*models.StepData, error)
	MatchActionToStep(actionIds *[]string, StepId string) error
	DeleteAllMatchActionToStep(stepId string) error
	GetNextStepByStepId(currentStepId string, currentTrust int) (*string, error)
	GetNextStepByAnswerId(answerId string, currentTrust int) (*string, error)

	CreateAction(actionType, messageId string, delay int) (string, error)
	EditAction(actionId, actionType, messageId string, delay *int) error
	GetActionById(actionId string) (*models.ActionData, error)

	CreateMessage(senderId, text *string, senderName string) (string, error)
	EditMessage(messageId string, senderId, senderName, text *string) error
	MatchAnswerToMessage(answerIds *[]string, messageId string) error
	DeleteAllMatchAnswerToMessage(messageId string) error
	DeleteAllMatchFileToMessage(messageId string) error
	MatchFileToMessage(filesIds *[]string, messageId string) error
	GetMessageById(messageId string) (*models.MessageData, error)

	CreateAnswer(addTrust int, errorText, text string) (string, error)
	EditAnswer(answerId string, errorText, text *string, addTrust *int) error
	GetAnswerById(answerId string) (*models.AnswerData, error)
	GetAnswersByMessageId(messageId string) (*[]models.AnswerData, error)
	RegisterTrustLevel(sessionId string, addTrust int) error

	CreateFile(filename string, isSafe bool, size int, fileError *string) (string, error)
	EditFile(fileId string, filename, fileError *string, isSafe *bool, size *int) error
	GetFilesByMessageId(messageId string) (*[]models.FileData, error)
	GetFileById(fileId string) (*models.FileData, error)
}

type simulatorService struct {
	log  *logrus.Logger
	repo SimulatorDataProvider
}

func New(
	log *logrus.Logger,
	repo SimulatorDataProvider,
) *simulatorService {
	return &simulatorService{
		log:  log,
		repo: repo,
	}
}
func (s *simulatorService) CreateSession(userId string, scenarioId string) error {
	scenario, err := s.repo.GetScenario(scenarioId)
	if err != nil {
		s.log.Error(err)
		return err
	}

	sessionId, err := s.repo.CreateSession(userId, scenario.FirstStep)
	if err != nil {
		s.log.Error(err)
		return err
	}

	return s.repo.SaveBeginTrustLevel(sessionId)
}
func (s *simulatorService) GetResults(userId string) (*simulator.SimulationFinal, error) {
	session, err := s.repo.GetCurrentSession(userId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	if err := s.repo.MarkSessionAsFinished(session.UUID); err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, err = s.repo.GetSessionBySessionId(session.UUID)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	userErrors, err := s.repo.GetErrors(session.UUID)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	trustGraph, err := s.repo.GetTrusts(session.UUID)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	duration := int(session.FinishedAt.Sub(session.CreatedAt).Minutes())

	return &simulator.SimulationFinal{
		Errors:       userErrors,
		GameDuration: &duration,
		TrustGraph:   trustGraph,
	}, nil

}

func (s *simulatorService) CreateScenario(data *simulator.ScenarioBaseRequired) (*simulator.ScenarioFull, error) {

	scenarioId, err := s.repo.CreateScenario(data.Title, data.Description, string(data.Difficulty))
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	strSlice := convert.AnyToStringSlice(data.ArticleIds, func(id types.UUID) string {
		return id.String()
	})
	if data.ArticleIds != nil {
		if err := s.repo.MatchArticleToScenario(strSlice, scenarioId); err != nil {
			s.log.Error(err)
			return nil, err
		}
	}
	return s.GetScenarioById(scenarioId)
}
func (s *simulatorService) EditScenario(scenarioId string, editedData *simulator.ScenarioBase) (*simulator.ScenarioFull, error) {

	if err := s.repo.EditScenario(scenarioId, editedData.Title, editedData.Description, (*string)(editedData.Difficulty)); err != nil {
		s.log.Error(err)
		return nil, err
	}

	if editedData.ArticleIds != nil {
		if err := s.repo.DeleteAllMatchArtilceToScenario(scenarioId); err != nil {
			s.log.Error(err)
			return nil, err
		}

		strSlice := convert.AnyToStringSlice(editedData.ArticleIds, func(id uuid.UUID) string {
			return id.String()
		})

		if err := s.repo.MatchArticleToScenario(strSlice, scenarioId); err != nil {
			s.log.Error(err)
			return nil, err
		}

	}
	return s.GetScenarioById(scenarioId)
}
func (s *simulatorService) GetAllScenarios() (*[]simulator.ScenarioFull, error) {
	scenarios, err := s.repo.GetAllScenarios()
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	return convert.AnyToAnySlice(scenarios, func(d models.ScenarioData) simulator.ScenarioFull {
		articleIdSlice := convert.StringToAnySlice(d.ArticleIds, func(s string) types.UUID {
			var articleUUID types.UUID
			if articleUUID.Scan(s) != nil {
				return types.UUID{}
			}
			return articleUUID
		})
		var scenarioUUID types.UUID
		scenarioUUID.Scan(d.UUID)

		return simulator.ScenarioFull{
			Id:          scenarioUUID,
			ArticleIds:  articleIdSlice,
			Description: d.Description,
			Difficulty:  simulator.ScenarioFullDifficulty(d.Difficulty),
			Title:       d.Title,
		}
	}), nil
}
func (s *simulatorService) GetScenarioById(scenarioId string) (*simulator.ScenarioFull, error) {
	scenario, err := s.repo.GetScenario(scenarioId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	var scenarioUUID types.UUID
	if scenarioUUID.Scan(scenario.UUID) != nil {
		s.log.Error(err)
		return nil, err
	}
	articleIdSlice := convert.StringToAnySlice(scenario.ArticleIds, func(s string) types.UUID {
		var uuid types.UUID
		if uuid.Scan(s) != nil {
			return types.UUID{}
		}
		return uuid
	})

	return &simulator.ScenarioFull{
		ArticleIds:  articleIdSlice,
		Description: scenario.Description,
		Difficulty:  simulator.ScenarioFullDifficulty(scenario.Difficulty),
		Id:          scenarioUUID,
		Title:       scenario.Title,
	}, nil
}

func (s *simulatorService) CreateStep(step *simulator.StepMetaRequired) (*simulator.StepMetaFull, error) {
	// Установка дефолтных значений
	var previousAnswer, previousStep *string
	if step.PreviousStep != nil {
		strId := step.PreviousStep.String()
		previousStep = &strId
	}
	if step.PreviousAnswer != nil {
		strId := step.PreviousAnswer.String()
		previousAnswer = &strId
	}

	// Создание шага
	stepId, err := s.repo.CreateStep(step.ScenarioId.String(), previousAnswer, previousStep, step.MaxTrust, step.MinTrust)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	// Привязка Action
	strSlice := convert.AnyToStringSlice(&step.Actions, func(id types.UUID) string {
		return id.String()
	})
	if err := s.repo.MatchActionToStep(strSlice, stepId); err != nil {
		s.log.Error(err)
		return nil, err
	}
	// Привязка первого шага к Scenario
	if step.PreviousAnswer == nil && step.PreviousStep == nil {
		s.repo.SetFirstStep(step.ScenarioId.String(), stepId)
	}

	// Получение созданного Step
	return s.GetStepMetaById(stepId)

}

func (s *simulatorService) EditStep(stepId string, editedStep *simulator.StepMetaMin) (*simulator.StepMetaFull, error) {
	// Редактирование Step
	if err := s.repo.EditStep(stepId, editedStep.MaxTrust, editedStep.MinTrust); err != nil {
		s.log.Error(err)
		return nil, err
	}

	// Обновление Actions (полная замена)
	if editedStep.Actions != nil {
		if err := s.repo.DeleteAllMatchActionToStep(stepId); err != nil {
			s.log.Error(err)
			return nil, err
		}
		strSlice := convert.AnyToStringSlice(editedStep.Actions, func(id uuid.UUID) string {
			return id.String()
		})

		if err := s.repo.MatchActionToStep(strSlice, stepId); err != nil {
			s.log.Error(err)
			return nil, err
		}
	}
	// Получение обновленного Step
	return s.GetStepMetaById(stepId)

}
func (s *simulatorService) GetStep(userId string) (*simulator.StepBase, error) {
	// Получаем сессию с  SessionId, CurrentTrust и CurrentStepId
	session, err := s.repo.GetCurrentSession(userId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	} else if session.CurrentStepId == nil {
		return nil, nil
	}

	// Получаем возможные Step
	nextStepId, err := s.repo.GetNextStepByStepId(*session.CurrentStepId, session.CurrentTrust)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	if err := s.repo.SetCurrentStep(session.UUID, nextStepId); err != nil {
		s.log.Error(err)
		return nil, err
	}

	//  Возвращаем текущий Step
	stepMeta, err := s.GetStepMetaById(*session.CurrentStepId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	return &simulator.StepBase{
		Actions: &stepMeta.Actions,
		Id:      stepMeta.Id,
	}, nil
}

func (s *simulatorService) CreateAction(action *simulator.ActionBaseRequired) (*simulator.ActionFull, error) {
	actionId, err := s.repo.CreateAction(string(action.Type), action.Message.String(), action.Delay)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	return s.GetActionById(actionId)
}
func (s *simulatorService) EditAction(actionId string, editedData *simulator.ActionBase) (*simulator.ActionFull, error) {
	var actionType string
	if editedData.Type != nil {
		actionType = string(*editedData.Type)
	}

	if err := s.repo.EditAction(actionId, actionType, editedData.Message.String(), editedData.Delay); err != nil {
		s.log.Error(err)
		return nil, err
	}

	return s.GetActionById(actionId)

}

func (s *simulatorService) CreateAnswer(answer *simulator.AnswerBaseRequired) (*simulator.AnswerFull, error) {
	var errorText string
	if answer.Error != nil {
		errorText = *answer.Error
	}

	answerId, err := s.repo.CreateAnswer(answer.AddTrust, errorText, answer.Text)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	answerData, err := s.repo.GetAnswerById(answerId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	var answerUUID types.UUID
	if answerUUID.Scan(answerData.UUID) != nil {
		s.log.Error(err)
		return nil, err
	}
	return &simulator.AnswerFull{
		AddTrust: answerData.AddTrust,
		Error:    answerData.Error,
		Id:       answerUUID,
		Text:     answerData.Text,
	}, nil
}
func (s *simulatorService) EditAnswer(answerId string, editedData *simulator.AnswerBase) (*simulator.AnswerFull, error) {

	if err := s.repo.EditAnswer(answerId, editedData.Error, editedData.Text, editedData.AddTrust); err != nil {
		s.log.Error(err)
		return nil, err
	}

	answerData, err := s.repo.GetAnswerById(answerId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	var answerUUID types.UUID
	if answerUUID.Scan(answerData.UUID) != nil {
		s.log.Error(err)
		return nil, err
	}
	return &simulator.AnswerFull{
		AddTrust: answerData.AddTrust,
		Error:    answerData.Error,
		Id:       answerUUID,
		Text:     answerData.Text,
	}, nil
}
func (s *simulatorService) SendUserAnswer(userId, answerId string) error {
	// Получаем SessionId
	session, err := s.repo.GetCurrentSession(userId)
	if err != nil {
		s.log.Error(err)
		return err
	}

	// Регистрируем nextStep
	nextStepId, err := s.repo.GetNextStepByAnswerId(answerId, session.CurrentTrust)
	if err != nil {
		s.log.Error(err)
		return err
	}

	if err := s.repo.SetCurrentStep(session.UUID, nextStepId); err != nil {
		s.log.Error(err)
		return err
	}

	// Получаем Answer и регистрируем ошибку если она есть
	answerData, err := s.repo.GetAnswerById(answerId)
	if err != nil {
		s.log.Error(err)
		return err
	}
	// Изменяем уровень trust
	if err := s.repo.RegisterTrustLevel(session.UUID, answerData.AddTrust); err != nil {
		s.log.Error(err)
		return err
	}

	if answerData.Error != nil && *answerData.Error != "" {
		if err := s.repo.SetSessionError(session.UUID, *answerData.Error); err != nil {
			s.log.Error(err)
			return err
		}
	}

	return nil
}

func (s *simulatorService) CreateFile(file *simulator.FileBaseRequired) (*simulator.FileFull, error) {
	var fileError string
	if file.Error != nil {
		fileError = *file.Error
	}

	fileId, err := s.repo.CreateFile(file.Filename, file.IsSafe, file.Size, &fileError)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	return s.getFileById(fileId)

}
func (s *simulatorService) EditFile(fileId string, editedData *simulator.FileBase) (*simulator.FileFull, error) {

	// Редактирование файла
	if err := s.repo.EditFile(fileId, editedData.Filename, editedData.Error, editedData.IsSafe, editedData.Size); err != nil {
		s.log.Error(err)
		return nil, err
	}

	// Получение обновленного файла
	return s.getFileById(fileId)
}
func (s *simulatorService) GetFileByFileId(fileId, userId string) (*simulator.FileFull, error) {
	// Получаем File
	file, err := s.getFileById(fileId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	// Если файл небезопасен - завершаем игру
	if !file.IsSafe {
		session, err := s.repo.GetCurrentSession(userId)
		if err != nil {
			s.log.Error(err)
			return nil, err
		}
		if err := s.repo.MarkSessionAsFinished(session.UUID); err != nil {
			s.log.Error(err)
			return nil, err
		}
	}

	if file.Error != nil && *file.Error != "" {
		session, err := s.repo.GetCurrentSession(userId)
		if err != nil {
			s.log.Error(err)
			return nil, err
		}
		if err := s.repo.SetSessionError(session.UUID, *file.Error); err != nil {
			s.log.Error(err)
			return nil, err
		}
	}
	// Возвращаем file
	return file, nil
}

func (s *simulatorService) CreateMessage(message *simulator.MessageBaseRequired) (*simulator.MessageFull, error) {
	var senderId *string
	if message.SenderId != nil {
		senderIdStr := message.SenderId.String()
		senderId = &senderIdStr
	}
	messageId, err := s.repo.CreateMessage(senderId, message.Text, message.SenderName)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	if message.Answers != nil {
		answerIdStrSlice := convert.AnyToStringSlice(message.Answers, func(a types.UUID) string {
			return a.String()
		})
		if err := s.repo.MatchAnswerToMessage(answerIdStrSlice, messageId); err != nil {
			s.log.Error(err)
			return nil, err
		}
	}

	if message.Files != nil {
		fileIdStrSlice := convert.AnyToStringSlice(message.Files, func(a types.UUID) string {
			return a.String()
		})

		if err := s.repo.MatchFileToMessage(fileIdStrSlice, messageId); err != nil {
			s.log.Error(err)
			return nil, err
		}
	}

	return s.GetMessageById(messageId)
}

func (s *simulatorService) EditMessage(messageId string, editedData *simulator.MessageBase) (*simulator.MessageFull, error) {
	var senderId *string
	if editedData.SenderId != nil {
		senderIdStr := editedData.SenderId.String()
		senderId = &senderIdStr
	}

	err := s.repo.EditMessage(messageId, senderId, editedData.SenderName, editedData.Text)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	if editedData.Answers != nil {
		if err := s.repo.DeleteAllMatchAnswerToMessage(messageId); err != nil {
			s.log.Error(err)
			return nil, err
		}
		answerIdStrSlice := convert.AnyToStringSlice(editedData.Answers, func(a types.UUID) string {
			return a.String()
		})
		if err := s.repo.MatchAnswerToMessage(answerIdStrSlice, messageId); err != nil {
			s.log.Error(err)
			return nil, err
		}
	}

	if editedData.Files != nil {
		if err := s.repo.DeleteAllMatchFileToMessage(messageId); err != nil {
			s.log.Error(err)
			return nil, err
		}
		fileIdStrSlice := convert.AnyToStringSlice(editedData.Files, func(a types.UUID) string {
			return a.String()
		})

		if err := s.repo.MatchFileToMessage(fileIdStrSlice, messageId); err != nil {
			s.log.Error(err)
			return nil, err
		}
	}

	return s.GetMessageById(messageId)
}

// Вспомогательные Getters
func (s *simulatorService) GetStepMetaById(stepId string) (*simulator.StepMetaFull, error) {
	stepData, err := s.repo.GetStep(stepId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	var stepUUID types.UUID
	if err := stepUUID.Scan(stepId); err != nil {
		s.log.Error(err)
		return nil, err
	}

	var previosAnswerUUID *types.UUID
	if stepData.PreviousAnswer != nil {
		if err := previosAnswerUUID.Scan(stepData.PreviousAnswer); err != nil {
			s.log.Error(err)
			return nil, err
		}
	}

	var previosStepUUID *types.UUID
	if stepData.PreviousStep != nil {
		if err := previosStepUUID.Scan(stepData.PreviousStep); err != nil {
			s.log.Error(err)
			return nil, err
		}
	}

	var scenarioUUID types.UUID
	if err := scenarioUUID.Scan(stepData.ScenarioId); err != nil {
		s.log.Error(err)
		return nil, err
	}

	var actions = make([]simulator.ActionFull, len(*stepData.ActionIds))
	for i, actionId := range *stepData.ActionIds {
		action, err := s.GetActionById(actionId)
		if err != nil {
			s.log.Error(err)
			return nil, err
		}
		actions[i] = *action
	}

	return &simulator.StepMetaFull{
		Id:             stepUUID,
		MaxTrust:       stepData.MaxTrust,
		MinTrust:       stepData.MinTrust,
		PreviousAnswer: previosAnswerUUID,
		PreviousStep:   previosStepUUID,
		ScenarioId:     scenarioUUID,
		Actions:        actions,
	}, nil
}

func (s *simulatorService) GetActionById(actionId string) (*simulator.ActionFull, error) {
	actionData, err := s.repo.GetActionById(actionId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	var actionUUID types.UUID
	if err := actionUUID.Scan(actionData.UUID); err != nil {
		s.log.Error(err)
		return nil, err
	}

	message, err := s.GetMessageById(actionData.MessageId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	return &simulator.ActionFull{
		Delay:   actionData.Delay,
		Id:      actionUUID,
		Message: *message,
		Type:    simulator.ActionFullType(actionData.Type),
	}, nil
}

func (s *simulatorService) GetMessageById(messageId string) (*simulator.MessageFull, error) {
	messageData, err := s.repo.GetMessageById(messageId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	var messageUUID types.UUID
	if err := messageUUID.Scan(messageData.UUID); err != nil {
		s.log.Error(err)
		return nil, err
	}

	var senderUUID types.UUID
	if err := senderUUID.Scan(messageData.UUID); err != nil {
		s.log.Error(err)
		return nil, err
	}

	answers, err := s.repo.GetAnswersByMessageId(messageId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	answerFullSlice := convert.AnyToAnySlice(answers, func(a models.AnswerData) simulator.AnswerFull {
		var answerUUID types.UUID
		if answerUUID.Scan(a.UUID) != nil {
			answerUUID = types.UUID{}
		}

		return simulator.AnswerFull{
			Id:       answerUUID,
			AddTrust: a.AddTrust,
			Error:    a.Error,
			Text:     a.Text,
		}
	})

	files, err := s.repo.GetFilesByMessageId(messageId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	fileFullSlice := convert.AnyToAnySlice(files, func(f models.FileData) simulator.FileFull {
		var FileUUID types.UUID
		if FileUUID.Scan(f.UUID) != nil {
			FileUUID = types.UUID{}
		}

		return simulator.FileFull{
			Id:       FileUUID,
			Filename: f.Filename,
			Error:    f.Error,
			Size:     f.Size,
			IsSafe:   f.IsSafe,
		}
	})

	return &simulator.MessageFull{
		Id:         messageUUID,
		Text:       messageData.Text,
		SenderId:   &senderUUID,
		SenderName: messageData.SenderName,
		Answers:    answerFullSlice,
		Files:      fileFullSlice,
	}, nil
}

func (s *simulatorService) getFileById(fileId string) (*simulator.FileFull, error) {
	fileData, err := s.repo.GetFileById(fileId)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	var fileUUID types.UUID
	if err := fileUUID.Scan(fileData.UUID); err != nil {
		s.log.Error(err)
		return nil, err
	}

	return &simulator.FileFull{
		Error:    fileData.Error,
		Filename: fileData.Filename,
		IsSafe:   fileData.IsSafe,
		Size:     fileData.Size,
		Id:       fileUUID,
	}, nil
}
