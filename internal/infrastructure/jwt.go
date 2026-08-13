package infrastructure

import (
	"github.com/TNJKL/bookmark-pkg/pkg/jwtutils"
	"github.com/TNJKL/bookmark-pkg/pkg/utils"
)

const (
	privateKeyPath = "./private.pem"
	publicKeyPath  = "./public.pem"
)

// CreateJWTProvider creates jwtutils.JWTGenerator and jwtutils.JWTValidator
func CreateJWTProvider() (jwtutils.JWTGenerator, jwtutils.JWTValidator) {
	jwtGenerator, err := jwtutils.NewJWTGenerator(privateKeyPath)

	utils.NoErr(err)
	jwtValidator, err := jwtutils.NewJWTValidator(publicKeyPath)

	utils.NoErr(err)

	return jwtGenerator, jwtValidator
}
