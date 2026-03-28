#jwt #open-policy-agent

### `rego/authentication.rego`
```rego
package service.rego

import rego.v1

default auth := false

auth if {
    [valid, _, _] := verify_jwt
    valid = true
}

verify_jwt := io.jwt.decode_verify(input.Token, {
    "cert": input.Key,
    "iss": input.ISS
})
```

### `main.go`
```go
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/vishal2098govind/service/foundations/keystore"
)

func main() {

	ks := keystore.New()
	ks.LoadKeys(os.DirFS("zarf/keys/"))

	// currently active kid is usually fetched from a central storage
	const kid = "8234c8c5-0508-4301-bfdf-da1d515399a1"

	token, err := GenToken(ks, kid)
	if err != nil {
		log.Fatalf("failed to gen key pair: %v", err)
	}
	err = ValidateTokenOPA(ks, token)
	if err != nil {
		log.Fatalf("failed to validate token: %v", err)
	}
	err = ValidateTokenParse(ks, token)
	if err != nil {
		log.Fatalf("failed to validate token parse: %v", err)
	}
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Roles []string
}

func GenToken(ks keystore.KeyStore, kid string) (string, error) {

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "vishal-user-id",
			Issuer:    "services-go service",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{
			"ADMIN",
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(jwt.SigningMethodRS256.Name), claims)
	token.Header["kid"] = kid

	privatePem, err := ks.PrivateKeyPEM(kid)
	if err != nil {
		return "", fmt.Errorf("private key not found: %w", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privatePem))
	if err != nil {
		return "", fmt.Errorf("parsing private pem bytes: %w", err)
	}

	tokenStr, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	fmt.Println("-------------TOKEN----------")
	fmt.Println(tokenStr)
	fmt.Println("----------------------------")

	return tokenStr, nil
}

//go:embed rego/authentication.rego
var opaAuthenticationPolicy string

func OpaEvaluation(input any) error {
	query := fmt.Sprintf("x = data.%s.%s", "service.rego", "auth")
	// fmt.Println(opaAuthenticationPolicy)
	ctx := context.Background()
	q, err := rego.New(
		rego.Query(query),
		rego.Module("policy.rego", opaAuthenticationPolicy),
	).PrepareForEval(ctx)
	if err != nil {
		return fmt.Errorf("prepare rego query: %w", err)
	}

	results, err := q.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return fmt.Errorf("rego eval: %w", err)
	}

	if len(results) == 0 {
		return errors.New("no results")
	}

	result, ok := results[0].Bindings["x"].(bool)
	if !ok {
		return fmt.Errorf("invalid rego result: %w", err)
	}

	if !result {
		fmt.Println(results[0].Bindings)
		return fmt.Errorf("invalid token")
	}

	return nil
}

func ValidateTokenOPA(ks keystore.KeyStore, tokenStr string) error {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))
	var claims CustomClaims
	token, _, err := parser.ParseUnverified(tokenStr, &claims)
	if err != nil {
		return fmt.Errorf("parse token string: %w", err)
	}
	kid, ok := token.Header["kid"]
	if !ok {
		return fmt.Errorf("invalid token (invalid kid)")
	}

	kidStr, ok := kid.(string)
	if !ok {
		return fmt.Errorf("invalid token (kid not string)")
	}

	publicKey, err := ks.PublicKeyPEM(kidStr)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	input := map[string]any{
		"Token": tokenStr,
		"Key":   publicKey,
		"ISS":   "services-go service",
	}

	err = OpaEvaluation(input)
	if err != nil {
		return fmt.Errorf("eval opa policy: %w", err)
	}

	fmt.Println("token is valid")

	return nil
}

func ValidateTokenParse(ks keystore.KeyStore, token string) error {

	var claims CustomClaims
	t, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		claims, ok := t.Claims.(*CustomClaims)
		if !ok {
			return nil, fmt.Errorf("invalid claims")
		}
		if !claims.VerifyIssuer("services-go service", true) {
			return nil, fmt.Errorf("invalid token")
		}

		if err := claims.Valid(); err != nil {
			return nil, fmt.Errorf("invalid token")
		}

		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid kid")
		}

		publicKeyPem, err := ks.PublicKeyPEM(kid)
		if err != nil {
			return nil, fmt.Errorf("reading public key pem file: %w", err)
		}

		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPem))
		if err != nil {
			return nil, fmt.Errorf("parse publick key bytes: %w", err)
		}

		return pubKey, nil
	}, jwt.WithValidMethods([]string{
		jwt.SigningMethodRS256.Name,
	}))
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		return err
	}
	if !t.Valid {
		return fmt.Errorf("invalid token")
	}

	fmt.Println(claims.Roles)
	fmt.Println(claims.Subject)

	return nil
}

func GenKey() error {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to rsa.GenerateKey: %w", err)
	}

	privatePem, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("failed to os.Create private.pem file: %w", err)
	}
	defer privatePem.Close()

	if err = pem.Encode(privatePem, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return fmt.Errorf("failed to write to private.pem file: %w", err)
	}

	publicKey := privateKey.PublicKey
	publicPem, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("failed to create os.Create public.pem file: %w", err)
	}
	defer publicPem.Close()

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return fmt.Errorf("marshal public key: %w", err)
	}
	if err = pem.Encode(publicPem, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}); err != nil {
		return fmt.Errorf("failed to write to public.pem file: %w", err)
	}

	return nil
}

```

