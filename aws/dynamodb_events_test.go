package aws

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pixie79/data-utils"
	"github.com/pixie79/data-utils/utils"
	"github.com/twmb/franz-go/pkg/kgo"
	"reflect"
	"testing"
	"time"
)

var (
	sampleKafkaRecords []*kgo.Record
	sampleKafkaRecord  = kgo.Record{Topic: "test", Value: json.RawMessage(`{"foo":"bar"}`)}
	timestamp          = events.SecondsEpochTime{Time: time.Now().UTC()}
)

func init() {
	sampleKafkaRecords = []*kgo.Record{&sampleKafkaRecord}
}

func getEvent(newItem json.RawMessage, oldItem json.RawMessage) data_utils.DynamoDBStreamEvent {

	return data_utils.DynamoDBStreamEvent{
		Records: []data_utils.DynamoDBStreamRecord{
			{
				AWSRegion:    "us-west-2",
				EventName:    "INSERT",
				EventSource:  "aws:dynamodb",
				EventID:      "sampleEventId",
				EventVersion: "1.1",
				Change: data_utils.DynamoDBStreamRecordChange{
					ApproximateCreationDateTime: timestamp,
					Keys:                        nil,
					NewImage:                    newItem,
					OldImage:                    oldItem,
					SequenceNumber:              "sampleSequenceNumber",
					SizeBytes:                   2,
					StreamViewType:              "NEW_AND_OLD_IMAGES",
				},
				EventSourceArn: "arn:aws:dynamodb:us-west-2:accountid:table/ExampleTableWithStream/stream/2015-06-27T00:48:05.899",
			},
		},
	}
}

func TestGetDynamoDbPayload(t *testing.T) {
	newItem := json.RawMessage(`{"email":"a@b.com", "state":"CA", "city":"San Francisco", "zipcode":"94107"}`)
	oldItem := json.RawMessage(`{"email":"a@example.com", "state":"CA", "city":"San Francisco", "zipcode":"94105"}`)

	event := getEvent(newItem, oldItem)

	key := []byte("key")
	actual := DynamoDbCreateEvent(event, key)

	// Prepare the expected value
	expected := []*kgo.Record{
		{
			Topic: `data-dynamodb-exampletablewithstream`,
			Value: utils.CreateBytes(event.Records[0]),
			Key:   key,
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestGetDynamoDBSource(t *testing.T) {
	var tests = []struct {
		name        string
		eventSource string
		want        string
	}{
		// the table itself
		{"test 1", "arn:aws:dynamodb:us-east-1:123456789012:table/MyTableWithStream/stream/2019-03-01T22:00:00.000", "mytablewithstream"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := GetDynamoDBSource(tt.eventSource)
			if source != tt.want {
				t.Errorf("GetDynamoDBSource() = %v, want %v", source, tt.want)
			}
		})
	}
}
