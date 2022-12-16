package staking_view

import (
	"context"
	"fmt"
	"log"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type DelegatorValidatorInfo struct {
	DelegatorAddress string
	TotalPower       *string
	ValidatorAddress string
	Share            string
	Amount           string // Delegated to this validator
}

type DelegatorExportInfo struct {
	Address                  string
	TotalPower               string
	ValidatorAddress         string
	DelegatedToThisValidator string
}

// var DelegatorExportInfos []DelegatorExportInfo
type DelegatorValidatorInfos []DelegatorValidatorInfo

//var DelegatorValidatorInfos []DelegatorValidatorInfo
type DelegatorInfoMapType map[string]DelegatorValidatorInfos

var DelegatorInfoMap = make(map[string]DelegatorValidatorInfos)

func OutDelegatorMap() {
	fmt.Printf("%v,%v,%v,%v,%v\n", "DelegatorAddress", "TotalPower", "ValidatorAddress", "Share", "Amount")

	for _, dvis := range DelegatorInfoMap {
		// fmt.Printf("Delegator : %v\n", k)
		for _, dvi := range dvis {
			fmt.Printf("%v,%v,%v,%v,%v\n", dvi.DelegatorAddress, *dvi.TotalPower, dvi.ValidatorAddress, dvi.Share, dvi.Amount)
		}
	}
}

func GetDelegatedInfos(sqc staking.QueryClient, vxis []ValidatorExportInfo) ([]DelegatorExportInfo, error) {
	// For store delegator address
	delegatorAddressMap := make(map[string]bool)

	for _, val := range vxis {
		valDelRes, err := sqc.ValidatorDelegations(
			context.Background(),
			&staking.QueryValidatorDelegationsRequest{
				ValidatorAddr: val.OperatorAddress,
				Pagination:    &query.PageRequest{Limit: 500, CountTotal: true},
			},
		)

		if err != nil {
			return nil, nil
		}

		var delResponses []staking.DelegationResponse = valDelRes.DelegationResponses
		for _, delRes := range delResponses {

			// var staking.Delegation
			del := delRes.Delegation
			delegatorAddressMap[del.DelegatorAddress] = true
		}
	}

	for delAddress, _ := range delegatorAddressMap {
		if DelegatorInfoMap[delAddress] == nil {
			var tmp DelegatorValidatorInfos = DelegatorValidatorInfos{}
			DelegatorInfoMap[delAddress] = tmp
		}

		fmt.Printf("Delegator Address: %v", delAddress)

		/// Process this address
		stakingRes, err := sqc.DelegatorDelegations(
			context.Background(),
			&staking.QueryDelegatorDelegationsRequest{
				DelegatorAddr: delAddress,
				Pagination:    &query.PageRequest{Limit: 500, CountTotal: true},
			},
		)
		if err != nil {
			log.Fatal("error while querying delegations, reason: ", err)
		}

		var delResponses []staking.DelegationResponse = stakingRes.DelegationResponses
		var total sdktypes.Coin = sdktypes.NewCoin("uqsr", sdktypes.ZeroInt())
		tmp := DelegatorInfoMap[delAddress]
		var totalPower = new(string)
		for _, d := range delResponses {
			var dvi DelegatorValidatorInfo
			dvi.DelegatorAddress = delAddress
			dvi.ValidatorAddress = d.Delegation.ValidatorAddress
			dvi.Amount = d.Balance.Amount.String()
			dvi.Share = d.Delegation.Shares.String()
			dvi.TotalPower = totalPower
			tmp = append(tmp, dvi)
			total = total.Add(d.Balance)
		}
		DelegatorInfoMap[delAddress] = tmp
		*totalPower = total.Amount.String() // Pointer will reflect to all keys in DelegatorInfoMap
		fmt.Printf(", Total Delegation Amount : %v,= %v \n", total.String(), *totalPower)

	}
	// OutDelegatorMap()
	return nil, nil
}

func UpdateSelfDelegation(sqc staking.QueryClient, vxis []ValidatorExportInfo) {
	for i, val := range vxis {
		stakingRes, err := sqc.DelegatorDelegations(
			context.Background(),
			&staking.QueryDelegatorDelegationsRequest{
				DelegatorAddr: val.AccountAddress,
				Pagination:    &query.PageRequest{Limit: 500, CountTotal: true},
			},
		)
		if err != nil {
			log.Fatal("error while querying delegations, reason: ", err)
		}

		var delRes []staking.DelegationResponse = stakingRes.DelegationResponses

		// for each validator, we look if their account address is in the set of delegators
		for _, del := range delRes {
			if del.Delegation.DelegatorAddress == val.AccountAddress {
				vxis[i].SelfDelegation = del.Balance.String() // append the self-delegation to the exportable slice
			} else {
				log.Fatal("delegtor address does not equal validator address")
			}
		}
	}
}
