package types

import "encoding/json"

// KafkaPartitionLogEvent is the event sent by Kafka to Lambda functions
type KafkaPartitionLogEvent struct {
	Partition int64           `json:"partition"`
	Offset    int64           `json:"offset"`
	SchemaId  int             `json:"schema_id"`
	Payload   json.RawMessage `json:"payload"`
}

// SourceKey is the key used to store the source in the context
type SourceKey struct{}

// TopicKey is the key used to store the topic in the context
type TopicKey struct{}
