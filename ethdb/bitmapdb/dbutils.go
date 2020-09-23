package bitmapdb

import (
	"bytes"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"time"
)

// PutMergeByOr - puts bitmap with recent changes into database by merging it with existing bitmap. Merge by OR.
func PutMergeByOr(c ethdb.Cursor, k []byte, delta *roaring.Bitmap) error {
	defer func(c uint64, t time.Time) { fmt.Printf("dbutils.go:13: %d %s\n", c, time.Since(t)) }(delta.GetCardinality(), time.Now())
	v, err := c.SeekExact(k)
	if err != nil {
		panic(err)
	}

	if len(v) > 0 {
		existing := roaring.New()
		_, err = existing.ReadFrom(bytes.NewReader(v))
		if err != nil {
			panic(err)
		}

		delta.Or(existing)
	}

	bufBytes, err := c.Reserve(k, int(delta.GetSerializedSizeInBytes()))
	if err != nil {
		panic(err)
	}

	_, err = delta.WriteTo(bytes.NewBuffer(bufBytes[:0]))
	if err != nil {
		return err
	}
	return nil
}

// RemoveRange - gets existing bitmap in db and call RemoveRange operator on it.
func RemoveRange(c ethdb.Cursor, k []byte, from, to uint64) error {
	v, err := c.SeekExact(k)
	if err != nil {
		panic(err)
	}

	bm := roaring.New()
	if len(v) > 0 {
		_, err = bm.ReadFrom(bytes.NewReader(v))
		if err != nil {
			panic(err)
		}

		bm.RemoveRange(from, to)
	}

	bufBytes, err := c.Reserve(k, int(bm.GetSerializedSizeInBytes()))
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(bufBytes[:0])
	_, err = bm.WriteTo(buf)
	if err != nil {
		return err
	}
	return nil
}

// Get - puts bitmap into database. If database already has such key - does merge by OR.
func Get(db ethdb.Getter, bucket string, k []byte) (*roaring.Bitmap, error) {
	bitmapBytes, err := db.Get(bucket, k)
	if err != nil {
		return nil, err
	}
	m := roaring.New()
	_, err = m.ReadFrom(bytes.NewReader(bitmapBytes))
	return m, err
}
