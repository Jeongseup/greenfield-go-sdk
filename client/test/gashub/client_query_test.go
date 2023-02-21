package gashub

import (
	"context"
	gnfdclient "github.com/bnb-chain/greenfield-go-sdk/client/chain"
	"github.com/bnb-chain/greenfield-go-sdk/client/test"
	gashubtypes "github.com/cosmos/cosmos-sdk/x/gashub/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGashubParams(t *testing.T) {
	client := gnfdclient.NewChainClient(test.TEST_GRPC_ADDR, test.TEST_CHAIN_ID)

	query := gashubtypes.QueryParamsRequest{}
	res, err := client.GashubQueryClient.Params(context.Background(), &query)
	assert.NoError(t, err)

	t.Log(res.String())
}
