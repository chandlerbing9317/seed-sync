package cookieCloud

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strings"
	"sync"
)

const (
	pkcs5SaltLen = 8
	aes256KeyLen = 32
)

type CookieCloudClient struct {
	config *CookieCloudConfig
	lock   sync.Mutex
}

// cookie-cloud返回的加密数据
type CookieCloudEncryptedData struct {
	Uuid      string `json:"uuid,omitempty"`
	Encrypted string `json:"encrypted,omitempty"`
}
type CookieData struct {
	Domain         string  `json:"domain"`
	ExpirationDate float64 `json:"expirationDate,omitempty"`
	HostOnly       bool    `json:"hostOnly"`
	HttpOnly       bool    `json:"httpOnly"`
	Name           string  `json:"name"`
	Path           string  `json:"path"`
	SameSite       string  `json:"sameSite"`
	Secure         bool    `json:"secure"`
	Session        bool    `json:"session"`
	StoreId        string  `json:"storeId"`
	Value          string  `json:"value"`
}

// cookie-cloud返回的数据
type CookieCloudResponse struct {
	CookieData       map[string][]CookieData `json:"cookie_data"`
	LocalStorageData map[string]interface{}  `json:"local_storage_data"`
	UpdateTime       string                  `json:"update_time"`
}

// 获取指定域名的cookie
func (data *CookieCloudResponse) GetCookieByDomain(domain string) (string, bool) {
	cookies, ok := data.CookieData[domain]
	if !ok {
		return "", false
	}
	cookieStr := ""
	for _, cookie := range cookies {
		cookieStr += fmt.Sprintf("%s=%s;", cookie.Name, cookie.Value)
	}
	return cookieStr, true
}

// 初始化cookie-cloud客户端 这个客户端由service层管理
func NewCookieCloudClient(cookieCloudConfig *CookieCloudConfig) (*CookieCloudClient, error) {
	return &CookieCloudClient{
		config: cookieCloudConfig,
		lock:   sync.Mutex{},
	}, nil
}

// 销毁cookie-cloud客户端
func (c *CookieCloudClient) Destroy() {
	c.config = nil
}

// 更新cookie-cloud客户端
func (c *CookieCloudClient) Update(config *CookieCloudConfig) (*CookieCloudClient, error) {
	c.config = config
	return c, nil
}

func (c *CookieCloudClient) GetConfig() *CookieCloudConfig {
	return c.config
}

// 获取cookie
func (c *CookieCloudClient) GetCookie() (*CookieCloudResponse, error) {
	if !strings.HasSuffix(c.config.Url, "/") {
		c.config.Url = c.config.Url + "/"
	}

	url := c.config.Url + "get/" + c.config.UserKey

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("cookie-cloud server return status %d", res.StatusCode)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var data *CookieCloudEncryptedData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cookie-cloud server response as json: %v", err)
	}
	keyPassword := Md5String(c.config.UserKey, "-", c.config.P2pPassword)[:16]
	decrypted, err := DecryptCryptoJsAesMsg(keyPassword, data.Encrypted)
	if err != nil {
		return nil, err
	}
	//解析结果转json
	var cookieCloudResponse CookieCloudResponse
	err = json.Unmarshal(decrypted, &cookieCloudResponse)
	if err != nil {
		return nil, err
	}
	return &cookieCloudResponse, nil
}

//cookie-cloud的解码 ref https://github.com/easychen/CookieCloud/blob/master/README.md

