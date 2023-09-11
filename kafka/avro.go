// Description: Avro utils
// Author: Pixie79
// ============================================================================
// package kafka

package kafka

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	sr "github.com/landoop/schema-registry"
	avro "github.com/linkedin/goavro/v2"
	"github.com/pixie79/data-utils/utils"
	srt "github.com/redpanda-data/redpanda/src/go/transform-sdk/sr"
)

// GetSchemaIdFromPayload returns the schema id from the payload
func GetSchemaIdFromPayload(msg []byte) int {
	schemaID := binary.BigEndian.Uint32(msg[1:5])
	return int(schemaID)
}

// GetSchemaTiny retrieves the schema with the given ID from the specified URL.
//
// Parameters:
// - id: The ID of the schema (as a string).
// - url: The URL of the Schema Registry.
//
// Returns:
// - The retrieved schema (as a string).
func GetSchemaTiny(id string, url string) string {
	registry := srt.NewClient()
	schemaIdInt, err := strconv.Atoi(id)
	utils.Logger.Debug(fmt.Sprintf("Schema ID: %s", id))
	utils.MaybeDie(err, fmt.Sprintf("SCHEMA_ID not an integer: %s", id))
	schema, err := registry.LookupSchemaById(schemaIdInt)
	utils.MaybeDie(err, fmt.Sprintf("Unable to retrieve schema for ID: %s", id))
	return schema.Schema
}

// GetSchema retrieves the schema with the given ID from the specified URL.
//
// Parameters:
// - id: The ID of the schema (as a string).
// - url: The URL of the Schema Registry.
//
// Returns:
// - The retrieved schema (as a string).
func GetSchema(id string, url string) string {
	registry, err := sr.NewClient(url)
	utils.MaybeDie(err, fmt.Sprintf("Cannot connect to Schema Registry: %+v", err))
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
// func GetSchemaCache(id string, url string) string {
// 	schemaCache := cache.New(2*time.Hour, 10*time.Minute)
// 	schema, found := schemaCache.Get(id)
// 	if found {
// 		return schema
// 	}
// 	schema = GetSchema(id)
// 	schemaCache.Set(id, schema, cache.DefaultExpiration)

// 	return schema
// }

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
	// utils.Logger.Debug(prettyPrint(native))
	nestedMap, ok := native.(map[string]interface{})
	if !ok {
		utils.Die("Unable to convert native to map[string]interface{}")
	}
	return nestedMap
}
