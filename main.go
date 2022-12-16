package main

import (

	/*
	   "log"
	   "google.golang.org/grpc"
	   "google.golang.org/grpc/credentials/insecure"
	   "github.com/cosmos/cosmos-sdk/codec"
	   staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	*/

	"fmt"
	"os"

	"github.com/akure/cosmo/connection"
	export "github.com/akure/cosmo/export"
	"github.com/akure/cosmo/staking_view"
	ini "gopkg.in/ini.v1"
)

func main() {

	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	grpcEndPoint := cfg.Section("").Key("grpc_node").String()
	fmt.Println("GRPC_NODE=", grpcEndPoint, "LEN=", len(grpcEndPoint))

	// cosmos-grpc.polkachu.com:14990
	// grpcEndPoint := string("cosmos-grpc.polkachu.com:14990")
	// quasar-testnet-rpc.polkachu.com:443
	//34.175.148.145:9090
	// grpcEndPoint := string("quasar-testnet-grpc.polkachu.com:18290") // quasar-testnet-rpc.polkachu.com:443
	// grpcEndPoint := string("34.175.148.145:9090")
	gc := connection.GrpcClient(grpcEndPoint)
	stakingClient := staking_view.GetStakingQueryClient(gc)
	activeValidators := staking_view.GetActiveValidators(stakingClient)
	validatorInfo := staking_view.TransformToValInfo(activeValidators)
	staking_view.UpdateSelfDelegation(stakingClient, validatorInfo)
	export.WriteValidatorInfo(validatorInfo)
	staking_view.GetDelegatedInfos(stakingClient, validatorInfo)
	export.OutDelegatorMap()
}
