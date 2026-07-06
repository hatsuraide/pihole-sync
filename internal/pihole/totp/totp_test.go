package totp

import (
	"testing"

	pquernaotp "github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"
)

func TestGenerateCode(t *testing.T) {
	const secret = "JBSWY3DPEHPK3PXP"

	code, err := GenerateCode(secret)
	require.NoError(t, err)
	require.Len(t, code, 6)

	valid := pquernaotp.Validate(code, secret)
	require.True(t, valid)
}
