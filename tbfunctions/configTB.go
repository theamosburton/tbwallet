package tbfunctions

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	RPCEndPoint  string `json:"RPCEndPoint"`
	WalletPath   string `json:"Walletpath"`
	SysWalletUID string `json:"SysWalletUID"`
	TxnBatch     string `json:"TxnBatch"`
}

func CreateDefaultConfig(filename string) error {
	defaultConfig := Config{
		RPCEndPoint:  "mainnet",
		WalletPath:   "/",
		SysWalletUID: "0",
		TxnBatch:     "1",
	}

	data, err := json.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	// Ensure the directory exists
	err = os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the YAML data to the file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func ChangeRPC(network string) {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	//Updating config
	config.RPCEndPoint = network
	// Savng config
	err = SaveConfig(config)
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}
	fmt.Println("Network configured to :", network)
}

func ChangeBatch(batch string) {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	batchConfigured := ""
	//Updating config
	batch = strings.ToLower(batch)
	if batch == "normal" {
		config.TxnBatch = "1"
		batchConfigured = "Normal"
	} else if batch == "hunter" {
		config.TxnBatch = "0"
		batchConfigured = "Hunter"
	} else {
		config.TxnBatch = "1"
		batchConfigured = "Normal"
	}
	// Savng config
	err = SaveConfig(config)
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}
	fmt.Println("Network configured to :", batchConfigured)
}
func ChangeWalletPath(walletpath string) {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	//Updating config
	config.WalletPath = walletpath
	// Savng config
	err = SaveConfig(config)
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}
	fmt.Println("Syatem Wallet configured to :", walletpath)
}

func LoadConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to get the user's home directory: %v", err)
	}
	dirName := filepath.Join(homeDir, ".config", "tulobyte")

	filename := dirName + "/config.json"

	var config Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func SaveConfig(config Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to get the user's home directory: %v", err)
	}
	dirName := filepath.Join(homeDir, ".config", "tulobyte")

	filename := dirName + "/config.yaml"

	data, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func InitConfig(dirName string) {

	configFile := dirName + "/config.json"

	// Check if config file exists, create one if not
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		err := CreateDefaultConfig(configFile)
		if err != nil {
			fmt.Println("Error creating default config:", err)
			return
		}
	}
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	//Updating config
	if config.WalletPath == "/" {
		config.WalletPath = dirName + "/wallet.tb"
	}

	// Savng config
	err = SaveConfig(config)
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}
}
