package tbwallet

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"tulobyte/tbfunctions"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/tyler-smith/go-bip32"
	"golang.org/x/term"
)

// RecoverWallet allows the user to select the recovery method (key or phrase)
func RecoverWallet(methodGiven bool, recoveryMethod string) {
	if methodGiven {
		switch recoveryMethod {
		case "key":
			RecoverWalletFromKey()
		case "phrase":
			RecoverWalletFromPhrase()
		default:
			tbfunctions.PrintRecoveryHelp()
		}
	} else {
		tbfunctions.PrintRecoveryHelp()
	}
}

// RecoverWalletFromKey recovers a wallet using the private key in hex format
func RecoverWalletFromKey() {
	var hexPrivateKey string
	fmt.Print("Enter your private key in hex format: ")
	_, err := fmt.Scanln(&hexPrivateKey)
	if err != nil {
		log.Fatal("Error reading input:", err)
	}

	hexPrivateKey = strings.TrimSpace(hexPrivateKey)

	// Validate hexadecimal format
	if !tbfunctions.IsValidHex(hexPrivateKey) {
		log.Fatal("Invalid input: Private key must be in hexadecimal format")
	}

	// Validate length of the private key
	if len(hexPrivateKey) < 64 {
		log.Fatal("Invalid input: Private key must be at least 64 hexadecimal characters long")
	}

	// Convert hex string to bytes
	privBytes, err := hex.DecodeString(hexPrivateKey)
	if err != nil {
		log.Fatal("Error decoding hex:", err)
	}

	// Ensure the private key is valid
	if len(privBytes) < 32 {
		log.Fatal("Private key must be at least 32 bytes long")
	}

	// Step 1: Hash the private key bytes to generate a valid private key
	privateKeyD := new(big.Int).SetBytes(privBytes[:])

	// Step 2: Generate the elliptic curve (P256 curve)
	curve := btcec.S256()

	// Step 3: Ensure the private key is valid (less than curve's order and greater than 0)
	if privateKeyD.Cmp(curve.Params().N) >= 0 || privateKeyD.Sign() <= 0 {
		log.Fatal("Invalid private key: Key out of range")
	}

	// Step 4: Create the private key object
	privateKey := &ecdsa.PrivateKey{
		D: privateKeyD,
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
		},
	}
	// Step 5: Derive the corresponding public key (using scalar multiplication)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(privateKey.D.Bytes())

	// Step 6: Serialize the public key (compressed)
	publicKey := SerializePublicKeyUncompressed(&privateKey.PublicKey)
	publicKeyHex := hex.EncodeToString(publicKey)

	// Step 7: Load configuration
	config, err := tbfunctions.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Step 8: Generate address based on the selected network (mainnet or testnet)
	address, err := tbfunctions.GenerateAddress(publicKey)
	if err != nil {
		log.Fatal("Error generating address:", err)
	}

	// Step 9: Print wallet information
	mnemonic := "RECOVERY PHRASE CAN'T BE RECOVERED BY ANY MEANS"
	tbfunctions.PrintWallet(mnemonic, hexPrivateKey, publicKeyHex, address, "RECOVERED")

	// Save wallet to system
	walletFile := config.WalletPath
	tbfunctions.SavePrivateKey(walletFile, hexPrivateKey)
}

// RecoverWalletFromPhrase recovers a wallet using the mnemonic phrase
func RecoverWalletFromPhrase() {
	fmt.Print("Enter your recovery phrase (Press Enter twice to finish): ")

	// Use bufio.NewReader to handle multi-line input
	reader := bufio.NewReader(os.Stdin)
	var mnemonicLines []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading input:", err)
		}

		// Trim spaces and newlines
		line = strings.TrimSpace(line)

		// If the line is empty, stop reading input
		if line == "" {
			break
		}

		// Append the line to the mnemonic slice
		mnemonicLines = append(mnemonicLines, line)
	}

	// Join all mnemonic lines into a single string
	mnemonic := strings.Join(mnemonicLines, " ")

	// Prompt for passphrase (optional)
	fmt.Print("Enter your passphrase (Optional): ")
	passphrase, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("Error reading passphrase:", err)
	}

	// Step 1: Derive the seed from mnemonic and passphrase
	seed := DeriveSeedFromMnemonic(mnemonic, string(passphrase))
	// Step 2: Derive private and public keys from the seed with derivation path
	privateKey, compressedPublicKey, err := DeriveKeyPair(seed)
	if err != nil {
		log.Fatal("Error deriving key pair:", err)
	}

	// Step 3: Load configuration
	config, err := tbfunctions.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Step 4: Generate the address using keccak256
	address, err := tbfunctions.GenerateAddress(compressedPublicKey)
	if err != nil {
		log.Fatal("Error generating address:", err)
	}

	// Step 5: Print wallet information
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())
	publicKeyHex := hex.EncodeToString(compressedPublicKey)
	tbfunctions.PrintWallet(mnemonic, privateKeyHex, publicKeyHex, address, "RECOVERED")

	// Save wallet to system
	walletFile := config.WalletPath
	tbfunctions.SavePrivateKey(walletFile, privateKeyHex)
}

// DeriveKeyPairWithPath derives a private and public key from the seed using BIP-44 path
func DeriveKeyPairWithPath(seed []byte) (*ecdsa.PrivateKey, []byte, error) {
	if len(seed) < 32 {
		return nil, nil, fmt.Errorf("seed must be at least 32 bytes long")
	}

	// Use BIP-32 to derive keys from the seed
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate master key: %v", err)
	}

	// BIP-44 path for Tulobyte: m/44'/202'/0'/0/0
	tbtath := []uint32{44 + 0x80000000, 202 + 0x80000000, 0 + 0x80000000, 0, 0}
	childKey := masterKey
	for _, index := range tbtath {
		childKey, err = childKey.NewChildKey(index)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to derive keys: %v", err)
		}
	}

	// Generate private key from the derived key
	privateKeyD := new(big.Int).SetBytes(childKey.Key)

	// Generate the public key using elliptic curve (P256)
	curve := btcec.S256()
	privateKey := &ecdsa.PrivateKey{
		D: privateKeyD,
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
		},
	}
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(privateKey.D.Bytes())

	// Serialize the public key (compressed)
	uncompressedPublicKey := SerializePublicKeyCompressed(&privateKey.PublicKey)

	return privateKey, uncompressedPublicKey, nil
}