// Decrypt a CryptoJS.AES.encrypt(msg, password) encrypted msg.
// ciphertext is the result of CryptoJS.AES.encrypt(), which is the base64 string of
// "Salted__" + [8 bytes random salt] + [actual ciphertext].
// actual ciphertext is padded (make it's length align with block length) using Pkcs7.
// CryptoJS use a OpenSSL-compatible EVP_BytesToKey to derive (key,iv) from (password,salt),
// using md5 as hash type and 32 / 16 as length of key / block.
// See: https://stackoverflow.com/questions/35472396/how-does-cryptojs-get-an-iv-when-none-is-specified ,
// https://stackoverflow.com/questions/64797987/what-is-the-default-aes-config-in-crypto-js
func DecryptCryptoJsAesMsg(password string, ciphertext string) ([]byte, error) {
	const keylen = 32
	const blocklen = 16
	rawEncrypted, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode Encrypted: %v", err)
	}
	if len(rawEncrypted) < 17 || len(rawEncrypted)%blocklen != 0 || string(rawEncrypted[:8]) != "Salted__" {
		return nil, fmt.Errorf("invalid ciphertext")
	}
	salt := rawEncrypted[8:16]
	encrypted := rawEncrypted[16:]
	key, iv := BytesToKey(salt, []byte(password), md5.New(), keylen, blocklen)
	newCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create aes cipher: %v", err)
	}
	cfbdec := cipher.NewCBCDecrypter(newCipher, iv)
	decrypted := make([]byte, len(encrypted))
	cfbdec.CryptBlocks(decrypted, encrypted)
	decrypted, err = pkcs7strip(decrypted, blocklen)
	if err != nil {
		return nil, fmt.Errorf("failed to strip pkcs7 paddings (password may be incorrect): %v", err)
	}
	return decrypted, nil
}

// From https://github.com/walkert/go-evp .
// BytesToKey implements the Openssl EVP_BytesToKey logic.
// It takes the salt, data, a hash type and the key/block length used by that type.
// As such it differs considerably from the openssl method in C.
func BytesToKey(salt, data []byte, h hash.Hash, keyLen, blockLen int) (key, iv []byte) {
	saltLen := len(salt)
	if saltLen > 0 && saltLen != pkcs5SaltLen {
		panic(fmt.Sprintf("Salt length is %d, expected %d", saltLen, pkcs5SaltLen))
	}
	var (
		concat   []byte
		lastHash []byte
		totalLen = keyLen + blockLen
	)
	for ; len(concat) < totalLen; h.Reset() {
		// concatenate lastHash, data and salt and write them to the hash
		h.Write(append(lastHash, append(data, salt...)...))
		// passing nil to Sum() will return the current hash value
		lastHash = h.Sum(nil)
		// append lastHash to the running total bytes
		concat = append(concat, lastHash...)
	}
	return concat[:keyLen], concat[keyLen:totalLen]
}

// BytesToKeyAES256CBC implements the SHA256 version of EVP_BytesToKey using AES CBC
func BytesToKeyAES256CBC(salt, data []byte) (key []byte, iv []byte) {
	return BytesToKey(salt, data, sha256.New(), aes256KeyLen, aes.BlockSize)
}

// BytesToKeyAES256CBCMD5 implements the MD5 version of EVP_BytesToKey using AES CBC
func BytesToKeyAES256CBCMD5(salt, data []byte) (key []byte, iv []byte) {
	return BytesToKey(salt, data, md5.New(), aes256KeyLen, aes.BlockSize)
}

// return the MD5 hex hash string (lower-case) of input string(s)
func Md5String(inputs ...string) string {
	keyHash := md5.New()
	for _, str := range inputs {
		io.WriteString(keyHash, str)
	}
	return hex.EncodeToString(keyHash.Sum(nil))
}

// from https://gist.github.com/nanmu42/b838acc10d393bc51cb861128ce7f89c .
// pkcs7strip remove pkcs7 padding
func pkcs7strip(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("pkcs7: Data is empty")
	}
	if length%blockSize != 0 {
		return nil, errors.New("pkcs7: Data is not block-aligned")
	}
	padLen := int(data[length-1])
	ref := bytes.Repeat([]byte{byte(padLen)}, padLen)
	if padLen > blockSize || padLen == 0 || !bytes.HasSuffix(data, ref) {
		return nil, errors.New("pkcs7: Invalid padding")
	}
	return data[:length-padLen], nil
}
