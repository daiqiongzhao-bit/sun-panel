// Package totp 提供最小可用的 TOTP(RFC 6238) 实现，用于两步验证(2FA)。
// 不依赖第三方库，避免引入新的模块依赖。
package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

// GenerateSecret 生成Base32编码的TOTP密钥
func GenerateSecret() (string, error) {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.TrimRight(base32.StdEncoding.EncodeToString(b), "="), nil
}

// OTPAuthURI 生成 otpauth:// URI，便于 authenticator 应用识别
func OTPAuthURI(account, issuer, secret string) string {
	account = strings.ReplaceAll(account, ":", "_")
	issuer = strings.ReplaceAll(issuer, ":", "_")
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, account, secret, issuer)
}

func pad(secret string) string {
	for len(secret)%8 != 0 {
		secret += "="
	}
	return secret
}

func hotp(secret string, counter uint64) (string, error) {
	key, err := base32.StdEncoding.DecodeString(pad(secret))
	if err != nil {
		return "", err
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)
	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	code := (binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7fffffff) % 1000000
	return fmt.Sprintf("%06d", code), nil
}

// GetTOTP 返回指定时间的6位动态码
func GetTOTP(secret string, t time.Time) (string, error) {
	return hotp(secret, uint64(t.Unix()/30))
}

// Validate 校验动态码，允许前后各一个时间步长(±30s)的误差
func Validate(secret, code string) bool {
	code = strings.TrimSpace(code)
	if len(code) != 6 {
		return false
	}
	now := time.Now()
	for _, t := range []time.Time{now.Add(-30 * time.Second), now, now.Add(30 * time.Second)} {
		c, err := GetTOTP(secret, t)
		if err == nil && c == code {
			return true
		}
	}
	return false
}
