package data_utils

import (
	"encoding/json"
	"time"
)

type CloudWatchEvent struct {
	Version    string          `json:"version"`
	ID         string          `json:"id"`
	DetailType string          `json:"detail-type"`
	Source     string          `json:"source"`
	AccountID  string          `json:"account"`
	Time       time.Time       `json:"time"`
	Region     string          `json:"region"`
	Resources  []string        `json:"resources"`
	Detail     json.RawMessage `json:"detail"`
}

type CredentialsType struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type KafkaPartitionLogEvent struct {
	Partition int64           `json:"partition"`
	Offset    int64           `json:"offset"`
	SchemaId  int             `json:"schema_id"`
	Payload   json.RawMessage `json:"payload"`
}

type SalesforceDetailEvent struct {
	Payload json.RawMessage `json:"payload"`
}

type SourceKey struct{}

type TagsType struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TopicKey struct{}
