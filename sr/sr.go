// Description: Schema Registry utils
// Author: Pixie79
// ============================================================================
// package sr

package avro

import (
	"fmt"
	"strconv"

	sr "github.com/landoop/schema-registry"
	tuUtils "github.com/pixie79/tiny-utils/utils"
)

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
	tuUtils.MaybeDie(err, fmt.Sprintf("Cannot connect to Schema Registry: %+v", err))
	schemaIdInt, err := strconv.Atoi(id)
	tuUtils.Print("DEBUG", fmt.Sprintf("Schema ID: %s", id))
	tuUtils.MaybeDie(err, fmt.Sprintf("SCHEMA_ID not an integer: %s", id))
	schema, err := registry.GetSchemaByID(schemaIdInt)
	tuUtils.MaybeDie(err, fmt.Sprintf("Unable to retrieve schema for ID: %s", id))
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
