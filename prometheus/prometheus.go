// Description: Prometheus utils to help parse prometheus metrics
// Author: Pixie79
// ============================================================================
// package prometheus

package prometheus

import (
	"fmt"
	"github.com/pixie79/data-utils/types"
	"strconv"
	"strings"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/pixie79/data-utils/utils"
)

// SplitTags splits the tags from a Prometheus metric into a slice of tags
func SplitTags(data string) []types.TagsType {
	var tags []types.TagsType
	sData := strings.Split(data, ",")

	for _, tagsData := range sData {
		tagsSplit := reTagSplit.FindAllSubmatch([]byte(tagsData), -1)

		for _, v := range tagsSplit {
			tag := types.TagsType{
				Name:  string(v[1]),
				Value: string(v[2]),
			}
			tags = append(tags, tag)
		}
	}
	tags = append(tags, additionalTags...) // Add additional tags

	return tags

}

func BuildMetrics(payload []string) []datadogV2.MetricSeries {
	var metrics []datadogV2.MetricSeries

	for _, line := range payload {
		splitLine := reInitialSplit.FindAllSubmatch([]byte(line), -1)

		if len(splitLine) > 0 {
			value, err := strconv.ParseFloat(string(splitLine[0][3]), 64)
			utils.MaybeDie(err, fmt.Sprintf("could not convert to float: %+q", splitLine[0][0]))

			var tags []string
			for _, tag := range SplitTags(string(splitLine[0][2])) {
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
