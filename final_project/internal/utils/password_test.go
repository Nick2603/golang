package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
		check    func(*testing.T, string)
	}{
		{
			name:     "hashes password successfully",
			password: "secret123",
			wantErr:  nil,
			check: func(t *testing.T, hash string) {
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, "secret123", hash)
			},
		},
		{
			name:     "hashes empty password",
			password: "",
			wantErr:  nil,
			check: func(t *testing.T, hash string) {
				assert.NotEmpty(t, hash)
			},
		},
		{
			name:     "hashes long password",
			password: "a_very_long_password_that_is_still_valid_for_bcrypt_123456",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				if tt.check != nil {
					tt.check(t, hash)
				}
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	const correctPassword = "secret123"

	hash, err := HashPassword(correctPassword)
	require.NoError(t, err)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "returns true for correct password",
			password: correctPassword,
			hash:     hash,
			want:     true,
		},
		{
			name:     "returns false for wrong password",
			password: "wrongpassword",
			hash:     hash,
			want:     false,
		},
		{
			name:     "returns false for empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "returns false for similar but different password",
			password: "secret124",
			hash:     hash,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPassword(tt.password, tt.hash)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestHashPasswordAndCheck(t *testing.T) {
	t.Run("hashed password matches original", func(t *testing.T) {
		const password = "test-password-123"

		hash, err := HashPassword(password)
		require.NoError(t, err)
		require.NotEmpty(t, hash)

		assert.NotEqual(t, password, hash)
		assert.True(t, CheckPassword(password, hash))
		assert.False(t, CheckPassword("wrong-password", hash))
	})

	t.Run("same password produces different hashes", func(t *testing.T) {
		const password = "test-password"

		hash1, err := HashPassword(password)
		require.NoError(t, err)

		hash2, err := HashPassword(password)
		require.NoError(t, err)

		assert.NotEqual(t, hash1, hash2)
		assert.True(t, CheckPassword(password, hash1))
		assert.True(t, CheckPassword(password, hash2))
	})
}
