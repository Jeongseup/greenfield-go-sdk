package chain

import (
	"github.com/bnb-chain/greenfield-go-sdk/client/test"
	"github.com/bnb-chain/greenfield-go-sdk/keys"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestSendTokenSucceedWithSimulatedGas(t *testing.T) {
	km, err := keys.NewPrivateKeyManager(test.TEST_PRIVATE_KEY)
	assert.NoError(t, err)
	gnfdClient := NewGreenfieldClient(test.TEST_GRPC_ADDR, test.TEST_CHAIN_ID, WithKeyManager(km),
		WithGrpcDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	to, err := sdk.AccAddressFromHexUnsafe(test.TEST_ADDR)
	assert.NoError(t, err)
	transfer := banktypes.NewMsgSend(km.GetAddr(), to, sdk.NewCoins(sdk.NewInt64Coin(test.TEST_DENOM, 12)))
	response, err := gnfdClient.BroadcastTx([]sdk.Msg{transfer}, nil)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), response.TxResponse.Code)
	t.Log(response.TxResponse.String())
}

func TestSendTokenWithTxOptionSucceed(t *testing.T) {
	km, err := keys.NewPrivateKeyManager(test.TEST_PRIVATE_KEY)
	assert.NoError(t, err)
	gnfdClient := NewGreenfieldClient(test.TEST_GRPC_ADDR, test.TEST_CHAIN_ID, WithKeyManager(km),
		WithGrpcDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	to, err := sdk.AccAddressFromHexUnsafe(test.TEST_ADDR)
	assert.NoError(t, err)
	transfer := banktypes.NewMsgSend(km.GetAddr(), to, sdk.NewCoins(sdk.NewInt64Coin(test.TEST_DENOM, 100)))
	payerAddr, err := sdk.AccAddressFromHexUnsafe(km.GetAddr().String())
	mode := tx.BroadcastMode_BROADCAST_MODE_ASYNC
	txOpt := &types.TxOption{
		Mode:      &mode,
		GasLimit:  123456,
		Memo:      "test",
		FeeAmount: sdk.Coins{{test.TEST_DENOM, sdk.NewInt(1)}},
		FeePayer:  payerAddr,
	}
	response, err := gnfdClient.BroadcastTx([]sdk.Msg{transfer}, txOpt)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), response.TxResponse.Code)
	t.Log(response.TxResponse.String())
}

func TestSimulateTx(t *testing.T) {
	km, err := keys.NewPrivateKeyManager(test.TEST_PRIVATE_KEY)
	assert.NoError(t, err)
	gnfdClient := NewGreenfieldClient(test.TEST_GRPC_ADDR, test.TEST_CHAIN_ID, WithKeyManager(km),
		WithGrpcDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	to, err := sdk.AccAddressFromHexUnsafe(test.TEST_ADDR)
	assert.NoError(t, err)
	transfer := banktypes.NewMsgSend(km.GetAddr(), to, sdk.NewCoins(sdk.NewInt64Coin(test.TEST_DENOM, 100)))
	simulateRes, err := gnfdClient.SimulateTx([]sdk.Msg{transfer}, nil)
	assert.NoError(t, err)
	t.Log(simulateRes.GasInfo.String())
}

func TestSendTokenSucceedWithSimulatedGas1(t *testing.T) {
	km, err := keys.NewPrivateKeyManager(test.TEST_PRIVATE_KEY)
	assert.NoError(t, err)
	gnfdClients := NewGnfdClients(
		[]string{test.TEST_GRPC_ADDR, test.TEST_GRPC_ADDR2, test.TEST_GRPC_ADDR3},
		[]string{test.TEST_RPC_ADDR, test.TEST_RPC_ADDR2, test.TEST_RPC_ADDR3},
		test.TEST_CHAIN_ID,
		WithKeyManager(km),
		WithGrpcDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	client, err := gnfdClients.GetClient()
	to, err := sdk.AccAddressFromHexUnsafe(test.TEST_ADDR)
	assert.NoError(t, err)
	transfer := banktypes.NewMsgSend(km.GetAddr(), to, sdk.NewCoins(sdk.NewInt64Coin(test.TEST_DENOM, 12)))
	response, err := client.BroadcastTx([]sdk.Msg{transfer}, nil)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), response.TxResponse.Code)
	t.Log(response.TxResponse.String())
}
