package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func GenRsaKey(bits int) error {
	//生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	//生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	if err = pem.Encode(file, block); err != nil {
		return err
	}
	return nil
}

// RsaEncrypt RSA加密
// plainText 要加密的数据
// path 公钥文件地址
func RsaEncryptFromPem(plainText []byte, path string) ([]byte, error) {
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//读取文件的内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	return RsaEncrypt(plainText, buf)
}

func RsaDecryptFromPem(cipherText []byte, path string) ([]byte, error) {
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//读取文件的内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	return RsaDecrypt(cipherText, buf)
}

func RsaEncrypt(plainText, pubKeyPem []byte) ([]byte, error) {
	//pem解码
	block, _ := pem.Decode(pubKeyPem)
	//x509解码
	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//类型断言
	pubKey := pubKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, plainText)
}

func RsaDecrypt(cipherText, priKeyPem []byte) ([]byte, error) {
	//pem解码
	block, _ := pem.Decode(priKeyPem)
	//x509解码
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//对密文进行解密
	return rsa.DecryptPKCS1v15(rand.Reader, priKey, cipherText)
}
