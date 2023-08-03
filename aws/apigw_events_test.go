package aws

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/pixie79/data-utils/utils"
	"strings"
	"testing"
)

var (
	event1 = `{"name": "Paul","Age": 25,"Location": "USA"}`
	event2 = `[
	{"name": "Paul","Age": 25,"Location": "USA"},
	{"name": "Fred","Age": 27,"Location": "UK"},
	{"name": "John","Age": 20,"Location": "UAE"}
	]`
	event3 = `{"name": "Paul","Age": 25,"Location": "USA"}
	{"name": "Fred","Age": 27,"Location": "UK"}
	{"name": "John","Age": 20,"Location": "UAE"}`
)

func getJsonStrings(data string) []string {
	var result []string
	obj, err := oj.ParseString(data)
	if err != nil {
		utils.Logger.Debug("going to try and parse as json lines")
		obj = convertJsonLines(data)
	}
	switch obj.(type) {
	case []interface{}:
		for _, v := range obj.([]interface{}) {
			parse := localParse(v)
			var b strings.Builder
			if err := oj.Write(&b, parse, &ojg.Options{Sort: true}); err != nil {
				panic(err)
			}
			result = append(result, b.String())
			utils.Logger.Info(fmt.Sprintf("parse: %s", b.String()))
		}
	case map[string]interface{}:
		parse := localParse(obj)
		var b strings.Builder
		if err := oj.Write(&b, parse, &ojg.Options{Sort: true}); err != nil {
			panic(err)
		}
		result = append(result, b.String())
		utils.Logger.Info(fmt.Sprintf("parse: %s", b.String()))
	default:
		panic("unknown type")
	}
	return result
}

func localParse(obj any) any {
	x, err := jp.ParseString("$.name")
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("error parsing json: %s", err))
	}
	utils.Logger.Info(fmt.Sprintf("obj: %T", obj))
	result := x.Get(obj)
	return result
}

func convertJsonLines(data string) any {
	var (
		payload []string
	)
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		payload = append(payload, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error occurred: %v\n", err)
	}
	result := strings.Join(payload, ",")
	obj, err := oj.ParseString("[" + result + "]")
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("error parsing jsonlines: %s", err))
	}
	return obj
}

func TestGetJson(t *testing.T) {
	var tests = []struct {
		name  string
		event string
		topic []string
	}{
		{"test1", event1, []string{"Paul"}},
		{"test2", event2, []string{"Paul", "Fred", "John"}},
		{"test3", event3, []string{"Paul", "Fred", "John"}},
	}
	{
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := getJsonStrings(tt.event)
				utils.Logger.Info(fmt.Sprintf("result: %s", result))

				if result[0] != tt.topic[0] {
					t.Errorf("got %s, want %s", result, tt.topic)
				}
			})
		}
	}
}

func _getApiGwEvent(event string) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{
		Body:            base64.URLEncoding.EncodeToString(json.RawMessage(event)),
		IsBase64Encoded: true,
		Headers:         map[string]string{"Content-Type": "application/vnd.kafka.v2+json"},
		PathParameters:  map[string]string{"proxy": "payment"},
	}

}

func TestFindJsonTopicFromBody(t *testing.T) {
	var tests = []struct {
		name     string
		apiEvent events.APIGatewayProxyRequest
		topic    []string
	}{
		//{"test1", _getApiGwEvent(event), []string{"CREDIT_COMPLETION_DELIVERED", "CREDIT_TRANSFER_RECEIVED", "CREDIT_AUTH_DELIVERY_PENDING", "VALIDATION_SUCCESSFUL"}},
		{"test2", _getApiGwEvent(event2), []string{"CREDIT_COMPLETION_DELIVERED"}},
	} // The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			utils.Logger.Info(fmt.Sprintf("Event: %s", tt.apiEvent.Body))
			payloads, _ := ApiGwCreateKafkaEvent(ctx, tt.apiEvent, []byte(""))
			for key, record := range payloads {
				utils.Logger.Info(fmt.Sprintf("Topic: %s", record.Topic))
				utils.Logger.Info(fmt.Sprintf("Value: %s", record.Value))
				utils.Logger.Info(fmt.Sprintf("Key: %d, Want Topic %s, Got %s", key, tt.topic[key], record.Topic))
				if record.Topic != tt.topic[key] {
					t.Errorf("got %s, want %s", record.Topic, tt.topic)
				}
			}
		})
	}

}
