package client

import (
	"context"

	"github.com/bnb-chain/greenfield-go-sdk/types"
	gnfdSdkTypes "github.com/bnb-chain/greenfield/sdk/types"
	paymentTypes "github.com/bnb-chain/greenfield/x/payment/types"
	storageTypes "github.com/bnb-chain/greenfield/x/storage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	types3 "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type Account interface {
	BuyQuotaForBucket(ctx context.Context, bucketName string, targetQuota uint64, opt types.BuyQuotaOption) (string, error)
	GetAccount(ctx context.Context, address string) (authTypes.AccountI, error)
	GetAccountBalance(ctx context.Context, address string) (*sdk.Coin, error)
	GetPaymentAccount(ctx context.Context, address string) (*paymentTypes.PaymentAccount, error)
	GetPaymentAccountsByOwner(ctx context.Context, owner string) ([]*paymentTypes.PaymentAccount, error)
	Transfer(ctx context.Context, toAddress string, amount int64, txOption gnfdSdkTypes.TxOption) (*sdk.TxResponse, error)
}

// BuyQuotaForBucket buy the target quota of the specific bucket
// targetQuota indicates the target quota to set for the bucket
func (c *client) BuyQuotaForBucket(ctx context.Context, bucketName string, targetQuota uint64, opt types.BuyQuotaOption) (string, error) {
	bucketInfo, err := c.HeadBucket(ctx, bucketName)
	if err != nil {
		return "", err
	}

	paymentAddr, err := sdk.AccAddressFromHexUnsafe(bucketInfo.PaymentAddress)
	if err != nil {
		return "", err
	}
	updateBucketMsg := storageTypes.NewMsgUpdateBucketInfo(c.defaultAccount.GetAddress(), bucketName, &targetQuota, paymentAddr, bucketInfo.Visibility)

	resp, err := c.chainClient.BroadcastTx(ctx, []sdk.Msg{updateBucketMsg}, opt.TxOpts)
	if err != nil {
		return "", err
	}

	return resp.TxResponse.TxHash, err
}

// GetAccount retrieves account information for a given address.
// It takes a context and an address as input and returns an AccountI interface and an error (if any).
func (c *client) GetAccount(ctx context.Context, address string) (authTypes.AccountI, error) {
	// Call the DefaultAccount method of the chain client with a QueryAccountRequest containing the address.
	response, err := c.chainClient.Account(ctx, &authTypes.QueryAccountRequest{Address: address})
	if err != nil {
		// Return an error if there was an issue retrieving the account.
		return nil, err
	}

	// Unmarshal the raw account data from the response into a BaseAccount object.
	baseAccount := authTypes.BaseAccount{}
	err = c.chainClient.GetCodec().Unmarshal(response.Account.GetValue(), &baseAccount)
	if err != nil {
		// Return an error if there was an issue unmarshalling the account data.
		return nil, err
	}

	// Return the BaseAccount object as an AccountI interface.
	return &baseAccount, err
}

// GetAccountBalance retrieves balance information of an account for a given address.
// It takes a context and an address as input and returns an sdk.Coin interface and an error (if any).
func (c *client) GetAccountBalance(ctx context.Context, address string) (*sdk.Coin, error) {
	response, err := c.chainClient.BankQueryClient.Balance(ctx, &types3.QueryBalanceRequest{Address: address, Denom: gnfdSdkTypes.Denom})
	if err != nil {
		return nil, err
	}

	return response.Balance, nil
}

// GetPaymentAccount function takes a context and an address string as parameters and returns a pointer to a paymentTypes.PaymentAccount struct and an error.
// This function uses the PaymentAccount method of the chainClient field of the client struct to query the payment account associated with the given address.
// If there is an error, the function returns nil and the error. If there is no error, the function returns a pointer to the PaymentAccount struct and nil.
func (c *client) GetPaymentAccount(ctx context.Context, address string) (*paymentTypes.PaymentAccount, error) {
	pa, err := c.chainClient.PaymentAccount(ctx, &paymentTypes.QueryGetPaymentAccountRequest{Addr: address})
	if err != nil {
		return nil, err
	}
	return &pa.PaymentAccount, nil
}

// GetPaymentAccountsByOwner retrieves all payment accounts owned by the given address
// and returns a slice of PaymentAccount pointers and an error (if any).
func (c *client) GetPaymentAccountsByOwner(ctx context.Context, owner string) ([]*paymentTypes.PaymentAccount, error) {
	// Call the GetPaymentAccountsByOwner method of the chain client with a QueryGetPaymentAccountsByOwnerRequest containing the owner address.
	accountsByOwnerResponse, err := c.chainClient.GetPaymentAccountsByOwner(ctx, &paymentTypes.QueryGetPaymentAccountsByOwnerRequest{Owner: owner})
	if err != nil {
		return nil, err
	}

	// Initialize a slice of PaymentAccount pointers.
	paymentAccounts := make([]*paymentTypes.PaymentAccount, 0, len(accountsByOwnerResponse.PaymentAccounts))

	// Iterate over each account address returned in the response.
	for _, accAddress := range accountsByOwnerResponse.PaymentAccounts {
		// Call the GetPaymentAccount method of the client to retrieve the PaymentAccount object for the given address.
		pa, err := c.GetPaymentAccount(ctx, accAddress)
		if err != nil {
			return nil, err
		}
		// Append the PaymentAccount object to the slice.
		paymentAccounts = append(paymentAccounts, pa)
	}

	// Return the slice of PaymentAccount pointers.
	return paymentAccounts, nil
}

// Transfer function takes a context, a toAddress string, an amount of type int64, and a txOption of
// type gnfdSdkTypes.TxOption as parameters and returns a pointer to an sdk.TxResponse struct and an error.
// This function first parses the toAddress parameter into an sdk.AccAddress object, and if there is an error,
// it returns nil and the error.
// Then it generates a MsgSend message using the NewMsgSend method of the types3 package and broadcasts the
// transaction to the chain by calling the BroadcastTx method of the chainClient field of the client struct.
// If there is an error during the broadcasting, the function returns nil and the error. If there is no error,
// the function returns a pointer to the TxResponse struct and nil.//
func (c *client) Transfer(ctx context.Context, toAddress string, amount int64, txOption gnfdSdkTypes.TxOption) (*sdk.TxResponse, error) {
	toAddr, err := sdk.AccAddressFromHexUnsafe(toAddress)
	if err != nil {
		return nil, err
	}
	msgSend := types3.NewMsgSend(c.defaultAccount.GetAddress(), toAddr, sdk.Coins{sdk.Coin{Denom: gnfdSdkTypes.Denom, Amount: sdk.NewInt(amount)}})
	tx, err := c.chainClient.BroadcastTx(ctx, []sdk.Msg{msgSend}, &txOption)
	if err != nil {
		return nil, err
	}
	return tx.TxResponse, nil
}
