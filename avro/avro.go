// Description: Avro utils
// Author: Pixie79
// ============================================================================
// package avro

package avro

import (
	"encoding/binary"
	"strings"

	avro "github.com/linkedin/goavro/v2"
	"github.com/pixie79/data-utils/utils"
)

// GetSchemaIdFromPayload returns the schema id from the payload
func GetSchemaIdFromPayload(msg []byte) int {
	schemaID := binary.BigEndian.Uint32(msg[1:5])
	return int(schemaID)
}

// decodeAvro decodes an Avro event using the provided schema and returns a nested map[string]interface{}.
//
// Parameters:
// - schema: The Avro schema used for decoding the event (string).
// - event: The Avro event to be decoded ([]byte).
//
// Returns:
// - nestedMap: The decoded event as a nested map[string]interface{}.
func DecodeAvro(schema string, event []byte) map[string]interface{} {
	sourceCodec, err := avro.NewCodec(schema)
	utils.MaybeDie(err, "Error creating Avro codec")

	strEvent := strings.Replace(string(event), "\"", "", -1)
	newEvent, err := utils.B64DecodeMsg(strEvent, 5)
	utils.MaybeDie(err, "Error decoding base64")
	native, _, err := sourceCodec.NativeFromBinary(newEvent)
	utils.MaybeDie(err, "Error creating native from binary")
	nestedMap, ok := native.(map[string]interface{})
	if !ok {
		utils.Die("Unable to convert native to map[string]interface{}")
	}
	return nestedMap
}
