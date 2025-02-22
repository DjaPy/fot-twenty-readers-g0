package command

import (
	"context"
	"github.com/DjaPy/fot-twenty-readers-go/decorator"
	"github.com/DjaPy/fot-twenty-readers-go/src/adapters"
	"github.com/DjaPy/fot-twenty-readers-go/src/domain"
	"github.com/sirupsen/logrus"
)

type CreatePsalmReaderTG struct {
	Username   string
	TelegramID int64
	Phone      string
}

type CreateCalendarOfReaderHandler decorator.CommandHandler[CreatePsalmReaderTG]

type createPsalmReaderTGHandler struct {
	repo                 adapters.PsalmReaderTGRepository
	psalmReaderTGService PsalmReaderTGService
}

func NewCreatePsalmReaderTGHandler(
	repo adapters.PsalmReaderTGRepository,
	psalmReaderTGService PsalmReaderTGService,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) decorator.CommandHandler[CreatePsalmReaderTG] {
	if repo == nil {
		panic("nil repo")
	}
	if userService == nil {
		panic("nil userService")
	}
	if trainerService == nil {
		panic("nil trainerService")
	}

	return decorator.ApplyCommandDecorators[CreatePsalmReaderTG](
		createPsalmReaderTGHandler{repo, psalmReaderTGService},
		logger,
		metricsClient,
	)
}

func (cpr createPsalmReaderTGHandler) Handle(ctx context.Context, cmd CreatePsalmReaderTG) error {
	defer func() {
		//logs.LogCommandExecution("ApproveTrainingReschedule", cmd, err)
	}()

	prTG, err := domain.NewPsalmReader(cmd.Username, cmd.TelegramID, cmd.Phone)
	if err != nil {
		return err
	}

	err = cpr.repo.CreatePsalmReaderTG(ctx, prTG)
	if err != nil {
		return err
	}
	return nil
}
