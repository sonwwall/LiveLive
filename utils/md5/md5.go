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
func GenerateSign(uri string, exp int64, secret string) string {
	data := fmt.Sprintf("%s:%d %s", uri, exp, secret)
	md5Sum := md5.Sum([]byte(data))
	return hex.EncodeToString(md5Sum[:])
}
