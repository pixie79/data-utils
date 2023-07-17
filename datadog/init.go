package datadog

import (
	"dataUtils/utils"
	"regexp"
	"strconv"
)

var (
	reInitialSplit        = regexp.MustCompile(`(.+){(.+)}\s(\d+\.?\d*)`)
	reTagSplit            = regexp.MustCompile(`(\w+)="(.+)"`)
	metricsBatchLength    string
	metricsBatchLengthInt int
)

func init() {
	metricsBatchLength = utils.GetEnv("METRICS_BATCH_LENGTH", "800")
	metricsBatchLengthInt, utils.Err = strconv.Atoi(metricsBatchLength)
	utils.MaybeDie(utils.Err, "cannot convert string metricsBatchLength to int")
}
