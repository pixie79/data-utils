// Description: Datadog setup
// Author: Pixie79
// ============================================================================
// package datadog

package datadog

import (
	"strconv"

	"github.com/pixie79/data-utils/utils"
)

var (
	metricsBatchLength    string // metricsBatchLength is the number of metrics to send to Datadog in a single batch
	metricsBatchLengthInt int    // metricsBatchLengthInt is the number of metrics to send to Datadog in a single batch integer
)

// init sets the metricsBatchLengthInt variable
func init() {
	metricsBatchLength = utils.GetEnv("METRICS_BATCH_LENGTH", "800")
	metricsBatchLengthInt, utils.Err = strconv.Atoi(metricsBatchLength)
	utils.MaybeDie(utils.Err, "cannot convert string metricsBatchLength to int")
}
