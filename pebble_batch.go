//go:build pebbledb

package db

import "github.com/cockroachdb/pebble"

type pebbleDBBatch struct {
	batch pebble.Batch
}

var _ Batch = (*pebbleDBBatch)(nil)

func newPebbleDBBatch(db *PebbleDB) *pebbleDBBatch {
	return &pebbleDBBatch{
		batch: *db.db.NewBatch(),
	}
}

// Set implements Batch.
func (b *pebbleDBBatch) Set(key, value []byte) error {
	if len(key) == 0 {
		return errKeyEmpty
	}
	if value == nil {
		return errValueNil
	}
	b.batch.Set(key, value, nil)
	return nil
}

// Delete implements Batch.
func (b *pebbleDBBatch) Delete(key []byte) error {
	if len(key) == 0 {
		return errKeyEmpty
	}
	b.batch.Delete(key, nil)
	return nil
}

// Write implements Batch.
func (b *pebbleDBBatch) Write() error {
	err := b.batch.Commit(pebble.NoSync)
	if err != nil {
		return err
	}
	// Make sure batch cannot be used afterwards. Callers should still call Close(), for errors.
	b.Close()
	return nil
}

// WriteSync implements Batch.
func (b *pebbleDBBatch) WriteSync() error {
	err := b.batch.Commit(pebble.Sync)
	if err != nil {
		return err
	}
	// Make sure batch cannot be used afterwards. Callers should still call Close(), for errors.
	return b.Close()
}

// Close implements Batch.
func (b *pebbleDBBatch) Close() error {
	err := b.batch.Close()
	if err != nil {
		return err
	}

	return nil
}
