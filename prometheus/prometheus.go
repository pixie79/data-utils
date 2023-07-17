package prometheus

import (
	"fmt"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"lambda"
	"lambda/utils"
	"strconv"
	"strings"
	"time"
)

func splitTags(data string) []lambda.tagsType {
	var tags []lambda.tagsType
	sData := strings.Split(data, ",")

	for _, tagsData := range sData {
		tagsSplit := lambda.reTagSplit.FindAllSubmatch([]byte(tagsData), -1)

		for _, v := range tagsSplit {
			tag := lambda.tagsType{
				Name:  string(v[1]),
				Value: string(v[2]),
			}
			tags = append(tags, tag)
		}

	}

	return tags

}

func buildMetrics(payload []string) []datadogV2.MetricSeries {
	var metrics []datadogV2.MetricSeries

	for _, line := range payload {
		splitLine := lambda.reInitialSplit.FindAllSubmatch([]byte(line), -1)

		if len(splitLine) > 0 {
			value, err := strconv.ParseFloat(string(splitLine[0][3]), 64)
			utils.maybeDie(err, fmt.Sprintf("could not convert to float: %+q", splitLine[0][0]))

			var tags []string
			for _, tag := range splitTags(string(splitLine[0][2])) {
				tags = append(tags, fmt.Sprintf("%s=%s", tag.Name, tag.Value))
			}

			metric := datadogV2.MetricSeries{
				Metric: fmt.Sprintf("redpanda.%s", string(splitLine[0][1])),
				Tags:   tags,
				Points: []datadogV2.MetricPoint{
					{
						Timestamp: datadog.PtrInt64(time.Now().Unix()),
						Value:     datadog.PtrFloat64(value),
					},
				},
			}

			metrics = append(metrics, metric)
		}

	}

	return metrics
}
