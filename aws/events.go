// Description: AWS utils
// Author: Pixie79
// ============================================================================
// package aws

package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	data_utils "github.com/pixie79/data-utils"
	"github.com/pixie79/data-utils/utils"
	"github.com/twmb/franz-go/pkg/kgo"
)

// GetSource retrieves the source from the event source
func GetSource(ctx context.Context, eventSource string) context.Context {
	r, _ := regexp.Compile(`^(aws.partner/)([a-zA-Z]*)`)
	sourceResult := r.FindStringSubmatch(eventSource)
	if len(sourceResult) >= 1 {
		return context.WithValue(ctx, data_utils.SourceKey{}, strings.ToLower(sourceResult[2]))
	} else if len(eventSource) > 0 {
		return context.WithValue(ctx, data_utils.SourceKey{}, strings.ToLower(eventSource))
	} else {
		utils.Die(fmt.Errorf("source result empty"), "no source found")
	}
	return ctx
}

// GetTopic retrieves the topic from the event detail type
func GetTopic(ctx context.Context, detailType string) context.Context {
	topic := fmt.Sprintf("%s-%s", ctx.Value(data_utils.SourceKey{}).(string),
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ToLower(detailType), "_", "-"),
			"--", "-"))
	return context.WithValue(ctx, data_utils.TopicKey{}, topic)
}

// GetPayload retrieves the payload from the event detail
func GetPayload(ctx context.Context, detail []byte, key []byte) []*kgo.Record {
	var (
		kafkaRecords []*kgo.Record
		keyValue     []byte
		topic        = ctx.Value(data_utils.TopicKey{}).(string)
		source       = ctx.Value(data_utils.SourceKey{}).(string)
	)

	// If key is empty, use hostname as key
	if len(key) < 1 {
		keyValue = []byte(utils.Hostname)
	} else {
		keyValue = key
	}
	if source == "salesforce" {
		// If source is salesforce, unmarshal the payload and use the payload as value
		customStructure := &data_utils.SalesforceDetailEvent{}
		_ = json.Unmarshal(detail, customStructure)
		payloadEvent := &kgo.Record{
			Topic: topic,
			Key:   keyValue,
			Value: customStructure.Payload,
		}
		kafkaRecords = append(kafkaRecords, payloadEvent)
	} else {
		// If source is not salesforce, use the payload as is
		payloadEvent := &kgo.Record{
			Topic: topic,
			Value: detail,
		}
		kafkaRecords = append(kafkaRecords, payloadEvent)
	}
	// Return the kafka records to be sent to kafka
	return kafkaRecords
}

// CreateEvent creates a kafka record from the event
func CreateEvent(ctx context.Context, event data_utils.CloudWatchEvent) ([]*kgo.Record, context.Context) {
	ctx = GetSource(ctx, event.Source)
	ctx = GetTopic(ctx, event.DetailType)
	kafkaRecords := GetPayload(ctx, event.Detail, nil)
	return kafkaRecords, ctx
}