package delivery

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/Vilinvil/task_messaggio/internal/message/message/usecases"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/Vilinvil/task_messaggio/pkg/responses"

	"github.com/google/uuid"
)

var _ MessageService = (*usecases.MessageService)(nil)

type MessageService interface {
	//GetMessagesByID(id ...int) ([]*models.Message, error)
	//ChangeStatusMessagesByID(status string, id ...int) error
	AddMessage(ctx context.Context, value string) (messageUUID uuid.UUID, err error)
}

type MessageHandler struct {
	service MessageService
	logger  *mylogger.MyLogger
}

func NewMessageHandler(service MessageService, logger *mylogger.MyLogger) *MessageHandler {
	return &MessageHandler{
		service: service,
		logger:  logger,
	}
}

var ResponseMessageAddSuccessful = responses.NewResponseSuccessful( //nolint:gochecknoglobals
	"Сообщение успешно добавлено")

// AddMessage godoc
//
//	@Summary  добавить сообщение в систему
//	@Tags message
//	@Produce    json
//	@Accept     x-www-form-urlencoded
//
// @Param orderChanges  body string true  "текст сообщения"
//
//	@Success    200  {object} responses.ResponseSuccessful
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    400  {object} myerrors.Error
//	@Router      /message [post]
func (m *MessageHandler) AddMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := m.logger.EnrichReqID(ctx)

	rawValueMessage, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		responses.SendErrResponse(w, logger, err)

		return
	}

	valueMessage, err := url.QueryUnescape(string(rawValueMessage))
	if err != nil {
		logger.Error(err)
		responses.SendErrResponse(w, logger, err)

		return
	}

	messageUUID, err := m.service.AddMessage(ctx, valueMessage)
	if err != nil {
		responses.SendErrResponse(w, logger, err)

		return
	}

	logger.Debugf("Message added: %s", messageUUID)

	responses.SendResponse(w, logger, ResponseMessageAddSuccessful)
}
