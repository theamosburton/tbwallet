// basicfunctions.go
package tbfunctions

import (
	"crypto/ecdsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

// DeriveSeedFromMnemonic derives the seed from the mnemonic phrase
func DeriveSeedFromMnemonic(mnemonic string) []byte {
	// Salt is simply "mnemonic" (no passphrase)
	salt := "mnemonic"

	// PBKDF2: HMAC-SHA512, 2048 iterations, 64 bytes key
	seed := pbkdf2.Key([]byte(mnemonic), []byte(salt), 2048, 64, sha512.New)

	return seed
}

func IsValidHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func SavePrivateKey(filename string, privateKey string) error {

	// Create the directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the private key to the file
	err = os.WriteFile(filename, []byte(privateKey), 0600) // 0600 for secure permissions
	if err != nil {
		return fmt.Errorf("failed to write private key to file: %w", err)
	}

	fmt.Println("Private key(hex) saved to:", filename)
	return nil
}

func InitDirs(isConfigInit bool) bool {
	homeDir, err := os.UserHomeDir()

	// Create tbwallet directory in ~/.config/tbwallet
	if err != nil {
		fmt.Println("Can't get user home directory")
		log.Fatalf("Reason: %v", err)
		return false
	}

	// Create config dir
	dirName := filepath.Join(homeDir, ".config", "tbwallet")
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(dirName, 0755)
		if err != nil {
			fmt.Println("Can't create ~/.config/tbwallet")
			log.Fatalf("Reason: %v", err)
			return false
		}
	}

	// Creating config Dir
	configDirName := filepath.Join(homeDir, ".config", "tbwallet")
	configFile := configDirName + "/config.json"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		CreateDefaultConfig(configFile, configDirName)
	}

	// Creating other dirs
	rootDirName := filepath.Join(homeDir, "tbwallet")
	if _, err := os.Stat(rootDirName); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(rootDirName, 0755)
		if err != nil {

			fmt.Println("Can't create ~/tbwallet")
			log.Fatalf("Reason: %v", err)
			return false
		}
	}
	//Mainnet directories
	mainnetDirName := filepath.Join(homeDir, "tbwallet", "mainnet")
	if _, err := os.Stat(mainnetDirName); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(mainnetDirName, 0755)
		if err != nil {
			fmt.Println("Can't create ~/tbwallet/mainnet")
			log.Fatalf("Reason: %v", err)
			return false
		}
	}

	//Testnet directories
	testnetDirName := filepath.Join(homeDir, "tbwallet", "testnet")
	if _, err := os.Stat(testnetDirName); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(testnetDirName, 0755)
		if err != nil {
			fmt.Println("Can't create ~/tbwallet/testnet")
			log.Fatalf("Reason: %v", err)
			return false
		}
	}
	// Initialiaze Configurations
	if isConfigInit {
		InitConfig(dirName)
	}
	return true

}

func ReadPrivateKey(filename string) ([]byte, error) {
	// Read the file content as a byte slice
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Convert hex string to bytes (assuming file contains valid hex characters)
	decodedData, err := hex.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	return decodedData, nil
}

