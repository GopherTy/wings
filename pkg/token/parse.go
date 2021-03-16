package token

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/gopherty/wings/common/conf"
)

//  var
var (
	AccessKeyFunc  = keyAcFunc
	RefreshKeyFunc = keyRfFunc

	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invaild token")
)

// ExtractTokenMetadata .
func ExtractTokenMetadata(ts string, keyFunc jwt.Keyfunc) (id uint64, err error) {
	token, err := jwt.Parse(ts, keyFunc)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		err = status.Errorf(codes.InvalidArgument, "token invaild")
		return
	}
	return strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
}

// ValidToken .
func ValidToken(ctx context.Context, keyFunc jwt.Keyfunc) error {
	token, err := VerifyToken(ctx, keyFunc)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return errInvalidToken
	}
	return nil
}

// VerifyToken .
func VerifyToken(ctx context.Context, keyFunc jwt.Keyfunc) (token *jwt.Token, err error) {
	ts, err := ExtractToken(ctx)
	if err != nil {
		return
	}
	token, err = jwt.Parse(ts, keyFunc)
	if err != nil {
		err = errInvalidToken
		return
	}
	return
}

// ExtractToken extract token from context
func ExtractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errMissingMetadata
	}

	authorization := md.Get("Authorization")
	if len(authorization) < 1 {
		return "", errInvalidToken
	}

	return strings.TrimPrefix(authorization[0], "Bearer "), nil
}

// access token keyFunc
func keyAcFunc(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, status.Errorf(codes.Unauthenticated, "sign alg method not match %v", t.Header["alg"])
	}
	return []byte(conf.Instance().Secret.Access), nil
}

// refresh token keyFunc
func keyRfFunc(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, status.Errorf(codes.Unauthenticated, "sign alg method not match %v", t.Header["alg"])
	}
	return []byte(conf.Instance().Secret.Refresh), nil
}
