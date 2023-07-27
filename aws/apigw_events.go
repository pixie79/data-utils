package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pixie79/data-utils/utils"
	"github.com/twmb/franz-go/pkg/kgo"
	"strings"
)

func ApiGwCreateKafkaEvent(ctx context.Context, event events.APIGatewayProxyRequest, key []byte) ([]*kgo.Record, context.Context) {
	var (
		kafkaRecords []*kgo.Record
		keyValue     = utils.CreateKey(key)
		source       string
		topic        string
	)

	if _, found := event.PathParameters["proxy"]; found {
		source = event.PathParameters["proxy"]
	} else {
		utils.MaybeDie(fmt.Errorf("no source found"), "api gw proxy path parameter not found")
	}

	utils.Logger.Debug(fmt.Sprintf("Source is: %s", source))
	topic = strings.ToLower(source)
	utils.Logger.Debug(fmt.Sprintf("Topic name: %s", topic))
	// Basic payload as is
	payloadEvent := &kgo.Record{
		Topic: topic,
		Value: []byte(event.Body),
		Key:   keyValue,
	}
	kafkaRecords = append(kafkaRecords, payloadEvent)

	utils.Logger.Debug(fmt.Sprintf("Kafka Payload: %+v", kafkaRecords))

	// Return the kafka records to be sent to kafka
	return kafkaRecords, ctx
}
