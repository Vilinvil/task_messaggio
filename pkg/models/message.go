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

type GeneratorUUID func() uuid.UUID

var (
	DummyUUID          = uuid.MustParse("52fdfc07-2182-454f-963f-5f0f9a621d72") //nolint:gochecknoglobals
	DummyGeneratorUUID = func() uuid.UUID {                                     //nolint:gochecknoglobals
		return DummyUUID
	}
)

type MessagePayload struct {
	ID    uuid.UUID
	Value string
}

func NewMessagePayload(value string, generatorUUID GeneratorUUID) *MessagePayload {
	return &MessagePayload{
		ID:    generatorUUID(),
		Value: value,
	}
}

const maxLenValueMessage = 4000

func (p *MessagePayload) Validate() error {
	if len(p.Value) == 0 || utf8.RuneCountInString(p.Value) > maxLenValueMessage {
		return ErrLenMessage
	}

	return nil
}
