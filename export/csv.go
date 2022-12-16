package export

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/akure/cosmo/staking_view"
)

func WriteValidatorInfo(vxis []staking_view.ValidatorExportInfo) {
	file, err := os.Create("active_validators.csv")
	if err != nil {
		log.Fatal("failed to open file", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	header := []string{"OperatorAddress", "Account Address", "Moniker", "Voting-Power", "Self-Delegation", "Total-Delegation"}
	err = w.Write(header)
	if err != nil {
		log.Fatal("error writing record to file", err)
	}
	fmt.Println(header)
	fmt.Println("=======================================================================")
	for _, val := range vxis {
		row := []string{val.OperatorAddress, val.AccountAddress, val.Moniker, fmt.Sprintf("%f", val.VotingPower), val.SelfDelegation, val.TotalDelegation}

		err := w.Write(row)
		if err != nil {
			log.Fatal("error writing record to file", err)
		}
		// If want to read from terminal.

		fmt.Println(row)
		// Continuous flush each record to disk
		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatal(err) // write file.csv: bad file descriptor
		}

	}
}

// NOT USED.
// TODO : Reduce code duplication using empty interfaces instead and generic write interface.
func WriteDelegatorsInfo(dxis []staking_view.DelegatorExportInfo) {
	file, err := os.Create("delegator_infos.csv")
	if err != nil {
		log.Fatal("failed to open file", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	header := []string{"DelegatorAddress", "TotalPower", "ValidatorAddress", "DelToThisVal"}
	err = w.Write(header)
	if err != nil {
		log.Fatal("error writing record to file", err)
	}
	fmt.Println(header)
	fmt.Println("=================================")
	for _, dxi := range dxis {
		row := []string{dxi.Address, dxi.TotalPower, dxi.ValidatorAddress, dxi.DelegatedToThisValidator}

		err := w.Write(row)
		if err != nil {
			log.Fatal("error writing record to file", err)
		}
		// If want to read from terminal.
		fmt.Println(row)
		// Continuous flush each record to disk
		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatal(err) // write file.csv: bad file descriptor
		}
	}
}

func OutDelegatorMap() {
	file, err := os.Create("delegator_validator_infos.csv")
	if err != nil {
		log.Fatal("failed to open file", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	header := []string{"DelegatorAddress", "DelTotalPower", "ValidatorAddress", "Share", "DelToThisVal"}
	err = w.Write(header)
	if err != nil {
		log.Fatal("error writing record to file", err)
	}
	fmt.Println(header)
	fmt.Println("==================================================================================")

	for _, dvis := range staking_view.DelegatorInfoMap {
		// fmt.Printf("Delegator : %v\n", k)
		for _, dvi := range dvis {
			row := []string{dvi.DelegatorAddress, *dvi.TotalPower, dvi.ValidatorAddress, dvi.Share, dvi.Amount}
			err := w.Write(row)
			if err != nil {
				log.Fatal("error writing record to file", err)
			}
			// If want to read from terminal.
			fmt.Println(dvi.DelegatorAddress, *dvi.TotalPower, dvi.ValidatorAddress, dvi.Share, dvi.Amount)
			// Continuous flush each record to disk
			w.Flush()
			if err := w.Error(); err != nil {
				log.Fatal(err) // write file.csv: bad file descriptor
			}

		}
	}
}
