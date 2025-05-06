package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ctxKey string

const (
	UserIDKey ctxKey = "id"
)

func AuthInterceptor(secret []byte) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata")
		}

		authHeaders := md.Get("Authorization")
    
		if len(authHeaders) == 0 {
			return nil, errors.New("missing auth header")
		}

		tokenStr := strings.TrimPrefix(authHeaders[0], "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			return nil, errors.New("invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("invalid claims")
		}

		userID, ok := claims["id"].(string)
		if !ok {
			return nil, errors.New("missing user ID in token")
		}

		ctx = context.WithValue(ctx, UserIDKey, userID)

		return handler(ctx, req)
	}
}
