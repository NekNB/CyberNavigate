package simulator

import (
	"context"
	"errors"

	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/http"
	"github.com/NekNB/CyberNavigate/backend/simulator-service/internal/storage"
	"github.com/sirupsen/logrus"
)

var _ http.SimulatorServiceInterface = (*simulatorService)(nil)

type simulatorContentProvider interface {
	simulatorTextById(ctx context.Context, id string) (string, error)
	SaveText(ctx context.Context, simulatorText string) (string, error)
	UpdateText(ctx context.Context, id string, simulatorText string) error
}

type simulatorMetaProvider interface {
	simulators() (*[]simulator.simulatorMetaData, error)
	simulatorByUUID(simulatorUUID string) (*simulator.simulatorMetaData, error)
	simulatorTextIDByUUID(simulatorUUID string) (textID string, err error)
	Createsimulator(simulatorName *string) (*simulator.simulatorMetaData, error)
	UpdatesimulatorByUUID(uuid, title, textID, status, videoUrl string) (*simulator.simulatorMetaData, error)
}

type simulatorService struct {
	log                   *logrus.Logger
	simulatorDataProvider simulatorContentProvider
}

func New(
	log *logrus.Logger,
	simulatorContentProvider simulatorContentProvider,
) *simulatorService {
	return &simulatorService{
		log:                   log,
		simulatorDataProvider: simulatorContentProvider,
	}
}

// Возвращает списко MetaData статей
func (a *simulatorService) simulators() (*[]simulator.simulatorMetaData, error) {
	metadata, err := a.simulatorMetaProvider.simulators()
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	return metadata, nil
}

func (a *simulatorService) simulatorByUUID(simulatorId string) (*simulator.simulatorMetaData, error) {
	metadata, err := a.simulatorMetaProvider.simulatorByUUID(simulatorId)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	return metadata, nil
}

func (a *simulatorService) simulatorTextByUUID(ctx context.Context, simulatorId string) (string, error) {
	mP := a.simulatorMetaProvider
	cP := a.simulatorDataProvider

	textId, err := mP.simulatorTextIDByUUID(simulatorId)
	if err != nil {
		if !errors.Is(err, storage.ErrsimulatorTextNotCreatedYet) {
			a.log.Error(err)
		}

		return "", err
	}

	text, err := cP.simulatorTextById(ctx, textId)
	if err != nil {
		a.log.Error(err)
		return text, err
	}

	return text, nil
}

func (a *simulatorService) Createsimulator(title *string) (*simulator.simulatorMetaData, error) {
	metadata, err := a.simulatorMetaProvider.Createsimulator(title)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	return metadata, err
}

func (a *simulatorService) SavesimulatorTextByUUID(ctx context.Context, simulatorId, text string) (*simulator.simulatorMetaData, error) {
	mP := a.simulatorMetaProvider
	cP := a.simulatorDataProvider

	textId, err := cP.SaveText(ctx, text)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}

	metadata, err := mP.UpdatesimulatorByUUID(
		simulatorId,
		"",
		textId,
		"",
		"",
	)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	return metadata, nil
}

func (a *simulatorService) UpdatesimulatorByUUID(simulatorId, text, title, videoURL, status string) (*simulator.simulatorMetaData, error) {
	metadata, err := a.simulatorMetaProvider.UpdatesimulatorByUUID(
		simulatorId,
		text,
		title,
		videoURL,
		status,
	)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	return metadata, err
}

func (a *simulatorService) UpdatesimulatorTextByUUID(ctx context.Context, simulatorId, text string) (*simulator.simulatorMetaData, error) {
	mP := a.simulatorMetaProvider
	cP := a.simulatorDataProvider

	textId, err := mP.simulatorTextIDByUUID(simulatorId)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	a.log.Debug(textId)
	if err := cP.UpdateText(ctx, textId, text); err != nil {
		a.log.Error(err)
		return nil, err
	}

	metadata, err := mP.simulatorByUUID(
		simulatorId,
	)
	if err != nil {
		a.log.Error(err)
		return nil, err
	}
	return metadata, nil
}
