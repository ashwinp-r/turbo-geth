package bitmapdb_test

import (
	"testing"

	"github.com/RoaringBitmap/roaring"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/ethdb/bitmapdb"
	"github.com/stretchr/testify/require"
)

func TestRoaringBitmapAddOffset(t *testing.T) {
	db := ethdb.NewMemDatabase()
	defer db.Close()

	{ // Simple encoding
		k := []byte{1}
		bm1 := roaring.NewBitmap()
		bm1.Add(940287)
		bm1.Add(940288)
		err := bitmapdb.Simple.PutMergeByOr(db, dbutils.LogIndex, k, bm1)
		require.NoError(t, err)

		bm2 := roaring.NewBitmap()
		bm2.Add(1_000_000)
		bm2.Add(1_000_001)
		bm2.Add(1_000_002)
		err = bitmapdb.Simple.PutMergeByOr(db, dbutils.LogIndex, k, bm2)
		require.NoError(t, err)
		err = bitmapdb.Simple.RemoveRange(db, dbutils.LogIndex, k, 1_000_001, 1_000_003) // [from, to)
		require.NoError(t, err)

		bm3, err := bitmapdb.Simple.Get(db, dbutils.LogIndex, k)
		require.NoError(t, err)
		arr := bm3.ToArray()
		require.Equal(t, 3, len(arr))
		require.Equal(t, uint32(940287), arr[0])
		require.Equal(t, uint32(940288), arr[1])
		require.Equal(t, uint32(1_000_000), arr[2])
	}

	{ // NoLeadingZeroes encoding
		k := []byte{2}
		bm1 := roaring.NewBitmap()
		bm1.Add(940287)
		bm1.Add(940288)
		err := bitmapdb.NoLeadingZeroes.PutMergeByOr(db, dbutils.LogIndex, k, bm1)
		require.NoError(t, err)

		bm2 := roaring.NewBitmap()
		bm2.Add(1_000_000)
		bm2.Add(1_000_001)
		bm2.Add(1_000_002)
		err = bitmapdb.NoLeadingZeroes.PutMergeByOr(db, dbutils.LogIndex, k, bm2)
		require.NoError(t, err)
		err = bitmapdb.NoLeadingZeroes.RemoveRange(db, dbutils.LogIndex, k, 1_000_001, 1_000_003) // [from, to)
		require.NoError(t, err)

		bm3, err := bitmapdb.NoLeadingZeroes.Get(db, dbutils.LogIndex, k)
		require.NoError(t, err)
		arr := bm3.ToArray()
		require.Equal(t, 3, len(arr))
		require.Equal(t, uint32(940287), arr[0])
		require.Equal(t, uint32(940288), arr[1])
		require.Equal(t, uint32(1_000_000), arr[2])
	}

}

func TestRemoveRange(t *testing.T) {
	db := ethdb.NewMemDatabase()
	defer db.Close()

	{ // Simple encoding
		k := []byte{1}
		bm1 := roaring.NewBitmap()
		bm1.Add(940287)
		bm1.Add(940288)
		err := bitmapdb.Simple.PutMergeByOr(db, dbutils.LogIndex, k, bm1)
		require.NoError(t, err)

		bm2 := roaring.NewBitmap()
		bm2.Add(1_000_000)
		bm2.Add(1_000_001)
		bm2.Add(1_000_002)
		err = bitmapdb.Simple.PutMergeByOr(db, dbutils.LogIndex, k, bm2)
		require.NoError(t, err)
		err = bitmapdb.Simple.RemoveRange(db, dbutils.LogIndex, k, 1_000_001, 1_000_002) // [from, to)
		require.NoError(t, err)

		bm3, err := bitmapdb.Simple.Get(db, dbutils.LogIndex, k)
		require.NoError(t, err)
		arr := bm3.ToArray()
		require.Equal(t, 4, len(arr))
		require.Equal(t, uint32(940287), arr[0])
		require.Equal(t, uint32(940288), arr[1])
		require.Equal(t, uint32(1_000_000), arr[2])
		require.Equal(t, uint32(1_000_002), arr[3])
	}

	{ // NoLeadingZeroes encoding
		k := []byte{2}
		bm1 := roaring.NewBitmap()
		bm1.Add(940287)
		bm1.Add(940288)
		err := bitmapdb.NoLeadingZeroes.PutMergeByOr(db, dbutils.LogIndex, k, bm1)
		require.NoError(t, err)

		bm2 := roaring.NewBitmap()
		bm2.Add(1_000_000)
		bm2.Add(1_000_001)
		bm2.Add(1_000_002)
		err = bitmapdb.NoLeadingZeroes.PutMergeByOr(db, dbutils.LogIndex, k, bm2)
		require.NoError(t, err)
		err = bitmapdb.NoLeadingZeroes.RemoveRange(db, dbutils.LogIndex, k, 1_000_001, 1_000_002) // [from, to)
		require.NoError(t, err)

		bm3, err := bitmapdb.NoLeadingZeroes.Get(db, dbutils.LogIndex, k)
		require.NoError(t, err)
		arr := bm3.ToArray()
		require.Equal(t, 4, len(arr))
		require.Equal(t, uint32(940287), arr[0])
		require.Equal(t, uint32(940288), arr[1])
		require.Equal(t, uint32(1_000_000), arr[2])
		require.Equal(t, uint32(1_000_002), arr[3])
	}
}
