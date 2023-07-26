package types

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type DynamoDBEvent struct {
	Records []DynamoDBEventRecord `json:"Records"`
}

type DynamoDBEventRecord struct {
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
