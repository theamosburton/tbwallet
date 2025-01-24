package txns

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tulobyte/tbfunctions"

	"github.com/chai2010/webp"
)

func VerifyTxn(rec_address string, tx_data string, amount_hb int) (bool, string, map[string]string) {
	var tx_nonce int
	var tx_amount = amount_hb
	var tx_raddress = rec_address
	aval_amount := 0
	addressVerified, txns, amt, networkType := VerifyAddress(rec_address)
	if !addressVerified {
		return false, "", nil
	}

	aval_amount = amt
	if amount_hb > aval_amount {
		returnError := `
+----------------------------------------+
| Error:  Insufficient TBYT Balance      |
| Reason: Transfer amount is larger then |
|         available balance              |
+----------------------------------------+
					`
		return false, returnError, nil
	}
	tx_sAddress, isFound := tbfunctions.ShowWalletInfo("address")
	if !isFound {
		return false, "", nil
	}
	tx_publicKey, isFound := tbfunctions.ShowWalletInfo("pubkey")
	if !isFound {
		return false, "", nil
	}
	tx_privateKey, isFound := tbfunctions.ShowWalletInfo("privatekey")
	if !isFound {
		return false, "", nil
	}
	isCreated, txnFolder := CreateTxnsDirs(networkType)
	if !isCreated {
		return false, "", nil
	}
	if txns == 0 {
		tx_nonce = 0
	} else {
		tx_nonce = txns - 1
	}
	inputs := map[string]string{
		"txnFolder":     txnFolder,
		"networkType":   networkType,
		"tx_sAddress":   tx_sAddress,
		"tx_raddress":   tx_raddress,
		"tx_publicKey":  tx_publicKey,
		"tx_privateKey": tx_privateKey,
		"tx_amount":     strconv.Itoa(tx_amount),
		"tx_nonce":      strconv.Itoa(tx_nonce),
		"tx_data":       tx_data,
	}
	return true, "", inputs
}

func CalculateFastTxnSize(tx_publicKey string, tx_sAddress string, tx_amount string, tx_nonce string, tx_raddress string) (int, bool) {
	return 0, true
}

func CalculateSlowTxnSize(tx_publicKey string, tx_sAddress string, tx_amount string, tx_nonce string, tx_raddress string, tx_image string) (int, bool) {
	return 0, true
}

func CreateTxnsDirs(networkType string) (bool, string) {
	dirsInitiliazed := tbfunctions.InitDirs(false)
	if !dirsInitiliazed {
		fmt.Println(`
+------------------------------------------------+
| Success: Can't Initiliaze Config Directories   |
+------------------------------------------------+
		`)
		return false, ""
	}
	noOfFolder := 0
	homeDir, err := os.UserHomeDir()

	// Create tulobyte directory in ~/.config/tulobyte
	if err != nil {
		fmt.Println("Can't get user home directory")
		log.Fatalf("   Reason: %v", err)
		return false, ""
	}
	txnDir := filepath.Join(homeDir, "tulobyte", networkType, "txns")
	files, err := os.ReadDir(txnDir)
	if err != nil {
		fmt.Println("Failed to read '", txnDir, "'")
		log.Fatalf("   Reason: %v", err)
		return false, ""
	}

	for _, file := range files {
		if file.IsDir() {
			noOfFolder++
		}
	}
	folderCountString := strconv.Itoa(noOfFolder)
	newTxDir := filepath.Join(txnDir, folderCountString)
	if _, err := os.Stat(newTxDir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.MkdirAll(newTxDir, 0755) // Use MkdirAll to ensure parent directories are created
		if err != nil {
			fmt.Println("Can't create a new txn directory")
			log.Fatalf("   Reason: %v", err)
			return false, ""
		}
	}

	return true, newTxDir
}

func VerifyAddress(rec_address string) (bool, int, int, string) {
	if !VerifyAddressFormat(rec_address) {
		fmt.Println(`
+-----------------------------------+
| Error: Invalid Recipent Address   |
+-----------------------------------+
					`)
		return false, 0, 0, ""
	}
	// check network type
	config, err := tbfunctions.LoadConfig()
	if err != nil {
		fmt.Println(`
+------------------------------------+
| Error: Problem with config file    |
+------------------------------------+
			`)
		return false, 0, 0, ""
	}
	networkType := config.RPCEndPoint
	if networkType == "mainnet" {
		amt, txns, walletChecked := CheckMyWalletMainnet()
		if !walletChecked {
			fmt.Println(`
+-----------------------------------------------------------+
| Error: Can't query RPC node to get wallet Transactions    |
+-----------------------------------------------------------+
						`)
			return false, 0, 0, ""
		} else {
			return true, txns, amt, networkType
		}
	} else if networkType == "testnet" {
		amt, txns, walletChecked := CheckMyWalletTestnet()
		if !walletChecked {
			fmt.Println(`
+-----------------------------------------------------------+
| Error: Can't query RPC node to get wallet Transactions    |
+-----------------------------------------------------------+
						`)
			return false, 0, 0, ""
		} else {
			return true, txns, amt, networkType
		}
	} else {
		return false, 0, 0, ""
	}

}
func VerifyAmount(amount_hb int) bool { return true }

func VerifyAddressFormat(rec_address string) bool {
	returnValue := true
	// Check if the address has the correct length (42 characters including '0x' prefix)
	if len(rec_address) != 42 {
		returnValue = false
	}

	// Check if the address starts with '0x'
	if !strings.HasPrefix(rec_address, "0x") {
		returnValue = false
	}

	// Check if the address contains only valid hexadecimal characters
	for _, char := range rec_address[2:] {
		if !strings.Contains("0123456789abcdefABCDEF", string(char)) {
			returnValue = false
		}
	}

	return returnValue
}

func CheckMyWalletMainnet() (int, int, bool) {
	return 10000, 0, true
}
func CheckMyWalletTestnet() (int, int, bool) {
	return 10000, 0, true
}

// convertImageToBase64 converts an image file (including WebP) to a Base64 string.
func ConvertImageToBase64(filePath string) (string, error) {
	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()
	// Decode the image to ensure it's a valid image file
	_, _, err = image.Decode(file)
	if err != nil {
		// Attempt WebP decoding if general decoding fails
		if _, err := webp.Decode(file); err == nil {
			// WebP format detected
		} else {
			return "", fmt.Errorf("could not decode image: %v", err)
		}
	}

	// Reset the file pointer to the beginning for reading raw bytes
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("could not reset file pointer: %v", err)
	}

	// Read the image file's bytes
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("could not read file bytes: %v", err)
	}

	// Encode the image bytes to a Base64 string
	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64String, nil
}
