package auth

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	secret := "super-secret-key"
	userID := uuid.New()

	validToken, _ := MakeJWT(userID, secret, time.Hour)
	expiredToken, _ := MakeJWT(userID, secret, -time.Hour)
	wrongSecretToken, _ := MakeJWT(userID, "wrong-secret", time.Hour)

	tests := map[string]struct {
		tokenString string
		tokenSecret string
		wantID      uuid.UUID
		wantErr     bool
	}{
		"Valid Token": {
			tokenString: validToken,
			tokenSecret: secret,
			wantID:      userID,
			wantErr:     false,
		},
		"Expired Token": {
			tokenString: expiredToken,
			tokenSecret: secret,
			wantID:      uuid.Nil,
			wantErr:     true,
		},
		"Wrong Secret": {
			tokenString: wrongSecretToken,
			tokenSecret: secret,
			wantID:      uuid.Nil,
			wantErr:     true,
		},
		"Malformed Token": {
			tokenString: "not.a.token",
			tokenSecret: secret,
			wantID:      uuid.Nil,
			wantErr:     true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotID, err := ValidateJWT(tc.tokenString, tc.tokenSecret)

			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidateJWT() error = %v, wantErr %v", err, tc.wantErr)
			}

			if diff := cmp.Diff(tc.wantID, gotID); diff != "" {
				t.Errorf("ValidateJWT() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

