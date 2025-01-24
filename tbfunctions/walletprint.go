// walletprint.go
package tbfunctions

import (
	"fmt"
	"strings"

	eastasianwidth "github.com/moznion/go-unicode-east-asian-width"
)

// PrintWallet prints the details of the wallet including the mnemonic, private key, public key, and Bech32 address.
func PrintWallet(mnemonic string, privateKeyHex string, publicKeyHex string, address string, purpose string) {
	// Set the box width for the wallet info
	boxWidth := 94
	border := "+" + strings.Repeat("-", boxWidth+2) + "+"

	// Print the wallet info inside a box
	var formattedInfo strings.Builder
	fmt.Print("\n \n")
	fmt.Fprintln(&formattedInfo, border)
	fmt.Fprintf(&formattedInfo, "|  %-90s    |\n", " ")
	fmt.Fprintf(&formattedInfo, "|  %-90s    |\n", "                                 WALLET "+purpose+" SUCCESSFULLY")
	fmt.Fprintf(&formattedInfo, "|  %-90s    |\n", " ")
	if purpose != "RECOVERED" {
		printWrappedLine(&formattedInfo, "| ", "         DON'T COPY/PASTE RECOVERY PHRASE, WRITE OR SAVE IT IN AN OFFLINE PLACE", boxWidth)
		fmt.Fprintf(&formattedInfo, "|  %-90s    |\n", " ")
		fmt.Fprintf(&formattedInfo, "|  %-90s    |\n", "Recovery Phrase:")
		printWrappedLine(&formattedInfo, "| ", mnemonic, boxWidth)
		fmt.Fprintf(&formattedInfo, "|  %-90s    |\n", " ")

	} else {
		printWrappedLine(&formattedInfo, "| ", "NOTE: "+"RECOVERY PHRASE CAN'T BE RECOVERED BY ANY MEAN", boxWidth)
		fmt.Fprintf(&formattedInfo, "|  %-90s    |\n", " ")
	}

	printWrappedLine(&formattedInfo, "| ", "Private Key (hex): "+privateKeyHex, boxWidth)
	fmt.Fprintf(&formattedInfo, "|  %-90s    |\n", " ")
	printWrappedLine(&formattedInfo, "| ", "Public Key (hex): "+publicKeyHex, boxWidth)
	fmt.Fprintf(&formattedInfo, "|  %-90s    |\n", " ")
	printWrappedLine(&formattedInfo, "| ", "Wallet Address: "+address, boxWidth)
	fmt.Fprintf(&formattedInfo, "|  %-90s    |\n", " ")

	fmt.Fprintln(&formattedInfo, border)

	// Output the formatted wallet info
	fmt.Println(formattedInfo.String())
}

// printWrappedLine takes care of wrapping the text to fit within the box
func printWrappedLine(writer *strings.Builder, prefix string, text string, maxWidth int) {
	// Wrap the text to fit within maxWidth
	words := wrapText(text, maxWidth-2) // reserve 2 spaces for the borders
	for i, line := range words {
		if i == len(words)-1 {
			// For the last line, do not add an extra space before the "|"
			fmt.Fprintf(writer, "%-90s       |\n", prefix+" "+line)
		} else {
			fmt.Fprintf(writer, "%-90s  |\n", prefix+" "+line)
		}
	}
}

// wrapText takes a string and breaks it into lines of maxWidth
func wrapText(text string, maxWidth int) []string {
	var wrappedLines []string
	var currentLine []rune

	// Iterate over each rune in the text
	for _, r := range text {
		currentLine = append(currentLine, r)

		// Check if the width of the line exceeds maxWidth
		if lineWidth(currentLine) > maxWidth {
			// If exceeded, add the current line to wrappedLines and reset it
			wrappedLines = append(wrappedLines, string(currentLine[:len(currentLine)-1]))
			currentLine = []rune{r} // Start new line with the current rune
		}
	}

	// Add any remaining text in currentLine
	if len(currentLine) > 0 {
		wrappedLines = append(wrappedLines, string(currentLine))
	}

	return wrappedLines
}

// lineWidth calculates the width of a string using eastasianwidth package
func lineWidth(runes []rune) int {
	width := 0
	for _, r := range runes {
		// Check the East Asian width of each rune
		if eastasianwidth.IsFullwidth(r) {
			width += 2
		} else {
			width++
		}
	}
	return width
}
