// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"tbwallet/tbfunctions"
	"tbwallet/tbwallet"
	"tbwallet/txns"
)

func main() {
	// Initliaze some configurations
	dirsInitiliazed := tbfunctions.InitDirs(false)

	if !dirsInitiliazed {
		return
	}

	// Capturing arguments
	if len(os.Args) < 2 {
		tbfunctions.NoArg()
	} else {
		FP := os.Args[1]
		if FP == "--help" || FP == "-h" {
			tbfunctions.PrintHelp()
		} else if FP == "-v" || FP == "--version" {
			tbfunctions.PrintVersion()
		} else if FP == "create" || FP == "create-wallet" {
			tbwallet.CreateWallet()
		} else if FP == "recover" {
			if len(os.Args) < 3 {
				tbwallet.RecoverWallet(false, "")
			} else {
				SP := os.Args[2]
				if SP == "-m" || SP == "--mnemonic" {
					tbwallet.RecoverWallet(true, "phrase")
				} else if SP == "-p" || SP == "--privatekey" {
					tbwallet.RecoverWallet(true, "key")
				} else if SP == "-h" || SP == "--help" {
					tbfunctions.PrintRecoveryHelp()
				} else {
					tbwallet.RecoverWallet(true, "")
				}
			}
		} else if FP == "config" {
			if len(os.Args) < 3 {
				tbfunctions.PrintConfigHelp()
			} else {
				SP := os.Args[2]
				if SP == "network" {
					if len(os.Args) >= 4 {
						TP := os.Args[3]
						if TP == "mainnet" {
							tbfunctions.ChangeRPC("mainnet")
						} else if TP == "testnet" {
							tbfunctions.ChangeRPC("testnet")
						} else if TP == "-d" {
							tbfunctions.ShowConfig("rpc")
						} else {
							fmt.Print("Usage: tulobyte config -rpc mainnet/testnet")
						}
					} else {
						tbfunctions.PrintConfigHelp()
					}
				} else if SP == "-wp" {
					if len(os.Args) >= 4 {
						TP := os.Args[3]
						if TP == "-d" {
							tbfunctions.ShowConfig("wp")
						} else {
							tbfunctions.ChangeWalletPath(TP)
						}
					} else {
						tbfunctions.PrintConfigHelp()
					}
				} else if SP == "-batch" {
					if len(os.Args) >= 4 {
						TP := os.Args[3]
						if TP == "-d" {
							tbfunctions.ShowConfig("batch")
						} else {
							tbfunctions.ChangeBatch(TP)
						}
					} else {
						tbfunctions.PrintConfigHelp()
					}
				} else if SP == "-h" || SP == "--help" {
					tbfunctions.PrintConfigHelp()
				} else {
					tbfunctions.PrintConfigHelp()
				}
			}
		} else if FP == "address" {
			address, isFound := tbfunctions.ShowWalletInfo("address")
			if !isFound {
				return
			}
			fmt.Println("Wallet Address:", address)
		} else if FP == "pubkey" {
			pubkey, isFound := tbfunctions.ShowWalletInfo("pubkey")
			if !isFound {
				return
			}
			fmt.Println("Public Key(hex): ", pubkey)
		} else if FP == "txn" {
			if len(os.Args) <= 2 {
				tbfunctions.PrintTxnHelp()
			} else if len(os.Args) > 2 {
				SP := os.Args[2]
				if SP == "-h" || SP == "--help" {
					tbfunctions.PrintTxnHelp()
				} else {
					startTxnsProcess()
				}
			}
		} else if FP == "-r" || FP == "--refresh" {
			dirsInitiliazed := tbfunctions.InitDirs(false)
			if !dirsInitiliazed {
				return
			} else {
				fmt.Println(`
+--------------------------------------------+
| Success: Tulobyte Refreshed Successfully   |
+--------------------------------------------+
										`)
			}

		} else {
			tbfunctions.NoArg()
		}
	}
}

func startTxnsProcess() {
	isError, isInputsVerified, dataMap := txns.VerifyTxnInputs()
	if !isInputsVerified && isError != "" {
		fmt.Println(isError)
		return

	} else if !isInputsVerified && isError == "" {
		fmt.Println("Transaction inputs cannot be verified")
		return
	}
	// amount converted back to int
	amount, err := strconv.Atoi(dataMap["amount_hb"])
	if err != nil {
		fmt.Println("Error converting amount_hb to int:", err)
		return
	}
	isTxnVerified, returnError, txnMap := txns.VerifyTxn(dataMap["rec_address"], dataMap["tx_data"], amount)
	if !isTxnVerified && returnError != "" {
		fmt.Println(returnError)
		return
	} else if !isTxnVerified && returnError == "" {
		fmt.Println("Transaction verification failed")
		return
	}
	x_sAddress := txnMap["tx_sAddress"]
	tx_amount := txnMap["tx_amount"]
	tx_nonce := txnMap["tx_nonce"]
	tx_rAddress := txnMap["tx_raddress"]
	tx_data := txnMap["tx_data"]
	tx_folder := txnMap["txnFolder"]
	txJsonFile := tx_folder + "/txn.json"

	isTxSigned, newTxnMap := txns.SignTxn(x_sAddress, tx_amount, tx_nonce, tx_rAddress, tx_data)
	if !isTxSigned {
		return
	}
	if !isTxSigned {
		fmt.Println("Failed to sign the transaction")
		return
	}
	file, err := os.Create(txJsonFile)
	if err != nil {
		fmt.Println("Error creating transaction file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(newTxnMap); err != nil {
		fmt.Println("Error encoding transaction to JSON:", err)
		return
	}

	_, err = file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	transactionFees := newTxnMap["f"]
	file, err = os.Create(txJsonFile)
	if err != nil {
		fmt.Println("Error creating transaction file:", err)
		return
	}
	defer file.Close()

	encoder = json.NewEncoder(file)
	if err := encoder.Encode(newTxnMap); err != nil {
		fmt.Println("Error encoding transaction to JSON:", err)
		return
	}
	// Transaction size is larger then 1KB or equals to 1KB
	printOutLine := `
  +-----------------------------------+
  |  Transaction Signed Successfully  |                                                                       
  +-----------------------------------+

  Hash : ` + newTxnMap["h"] + `                                                                                                               
  Estimated Fees : ` + transactionFees + ` Hanas                                                                                                  
`
	fmt.Println(printOutLine)
	var isBroadCast string
	fmt.Print("  Broadcast Transaction (Y/N): ")
	_, err = fmt.Scanln(&isBroadCast)
	if err != nil {
		log.Fatal("Error reading input:", err)
	}

	if isBroadCast == "Y" || isBroadCast == "y" {
		printOutLine := `
  +----------------------------+
  |  Transaction Broadcasted   |                                                                       
  +----------------------------+`
		fmt.Println(printOutLine)
	} else {
		printOutLine := `
  +-------------------------+
  |  Transaction Declined   |                                                                       
  +-------------------------+`
		fmt.Println(printOutLine)
	}
}
