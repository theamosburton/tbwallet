package tbwallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"tbwallet/tbfunctions"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/term"
)

// DeriveSeedFromMnemonic derives a seed from the given mnemonic and passphrase
func DeriveSeedFromMnemonic(mnemonic string, passphrase string) []byte {
	// Convert mnemonic to seed using BIP-39
	seed := bip39.NewSeed(mnemonic, passphrase)
	return seed
}

// DeriveKeyPair derives a private and public key from a seed based on BIP-32 and BIP-44 path
func DeriveKeyPair(seed []byte) (*ecdsa.PrivateKey, []byte, error) {
	if len(seed) < 32 {
		return nil, nil, errors.New("seed must be at least 32 bytes long")
	}

	// Use BIP-32 to derive keys from the seed
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate master key: %v", err)
	}

	// BIP-44 path for Tulobyte: m/44'/202'/0'/0/0
	tbtpath := []uint32{44 + 0x80000000, 202 + 0x80000000, 0 + 0x80000000, 0, 0}
	childKey := masterKey
	for _, index := range tbtpath {
		childKey, err = childKey.NewChildKey(index)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to derive Ethereum keys: %v", err)
		}
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive Ethereum keys: %v", err)
	}

	// Generate private key from the derived key
	privateKeyD := new(big.Int).SetBytes(childKey.Key)

	curve := btcec.S256() // Use secp256k1
	privateKey := &ecdsa.PrivateKey{
		D: privateKeyD,
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
		},
	}
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(privateKey.D.Bytes())

	// Compress the public key
	uncompressedPublicKey := SerializePublicKeyUncompressed(&privateKey.PublicKey)

	// Return the private key and compressed public key
	return privateKey, uncompressedPublicKey, nil
}

// SerializePublicKeyUncompressed serializes the public key in uncompressed format (0x04 + X + Y)
func SerializePublicKeyUncompressed(pubKey *ecdsa.PublicKey) []byte {
	// Uncompressed public key starts with 0x04 followed by X and Y coordinates
	pubKeyBytes := pubKey.X.Bytes()
	pubKeyBytes = append(pubKeyBytes, pubKey.Y.Bytes()...)
	return pubKeyBytes
}

// SerializePublicKeyCompressed serializes the public key in compressed format for secp256k1
func SerializePublicKeyCompressed(pub *ecdsa.PublicKey) []byte {
	xBytes := pub.X.Bytes()

	var prefix byte
	if pub.Y.Bit(0) == 0 {
		// y is even
		prefix = 0x02
	} else {
		// y is odd
		prefix = 0x03
	}

	// Compress the x-coordinate to 32 bytes
	paddedX := make([]byte, 32)
	copy(paddedX[32-len(xBytes):], xBytes)

	// The compressed public key consists of the prefix + x-coordinate
	serializedKey := append([]byte{prefix}, paddedX...)
	return serializedKey
}

// CreateWallet creates a wallet with a mnemonic, derives keys, and generates an Ethereum address
func CreateWallet() {
	// Step 1: Input for mnemonic length
	var mnemonicLength int
	fmt.Print("Enter mnemonic length (default 12): ")
	_, err := fmt.Scanln(&mnemonicLength)
	if err != nil || mnemonicLength < 12 || mnemonicLength > 24 || mnemonicLength%3 != 0 {
		mnemonicLength = 12
	}

	// Step 2: Input for passphrase (hidden input)
	fmt.Print("Enter passphrase (Optional): ")
	passphrase, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading passphrase:", err)
		return
	}

	// Step 3: Generate mnemonic
	mnemonic, err := GenerateMnemonic(mnemonicLength)
	if err != nil {
		fmt.Println("Error generating mnemonic:", err)
		return
	}
	// Step 4: Derive seed
	seed := DeriveSeedFromMnemonic(mnemonic, string(passphrase))

	// Step 5: Derive private and public keys
	privateKey, uncompressedPublicKey, err := DeriveKeyPair(seed)
	if err != nil {
		fmt.Println("Error deriving key pair:", err)
		return
	}

	// Generate Ethereum address from the public key
	tbtAddress, err := tbfunctions.GenerateAddress(uncompressedPublicKey)
	if err != nil {
		log.Fatal("Error generating address:", err)
	}

	// Step 6: Load configuration
	config, err := tbfunctions.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Step 7: Print wallet information
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())
	publicKeyHex := hex.EncodeToString(uncompressedPublicKey)
	tbfunctions.PrintWallet(mnemonic, privateKeyHex, publicKeyHex, tbtAddress, "CREATED")

	// Save wallet to system
	walletFile := config.WalletPath
	tbfunctions.SavePrivateKey(walletFile, privateKeyHex)
}
