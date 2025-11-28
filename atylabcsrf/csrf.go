package atylabcsrf

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type CsrfPkgInterface interface {
	GenerateNonceString() string
	GenerateCSRFCookieToken(secret string, timestamp int64, nonceStr string) string
	ValidateCSRFCookieToken(token string, secret string, timestamp int64) error
}

type CsrfPkgStruct struct{}

func NewCsrfPkgStruct() CsrfPkgInterface {
	return &CsrfPkgStruct{}
}

func (c *CsrfPkgStruct) GenerateNonceString() string {
	nonce := make([]byte, 32)
	rand.Read(nonce)
	return base64.StdEncoding.EncodeToString(nonce)
}

func (c *CsrfPkgStruct) GenerateCSRFCookieToken(secret string, timestamp int64, nonceStr string) string {

	data := fmt.Sprintf("%d:%s", timestamp, nonceStr)
	sig := c.hmacSha256(data, secret)

	return fmt.Sprintf("%s:%s", data, sig)
}

func (c *CsrfPkgStruct) ValidateCSRFCookieToken(token string, secret string, timestamp int64) error {
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		return fmt.Errorf("invalid token format")
	}

	timestampStr, nonce, sig := parts[0], parts[1], parts[2]
	data := fmt.Sprintf("%s:%s", timestampStr, nonce)
	expectedSig := c.hmacSha256(data, secret)

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return fmt.Errorf("invalid token signature")
	}

	// 有効期限を 10分（600秒）にする場合：
	expTimestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp")
	}
	if timestamp-expTimestamp > 600 {
		return fmt.Errorf("token expired")
	}

	return nil
}

func (c *CsrfPkgStruct) hmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
