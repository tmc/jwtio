package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	tokenBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to read JWT token from stdin: %v", err)
	}
	token := string(tokenBytes)
	token = strings.TrimSpace(token)

	token = clean(token)

	// Split the token into its parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		fmt.Fprintln(os.Stderr, "issue: token does not have 3 parts")
	}

	colors := []string{colorLightRed, colorLightPurple, colorLightBlue}
	for i, part := range parts {
		fmt.Printf("%s", colors[i%len(colors)])
		printPart(part)
	}
	return nil
}

// clean up input such as stripping Bearer prefix
func clean(token string) string {
	token = strings.TrimPrefix(token, "Bearer ")
	return token
}

const (
	colorReset       = "\033[0m"
	colorRed         = "\033[31m"
	colorGreen       = "\033[32m"
	colorYellow      = "\033[33m"
	colorBlue        = "\033[34m"
	colorLightRed    = "\033[91m"
	colorLightPurple = "\033[94m"
	colorLightBlue   = "\033[96m"
)

func printPart(part string) {
	// Decode the part
	decoded, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(decoded, &m); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
	pretty, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
	fmt.Println(string(pretty))
}