// SerializePublicKeyCompressed serializes the public key in compressed format
func CsSerializePublicKeyCompressed(pub *ecdsa.PublicKey) []byte {
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

func ShowWalletInfo(infoType string) (string, bool) {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return "", false
	}

	walletFilePath := config.WalletPath
	_, statErr := os.Stat(walletFilePath) // Use a new variable 'statErr'

	// Check if the file doesn't exist
	if os.IsNotExist(statErr) {
		fmt.Println("Wallet doesn't exist or file removed")
		return "", false
	}

	// Read private key from wallet file
	privateKeyBytes, err := ReadPrivateKey(walletFilePath)
	if err != nil {
		fmt.Println("Can't read wallet file")
		return "", false
	}

	// Ensure the private key bytes are valid (at least 32 bytes long)
	if len(privateKeyBytes) < 32 {
		fmt.Println("Error: Private key must be at least 32 bytes long")
		return "", false
	}

	// Step 1: Hash the private key bytes to generate a valid private key
	privateKeyD := new(big.Int).SetBytes(privateKeyBytes[:])

	// Step 2: Generate the elliptic curve (P256 curve)
	curve := btcec.S256() // Use secp256k1

	// Step 3: Ensure the private key is valid (less than curve's order and greater than 0)
	if privateKeyD.Cmp(curve.Params().N) >= 0 || privateKeyD.Sign() <= 0 {
		fmt.Println("Error: Invalid private key provided")
		return "", false
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

	// Step 8: Generate the Ethereum-like address using Keccak-256
	address, err := GenerateAddress(publicKey)
	if err != nil {
		fmt.Println("Error generating address:", err)
		return "", false
	}

	// Return the appropriate wallet information based on the infoType
	switch infoType {
	case "address":
		return address, true
	case "pubkey":
		return publicKeyHex, true
	case "privatekey":
		return hex.EncodeToString(privateKeyBytes), true
	default:
		return "", false
	}
}

// GenerateEthereumAddress generates an Ethereum-like address using Keccak-256 hash of the public key
func GenerateAddress(pubKey []byte) (string, error) {
	// Perform Keccak-256 hash on the public key
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubKey)
	addressHash := hash.Sum(nil)

	// Ethereum address is the last 20 bytes of the Keccak-256 hash
	address := addressHash[12:]

	// Format the address as a hex string
	tbtAddress := "0x" + hex.EncodeToString(address)
	return tbtAddress, nil
}
func SerializePublicKeyUncompressed(pubKey *ecdsa.PublicKey) []byte {
	// Uncompressed public key starts with 0x04 followed by X and Y coordinates
	pubKeyBytes := pubKey.X.Bytes()
	pubKeyBytes = append(pubKeyBytes, pubKey.Y.Bytes()...)
	return pubKeyBytes
}
func ShowConfig(display string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	walletPath := config.WalletPath
	rpcNetwork := config.Network
	batchChoice := config.TxnBatch
	if display == "wp" {
		fmt.Print("Wallet Path: ", walletPath)
	} else if display == "rpc" {
		fmt.Println("RPC Network: ", rpcNetwork)
	} else if display == "batch" {
		if batchChoice == "0" {
			fmt.Println("Batch Choice: ", "Nromal")
		} else {
			fmt.Println("Batch Choice: ", "Hunter")
		}

	}
}

// Function to check if the file exists and if it's an image
func CheckValidImage(filePath string) (bool, string) {
	// Check if the file exists
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist
			returnError := `
			Error: Entered file path doesn't exists  
			   File Path: ` + filePath
			return false, returnError
		}
		// Some other error occurred (e.g., permission error)
		returnError := `
+---------------------------------------------------------------+
| Error: Some other error occurred (e.g., permission denied)    |
+---------------------------------------------------------------+
					`
		return false, returnError
	}

	// Try opening the file to check if it is an image
	file, err := os.Open(filePath)
	if err != nil {
		returnError := `
+--------------------------------------------+
| Error: We can't open and read your file    |
+--------------------------------------------+
			`
		return false, returnError

	}
	defer file.Close()
	// Decode the file to check if it's a valid image
	_, _, err = image.Decode(file)
	if err != nil {
		// Corrected the error formatting to include filePath
		returnError := `
  Error: Entered file is not a valid image   
     File Path:` + filePath
		return false, returnError
	}

	return true, ""
}

// Encode image to base64
func EncodeImageToBase64(filePath string) (string, bool) {
	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open file '", filePath, "'")
		log.Fatalf("   Reason: %v", err)
		return "", false
	}
	defer file.Close()

	// Determine the image format
	_, format, err := image.Decode(file) // Decoding is still needed to detect format
	if err != nil {
		fmt.Println("Failed to decode image")
		log.Fatalf("   Reason: %v", err)
		return "", false
	}

	// Reopen the file to get raw bytes (since image.Decode reads the stream)
	file.Seek(0, 0)
	imgBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Failed to read file'", filePath, "'")
		log.Fatalf("   Reason: %v", err)
		return "", false
	}

	// Encode the image to Base64
	base64Data := base64.StdEncoding.EncodeToString(imgBytes)

	// Build the data URL
	dataURL := fmt.Sprintf("data:image/%s;base64,%s", strings.ToLower(format), base64Data)
	return dataURL, true
}

func GetPrivateKey() (string, bool) {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return "", false
	}
	walletPath := config.WalletPath
	data, err := os.ReadFile(walletPath)
	if err != nil {
		fmt.Println("Failed to read file: %w", err)
		return "", false
	}
	privateKeyHex := string(data)
	return privateKeyHex, true

}
