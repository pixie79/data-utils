// Description: Generic utils functions
// Author: Pixie79
// ============================================================================
// package utils

package utils

import (
	"bytes"
	"encoding/gob"

	tuUtils "github.com/pixie79/tiny-utils/utils"
)

func CreateBytes(data any) []byte {
	var envBuffer bytes.Buffer
	encData := gob.NewEncoder(&envBuffer)
	err := encData.Encode(data)
	tuUtils.MaybeDie(err, "encoding to bytes failed")
	return envBuffer.Bytes()
}
