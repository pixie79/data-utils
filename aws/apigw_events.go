package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pixie79/data-utils/utils"
	"github.com/tidwall/gjson"
	"github.com/twmb/franz-go/pkg/kgo"
)

func ApiGwCreateKafkaEvent(ctx context.Context, event events.APIGatewayProxyRequest, key []byte) ([]*kgo.Record, context.Context) {
	var (
		kafkaRecords      []*kgo.Record
		keyValue          = utils.CreateKey(key)
		source            string
		topic             string
		payloadKey        string
		topicInBody       = true
		partialPayloadKey string
		partialPayload    = false
		value             []byte
	)

	if _, found := event.PathParameters["proxy"]; found {
		source = event.PathParameters["proxy"]
		source = strings.ToLower(source)
	} else {
		utils.MaybeDie(fmt.Errorf("no source found"), "api gw proxy path parameter not found")
	}

	if source == "electrum" {
		payloadKey = "type"
		topicInBody = true
		partialPayload = false
		partialPayloadKey = "tranInfo"
	}

	utils.Print("DEBUG", fmt.Sprintf("Source is: %s", source))

	decodedPayload := make([]byte, base64.StdEncoding.DecodedLen(len(event.Body)))
	n, err := base64.StdEncoding.Decode(decodedPayload, []byte(event.Body))
	utils.MaybeDie(err, "unable to decode base64 payload")

	//payloads := ReturnListFromString(string(decodedPayload[:n]))
	if !gjson.Valid(string(decodedPayload[:n])) {
		utils.MaybeDie(errors.New("invalid json"), "error parsing json")
	}
	//value1 := gjson.Get(string(decodedPayload[:n]), ".")
	//m, ok := gjson.Parse(string(decodedPayload[:n])).Value().(map[string]interface{})
	m := gjson.GetManyBytes(decodedPayload[:n], "")
	//if !ok {
	//	utils.Print("INFO", "Not a map")
	//}
	fmt.Printf("%+v\n", m)
	fmt.Printf("length: %d\n\n", len(m))
	for _, payload := range m {
		payloadJson, err := json.Marshal(payload)
		utils.MaybeDie(err, "unable to marshal payload")

		if topicInBody {
			topic = strings.ToLower(source) + "-" + gjson.Get(
				string(payloadJson),
				payloadKey,
			).String()
		} else {
			topic = strings.ToLower(source)
		}

		if partialPayload {
			value = []byte(gjson.Get(string(payloadJson), partialPayloadKey).String())
		} else {
			value = payloadJson
		}

		utils.Print("DEBUG", fmt.Sprintf("topic: %s, value: %s, key: %s", topic, value, key))

		payloadEvent := &kgo.Record{
			Topic: topic,
			Value: value,
			Key:   keyValue,
		}
		kafkaRecords = append(kafkaRecords, payloadEvent)
	}

	// Return the kafka records to be sent to kafka
	return kafkaRecords, ctx
}

func ReturnListFromString(body string) []map[string]interface{} {
	var (
		payloads []map[string]interface{}
		payload  interface{}
	)

	// Unmarshal or Decode the JSON to the interface.
	err := json.Unmarshal([]byte(body), &payloads)
	if err != nil {
		err = json.Unmarshal([]byte(body), &payload)
		utils.MaybeDie(err, "unable to unmarshal json body")
		payloads = append(payloads, payload.(map[string]interface{}))
	}
	utils.MaybeDie(err, "unable to unmarshal json body")

	utils.Print("INFO", fmt.Sprintf("number of items in payload: %d", len(payloads)))
	return payloads
}
