package models

import (
	"time"
	"unicode/utf8"

	"github.com/Vilinvil/task_messaggio/pkg/myerrors"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	Value     string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrLenMessage = myerrors.NewBadRequestError("Длина сообщения должна быть больше 0 и меньше 4000 символов utf8")

type PreMessage struct {
	ID    uuid.UUID
	Value string
}

func NewPreMessage(value string) *PreMessage {
	return &PreMessage{
		ID:    uuid.New(),
		Value: value,
	}
}

const maxLenValueMessage = 4000

func (p *PreMessage) Validate() error {
	if len(p.Value) == 0 || utf8.RuneCountInString(p.Value) > maxLenValueMessage {
		return ErrLenMessage
	}

	return nil
}
