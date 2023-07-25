package aws

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/pixie79/data-utils"
	"github.com/twmb/franz-go/pkg/kgo"
)

// TestGetCloudWatchTopic - Checks source and topicDetail, switches underscore for hyphen, multiple hyphens for single and lowercase
func TestGetCloudWatchTopic(t *testing.T) {
	// Defining the columns of the table
	var tests = []struct {
		name       string
		source     string
		detailType string
		want       string
	}{
		// the table itself
		{"topicname, topictype should be topicname-topictype", "topicname", "topictype", "topicname-topictype"},
		{"salesforce, DD_chatTranscript__c should be salesforce-dd-chattranscript_c", "salesforce", "DD_chatTranscript__c", "salesforce-dd-chattranscript-c"},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, data_utils.SourceKey{}, tt.source)
			ctx = GetCloudWatchTopic(ctx, tt.detailType)
			topic := ctx.Value(data_utils.TopicKey{}).(string)
			if topic != tt.want {
				t.Errorf("got %s, want %s", topic, tt.want)
			}
		})
	}
}

// TestGetCloudWatchSource - Checks source and filters for the correct source name
func TestGetCloudWatchSource(t *testing.T) {
	// Defining the columns of the table
	var tests = []struct {
		name   string
		source string
		want   string
	}{
		// the table itself
		{"source should be topicname", "topicname", "topicname"},
		{"source should be salesforce", "aws.partner/salesforce.com/00A0A000000AA0aBCD/ChangeEvents", "salesforce"},
		{"source should be test", "aws.partner/test/00A0A000000AA0aBCD/ChangeEvents", "test"},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = GetCloudWatchSource(ctx, tt.source)
			source := ctx.Value(data_utils.SourceKey{}).(string)
			if source != tt.want {
				t.Errorf("got %s, want %s", source, tt.want)
			}
		})
	}
}

// TestCloudWatchCreateEvent - Creates a sample event in the correct format
func TestCloudWatchCreateEvent(t *testing.T) {
	testArn := []string{"arn:aws:events:eu-west-2:123456789012:event-bus/123456789012-eu-west-2-testapp"}
	detailPayload := json.RawMessage(`{"email":"a@b.com", "state":"CA", "city":"San Francisco", "zipcode":"94107"}`)
	event1 := data_utils.CloudWatchEvent{
		Version:    "0",
		ID:         "972723a0-69b8-4ddf-8729-5b0b4fb4af15",
		DetailType: "mobile",
		Source:     "testapp",
		AccountID:  "123456789012",
		Time:       time.Now(),
		Region:     "eu-west-2",
		Resources:  testArn,
		Detail:     detailPayload,
	}
	kafkaEvent1 := kgo.Record{
		Topic: "testapp-mobile",
		Value: json.RawMessage{},
	}
	event2 := data_utils.CloudWatchEvent{
		Version:    "0",
		ID:         "972723a0-69b8-4ddf-8729-5b0b4fb4af15",
		DetailType: "CaseChangeEvent",
		Source:     "aws.partner/salesforce.com/00A0A000000AA0aBCD/ChangeEvents",
		AccountID:  "123456789012",
		Time:       time.Now(),
		Region:     "eu-west-2",
		Resources:  testArn,
		Detail:     detailPayload,
	}
	kafkaEvent2 := kgo.Record{
		Topic: "salesforce-casechangeevent",
		Value: json.RawMessage{},
	}

	// Defining the columns of the table
	var tests = []struct {
		name  string
		event data_utils.CloudWatchEvent
		want  kgo.Record
	}{
		// the table itself
		{"Simple generic test", event1, kafkaEvent1},
		{"Simple salesforce match", event2, kafkaEvent2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, ctx := CloudWatchCreateEvent(ctx, tt.event, []byte(""))
			if result[0].Topic != tt.want.Topic {
				t.Errorf("got: %+v, want: %+v, context: %+v", result[0], tt.want, ctx)
			}
		})
	}
}
