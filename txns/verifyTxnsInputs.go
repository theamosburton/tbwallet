package txns

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"tbwallet/tbfunctions"
)

func VerifyTxnInputs() (string, bool, map[string]string) {
	argsReq := 5

	if len(os.Args) == argsReq {
		var tx_data string
		rec_address := os.Args[2]
		localAddress, isFound := tbfunctions.ShowWalletInfo("address")
		if isFound {
			rec_address = strings.ToLower(rec_address)
			localAddress = strings.ToLower(localAddress)
			if localAddress == rec_address {
				returnError := `
+---------------------------------------------------+
| Error: Reciever and sender address can't be same  |
+---------------------------------------------------+
												`
				return returnError, false, nil
			}
		} else {
			returnError := `
+--------------------------------------------+
| Error: Failed to load local wallet address |
+--------------------------------------------+
								`
			return returnError, false, nil
		}
		amount_hb, amount_error := strconv.Atoi(os.Args[3])
		if amount_error != nil {
			returnError := `
+---------------------------------------+
| Error: Please enter amount in digits  |
|        e.g 100, 23, 1000 etc          |
+---------------------------------------+
					`
			return returnError, false, nil
		}
		tx_data = os.Args[4]
		verifiedInput := CheckInputErrors(rec_address, amount_hb)
		if verifiedInput { // Inputs have no problem

			inputs := map[string]string{
				"rec_address": rec_address,
				"tx_data":     tx_data,
				"amount_hb":   strconv.Itoa(amount_hb), // amount converted to string
			}
			return "", true, inputs
		}

	} else {
		tbfunctions.PrintTxnHelp()
		return "", false, nil
	}
	return "", true, nil
}

func CheckInputErrors(rec_address string, amount_hb int) bool {

	// Check recipient address length
	if len(rec_address) != 42 {
		fmt.Println(`
+--------------------------------------------------+
| Error: Recipient address should be 42 chars long |
+--------------------------------------------------+
			`)
		return false // Return false immediately
	}

	// Check if the amount is greater than zero
	if amount_hb <= 0 {
		fmt.Println(`
+------------------------------------------+
| Error: Amount cannot be 0 or negative    |
+------------------------------------------+
			`)
		return false // Return false immediately
	}
	return true
}
