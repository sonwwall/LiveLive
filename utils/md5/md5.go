package md5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// MD5 use md5 to encrypt strings
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateSign 生成带签名的url路径
func GenerateSign(uri, uid, exp, secret string) string {
	// 拼接签名字符串
	raw := fmt.Sprintf("%s:%s:%s %s", uri, uid, exp, secret)
	hash := md5.Sum([]byte(raw))
	return hex.EncodeToString(hash[:])
}
