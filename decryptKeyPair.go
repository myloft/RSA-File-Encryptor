package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

//从磁盘读取私钥并解码
func decryptPEMKey(keyPath string) *rsa.PrivateKey {
	//打开输入流
	file, err := os.Open(keyPath)
	checkError(err)
	defer file.Close()

	//读取文件描述符
	info, err := file.Stat()
	checkError(err)
	//开辟空间
	b := make([]byte, info.Size())
	//读取密钥
	file.Read(b)

	//解码密钥
	pemKey, _ := pem.Decode(b)
	privateKey, err := x509.ParsePKCS1PrivateKey(pemKey.Bytes)
	checkError(err)
	return privateKey
}

//从磁盘读取公钥并解码
func decryptPublicPEMKey(keyPath string) *rsa.PublicKey {
	//打开输入流
	file, err := os.Open(keyPath)
	checkError(err)
	defer file.Close()

	//读取文件描述符
	info, err := file.Stat()
	checkError(err)
	//开辟空间
	b := make([]byte, info.Size())
	//读取密钥
	file.Read(b)

	//解码密钥
	pemKey, _ := pem.Decode(b)
	publicKey, err := x509.ParsePKCS1PublicKey(pemKey.Bytes)
	checkError(err)
	return publicKey
}
