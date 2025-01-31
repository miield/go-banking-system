package utils

import (
	"math/rand"
	"time"
)

// GenerateTransactionID creates a unique 16-character alphanumeric ID
func GenerateTransactionID() string {
	charRange := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charSlice := make([]byte, 16)
	seed := rand.NewSource(time.Now().UTC().UnixNano())
	source := rand.New(seed)
	
	for i := range charSlice { 
		charSlice[i] = charRange[source.Intn(len(charRange))]
	}

	// create an instance of the randomID
	randomId := string(charSlice)

	return randomId
}

// GenerateAccountNumber creates a random 10-digit account number
func GenerateAccountNumber() int64 {
	rand.NewSource(time.Now().UTC().UnixNano())
	return 1000000000 + rand.Int63n(9000000000)
}
