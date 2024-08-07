package repository

import (
	"context"

	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/segmentio/kafka-go"
)

const (
	TopicMessageName  = "messages"
	NumPartitions     = 1
	ReplicationFactor = 1
)

// BrokerMessageKafka uses a long-lived connection to Kafka. The user of BrokerMessageKafka must
// call method BrokerMessageKafka.Close() to gracefully shutdown the work.
type BrokerMessageKafka struct {
	conn   *kafka.Conn
	writer *kafka.Writer
	logger *mylogger.MyLogger
}

func NewBrokerMessageKafka(brokerAddr string, logger *mylogger.MyLogger) (*BrokerMessageKafka, error) {
	preBroker := &BrokerMessageKafka{ //nolint:exhaustruct
		logger: logger,
	}

	err := preBroker.initConn(logger, brokerAddr)
	if err != nil {
		return nil, err
	}

	preBroker.initWriter(brokerAddr)

	err = preBroker.initTopic(logger)
	if err != nil {
		return nil, err
	}

	return preBroker, nil
}

// initConn init start connection to Kafka and should be call first (before another init functions).
func (b *BrokerMessageKafka) initConn(logger *mylogger.MyLogger, brokerAddr string) error {
	var err error

	b.conn, err = kafka.Dial("tcp", brokerAddr)
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

func (b *BrokerMessageKafka) initTopic(logger *mylogger.MyLogger) error {
	topicConfigs := []kafka.TopicConfig{
		{ //nolint:exhaustruct
			Topic:             TopicMessageName,
			NumPartitions:     NumPartitions,
			ReplicationFactor: ReplicationFactor,
		},
	}

	err := b.conn.CreateTopics(topicConfigs...)
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

func (b *BrokerMessageKafka) initWriter(brokerAddr string) {
	b.logger.Debugln(brokerAddr, b.conn.RemoteAddr())

	b.writer = &kafka.Writer{
		Addr:     kafka.TCP(b.conn.RemoteAddr().String()),
		Topic:    TopicMessageName,
		Balancer: &kafka.Hash{}, //nolint:exhaustruct
	}
}

func (b *BrokerMessageKafka) Close() error {
	err := b.conn.Close()
	if err != nil {
		b.logger.Error(err)

		return err
	}

	err = b.writer.Close()
	if err != nil {
		b.logger.Error(err)

		return err
	}

	return nil
}

func (b *BrokerMessageKafka) WriteMessage(ctx context.Context, msgPayload *models.MessagePayload) error {
	logger := b.logger.EnrichReqID(ctx)

	err := b.writer.WriteMessages(ctx,
		kafka.Message{Key: msgPayload.ID[:], Value: []byte(msgPayload.Value)}, //nolint:exhaustruct
	)
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}
