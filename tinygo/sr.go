// Description: Avro utils
// Author: Pixie79
// ============================================================================
// package tinygo

package tinygo

import (
	"fmt"
	"strconv"

	"github.com/pixie79/data-utils/utils"
	sr "github.com/redpanda-data/redpanda/src/go/transform-sdk/sr"
)

// GetSchemaTiny retrieves the schema with the given ID from the specified URL.
//
// Parameters:
// - id: The ID of the schema (as a string).
// - url: The URL of the Schema Registry.
//
// Returns:
// - The retrieved schema (as a string).
func GetSchemaTiny(id string, url string) string {
	registry := sr.NewClient()
	schemaIdInt, err := strconv.Atoi(id)
	utils.Logger.Debug(fmt.Sprintf("Schema ID: %s", id))
	utils.MaybeDie(err, fmt.Sprintf("SCHEMA_ID not an integer: %s", id))
	schema, err := registry.LookupSchemaById(schemaIdInt)
	utils.MaybeDie(err, fmt.Sprintf("Unable to retrieve schema for ID: %s", id))
	return schema.Schema
}
