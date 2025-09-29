package utils

import (
	"fmt"
	"sync"
	"time"
	"website-api/common"
)

var (
	verificationAttempts = make(map[string][]time.Time)
	verificationMutex    = sync.RWMutex{}
)

// CheckVerificationRateLimit checks if verification attempts are within limit
func CheckVerificationRateLimit(identifier string) error {
	verificationMutex.Lock()
	defer verificationMutex.Unlock()

	now := time.Now()
	attempts := verificationAttempts[identifier]

	// Remove attempts older than 1 hour
	var validAttempts []time.Time
	for _, attempt := range attempts {
		if now.Sub(attempt) < time.Duration(common.RateLimitWindowHours)*time.Hour {
			validAttempts = append(validAttempts, attempt)
		}
	}

	// Check if exceeded max attempts
	if len(validAttempts) >= common.MaxVerificationAttempts {
		return fmt.Errorf(common.TooManyVerificationAttempts)
	}

	// Add current attempt
	validAttempts = append(validAttempts, now)
	verificationAttempts[identifier] = validAttempts

	return nil
}

// ResetVerificationRateLimit resets rate limit for identifier
func ResetVerificationRateLimit(identifier string) {
	verificationMutex.Lock()
	defer verificationMutex.Unlock()
	delete(verificationAttempts, identifier)
}