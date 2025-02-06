// help.go
package tbfunctions

import (
	"fmt"
)

// NoArg function shows the default help text
func NoArg() {
	printText := `
+--------- Welcome To Tulobyte CLI ------------------+
| Version: 1.0.0                                     |
| Usage: tbwallet <tags/commands>                    |
|                                                    |
|  -h, --help                  Display help options  |
|  -v, --version               Display app version   |
|  -a, --about                 Summary of our vision |
|                              and work              |
|  -r, --refresh               Initiliaze important  |
|                              directory and files   |
+----------------------------------------------------+
`
	fmt.Println(printText)
}

// PrintHelp function shows all available commands
func PrintHelp() {
	helpText := `

Usage: tbwallet <flags> or <SUBCOMMANDS> <sub-flags>

SUBCOMMANDS:
    create-wallet, create                Create a Tulobyte SegWit Bech32 TB wallet.
    recover                              Recover your wallet using a private key or recovery phrase.
    address                              Display your wallet address.
    pubkey                               Display your wallet's public key.
    balance                              Check your wallet balance.
    config                               Manage Tulobyte command-line tool configuration settings.
    txn                                  Calculate transaction size, fees, and perform actual transfers.

FLAGS:
    -p                                   To recover wallet from private key(hex)
    -m                                   To recover wallet from Recovery Phrase
    network                                 Manage RPC configurations
    -wp                                  System wallet configurations
    -h, --help                           Display help options.
    -v, --version                        Display the application version.
    -r, --refresh                        Initiliaze important directory and files 

SUB-FLAGS:
    -h, --help                           Display help options for subcommands.

`
	fmt.Println(helpText)
}

// Print configuration help
func PrintConfigHelp() {
	helpText := `
Usage: tbwallet config <tags> <values>

Tags:
    network                       Manage Network configurations.
    -wp                           Manage the system wallet file path.

Values:
    network mainnet               Switch the network to mainnet.
    network testnet               Switch the network to testnet.
    network -d                    Display the current network configuration.
    -wp <filepath>                Update the wallet file path.
                                  Example usage:
                                  -  tulobyte config -wp /path/to/wallet.tb
    -wp -d                        Display the current wallet's file path from configuration.
    -batch -d                     Display the current batch choice.
                                  -  Normal: Suitable for fast, light, and cost-effective transactions.
                                  -  Hunter: Typically slower, heavier, and more expensive transactions.
                                             A block from every batch is chosen which has a high size
                                             to disburse 10% of transaction fees among senders.
                                             When you transact from this batch, you have a chance
                                             to win.
    -batch normal                 Update batch to Normal
    -batch hunter                 Update batch to Hunter
`
	fmt.Println(helpText)
}

func PrintTxnHelp() {
	helpText := `
Usage: tbwallet txn <RECIPIENT_ADDRESS> <AMOUNT> <DATA>

Arguments:
    RECIPIENT_ADDRESS           The address to which you want to send TBS.

    AMOUNT                      The amount to transfer, specified in Hanabytes.
                                Note: 1 TBT = 10,000,000 Hanabytes.

    DATA                        Additional data to include in the transaction for
                                verification purposes 
                                                    OR
                                This can be any text, such as a base64-encoded image or a quote.

    NOTE: The maximum size limit of a transaction is 1MB (1024KB).

`
	fmt.Println(helpText)
}

func PrintTxnHelpShort(argsReq string) {
	fmt.Println(`
 +---------------------------------------+
 | Error: Atleast `, argsReq, ` arguments required |
 +---------------------------------------+

 Usage: tulobyte txn <PURPOSE> <RECIPIENT_ADDRESS> <AMOUNT> <DATA>


 or type  "tbwallet txn -h" for help and more commands
    `)
}

// PrintVersion function shows the version of the application
func PrintVersion() {
	printText := `
Name: Tulobyte CLI App
Root Command: tbwallet
Version: 1.0.0

`
	fmt.Println(printText)
}

// PrintRecovery function shows the subcommands
func PrintRecoveryHelp() {
	helpText := `
Usage: tbwallet recover <flags> 

flags:
    -h, --help                       Display help options
    -m, --mnemonic                   To recover wallet using mnemonic phrase or recovery phrase
    -p, --privatekey                 To recover wallet using private key
`
	fmt.Println(helpText)
}
