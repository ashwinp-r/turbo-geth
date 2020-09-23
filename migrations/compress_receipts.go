package migrations

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/common/etl"
	"github.com/ledgerwatch/turbo-geth/core/rawdb"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/log"
	"github.com/ledgerwatch/turbo-geth/rlp"
	"time"
)

var receiptLeadingZeroes = Migration{
	Name: "receipt_leading_zeroes_1",
	Up: func(db ethdb.Database, datadir string, OnLoadCommit etl.LoadCommitHandler) error {
		if exists, err := db.(ethdb.BucketsMigrator).BucketExists(dbutils.BlockReceiptsPrefixOld1); err != nil {
			return err
		} else if !exists {
			return OnLoadCommit(db, nil, true)
		}

		if err := db.(ethdb.BucketsMigrator).ClearBuckets(dbutils.BlockReceiptsPrefix); err != nil {
			return err
		}

		logEvery := time.NewTicker(30 * time.Second)
		defer logEvery.Stop()

		tx := db.(ethdb.HasTx).Tx()
		receiptsOld := tx.Cursor(dbutils.BlockReceiptsPrefixOld1)
		receipts := tx.Cursor(dbutils.BlockReceiptsPrefix)
		for k, v, err := receiptsOld.First(); k != nil; k, v, err = receiptsOld.Next() {
			if err != nil {
				return err
			}

			blockHashBytes := k[len(k)-32:]
			blockNum64Bytes := k[:len(k)-32]
			blockNum32Bytes := k[4:8]
			blockNum := binary.BigEndian.Uint64(blockNum64Bytes)
			canonicalHash := rawdb.ReadCanonicalHash(db, blockNum)
			if !bytes.Equal(blockHashBytes, canonicalHash[:]) {
				return nil
			}

			// Decode the receipts by legacy data type
			storageReceiptsLegacy := []*types.DeprecatedReceiptForStorage1{}
			if err := rlp.DecodeBytes(v, &storageReceiptsLegacy); err != nil {
				return fmt.Errorf("invalid receipt array RLP: %w, blockNum=%d", err, blockNum)
			}

			// Encode by new data type
			storageReceipts := make([]*types.ReceiptForStorage, len(storageReceiptsLegacy))
			for i, r := range storageReceiptsLegacy {
				storageReceipts[i] = (*types.ReceiptForStorage)(r)
			}
			newV, err := rlp.EncodeToBytes(storageReceipts)
			if err != nil {
				log.Crit("Failed to encode block receipts", "err", err)
			}

			if err := rlp.DecodeBytes(v, &storageReceipts); err != nil {
				return fmt.Errorf("invalid receipt array RLP: %w, blockNum=%d", err, blockNum)
			}

			if err := receipts.Append(blockNum32Bytes, newV); err != nil {
				return err
			}

			select {
			default:
			case <-logEvery.C:
				sz, _ := tx.BucketSize(dbutils.BlockReceiptsPrefix)
				log.Info("Progress", "blockNum", blockNum, "bucketSize", common.StorageSize(sz))
			}
		}

		//if err := db.(ethdb.BucketsMigrator).DropBuckets(dbutils.BlockReceiptsPrefixOld1); err != nil {
		//	return err
		//}
		return OnLoadCommit(db, nil, true)
	},
}
