package chain

import (
	"context"
	"testing"

	"github.com/bnb-chain/greenfield-go-sdk/client/test"
	"github.com/bnb-chain/greenfield-go-sdk/keys"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
)

func TestSendTokenSucceedWithSimulatedGas(t *testing.T) {
	km, err := keys.NewPrivateKeyManager(test.TEST_PRIVATE_KEY)
	assert.NoError(t, err)
	gnfdClient, err := NewGreenfieldClient(test.TEST_GRPC_ADDR, test.TEST_CHAIN_ID,
		WithKeyManager(km),
	)
	assert.NoError(t, err)
	to, err := sdk.AccAddressFromHexUnsafe(test.TEST_ADDR)
	assert.NoError(t, err)
	transfer := banktypes.NewMsgSend(km.GetAddr(), to, sdk.NewCoins(sdk.NewInt64Coin(test.TEST_DENOM, 12)))
	response, err := gnfdClient.BroadcastTx(context.Background(), []sdk.Msg{transfer}, nil)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), response.TxResponse.Code)
	t.Log(response.TxResponse.String())
}

func TestSendTokenWithTxOptionSucceed(t *testing.T) {
	km, err := keys.NewPrivateKeyManager(test.TEST_PRIVATE_KEY)
	assert.NoError(t, err)
	gnfdClient, err := NewGreenfieldClient(test.TEST_GRPC_ADDR, test.TEST_CHAIN_ID, WithKeyManager(km))
	assert.NoError(t, err)
	to, err := sdk.AccAddressFromHexUnsafe(test.TEST_ADDR)
	assert.NoError(t, err)
	transfer := banktypes.NewMsgSend(km.GetAddr(), to, sdk.NewCoins(sdk.NewInt64Coin(test.TEST_DENOM, 100)))
	payerAddr, err := sdk.AccAddressFromHexUnsafe(km.GetAddr().String())
	assert.NoError(t, err)
	mode := tx.BroadcastMode_BROADCAST_MODE_ASYNC
	txOpt := &types.TxOption{
		Mode:     &mode,
		GasLimit: 123456,
		Memo:     "test",
		FeePayer: payerAddr,
	}
	response, err := gnfdClient.BroadcastTx(context.Background(), []sdk.Msg{transfer}, txOpt)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), response.TxResponse.Code)
	t.Log(response.TxResponse.String())
}

func TestSimulateTx(t *testing.T) {
	km, err := keys.NewPrivateKeyManager(test.TEST_PRIVATE_KEY)
	assert.NoError(t, err)
	gnfdClient, err := NewGreenfieldClient(test.TEST_GRPC_ADDR, test.TEST_CHAIN_ID, WithKeyManager(km))
	assert.NoError(t, err)
	to, err := sdk.AccAddressFromHexUnsafe(test.TEST_ADDR)
	assert.NoError(t, err)
	transfer := banktypes.NewMsgSend(km.GetAddr(), to, sdk.NewCoins(sdk.NewInt64Coin(test.TEST_DENOM, 100)))
	simulateRes, err := gnfdClient.SimulateTx(context.Background(), []sdk.Msg{transfer}, nil)
	assert.NoError(t, err)
	t.Log(simulateRes.GasInfo.String())
}
