package Apush

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 生成查询条件。
func GetQueryString() ([]string, string) {
	token := base64.StdEncoding.EncodeToString(AesEncrypt([]byte("token=" + fmt.Sprintf("%d", time.Now().UnixNano() / 1e6 + 100000) + ",P_292$0$"), GetKeyByte("YouMustModifyThisToAnOtherString")))
	urlToken := url.QueryEscape(token)
	queryStr := "?param=" + fmt.Sprintf("%x", md5.Sum([]byte("P_292$0$"))) + ",pmsid1," + urlToken
	resp, err := http.Get("http://10.125.8.90:6080/apush/1/" + queryStr)

	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return strings.Split(string(body), ":"), queryStr
}

// 加密
func AesEncrypt(content, key []byte) (crypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		// TODO
	}
	ecb := NewECBEncrypter(block)
	content = PKCS5Padding(content, block.BlockSize())
	crypted = make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return
}

// 解密
func AesDecrypt(crypted, key []byte) (content []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		// TODO
	}
	blockMode := NewECBDecrypter(block)
	content = make([]byte, len(crypted))
	blockMode.CryptBlocks(content, crypted)
	content = PKCS5UnPadding(content)
	return
}

// 公用的数据结构。
type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

// 补全数据。
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 去除补全的数据。
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length - 1])
	return origData[:(length - unpadding)]
}

// ECB模式的加密实现。
type ecbEncrypter ecb

func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int {
	return x.blockSize
}
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src) % x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// ECB模式的解密实现。
type ecbDecrypter ecb

func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int {
	return x.blockSize
}
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src) % x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// 生成对应的key的数组。
func GetKeyByte(key string) (ret []byte) {
	ret = []byte(key)

	if len(ret) < 16 {
		ret = append(ret, bytes.Repeat([]byte{0}, 16 - len(ret))...)
	} else if len(ret) > 16 {
		ret = ret[:16]
	}

	return
}
