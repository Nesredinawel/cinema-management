package utils

import (
	"context"
	"os"

	"google.golang.org/api/idtoken"
)

func VerifyGoogleToken(idToken string) (*idtoken.Payload, error) {
	ctx := context.Background()
	payload, err := idtoken.Validate(ctx, idToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, err
	}
	return payload, nil
}
