// Description: AWS utils
// Author: Pixie79
// ============================================================================
// package aws

package aws

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pixie79/data-utils/types"

	"github.com/pixie79/data-utils/utils"
	tuUtils "github.com/pixie79/tiny-utils/utils"
	"github.com/twmb/franz-go/pkg/kgo"
)

// GetDynamoDBSource retrieves the source from the event source
func GetDynamoDBSource(eventSourceArn string) string {
	r, _ := regexp.Compile(`^.*:table/(.*)/stream`) // arn:aws:dynamodb:us-east-1:123456789012:table/MyTableWithStream/stream/2019-03-01T22:00:00.000
	sourceResult := r.FindStringSubmatch(eventSourceArn)
	if len(sourceResult) >= 1 {
		return strings.ToLower(sourceResult[1])
	} else if len(eventSourceArn) > 0 {
		return strings.ToLower(eventSourceArn)
	} else {
		tuUtils.Die("no source found")
	}
	return ""
}

// DynamoDbCreateKafkaEvent retrieves the event
func DynamoDbCreateKafkaEvent(ctx context.Context, event types.DynamoDBEvent, key []byte) ([]*kgo.Record, context.Context) {
	var (
		kafkaRecords []*kgo.Record
		keyValue     = tuUtils.CreateKey(key)
		source       = GetDynamoDBSource(event.Records[0].EventSourceArn)
		topic        = strings.ToLower(utils.Prefix) + `-dynamodb-` + strings.ToLower(source)
	)

	for _, v := range event.Records {
		payloadEvent := &kgo.Record{
			Topic: topic,
			Value: utils.CreateBytes(v),
			Key:   keyValue,
		}
		kafkaRecords = append(kafkaRecords, payloadEvent)
	}

	// Return the kafka records to be sent to kafka
	return kafkaRecords, context.WithValue(ctx, types.TopicKey{}, topic)
}

func MarshalDynamoDBEventToLocal(event events.DynamoDBEvent) types.DynamoDBEvent {
	holdingEvent, err := json.Marshal(event)
	tuUtils.MaybeDie(err, "unable to parse raw event: %+v")
	newEvent := types.DynamoDBEvent{}
	err = json.Unmarshal(holdingEvent, &newEvent)
	tuUtils.MaybeDie(err, "unable to load event: %+v")
	return newEvent
}
