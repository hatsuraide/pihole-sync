package totp

import (
	"time"

	"github.com/pquerna/otp/totp"
)

// GenerateCode returns the current TOTP passcode for the given base32 secret,
// matching the code a user's authenticator app would produce for Pi-hole 2FA.
func GenerateCode(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now())
}
