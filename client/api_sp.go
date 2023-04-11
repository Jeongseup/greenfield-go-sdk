package client

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/bnb-chain/greenfield-go-sdk/pkg/utils"
	spTypes "github.com/bnb-chain/greenfield/x/sp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SP interface {
	// ListSP return the storage provider info on chain
	// isInService indicates if only display the sp with STATUS_IN_SERVICE status
	ListSP(ctx context.Context, isInService bool) ([]spTypes.StorageProvider, error)
	// GetSPInfo return the sp info the sp chain address
	GetSPInfo(ctx context.Context, SPAddr sdk.AccAddress) (*spTypes.StorageProvider, error)
	// GetSpAddrFromEndpoint return the chain addr according to the SP endpoint
	GetSpAddrFromEndpoint(ctx context.Context, spEndpoint string) (sdk.AccAddress, error)
	// GetStoragePrice returns the storage price for a particular storage provider, including update time, read price, store price and .etc.
	GetStoragePrice(ctx context.Context, SPAddr sdk.AccAddress) (*spTypes.SpStoragePrice, error)
}

func (c *client) GetStoragePrice(ctx context.Context, SPAddr sdk.AccAddress) (*spTypes.SpStoragePrice, error) {
	resp, err := c.chainClient.QueryGetSpStoragePriceByTime(ctx, &spTypes.QueryGetSpStoragePriceByTimeRequest{
		SpAddr:    SPAddr.String(),
		Timestamp: 0,
	})
	if err != nil {
		return nil, err
	}
	return &resp.SpStoragePrice, nil
}

// ListSP return the storage provider info on chain
// isInService indicates if only display the sp with STATUS_IN_SERVICE status
func (c *client) ListSP(ctx context.Context, isInService bool) ([]spTypes.StorageProvider, error) {
	request := &spTypes.QueryStorageProvidersRequest{}
	gnfdRep, err := c.chainClient.StorageProviders(ctx, request)
	if err != nil {
		return nil, err
	}

	spList := gnfdRep.GetSps()
	spInfoList := make([]spTypes.StorageProvider, 0)
	for _, info := range spList {
		if isInService && info.Status != spTypes.STATUS_IN_SERVICE {
			continue
		}
		spInfoList = append(spInfoList, *info)
	}

	return spInfoList, nil
}

// GetSpAddrFromEndpoint return the chain addr according to the SP endpoint
func (c *client) GetSpAddrFromEndpoint(ctx context.Context, spEndpoint string) (sdk.AccAddress, error) {
	spList, err := c.ListSP(ctx, false)
	if err != nil {
		return nil, err
	}

	if strings.Contains(spEndpoint, "http") {
		s := strings.Split(spEndpoint, "//")
		spEndpoint = s[1]
	}

	for _, spInfo := range spList {
		endpoint := spInfo.GetEndpoint()
		if strings.Contains(endpoint, "http") {
			s := strings.Split(endpoint, "//")
			endpoint = s[1]
		}
		if endpoint == spEndpoint {
			addr := spInfo.GetOperatorAddress()
			if addr == "" {
				return nil, errors.New("fail to get addr")
			}
			return sdk.MustAccAddressFromHex(spInfo.GetOperatorAddress()), nil
		}
	}
	return nil, errors.New("fail to get addr")
}

// GetSPInfo return the sp info the sp chain address
func (c *client) GetSPInfo(ctx context.Context, SPAddr sdk.AccAddress) (*spTypes.StorageProvider, error) {
	request := &spTypes.QueryStorageProviderRequest{
		SpAddress: SPAddr.String(),
	}

	gnfdRep, err := c.chainClient.StorageProvider(ctx, request)
	if err != nil {
		return nil, err
	}

	return gnfdRep.StorageProvider, nil
}

func (c *client) getSPUrlInfo() (map[string]*url.URL, error) {
	ctx := context.Background()
	spInfo := make(map[string]*url.URL, 0)
	request := &spTypes.QueryStorageProvidersRequest{}
	gnfdRep, err := c.chainClient.StorageProviders(ctx, request)
	if err != nil {
		return nil, err
	}
	spList := gnfdRep.GetSps()
	for _, info := range spList {
		endpoint := info.Endpoint
		urlInfo, urlErr := utils.GetEndpointURL(endpoint, c.secure)
		if urlErr != nil {
			return nil, urlErr
		}
		spInfo[info.GetOperator().String()] = urlInfo
	}

	return spInfo, nil
}
