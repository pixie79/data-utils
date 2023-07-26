// Description: AWS utils
// Author: Pixie79
// ============================================================================
// package aws

package aws

import (
	"context"
	"encoding/json"
	"fmt"
	data_utils "github.com/pixie79/data-utils/types"
	"regexp"
	"strings"

	"github.com/pixie79/data-utils/utils"
	"github.com/twmb/franz-go/pkg/kgo"
)

// GetCloudWatchSource retrieves the source from the event source
func GetCloudWatchSource(ctx context.Context, eventSource string) context.Context {
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

// GetCloudWatchTopic retrieves the topic from the event detail type
func GetCloudWatchTopic(ctx context.Context, detailType string) context.Context {
	topic := fmt.Sprintf("%s-%s", ctx.Value(data_utils.SourceKey{}).(string),
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ToLower(detailType), "_", "-"),
			"--", "-"))
	return context.WithValue(ctx, data_utils.TopicKey{}, topic)
}

// CloudWatchCreateEvent retrieves the payload from the event detail
func CloudWatchCreateEvent(ctx context.Context, event data_utils.CloudWatchEvent, key []byte) ([]*kgo.Record, context.Context) {
	utils.Logger.Info(fmt.Sprintf("Running CloudWatchCreateEvent %+v", event))
	ctx = GetCloudWatchSource(ctx, event.Source)
	ctx = GetCloudWatchTopic(ctx, event.DetailType)

	var (
		kafkaRecords []*kgo.Record
		keyValue     []byte
		detail       = event.Detail
		source       = ctx.Value(data_utils.SourceKey{}).(string)
		topic        = ctx.Value(data_utils.TopicKey{}).(string)
	)

	utils.Logger.Info(fmt.Sprintf("Running topic: %s, source: %s", topic, source))
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
			Key:   keyValue,
		}
		kafkaRecords = append(kafkaRecords, payloadEvent)
	}
	// Return the kafka records to be sent to kafka
	return kafkaRecords, ctx
}
