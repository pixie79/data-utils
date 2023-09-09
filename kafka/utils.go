// Description: Kafka utils
// Author: Pixie79
// ============================================================================
// package kafka

package kafka

import (
	"binary"
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

// ProduceMessages produces messages to Kafka
func ProduceMessages(ctx context.Context, client *kgo.Client, record []*kgo.Record) error {
	var (
		errPromise kgo.FirstErrPromise
	)

	for _, s := range record {
		client.Produce(ctx, s, errPromise.Promise())
	}
	// Wait for all the records to be flushed or for an error to be returned.
	return errPromise.Err()
}

// RollbackTransaction rolls back a transaction
func RollbackTransaction(client *kgo.Client) error {
	// Background context is used because cancelling either of these operations can result
	// in buffered messages being added to the next transaction.
	ctx := context.Background()
	// Remove any records that have not yet been flushed.
	err := client.AbortBufferedRecords(ctx)
	if err != nil {
		return err
	}
	// End the transaction itself so that flushed records will be committed.
	if err := client.EndTransaction(ctx, kgo.TryAbort); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func IntToBytes(n int) []byte {
	byteArray := make([]byte, 4)
	binary.BigEndian.PutUint32(byteArray, uint32(n))
	return byteArray
}

func AddZeroToStart(byteArray []byte) []byte {
	return append([]byte{0}, byteArray...)
}

func EncodedBuffer(i) []byte {
	return addZeroToStart(intToBytes(i))
}