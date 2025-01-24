package txns

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"tbwallet/tbfunctions"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

// SignTxns signs a transaction using the private key derived from the wallet.
func SignTxn(txSenderAddress, txAmount, txNonce, txReceiverAddress, tx_data string) (bool, map[string]string) {
	// Retrieve private key from the wallet file
	privateKeyHex, isKeyFound := tbfunctions.GetPrivateKey()
	if !isKeyFound {
		fmt.Println("Private key not found in the wallet file.")
		return false, nil
	}

	// Decode the private key from hex (64 chars = 32 bytes)
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		fmt.Println("Failed to decode private key:", err)
		return false, nil
	}

	// Ensure the private key is the correct length (32 bytes)
	if len(privateKeyBytes) != 32 {
		fmt.Println("Invalid private key length.")
		return false, nil
	}
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to convert private key bytes to ECDSA: %v", err)
	}

	// Generate the current Unix timestamp
	txTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
	transaction := map[string]interface{}{
		"n": txNonce,           // Nonce
		"s": txSenderAddress,   // Sender address
		"r": txReceiverAddress, // Receiver address
		"t": txTimestamp,       // Timestamp
		"a": txAmount,          // Amount
		"b": 0,                 // 0 for even 1 for odd
		"d": tx_data,
	}

	txn, err := json.Marshal(transaction)
	if err != nil {
		return false, nil
	}

	txHash := crypto.Keccak256Hash(txn)
	// Sign the transaction hash
	signature, err := crypto.Sign(txHash.Bytes(), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign the transaction: %v", err)
	}

	// Verify the signature and recover the sender's address
	senderAddress, err := recoverAddress(txHash.Bytes(), signature)
	if err != nil {
		log.Fatalf("Failed to recover address: %v", err)
	}
	// check batch type
	config, err := tbfunctions.LoadConfig()
	if err != nil {
		fmt.Println(`
+------------------------------------+
| Error: Problem with config file    |
+------------------------------------+
				`)
	}
	TxnBatch := config.TxnBatch
	result := map[string]string{
		"n":  txNonce,
		"s":  txSenderAddress,
		"r":  txReceiverAddress,
		"t":  txTimestamp,
		"a":  txAmount,
		"b":  TxnBatch,
		"d":  tx_data,
		"sg": hex.EncodeToString(signature),
		"f":  "0000000000",
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error marshalling map to JSON:", err)
		return false, nil
	}

	// Get the size of the resulting JSON string
	transactionSize := len(jsonData)
	fees := transactionSize * 10
	result["f"] = strconv.Itoa(fees)
	result["h"] = txHash.Hex()
	// Convert the transaction to map[string]string
	senderAddress = strings.ToLower(senderAddress)

	// Example expected address for verification
	expectedAddress, isFound := tbfunctions.ShowWalletInfo("address")
	if !isFound {
		fmt.Println("Cannot get local wallet address.")
		return false, nil
	}
	if senderAddress == expectedAddress {
		return true, result
	} else {
		fmt.Println("Signature verification failed.")
		return false, nil
	}
}

func recoverAddress(hash []byte, signature []byte) (string, error) {
	// Ensure the signature length is 65 bytes (R, S, V)
	if len(signature) != 65 {
		return "", fmt.Errorf("invalid signature length: %d", len(signature))
	}

	// Recover the public key
	pubKey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return "", err
	}

	// Derive the Ethereum address from the public key
	address := crypto.PubkeyToAddress(*pubKey)
	return address.Hex(), nil
}
