package aws

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	data_utils "github.com/pixie79/data-utils"
	"github.com/twmb/franz-go/pkg/kgo"
)

// TestGetTopic - Checks source and topicDetail, switches underscore for hyphen, multiple hyphens for single and lowercase
func TestGetTopic(t *testing.T) {
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
			ctx = GetTopic(ctx, tt.detailType)
			topic := ctx.Value(data_utils.TopicKey{}).(string)
			if topic != tt.want {
				t.Errorf("got %s, want %s", topic, tt.want)
			}
		})
	}
}

// TestGetPayload - Checks source and topicDetail, getting the correct payload for the source
func TestGetPayload(t *testing.T) {
	// Defining the columns of the table
	var tests = []struct {
		name   string
		source string
		topic  string
		detail []byte
		want   *kgo.Record
	}{
		// the table itself
		{"Simple source", "topicname", "topicname", json.RawMessage(`{"foo":"bar"}`), &kgo.Record{Topic: "topicname", Value: json.RawMessage(`{"foo":"bar"}`)}},
		{"Salesforce source", "salesforce", "salesforce", json.RawMessage(`{"Payload": {"foo":"bar"}}`), &kgo.Record{Topic: "salesforce", Value: json.RawMessage(`{"foo":"bar"}`)}},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, data_utils.SourceKey{}, tt.source)
			ctx = context.WithValue(ctx, data_utils.TopicKey{}, tt.topic)
			result := GetPayload(ctx, tt.detail, nil)
			if cmp.Equal(result[0].Topic, tt.want.Topic) == false {
				t.Errorf("got %+v, want %+v", result[0], tt.want)
			}
			if cmp.Equal(result[0].Value, tt.want.Value) == false {
				t.Errorf("got %+v, want %+v", result[0].Value, tt.want.Value)
			}
		})
	}
}

// TestGetSource - Checks source and filters for the correct source name
func TestGetSource(t *testing.T) {
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
			ctx = GetSource(ctx, tt.source)
			source := ctx.Value(data_utils.SourceKey{}).(string)
			if source != tt.want {
				t.Errorf("got %s, want %s", source, tt.want)
			}
		})
	}
}

// TestCreateEvent - Creates a sample event in the correct format
func TestCreateEvent(t *testing.T) {
	testArn := []string{"arn:aws:events:eu-west-2:123456789012:event-bus/123456789012-eu-west-2-testapp"}
	detailPayload := json.RawMessage{}
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
			result, ctx := CreateEvent(ctx, tt.event)
			if result[0].Topic != tt.want.Topic {
				t.Errorf("got: %+v, want: %+v, context: %+v", result[0], tt.want, ctx)
			}
		})
	}
}
