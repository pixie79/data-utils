// Description: Kafka utils
// Author: Pixie79
// ============================================================================
// package kafka

package kafka

import (
	"context"
	"crypto/tls"
	"fmt"

	data_utils "github.com/pixie79/data-utils"
	"github.com/pixie79/data-utils/utils"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

// CreateConnectionAndSubmitRecords creates a connection to Kafka and submits records
func CreateConnectionAndSubmitRecords(ctx context.Context, kafkaRecords []*kgo.Record, credentials data_utils.CredentialsType) *kgo.Client {
	var (
		opts          []kgo.Opt
		transactionId = fmt.Sprintf("eventbridge-%s", utils.Hostname)
		topic         = ctx.Value(data_utils.TopicKey{}).(string)
	)
	// Set up the kgo Client, which handles all the broker communication
	// and underlies any producer/consumer actions.
	seedEnv := utils.GetEnvOrDie("KAFKA_SEEDS")
	seeds := []string{seedEnv}
	opts = append(opts,
		kgo.SeedBrokers(seeds...),
		kgo.TransactionalID(transactionId),
		kgo.DefaultProduceTopic(topic),
		kgo.RecordPartitioner(kgo.RoundRobinPartitioner()),
		kgo.RecordRetries(4),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.AllowAutoTopicCreation(),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
	)
	// Initialize public CAs for TLS
	opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))

	//// Initializes SASL/SCRAM 256
	opts = append(opts, kgo.SASL(scram.Auth{
		User: credentials.Username,
		Pass: credentials.Password,
	}.AsSha256Mechanism()))

	client, err := kgo.NewClient(opts...)
	utils.MaybeDie(err, "could not connect to Kafka")
	defer client.Close()

	return client
}

// SubmitRecords submits records to Kafka
func SubmitRecords(ctx context.Context, client *kgo.Client, kafkaRecords []*kgo.Record) error {
	var (
		topic = ctx.Value(data_utils.TopicKey{}).(string)
	)
	// Start the transaction so that we can start buffering records.
	if err := client.BeginTransaction(); err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	// Write a message for each log file to Kafka.
	if err := ProduceMessages(ctx, client, kafkaRecords); err != nil {
		if rollbackErr := RollbackTransaction(client); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	// we're running in autocommit mode by default, which will flush all the
	// buffered messages before attempting to commit the transaction.
	if err := client.EndTransaction(ctx, kgo.TryCommit); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	utils.Logger.Info(fmt.Sprintf("kafka topic: %s, records produced %d", topic, len(kafkaRecords)))

	//TODO Update Metric production
	//ProduceMetric(topic, len(kafkaRecords))

	return nil
}
