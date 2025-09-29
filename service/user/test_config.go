package user

import (
	"os"
	"testing"
)

// TestMain runs before all tests in this package
func TestMain(m *testing.M) {
	// Setup test environment
	setupTestEnvironment()

	// Run tests
	code := m.Run()

	// Cleanup
	cleanupTestEnvironment()

	// Exit with test result code
	os.Exit(code)
}

func setupTestEnvironment() {
	os.Setenv("APP_BASE_URL", "http://localhost:8080")
	os.Setenv("JWT_SECRET", "test_secret_key_for_testing")
	os.Setenv("SMTP_HOST", "smtp.test.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "testpassword")
}

func cleanupTestEnvironment() {
	os.Unsetenv("APP_BASE_URL")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USERNAME")
	os.Unsetenv("SMTP_PASSWORD")
}
