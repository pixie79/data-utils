// Description: Avro utils
// Author: Pixie79
// ============================================================================
// package kafka

package kafka

import (
	"context"
	"encoding/binary"

	"github.com/hamba/avro/v2"
	"github.com/pixie79/data-utils/utils"
	"github.com/twmb/franz-go/pkg/sr"
)

// FetchSchema fetches the schema from the schema registry
func FetchSchema(schemaId int) avro.Schema {
	var url = "http://localhost:18081"
	rcl, err := sr.NewClient(sr.URLs(url))
	utils.MaybeDie(err, "unable to create schema registry client")
	ss, err := rcl.SchemaByID(context.Background(), schemaId)
	utils.MaybeDie(err, "unable to parse avro schema")
	kafkaSchema, err := avro.Parse(ss.Schema)
	utils.MaybeDie(err, "unable to parse avro schema")
	return kafkaSchema
}

// GetSchemaIdFromPayload returns the schema id from the payload
func GetSchemaIdFromPayload(msg []byte) int {
	schemaID := binary.BigEndian.Uint32(msg[1:5])
	return int(schemaID)
}