### `keystore/keystore.go`
```go
package keystore

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

type key struct {
	privatePEM string
	publicPEM  string
}

type KeyStore struct {
	store map[string]key
}

func New() KeyStore {
	return KeyStore{
		store: make(map[string]key),
	}
}

func (ks *KeyStore) PrivateKeyPEM(kid string) (string, error) {
	key, ok := ks.store[kid]
	if !ok {
		return "", fmt.Errorf("kid not found")
	}
	return key.privatePEM, nil
}

func (ks *KeyStore) PublicKeyPEM(kid string) (string, error) {
	key, ok := ks.store[kid]
	if !ok {
		return "", fmt.Errorf("kid not found")
	}
	return key.publicPEM, nil
}

// this is usually with the auth service
// and done when the auth service boots up
// best practice is to do initial fetch
// from a KMS like AWS KMS, and cache in
// memory using [KeyStore.store]
func (ks *KeyStore) LoadKeys(fsfs fs.FS) {
	fs.WalkDir(fsfs, ".", func(fileName string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if path.Ext(fileName) != ".pem" {
			return nil
		}

		file, err := fsfs.Open(fileName)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("open pem file: %w", err)
		}

		pem, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("read pem: %w", err)
		}

		privatePem := string(pem)
		publicPem, err := toPublicPEM(privatePem)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("private to public pem: %w", err)
		}
		ks.store[strings.Split(fileName, ".pem")[0]] = key{
			privatePEM: privatePem,
			publicPEM:  publicPem,
		}

		return nil
	})
}

func toPublicPEM(privatePEM string) (string, error) {
	pb, _ := pem.Decode([]byte(privatePEM))
	if pb == nil {
		return "", fmt.Errorf("error decoding private pem. it should of PKCS1 or PKCS8")
	}

	var pvtKey any
	pvtKey, err := x509.ParsePKCS1PrivateKey(pb.Bytes)
	if err != nil {
		pvtKey, err = x509.ParsePKCS8PrivateKey(pb.Bytes)
		if err != nil {
			return "", fmt.Errorf("parse private key: %w", err)
		}
	}

	privateKey, ok := pvtKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("invalid private key")
	}

	publicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", fmt.Errorf("marshal public key: %w", err)
	}
	publicLem := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKey,
	}

	var b bytes.Buffer
	err = pem.Encode(&b, &publicLem)
	if err != nil {
		return "", fmt.Errorf("encoding public pem: %w", err)
	}

	return b.String(), nil
}

```