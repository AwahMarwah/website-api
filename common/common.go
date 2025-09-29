package common

const (
	DefaultPassword            = "P4ssword"
	EmailAlreadyExists         = "email already exists"
	EmailOrPasswordIsIncorrect = "email or password is incorrect"
	SuccessfullyChecked        = "successfully checked"
	SuccessfullyCreated        = "successfully created"
	SuccessfullySignedIn       = "successfully signed in"
	SuccessfullySignedOut      = "successfully signed out"
	SuccessfullyUpdated        = "successfully updated"
	SuccessfullyDeleted        = "successfully deleted"
	UnexpectedNewline          = "unexpected newline"

	// Security constants
	InvalidVerificationToken    = "invalid verification token"
	VerificationTokenExpired    = "verification token expired, please request a new one"
	EmailAlreadyVerified        = "email already verified"
	TooManyVerificationAttempts = "too many verification attempts, please try again later"
	VerificationFailed          = "verification failed"
	VerificationSuccess         = "email successfully verified"

	// Password Reset
	PasswordResetTokenInvalid = "password reset token is invalid"
	PasswordResetTokenExpired = "password reset token has expired"
	PasswordResetTokenUsed    = "password reset token has already been used"
	PasswordResetSuccess      = "password has been reset successfully"
	PasswordResetRequestSent  = "if your email is registered, you will receive a password reset link"

	// Simple rate limiting
	MaxVerificationAttempts = 5
	RateLimitWindowHours    = 1
)
