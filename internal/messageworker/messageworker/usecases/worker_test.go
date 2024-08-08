package usecases_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/Vilinvil/task_messaggio/internal/messageworker/messageworker/mocks"
	"github.com/Vilinvil/task_messaggio/internal/messageworker/messageworker/repository"
	"github.com/Vilinvil/task_messaggio/internal/messageworker/messageworker/usecases"
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"go.uber.org/mock/gomock"
)

func NewMessageWorker(
	ctrl *gomock.Controller,
	behaviorMessageRepository func(mock *mocks.MockMessageRepository),
	behaviorBrokerMessage func(mock *mocks.MockBrokerMessage),
	logger *mylogger.MyLogger,
) *usecases.MessageWorker {
	mockBrokerMessage := mocks.NewMockBrokerMessage(ctrl)
	mockMessageRepository := mocks.NewMockMessageRepository(ctrl)

	behaviorMessageRepository(mockMessageRepository)
	behaviorBrokerMessage(mockBrokerMessage)

	return usecases.NewMessageWorker(mockBrokerMessage, mockMessageRepository, logger)
}

func TestJobMessageErrWhenHandleConsumedMessage(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	chConsumptionMessages := make(chan repository.MessagePayloadWithCommitFunc)
	chErrConsumptionMessage := make(chan error)

	defer close(chConsumptionMessages)
	defer close(chErrConsumptionMessage)

	behaviorMessageRepository := func(_ *mocks.MockMessageRepository) {}
	behaviorBrokerMessage := func(mock *mocks.MockBrokerMessage) {
		mock.EXPECT().StartConsumption(gomock.Any()).Return(chConsumptionMessages, chErrConsumptionMessage)
	}

	errHandler := fmt.Errorf("my err handle") //nolint
	handlerFunc := func(_ context.Context, _ *models.MessagePayload) error {
		return errHandler
	}

	expectedErr := errHandler

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageWorker := NewMessageWorker(ctrl, behaviorMessageRepository, behaviorBrokerMessage, nopLogger)

	chErr := messageWorker.JobMessages(context.Background(), handlerFunc)

	chConsumptionMessages <- repository.MessagePayloadWithCommitFunc{
		MessagePayload: *models.NewMessagePayload("", models.DummyGeneratorUUID),
		CommitFunc: func() error {
			return nil
		},
	}

	if receivedErr := <-chErr; !errors.Is(receivedErr, expectedErr) {
		t.Errorf("Not equal receivedeErr: %s expectedErr: %s", receivedErr.Error(), expectedErr.Error())
	}
}

func TestJobMessageErrMessageRepositorySetStatusMessage(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	chConsumptionMessages := make(chan repository.MessagePayloadWithCommitFunc)
	chErrConsumptionMessage := make(chan error)

	defer close(chConsumptionMessages)
	defer close(chErrConsumptionMessage)

	errSetStatusMessage := fmt.Errorf("internal error set status") //nolint

	behaviorMessageRepository := func(mock *mocks.MockMessageRepository) {
		mock.EXPECT().
			SetStatusMessage(gomock.Any(), &models.DummyUUID, repository.StatusMessageDone).Return(errSetStatusMessage)
	}
	behaviorBrokerMessage := func(mock *mocks.MockBrokerMessage) {
		mock.EXPECT().StartConsumption(gomock.Any()).Return(chConsumptionMessages, chErrConsumptionMessage)
	}

	handlerFunc := func(_ context.Context, _ *models.MessagePayload) error {
		return nil
	}

	expectedErr := errSetStatusMessage

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageWorker := NewMessageWorker(ctrl, behaviorMessageRepository, behaviorBrokerMessage, nopLogger)

	chErr := messageWorker.JobMessages(context.Background(), handlerFunc)

	chConsumptionMessages <- repository.MessagePayloadWithCommitFunc{
		MessagePayload: *models.NewMessagePayload("", models.DummyGeneratorUUID),
		CommitFunc: func() error {
			return nil
		},
	}

	if receivedErr := <-chErr; !errors.Is(receivedErr, expectedErr) {
		t.Errorf("Not equal receivedeErr: %v expectedErr: %v", receivedErr, expectedErr)
	}
}

func TestJobMessageErrCommitFunc(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	chConsumptionMessages := make(chan repository.MessagePayloadWithCommitFunc)
	chErrConsumptionMessage := make(chan error)

	defer close(chConsumptionMessages)
	defer close(chErrConsumptionMessage)

	behaviorMessageRepository := func(mock *mocks.MockMessageRepository) {
		mock.EXPECT().SetStatusMessage(gomock.Any(), &models.DummyUUID, repository.StatusMessageDone).Return(nil)
	}
	behaviorBrokerMessage := func(mock *mocks.MockBrokerMessage) {
		mock.EXPECT().StartConsumption(gomock.Any()).Return(chConsumptionMessages, chErrConsumptionMessage)
	}

	handlerFunc := func(_ context.Context, _ *models.MessagePayload) error { return nil }

	errCommitFunc := fmt.Errorf("internal error commit func") //nolint
	msgPayloadWithCommitFunc := repository.MessagePayloadWithCommitFunc{
		MessagePayload: *models.NewMessagePayload("", models.DummyGeneratorUUID),
		CommitFunc: func() error {
			return errCommitFunc
		},
	}
	expectedErr := errCommitFunc

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageWorker := NewMessageWorker(ctrl, behaviorMessageRepository, behaviorBrokerMessage, nopLogger)

	chErr := messageWorker.JobMessages(context.Background(), handlerFunc)

	chConsumptionMessages <- msgPayloadWithCommitFunc

	if receivedErr := <-chErr; !errors.Is(receivedErr, expectedErr) {
		t.Errorf("Not equal receivedeErr: %v expectedErr: %v", receivedErr, expectedErr)
	}
}

