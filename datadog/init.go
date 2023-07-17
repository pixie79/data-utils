package datadog

import (
	"dataUtils/utils"
	"strconv"
)

var (
	metricsBatchLength    string
	metricsBatchLengthInt int
)

func init() {
	metricsBatchLength = utils.GetEnv("METRICS_BATCH_LENGTH", "800")
	metricsBatchLengthInt, utils.Err = strconv.Atoi(metricsBatchLength)
	utils.MaybeDie(utils.Err, "cannot convert string metricsBatchLength to int")
}
