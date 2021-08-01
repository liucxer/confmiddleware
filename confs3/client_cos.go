package confs3

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"hash"
	"net/url"
	"strings"
	"time"
)

func COSPresignedValues(db *ObjectDB, objectKey string, expiresIn time.Duration) url.Values {
	authTime := NewAuthTime(expiresIn)
	signTime := authTime.signString()
	keyTime := authTime.keyString()
	signKey := calSignKey(db.SecretAccessKey.String(), keyTime)
	formatString := genFormatString("get", "/"+objectKey, "", "")
	stringToSign := calStringToSign(sha1SignAlgorithm, keyTime, formatString)
	signature := calSignature(signKey, stringToSign)
	signedHeaderList := make([]string, 0)
	signedParameterList := make([]string, 0)

	values := url.Values{}

	values.Set("q-sign-algorithm", sha1SignAlgorithm)
	values.Set("q-ak", db.AccessKeyID)
	values.Set("q-sign-time", signTime)
	values.Set("q-objectKey-time", keyTime)
	values.Set("q-header-list", strings.Join(signedHeaderList, ";"))
	values.Set("q-url-param-list", strings.Join(signedParameterList, ";"))
	values.Set("q-signature", signature)

	return values
}

// NewAuthTime 生成 AuthTime 的便捷函数
//   expire: 从现在开始多久过期.
func NewAuthTime(expire time.Duration) *AuthTime {
	if expire == time.Duration(0) {
		expire = defaultAuthExpire
	}
	signStartTime := time.Now()
	keyStartTime := signStartTime
	signEndTime := signStartTime.Add(expire)
	keyEndTime := signEndTime
	return &AuthTime{
		SignStartTime: signStartTime,
		SignEndTime:   signEndTime,
		KeyStartTime:  keyStartTime,
		KeyEndTime:    keyEndTime,
	}
}

const sha1SignAlgorithm = "sha1"
const defaultAuthExpire = time.Hour

// AuthTime 用于生成签名所需的 q-signSearch-time 和 q-key-time 相关参数
type AuthTime struct {
	SignStartTime time.Time
	SignEndTime   time.Time
	KeyStartTime  time.Time
	KeyEndTime    time.Time
}

// signString return q-signSearch-time string
func (a *AuthTime) signString() string {
	return fmt.Sprintf("%d;%d", a.SignStartTime.Unix(), a.SignEndTime.Unix())
}

// keyString return q-key-time string
func (a *AuthTime) keyString() string {
	return fmt.Sprintf("%d;%d", a.KeyStartTime.Unix(), a.KeyEndTime.Unix())
}

// calSignKey 计算 SignKey
func calSignKey(secretKey, keyTime string) string {
	digest := HMAC(secretKey, keyTime, sha1SignAlgorithm)
	return fmt.Sprintf("%x", digest)
}

// calStringToSign 计算 StringToSign
func calStringToSign(signAlgorithm, signTime, formatString string) string {
	h := sha1.New()
	h.Write([]byte(formatString))
	return fmt.Sprintf("%s\n%s\n%x\n", signAlgorithm, signTime, h.Sum(nil))
}

// calSignature 计算 Signature
func calSignature(signKey, stringToSign string) string {
	digest := HMAC(signKey, stringToSign, sha1SignAlgorithm)
	return fmt.Sprintf("%x", digest)
}

// genFormatString 生成 FormatString
func genFormatString(method string, url string, formatParameters, formatHeaders string) string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n", method, url,
		formatParameters, formatHeaders,
	)
}

// HMAC 签名
func HMAC(key, msg, signMethod string) []byte {
	var hashFunc func() hash.Hash
	switch signMethod {
	case "sha1":
		hashFunc = sha1.New
	default:
		hashFunc = sha1.New
	}
	h := hmac.New(hashFunc, []byte(key))
	h.Write([]byte(msg))
	return h.Sum(nil)
}
