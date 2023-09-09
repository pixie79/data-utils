// Description: Avro utils
// Author: Pixie79
// ============================================================================
// package kafka

package kafka

import (
	"context"
	"encoding/binary"
	sr "github.com/landoop/schema-registry"
	avro "github.com/linkedin/goavro/v2"
	"github.com/pixie79/data-utils/utils"
	"github.com/twmb/franz-go/pkg/sr"
)


// GetSchemaIdFromPayload returns the schema id from the payload
func GetSchemaIdFromPayload(msg []byte) int {
	schemaID := binary.BigEndian.Uint32(msg[1:5])
	return int(schemaID)
}

// GetSchema retrieves the schema with the given ID.
//
// Parameters:
// - id: the ID of the schema to retrieve.
//
// Returns:
// - *sr.Schema: the retrieved schema.
func GetSchema(id string) string {
	registry, err := sr.NewClient()
	maybeDie(err, fmt.Sprintf("Cannot connect to Schema Registry: %+v", err))
	schemaIdInt, err := strconv.Atoi(id)
	utils.Logger.Debug(fmt.Sprintf("Schema ID: %s", id))
	utils.MaybeDie(err, fmt.Sprintf("SCHEMA_ID not an integer: %s", id))
	schema, err := registry.GetSchemaByID(schemaIdInt)
	utils.MaybeDie(err, fmt.Sprintf("Unable to retrieve schema for ID: %s", id))
	return schema
}

// GetSchemaCache retrieves the schema with the given ID.
//
// Parameters:
// - id: the ID of the schema to retrieve.
//
// Returns:
// - *sr.Schema: the retrieved schema.
func GetSchemaCache(id string) string {	
	schema, found := schemaCache.Get(id)
	if found {
		codec, err := avro.NewCodec(schema.(string))
		utils.MaybeDie(err, fmt.Sprintf("Error creating Avro codec: %+v", err))
		return codec
	}
	schema = GetSchema(id)
	schemaCache.Set(id, schema, cache.DefaultExpiration)
	
	return schema
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
	newEvent, err := B64DecodeMsg(strEvent, 5)
	utils.MaybeDie(err, "Error decoding base64")
	native, _, err := sourceCodec.NativeFromBinary(newEvent)
	utils.MaybeDie(err, "Error creating native from binary")
	// utils.Logger.Debug(prettyPrint(native))
	nestedMap, ok := native.(map[string]interface{})
	if !ok {
		utils.Die("Unable to convert native to map[string]interface{}")
	}
	return nestedMap
}