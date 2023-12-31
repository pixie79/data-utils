package aws

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pixie79/data-utils/types"
	"github.com/pixie79/data-utils/utils"
	"github.com/twmb/franz-go/pkg/kgo"
)

var (
	timestamp = events.SecondsEpochTime{Time: time.Now().UTC()}
)

func _getDynamoDbEvent(newItem json.RawMessage, oldItem json.RawMessage) types.DynamoDBEvent {

	return types.DynamoDBEvent{
		Records: []types.DynamoDBEventRecord{
			{
				AWSRegion:    "us-west-2",
				EventName:    "INSERT",
				EventSource:  "aws:dynamodb",
				EventID:      "sampleEventId",
				EventVersion: "1.1",
				Change: types.DynamoDBStreamRecord{
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

func TestDynamoDbCreateKafkaEvent(t *testing.T) {
	newItem := json.RawMessage(`{"email":"a@b.com", "state":"CA", "city":"San Francisco", "zipcode":"94107"}`)
	oldItem := json.RawMessage(`{"email":"a@example.com", "state":"CA", "city":"San Francisco", "zipcode":"94105"}`)
	ctx := context.Background()
	event := _getDynamoDbEvent(newItem, oldItem)

	key := []byte("key")
	actual, ctx := DynamoDbCreateKafkaEvent(ctx, event, key)

	// Prepare the expected value
	expected := []*kgo.Record{
		{
			Topic: `data-dynamodb-exampletablewithstream`,
			Value: utils.CreateBytes(event.Records[0]),
			Key:   key,
		},
	}
	if ctx.Value(types.TopicKey{}).(string) != expected[0].Topic {
		t.Errorf("Expected %v, got %v", expected[0].Topic, ctx.Value(types.TopicKey{}).(string))
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
