package datadog

import (
	"context"
	"fmt"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/pixie79/data-utils/utils"
)

func submitMetrics(metrics []datadogV2.MetricSeries) {
	body := datadogV2.MetricPayload{
		Series: metrics,
	}

	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewMetricsApi(apiClient)

	_, r, err := api.SubmitMetrics(ctx, body, *datadogV2.NewSubmitMetricsOptionalParameters())
	utils.MaybeDie(err, fmt.Sprintf("Error when calling `MetricsApi.SubmitMetrics`: %v, payload length: %d", r, len(metrics)))

	utils.Logger.Debug(fmt.Sprintf("metrics submitted: %d: error code: %d", len(metrics), r.StatusCode))
}

func ChunkMetrics(metrics []datadogV2.MetricSeries) {
	var (
		counter       = 0
		metricsLength = len(metrics)
	)

	if metricsLength < metricsBatchLengthInt {
		submitMetrics(metrics)
	} else {

		var (
			chunkedMetrics       = utils.ChunkBy(metrics, metricsBatchLengthInt)
			chunkedMetricsLength = len(chunkedMetrics)
		)

		utils.Logger.Debug(fmt.Sprintf("payload to large splitting current length: %d, total number of new batches: %d", metricsLength, chunkedMetricsLength))
		for counter < chunkedMetricsLength {
			utils.Logger.Debug(fmt.Sprintf("Submitting Batch: %d of %d records", counter, len(chunkedMetrics[counter])))
			submitMetrics(chunkedMetrics[counter])
			counter++
		}
	}
}
