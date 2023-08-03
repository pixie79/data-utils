package types

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type DynamoDBEvent struct {
	Records []DynamoDBEventRecord `json:"Records"`
}

type DynamoDBUserIdentity struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
}

type DynamoDBEventRecord struct {
	EventID        string                       `json:"eventID"`   // e.g. 1
	EventName      string                       `json:"eventName"` // e.g. 1.1
	Change         DynamoDBStreamRecord         `json:"dynamodb"`
	AWSRegion      string                       `json:"awsRegion"`
	EventSourceArn string                       `json:"eventSourceARN"` // e.g. arn:aws:dynamodb:us-east-1:123456789012:table/MyTableWithStream/stream/2019-03-01T22:00:00.000
	EventSource    string                       `json:"eventSource"`
	EventVersion   string                       `json:"eventVersion"`
	UserIdentity   *events.DynamoDBUserIdentity `json:"userIdentity,omitempty"`
}

type DynamoDBStreamRecord struct {
	ApproximateCreationDateTime events.SecondsEpochTime `json:"ApproximateCreationDateTime,omitempty"`
	Keys                        json.RawMessage         `json:"Keys,omitempty"`
	NewImage                    json.RawMessage         `json:"NewImage,omitempty"`
	OldImage                    json.RawMessage         `json:"OldImage,omitempty"`
	StreamViewType              string                  `json:"StreamViewType"`
	SequenceNumber              string                  `json:"SequenceNumber"`
	SizeBytes                   int64                   `json:"SizeBytes"`
}
