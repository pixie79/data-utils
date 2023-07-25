// Description: A package containing useful reusable snippets of GO code
// Author: Pixie79
// ============================================================================
// package data_utils

package data_utils

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
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

type DynamoDBStreamEvent struct {
	Records []DynamoDBStreamRecord
}

type DynamoDBStreamRecord struct {
	EventID        string                     `json:"event-id"`      // e.g. 1
	EventVersion   string                     `json:"event-version"` // e.g. 1.1
	Change         DynamoDBStreamRecordChange `json:"change"`
	AWSRegion      string                     `json:"aws-region"`
	EventName      string                     `json:"event-name"`
	EventSourceArn string                     `json:"event-source-arn"` // e.g. arn:aws:dynamodb:us-east-1:123456789012:table/MyTableWithStream/stream/2019-03-01T22:00:00.000
	EventSource    string                     `json:"event-source"`
}

type DynamoDBStreamRecordChange struct {
	ApproximateCreationDateTime events.SecondsEpochTime `json:"approximate-creation-date-time"`
	Keys                        json.RawMessage         `json:"keys"`
	NewImage                    json.RawMessage         `json:"new-image"`
	StreamViewType              string                  `json:"stream-view-type"`
	SequenceNumber              string                  `json:"sequence-number"`
	SizeBytes                   int                     `json:"size-bytes"`
	OldImage                    json.RawMessage         `json:"old-image,omitempty"`
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
