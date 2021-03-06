package commands

import (
	"context"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/turbo/rpchelper"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/hexutil"
	"github.com/ledgerwatch/turbo-geth/rpc"
)

func (api *APIImpl) GetBalance(_ context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	blockNumber, _, err := rpchelper.GetBlockNumber(blockNrOrHash, api.dbReader)
	if err != nil {
		return nil, err
	}

	acc, err := rpchelper.GetAccount(api.db, blockNumber, address)
	if err != nil {
		return nil, fmt.Errorf("cant get a balance for account %q for block %v", address.String(), blockNumber)
	}

	return (*hexutil.Big)(acc.Balance.ToBig()), nil
}
