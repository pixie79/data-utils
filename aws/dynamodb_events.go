// Description: AWS utils
// Author: Pixie79
// ============================================================================
// package aws

package aws

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pixie79/data-utils"
	"github.com/pixie79/data-utils/utils"
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
		utils.Die(fmt.Errorf("source result empty"), "no source found")
	}
	return ""
}

// DynamoDbCreateEvent retrieves the event
func DynamoDbCreateEvent(event data_utils.DynamoDBStreamEvent, key []byte) []*kgo.Record {
	var (
		kafkaRecords []*kgo.Record
		keyValue     []byte
		source       = GetDynamoDBSource(event.Records[0].EventSourceArn)
		topic        = strings.ToLower(utils.Prefix) + `-dynamodb-` + strings.ToLower(source)
	)

	// If key is empty, use hostname as key
	if len(key) < 1 {
		keyValue = []byte(utils.Hostname)
	} else {
		keyValue = key
	}

	for _, v := range event.Records {
		payloadEvent := &kgo.Record{
			Topic: topic,
			Value: utils.CreateBytes(v),
			Key:   keyValue,
		}
		kafkaRecords = append(kafkaRecords, payloadEvent)
	}

	// Return the kafka records to be sent to kafka
	return kafkaRecords
}