func TestJobMessageErrInChErrConsumptionMessage(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	chConsumptionMessages := make(chan repository.MessagePayloadWithCommitFunc)
	chErrConsumptionMessage := make(chan error)

	defer close(chErrConsumptionMessage)

	behaviorMessageRepository := func(_ *mocks.MockMessageRepository) {}
	behaviorBrokerMessage := func(mock *mocks.MockBrokerMessage) {
		mock.EXPECT().StartConsumption(gomock.Any()).Return(chConsumptionMessages, chErrConsumptionMessage)
	}

	handlerFunc := func(_ context.Context, _ *models.MessagePayload) error {
		return nil
	}

	expectedErr := fmt.Errorf("internal err in chErrConsumptionMessage") //nolint

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageWorker := NewMessageWorker(ctrl, behaviorMessageRepository, behaviorBrokerMessage, nopLogger)

	chErr := messageWorker.JobMessages(context.Background(), handlerFunc)

	close(chConsumptionMessages)
	chErrConsumptionMessage <- expectedErr

	if receivedErr := <-chErr; !errors.Is(receivedErr, expectedErr) {
		t.Errorf("Not equal receivedeErr: %v expectedErr: %v", receivedErr, expectedErr)
	}
}

func multipleWriteInChConsumptionMessage(
	amountWrites int, chConsumptionMessage chan<- repository.MessagePayloadWithCommitFunc,
) {
	msgPayloadWithCommitFunc := repository.MessagePayloadWithCommitFunc{
		MessagePayload: *models.NewMessagePayload("", models.DummyGeneratorUUID),
		CommitFunc: func() error {
			return nil
		},
	}

	for range amountWrites {
		chConsumptionMessage <- msgPayloadWithCommitFunc
	}

	close(chConsumptionMessage)
}

func TestJobMessageMultipleWriteInConsumption(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	chConsumptionMessages := make(chan repository.MessagePayloadWithCommitFunc)
	chErrConsumptionMessage := make(chan error)

	defer close(chErrConsumptionMessage)

	amountWritesConsumption := 100

	behaviorMessageRepository := func(mock *mocks.MockMessageRepository) {
		mock.EXPECT().SetStatusMessage(gomock.Any(), &models.DummyUUID, repository.StatusMessageDone).Return(nil).
			Times(amountWritesConsumption)
	}
	behaviorBrokerMessage := func(mock *mocks.MockBrokerMessage) {
		mock.EXPECT().StartConsumption(gomock.Any()).Return(chConsumptionMessages, chErrConsumptionMessage)
	}

	handlerFunc := func(_ context.Context, _ *models.MessagePayload) error {
		return nil
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageWorker := NewMessageWorker(ctrl, behaviorMessageRepository, behaviorBrokerMessage, nopLogger)

	chErr := messageWorker.JobMessages(context.Background(), handlerFunc)

	multipleWriteInChConsumptionMessage(amountWritesConsumption, chConsumptionMessages)

	chErrConsumptionMessage <- nil

	if receivedErr := <-chErr; !errors.Is(receivedErr, nil) {
		t.Errorf("Not equal receivedeErr: %v expected: nil", receivedErr)
	}
}

func multipleCallJobMessage(amountCall int, messageWorker *usecases.MessageWorker) error {
	chErr := make(chan error)

	handlerFunc := func(_ context.Context, _ *models.MessagePayload) error { return nil }

	ctx := context.Background()

	expectedChErr := messageWorker.JobMessages(ctx, handlerFunc)

	amountCall--

	wgJobMessages := sync.WaitGroup{}
	wgJobMessages.Add(amountCall)

	for range amountCall {
		go func() {
			defer wgJobMessages.Done()

			receivedChErr := messageWorker.JobMessages(ctx, handlerFunc)
			if receivedChErr != expectedChErr {
				chErr <- fmt.Errorf("receivedChErr and expectedChErr are different") //nolint
			}
		}()
	}

	go func() {
		wgJobMessages.Wait()
		close(chErr)
	}()

	return <-chErr
}

func TestJobMessageMultipleCallJobMessage(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	chConsumptionMessages := make(chan repository.MessagePayloadWithCommitFunc)
	chErrConsumptionMessage := make(chan error)

	defer close(chConsumptionMessages)
	defer close(chErrConsumptionMessage)

	amountCallJobMessage := 100

	behaviorMessageRepository := func(_ *mocks.MockMessageRepository) {}
	behaviorBrokerMessage := func(mock *mocks.MockBrokerMessage) {
		mock.EXPECT().StartConsumption(gomock.Any()).Return(chConsumptionMessages, chErrConsumptionMessage)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageWorker := NewMessageWorker(ctrl, behaviorMessageRepository, behaviorBrokerMessage, nopLogger)

	receivedErr := multipleCallJobMessage(amountCallJobMessage, messageWorker)
	if receivedErr != nil {
		t.Errorf("receivedErr not equal nil: %v", receivedErr)
	}
}
