// mnemonic.go
package tbwallet

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"
	"strings"
)

// LoadWordlist loads a plain array of words from the JSON file
func LoadWordlist(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var wordlist []string
	err = json.Unmarshal(data, &wordlist)
	if err != nil {
		return nil, err
	}

	return wordlist, nil
}

// CreateMnemonic generates a random mnemonic from the wordlist using crypto/rand
func CreateMnemonic(wordlist []string, wordCount int) (string, error) {
	var mnemonic []string
	for i := 0; i < wordCount; i++ {
		// Generate a cryptographically secure random index
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(wordlist))))
		if err != nil {
			return "", err
		}

		mnemonic = append(mnemonic, wordlist[randomIndex.Int64()])
	}
	// Join the words together with spaces and return it as a string
	return strings.Join(mnemonic, " "), nil
}

// GenerateMnemonic generates a random mnemonic and returns it as a string
func GenerateMnemonic(wordCount int) (string, error) {
	// Load wordlist from JSON file
	wordlist, err := LoadWordlist("wordlist.json")
	if err != nil {
		return "", err
	}

	// Generate a random mnemonic
	mnemonic, err := CreateMnemonic(wordlist, wordCount) // 12-word mnemonic
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}
