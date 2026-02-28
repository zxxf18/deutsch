package jwt

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
)

const (
	JWTKey        = "jwt"
	JWTExpore     = "jwt_expire"
	JWTTimeOrigin = "jwt_origin"
)

var (
	// ErrMissingSecretKey indicates Secret key is required
	ErrMissingSecretKey = errors.New("filed to get secret key")

	// ErrFailedTokenCreation indicates JWT Token failed to create, reason unknown
	ErrFailedTokenCreation = errors.New("failed to create JWT Token")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrExpiredToken = errors.New("failed to get valid token")

	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = errors.New("failed to get auth header")

	// ErrInvalidToken can be thrown if get token from context error
	ErrInvalidToken = errors.New("failed to get token")

	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	ErrInvalidSigningAlgorithm = errors.New("failed to get get signing algorithm")

	// ErrNoPrivKeyFile indicates that the given private key is unreadable
	ErrNoPrivKeyFile = errors.New("failed to get private key file")

	// ErrNoPubKeyFile indicates that the given public key is unreadable
	ErrNoPubKeyFile = errors.New("failed to get public key file")

	// ErrInvalidPrivKey indicates that the given private key is invalid
	ErrInvalidPrivKey = errors.New("failed to get private key")

	// ErrInvalidPubKey indicates the the given public key is invalid
	ErrInvalidPubKey = errors.New("failed to get public key")
)

type JWTHelper struct {
	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration
	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is TokenTime + MaxRefresh.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// signing algorithm - possible values are
	// HS256, HS384, HS512 (对称的，使用key作为密钥)
	// RS256, RS384 or RS512 （非对称的，需要使用rsa证书）
	// Optional, default is HS256.
	SigningAlgorithm string
	// Secret key used for signing. Required.
	Key string
	// Private key file for asymmetric algorithms
	PrivKeyFile string
	// Public key file for asymmetric algorithms
	PubKeyFile string

	// Private key
	privKey *rsa.PrivateKey
	// Public key
	pubKey *rsa.PublicKey

	// 采用 go-zero 框架的 parser ，逻辑为从 header 的 Authorization 中获取 token 并解析
	parser *token.TokenParser

	ctx context.Context
}

func NewJWTHelper(cfg JWTConfig) (*JWTHelper, error) {
	helper := &JWTHelper{
		Timeout:          cfg.Timeout,
		MaxRefresh:       cfg.MaxRefresh,
		SigningAlgorithm: cfg.SigningAlgorithm,
		Key:              cfg.Key,
		PrivKeyFile:      cfg.PrivKeyFile,
		PubKeyFile:       cfg.PubKeyFile,
		parser:           token.NewTokenParser(),
		ctx:              logx.ContextWithFields(context.Background(), logx.Field("pkg", "jwt")),
	}
	if helper.usingPublicKeyAlgo() {
		err := helper.readKeys()
		if err != nil {
			return nil, err
		}
	} else if cfg.Key == "" {
		return nil, ErrMissingSecretKey
	}
	return helper, nil
}

// ParseToken parse jwt token from context
func (j *JWTHelper) parseToken(r *http.Request) (*jwt.Token, error) {
	return j.parser.ParseToken(r, j.Key, "")
}

func (j *JWTHelper) keyFunc(t *jwt.Token) (any, error) {
	if jwt.GetSigningMethod(j.SigningAlgorithm) != t.Method {
		return nil, ErrInvalidSigningAlgorithm
	}
	if j.usingPublicKeyAlgo() {
		return j.pubKey, nil
	}
	// save token string if valid
	return j.Key, nil
}

func (j *JWTHelper) readKeys() error {
	err := j.privateKey()
	if err != nil {
		return err
	}
	err = j.publicKey()
	if err != nil {
		return err
	}
	return nil
}

func (j *JWTHelper) privateKey() error {
	keyData, err := os.ReadFile(j.PrivKeyFile)
	if err != nil {
		return ErrNoPrivKeyFile
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPrivKey
	}
	j.privKey = key
	return nil
}

func (j *JWTHelper) publicKey() error {
	keyData, err := os.ReadFile(j.PubKeyFile)
	if err != nil {
		return ErrNoPubKeyFile
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPubKey
	}
	j.pubKey = key
	return nil
}

func (j *JWTHelper) usingPublicKeyAlgo() bool {
	switch j.SigningAlgorithm {
	case jwt.SigningMethodRS256.Alg(), jwt.SigningMethodRS384.Alg(), jwt.SigningMethodRS512.Alg():
		return true
	}
	return false
}

func (j *JWTHelper) Generate(claims map[string]any) (token string, expire time.Time, err error) {
	now := time.Now()
	expire = now.Add(j.Timeout)
	claims[JWTExpore] = expire.Unix()
	claims[JWTTimeOrigin] = now.Unix()

	builder := jwt.New(jwt.GetSigningMethod(j.SigningAlgorithm))
	data := builder.Claims.(jwt.MapClaims)
	for k, v := range claims {
		data[k] = v
	}

	if j.usingPublicKeyAlgo() {
		token, err = builder.SignedString(j.privKey)
	} else {
		token, err = builder.SignedString([]byte(j.Key))
	}
	return
}

func (j *JWTHelper) Refresh(r *http.Request) (string, time.Time, error) {
	claims, err := j.CheckMaxRefreshAndParse(r)
	if err != nil {
		return "", time.Now(), err
	}
	return j.Generate(claims)
}

func (j *JWTHelper) CheckExpireAndParse(r *http.Request) (map[string]any, error) {
	return j.checkAndParse(r, j.Timeout)
}

func (j *JWTHelper) CheckMaxRefreshAndParse(r *http.Request) (map[string]any, error) {
	return j.checkAndParse(r, j.MaxRefresh)
}

// parseInt64 兼容 float64 和 json.Number，避免 interface conversion panic
func parseInt64(v any) (int64, error) {
	if v == nil {
		return 0, ErrInvalidToken
	}
	switch val := v.(type) {
	case float64:
		return int64(val), nil
	case json.Number:
		n, err := val.Int64()
		if err != nil {
			return 0, ErrInvalidToken
		}
		return n, nil
	default:
		return 0, ErrInvalidToken
	}
}

// GetExpireFromToken 从 token 字符串解析出过期时间戳（秒），用于黑名单 TTL
func GetExpireFromToken(tokenString string) (int64, error) {
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}
	return parseInt64(claims[JWTExpore])
}

func (j *JWTHelper) checkAndParse(r *http.Request, offset time.Duration) (map[string]any, error) {
	jtoken, err := j.parser.ParseToken(r, j.Key, "")
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if !ok || validationErr.Errors != jwt.ValidationErrorExpired {
			return nil, err
		}
	}
	if jtoken == nil {
		return nil, ErrFailedTokenCreation
	}

	claims := jtoken.Claims.(jwt.MapClaims)
	exp, err := parseInt64(claims[JWTExpore])
	if err != nil {
		return nil, ErrInvalidToken
	}
	if exp < time.Now().Add(-offset).Unix() {
		return nil, ErrExpiredToken
	}
	return claims, nil
}
