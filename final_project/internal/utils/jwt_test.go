package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		secret  string
		expiry  time.Duration
		wantErr error
		check   func(*testing.T, string)
	}{
		{
			name:    "generates token successfully",
			userID:  "507f1f77bcf86cd799439011",
			secret:  "test-secret",
			expiry:  24 * time.Hour,
			wantErr: nil,
			check: func(t *testing.T, token string) {
				assert.NotEmpty(t, token)
			},
		},
		{
			name:    "generates token with short expiry",
			userID:  "000000000000000000000001",
			secret:  "another-secret",
			expiry:  1 * time.Second,
			wantErr: nil,
		},
		{
			name:    "generates token with different user id",
			userID:  "507f191e810c19729de860ea",
			secret:  "test-secret",
			expiry:  time.Hour,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.secret, tt.expiry)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				if tt.check != nil {
					tt.check(t, token)
				}
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	const secret = "test-secret"
	const userID = "507f1f77bcf86cd799439011"

	validToken, err := GenerateToken(userID, secret, time.Hour)
	require.NoError(t, err)

	expiredToken, err := GenerateToken(userID, secret, -time.Second)
	require.NoError(t, err)

	tests := []struct {
		name      string
		token     string
		secret    string
		wantErr   error
		checkFunc func(*testing.T, *JWTClaims)
	}{
		{
			name:    "parses valid token successfully",
			token:   validToken,
			secret:  secret,
			wantErr: nil,
			checkFunc: func(t *testing.T, claims *JWTClaims) {
				assert.Equal(t, userID, claims.UserID)
				assert.True(t, claims.ExpiresAt.After(time.Now()))
			},
		},
		{
			name:    "returns error for expired token",
			token:   expiredToken,
			secret:  secret,
			wantErr: jwt.ErrTokenExpired,
		},
		{
			name:    "returns error for wrong secret",
			token:   validToken,
			secret:  "wrong-secret",
			wantErr: jwt.ErrTokenSignatureInvalid,
		},
		{
			name:    "returns error for malformed token",
			token:   "not.a.valid.token",
			secret:  secret,
			wantErr: jwt.ErrTokenMalformed,
		},
		{
			name:    "returns error for empty token",
			token:   "",
			secret:  secret,
			wantErr: jwt.ErrTokenMalformed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.token, tt.secret)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, claims)
				if tt.checkFunc != nil {
					tt.checkFunc(t, claims)
				}
			}
		})
	}
}

func TestGenerateAndParseToken(t *testing.T) {
	t.Run("generated token can be parsed back", func(t *testing.T) {
		const userID = "507f1f77bcf86cd799439011"
		const secret = "integration-secret"

		token, err := GenerateToken(userID, secret, 2*time.Hour)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		claims, err := ParseToken(token, secret)
		require.NoError(t, err)
		require.NotNil(t, claims)

		assert.Equal(t, userID, claims.UserID)
		assert.True(t, claims.ExpiresAt.After(time.Now()))
	})
}
