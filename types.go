// Description: A package containing useful reusable snippets of GO code
// Author: Pixie79
// ============================================================================
// package data_utils

package data_utils

import (
	"encoding/json"
	"time"
)

// CloudWatchEvent is the event sent by CloudWatch to Lambda functions
type CloudWatchEvent struct {
	Version    string          `json:"version"`     // e.g. 0
	ID         string          `json:"id"`          // e.g. 12345678-1234-1234-1234-123456789012
	DetailType string          `json:"detail-type"` // e.g. AWS API Call via CloudTrail or mobile
	Source     string          `json:"source"`      // aws.partner/salesforce.com/... or aws.ec2 or rudderstack
	AccountID  string          `json:"account"`     // e.g. 123456789012
	Time       time.Time       `json:"time"`        // e.g. 2019-03-01T21:49:13Z
	Region     string          `json:"region"`      // e.g. us-east-1
	Resources  []string        `json:"resources"`   // ARNs of resources
	Detail     json.RawMessage `json:"detail"`      // this is the raw JSON event
}

// CredentialsType is a struct to hold credentials
type CredentialsType struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// KafkaPartitionLogEvent is the event sent by Kafka to Lambda functions
type KafkaPartitionLogEvent struct {
	Partition int64           `json:"partition"`
	Offset    int64           `json:"offset"`
	SchemaId  int             `json:"schema_id"`
	Payload   json.RawMessage `json:"payload"`
}

// SalesforceDetailEvent is the event sent by Salesforce to Lambda functions
type SalesforceDetailEvent struct {
	Payload json.RawMessage `json:"payload"`
}

// SourceKey is the key used to store the source in the context
type SourceKey struct{}

// TagsType is a struct to hold tags
type TagsType struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// TopicKey is the key used to store the topic in the context
type TopicKey struct{}
