// Description: Datadog utils
// Author: Pixie79
// ============================================================================
// package datadog

package datadog

import (
	"context"
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	tuUtils "github.com/pixie79/tiny-utils/utils"
)

// SubmitMetrics submits the metrics to Datadog
func submitMetrics(metrics []datadogV2.MetricSeries) {
	body := datadogV2.MetricPayload{
		Series: metrics,
	}

	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewMetricsApi(apiClient)

	_, r, err := api.SubmitMetrics(ctx, body, *datadogV2.NewSubmitMetricsOptionalParameters())
	tuUtils.MaybeDie(err, fmt.Sprintf("calling `MetricsApi.SubmitMetrics`: %v, payload length: %d", r, len(metrics)))

	tuUtils.Print("DEBUG", fmt.Sprintf("metrics submitted: %d: error code: %d", len(metrics), r.StatusCode))
}

// ChunkMetrics splits the metrics into chunks of X and submits them to Datadog
func ChunkMetrics(metrics []datadogV2.MetricSeries) {
	var (
		counter       = 0
		metricsLength = len(metrics)
	)

	if metricsLength < metricsBatchLengthInt {
		submitMetrics(metrics)
	} else {

		var (
			chunkedMetrics       = tuUtils.ChunkBy(metrics, metricsBatchLengthInt)
			chunkedMetricsLength = len(chunkedMetrics)
		)

		tuUtils.Print("DEBUG", fmt.Sprintf("payload to large splitting current length: %d, total number of new batches: %d", metricsLength, chunkedMetricsLength))
		// submit the metrics in batches
		for counter < chunkedMetricsLength {
			tuUtils.Print("DEBUG", fmt.Sprintf("submitting Batch: %d of %d records", counter, len(chunkedMetrics[counter])))
			submitMetrics(chunkedMetrics[counter])
			counter++
		}
	}
}
