package staking_view

import (
	"context"
	"errors"
	"log"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc"
)

type ValidatorI interface {
	GetConnection() *grpc.ClientConn
	GetSelfDelegation() (float64, error)
	Output() error
	SortValidatorsByPower(vxis []ValidatorExportInfo)
	ValAddressFromBech32(address, prefix string) (valAddr sdk.ValAddress, err error)
	Bech32ifyAddressBytes(prefix string, address sdk.AccAddress) (string, error)
}

type ValidatorExportInfo struct {
	Moniker         string
	OperatorAddress string
	AccountAddress  string
	TotalDelegation string  // Total abs vting power
	SelfDelegation  string  // Amount of delegation to self.
	VotingPower     float64 // Percentage of total voting power
}

func NewValidatorExportInfo() ValidatorExportInfo {
	return ValidatorExportInfo{}
}

func GetStakingQueryClient(grpcConn *grpc.ClientConn) staking.QueryClient {
	stakingClient := staking.NewQueryClient(grpcConn)
	return stakingClient
}

// func GetActiveValidators(grpcConn *grpc.ClientConn) []staking.Validator {
func GetActiveValidators(sqc staking.QueryClient) []staking.Validator {
	// This creates a gRPC client to query the x/staking/types service.
	// stakingClient := staking.NewQueryClient(grpcConn)
	// stakingClient := GetStakingQueryClient(grpcConn)
	stakingRes, err := sqc.Validators(
		context.Background(),
		&staking.QueryValidatorsRequest{
			Status:     "BOND_STATUS_BONDED", // Get all validators that are active and bonded
			Pagination: &query.PageRequest{Limit: 500, CountTotal: true}},
	)
	if err != nil {
		log.Fatal("error while querying validators, reason: ", err)
	}
	vals := stakingRes.Validators
	sort.SliceStable(vals, func(i, j int) bool {
		return vals[i].Tokens.GT(vals[j].Tokens)
	})
	return vals
}

// Implementation to be improved for time complexity for self delegation
func TransformToValInfo(vs []staking.Validator) []ValidatorExportInfo {
	var vxis = []ValidatorExportInfo{}
	var c sdk.Coin
	var totalPower = getTotalBondedTokens(vs)

	for _, val := range vs {
		accAddStr := deriveValAccAddress(val)
		//accAdd := deriveValAccAddress(val)
		//accAddStr := accAdd.String()                // Added for easy debugging
		bondedTokens := val.BondedTokens().String() // Added for easy debugging
		votingPower := getVotingPower(val, totalPower)
		vxis = append(vxis, ValidatorExportInfo{val.GetMoniker(),
			val.OperatorAddress,
			accAddStr,
			bondedTokens,
			c.Amount.String(),
			votingPower,
		})
	}

	return vxis
}

// This method loops over the validator's bonded tokens and sums them up
func getTotalBondedTokens(vs []staking.Validator) (total int64) {
	for _, val := range vs {
		total += val.GetBondedTokens().BigInt().Int64()
	}
	return
}

// This method calculates the voting power of a validator
// TODO - Use sdk.Dec or math instead
func getVotingPower(val staking.Validator, totalTokens int64) float64 {
	var vp = float64(val.BondedTokens().Int64()) / float64(totalTokens)
	return vp
}

// Conversion of validator address from operatorAddress to validator self delegation address
// func deriveValAccAddress(val staking.Validator) sdk.AccAddress {
func deriveValAccAddress(val staking.Validator) string {
	// valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
	valAddr, err := ValAddressFromBech32(val.OperatorAddress, "quasarvaloper")
	if err != nil {
		log.Fatal("could not get validator address, reason: ", err)
	}
	// accAddr, err := sdk.AccAddressFromHexUnsafe(hex.EncodeToString(valAddr.Bytes()))
	accAddr, err := Bech32ifyAddressBytes("quasar", valAddr.Bytes())
	if err != nil {
		log.Fatal("could not get account address, reason: ", err)
	}
	return accAddr
}

// ValAddressFromBech32 creates a ValAddress from a Bech32 string.
func ValAddressFromBech32(address, prefix string) (valAddr sdk.ValAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.ValAddress{}, errors.New("empty address string is not allowed")
	}

	bz, err := sdk.GetFromBech32(address, prefix)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// Bech32ifyAddressBytes returns a bech32 representation of address bytes.
// Returns an empty sting if the byte slice is 0-length.
// Returns an error if the bech32 conversion
// fails or the prefix is empty.
func Bech32ifyAddressBytes(prefix string, address sdk.AccAddress) (string, error) {
	if address.Empty() {
		return "", nil
	}
	if len(address.Bytes()) == 0 {
		return "", nil
	}
	if len(prefix) == 0 {
		return "", errors.New("prefix cannot be empty")
	}
	return bech32.ConvertAndEncode(prefix, address.Bytes())
}
