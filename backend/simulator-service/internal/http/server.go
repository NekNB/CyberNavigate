package http

import (
	"encoding/json"
	"errors"

	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/models"
	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/storage"
	"github.com/NekNB/CyberNavigate/swagger/gen/simulator"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

// Здесь реализуем все методы обработки APIServer

// Проверяем, соответствует ли APIServer сгенерированному ServerInterface
var _ simulator.ServerInterface = (*APIServer)(nil)

type APIServer struct {
	log              *logrus.Logger
	simulatorService SimulatorServiceInterface
}

type SimulatorServiceInterface interface {
	CreateSession(userId string) error
	GetResults(userId string) (*simulator.SimulationFinal, error)

	CreateScenario(data *simulator.ScenarioBaseRequired) (*simulator.ScenarioFull, error)
	EditScenario(scenarioId string, editedData *simulator.ScenarioBase) (*simulator.ScenarioFull, error)
	GetAllScenarios() (*[]simulator.ScenarioFull, error)
	GetScenarioById(scenarioId string) (*simulator.ScenarioFull, error)

	CreateStep(step *simulator.StepMetaRequired) (models.StepMetaResponse, error)
	EditStep(stepId string, editedData *simulator.StepMeta) (*simulator.StepMetaFull, error)
	GetStep(userId string) (*simulator.StepBase, error)

	CreateAction(action *simulator.ActionBaseRequired) (*simulator.ActionFull, error)
	EditAction(actionId string, editedData *simulator.ActionBase) (*simulator.ActionFull, error)

	CreateAnswer(answer *simulator.AnswerBaseRequired) (*simulator.AnswerFull, error)
	EditAnswer(answerId string, editedData *simulator.AnswerBase) (*simulator.AnswerFull, error)
	SendUserAnswer(answerId string) error

	CreateFile(file *simulator.FileBaseRequired) (*simulator.FileFull, error)
	EditFile(fileId string, editedData *simulator.FileBase) (*simulator.FileFull, error)
	GetFileByFileId(fileId, userId string) (*simulator.FileFull, error)

	CreateMessage(Message *simulator.MessageBaseRequired) (*simulator.MessageFull, error)
	EditMessage(messageId string, editedData *simulator.MessageBase) (*simulator.MessageFull, error)
}

func New(log *logrus.Logger, simulatorService SimulatorServiceInterface) *APIServer {

	return &APIServer{log: log, simulatorService: simulatorService}
}

func (a *APIServer) CreateScenario(c fiber.Ctx) error {
	// Парсим Body
	req := &simulator.ScenarioBaseRequired{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	scenario, err := a.simulatorService.CreateScenario(req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.Status(400).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(201).JSON(scenario)
}

func (a *APIServer) GetScenarioById(c fiber.Ctx, scenarioId string) error {

	file, err := a.simulatorService.GetScenarioById(scenarioId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(404).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(200).JSON(file)
}

func (a *APIServer) CreateSimulatorSession(c fiber.Ctx, params simulator.CreateSimulatorSessionParams) error {

	if err := a.simulatorService.CreateSession(params.XUserId); err != nil {
		a.log.Error(err)
		c.SendStatus(500)
	}

	return c.SendStatus(204)
}

func (a *APIServer) GetResults(c fiber.Ctx, params simulator.GetResultsParams) error {
	final, err := a.simulatorService.GetResults(params.XUserId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(400).JSON(simulator.ErrorResponse{Message: err.Error()})
		}

		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(200).JSON(final)
}

func (a *APIServer) CreateStep(c fiber.Ctx) error {
	req := &simulator.StepMetaRequired{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	stepMeta, err := a.simulatorService.CreateStep(req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.Status(400).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(201).JSON(stepMeta)
}

func (a *APIServer) EditStep(c fiber.Ctx, stepId string) error {
	// Парсим Body
	req := &simulator.StepMeta{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	step, err := a.simulatorService.EditStep(stepId, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.Status(400).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(404).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(200).JSON(step)
}

func (a *APIServer) GetStep(c fiber.Ctx, params simulator.GetStepParams) error {
	step, err := a.simulatorService.GetStep(params.XUserId.String())
	if err != nil {
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(200).JSON(step)
}

func (a *APIServer) EditScenario(c fiber.Ctx, scenarioId string) error {
	// Парсим Body
	req := &simulator.EditScenarioJSONRequestBody{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	scenario, err := a.simulatorService.EditScenario(scenarioId, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.Status(400).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(404).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(201).JSON(scenario)
}

func (a *APIServer) GetAllScenarios(c fiber.Ctx) error {
	scenarios, err := a.simulatorService.GetAllScenarios()
	if err != nil {
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(200).JSON(scenarios)
}

func (a *APIServer) CreateAction(c fiber.Ctx) error {
	req := &simulator.ActionBaseRequired{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	action, err := a.simulatorService.CreateAction(req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.Status(400).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(201).JSON(action)
}
func (a *APIServer) EditAction(c fiber.Ctx, actionId string) error {
	// Парсим Body
	req := &simulator.ActionBase{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	action, err := a.simulatorService.EditAction(actionId, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.Status(400).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(404).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(200).JSON(action)
}

func (a *APIServer) CreateAnswer(c fiber.Ctx) error {
	req := &simulator.AnswerBaseRequired{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	answer, err := a.simulatorService.CreateAnswer(req)
	if err != nil {
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(201).JSON(answer)
}
func (a *APIServer) EditAnswer(c fiber.Ctx, answerId string) error {
	// Парсим Body
	req := &simulator.AnswerBase{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	Answer, err := a.simulatorService.EditAnswer(answerId, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(404).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(200).JSON(Answer)
}
func (a *APIServer) SendUserAnswer(c fiber.Ctx, params simulator.SendUserAnswerParams) error {

	req := &simulator.SendUserAnswerJSONBody{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	if err := a.simulatorService.SendUserAnswer(req.AnswerId.String()); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.Status(400).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}
	return c.SendStatus(204)
}

func (a *APIServer) CreateFile(c fiber.Ctx) error {
	req := &simulator.FileBaseRequired{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	file, err := a.simulatorService.CreateFile(req)
	if err != nil {
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(201).JSON(file)
}

func (a *APIServer) EditFile(c fiber.Ctx, fileId string) error {
	// Парсим Body
	req := &simulator.FileBase{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	file, err := a.simulatorService.EditFile(fileId, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(404).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(200).JSON(file)
}

func (a *APIServer) GetFileByFileId(c fiber.Ctx, fileId string, params simulator.GetFileByFileIdParams) error {

	file, err := a.simulatorService.GetFileByFileId(fileId, params.XUserId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(404).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(200).JSON(file)
}

func (a *APIServer) CreateMessage(c fiber.Ctx) error {
	req := &simulator.MessageBaseRequired{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	Message, err := a.simulatorService.CreateMessage(req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.Status(400).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		if errors.Is(err, storage.ErrDataConfllict) {
			c.Status(409).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(201).JSON(Message)
}
func (a *APIServer) EditMessage(c fiber.Ctx, messageId string) error {
	// Парсим Body
	req := &simulator.MessageBase{}
	if json.Unmarshal(c.Body(), req) != nil {
		return c.SendStatus(422)
	}

	message, err := a.simulatorService.EditMessage(messageId, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.Status(400).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		if errors.Is(err, storage.ErrNotFound) {
			return c.Status(404).JSON(simulator.ErrorResponse{Message: err.Error()})
		}
		a.log.Error(err)
		return c.SendStatus(500)
	}

	return c.Status(200).JSON(message)
}
