package repository

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type MessagePayloadWithCommitFunc struct {
	models.MessagePayload
	CommitFunc func() error
}

func newCommitFunc(generation *kafka.Generation, partition int, offset int64) func() error {
	return func() error {
		return generation.CommitOffsets(map[string]map[int]int64{TopicName: {partition: offset}})
	}
}

// BrokerMessageKafka uses a long-lived connection to Kafka. The user of BrokerMessageKafka must
// call method BrokerMessageKafka.Close() to gracefully shutdown the work.
type BrokerMessageKafka struct {
	brokerAddr               string
	consumerGroup            *kafka.ConsumerGroup
	logger                   *mylogger.MyLogger
	chConsumptionMessages    chan MessagePayloadWithCommitFunc
	chErrConsumptionMessages chan error

	// protect channels chConsumptionMessages, chErrConsumptionMessages
	onceChannels *sync.Once
}

func NewBrokerMessageKafka(brokerAddr string, logger *mylogger.MyLogger) (*BrokerMessageKafka, error) {
	preBroker := &BrokerMessageKafka{ //nolint:exhaustruct
		brokerAddr:   brokerAddr,
		logger:       logger,
		onceChannels: &sync.Once{},
	}

	err := preBroker.initConsumerGroup()
	if err != nil {
		return nil, err
	}

	return preBroker, nil
}

const (
	ConsumerGroupID = "group_consumer_1"
	TopicName       = "messages"
)

func (b *BrokerMessageKafka) initConsumerGroup() error {
	var err error

	b.consumerGroup, err = kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      ConsumerGroupID,
		Brokers: []string{b.brokerAddr},
		Topics:  []string{TopicName},
	})
	if err != nil {
		b.logger.Error(err)

		return err
	}

	return nil
}

func (b *BrokerMessageKafka) Close() error {
	err := b.consumerGroup.Close()
	if err != nil {
		b.logger.Error(err)

		return err
	}

	return nil
}

func (b *BrokerMessageKafka) startConsumptionGeneration(ctx context.Context,
	generation *kafka.Generation, partition int, offset int64,
) {
	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:   []string{b.brokerAddr},
			Topic:     TopicName,
			Partition: partition,
		})
	defer func() {
		err := reader.Close()
		if err != nil {
			b.logger.Error(err)
			close(b.chConsumptionMessages)
			b.chErrConsumptionMessages <- err
		}
	}()

	err := reader.SetOffset(offset)
	if err != nil {
		b.logger.Error(err)
		close(b.chConsumptionMessages)
		b.chErrConsumptionMessages <- err

		return
	}

	for {
		b.logger.Debugln("wait read message")

		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			b.logger.Error(err)
			close(b.chConsumptionMessages)
			b.chErrConsumptionMessages <- err

			return
		}

		b.logger.Debugf("received message %s/%d/%d : %s\n", msg.Topic,
			msg.Partition, msg.Offset, string(msg.Value))

		var msgKey uuid.UUID

		copy(msgKey[:], msg.Key[:len(msgKey)])

		b.chConsumptionMessages <- MessagePayloadWithCommitFunc{
			MessagePayload: models.MessagePayload{
				ID:    msgKey,
				Value: string(msg.Value),
			},
			CommitFunc: newCommitFunc(generation, partition, msg.Offset+1),
		}
	}
}

// StartConsumption one time starts consuming messages from Kafka. At the next call of this method only returning
// existing channels happens. User of StartConsumption may commit offset of consumed messages by call
// MessagePayloadWithCommitFunc.CommitFunc. It may do after reading every message, after reading a batch of messages,
// or after time plus reading message (user self-control this process).
//
// In case error channel chConsumptionMesses will be close, and err will be sent to chErrConsumptionMessage.
func (b *BrokerMessageKafka) StartConsumption(ctx context.Context) ( //nolint:nonamedreturns
	chConsumptionMessages <-chan MessagePayloadWithCommitFunc, chErrConsumptionMessage <-chan error,
) {
	if b.chConsumptionMessages != nil || b.chErrConsumptionMessages != nil {
		return b.chConsumptionMessages, b.chErrConsumptionMessages
	}

	wgCreationChannels := sync.WaitGroup{}

	wgCreationChannels.Add(1)

	b.onceChannels.Do(func() {
		b.chConsumptionMessages = make(chan MessagePayloadWithCommitFunc)
		b.chErrConsumptionMessages = make(chan error)

		wgCreationChannels.Done()

		go func() {
			for {
				gen, err := b.consumerGroup.Next(ctx)
				if err != nil {
					if errors.Is(err, io.EOF) {
						continue
					}

					b.logger.Debugf("%+v", b.consumerGroup)
					b.logger.Error(err)
					close(b.chConsumptionMessages)
					b.chErrConsumptionMessages <- err

					return
				}

				assignments := gen.Assignments[TopicName]
				for _, assignment := range assignments {
					partition, offset := assignment.ID, assignment.Offset

					b.logger.Debugf("partition/offset: %d/%d\n", partition, offset)

					gen.Start(func(ctx context.Context) {
						b.startConsumptionGeneration(ctx, gen, partition, offset)
					})
				}
			}
		}()
	})

	wgCreationChannels.Wait()

	return b.chConsumptionMessages, b.chErrConsumptionMessages
}
