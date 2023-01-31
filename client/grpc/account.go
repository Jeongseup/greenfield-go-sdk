package client

import (
	"context"
	"github.com/bnb-chain/gnfd-go-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (c *GreenfieldClient) Account(address string) (authtypes.AccountI, error) {
	acct, err := c.AuthQueryClient.Account(context.Background(), &authtypes.QueryAccountRequest{Address: address})
	if err != nil {
		return nil, err
	}
	var account authtypes.AccountI
	if err := types.Cdc().InterfaceRegistry().UnpackAny(acct.Account, &account); err != nil {
		return nil, err
	}
	return account, nil
}
