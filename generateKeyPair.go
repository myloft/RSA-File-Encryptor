package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

//生成密钥对
func generateKey() {
	//使用随机数生成指定长度的 RSA 密钥对
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	checkError(err)
	//导出公钥
	publicKey := key.PublicKey

	//保存 PEM 格式的密钥对
	savePEMKey("private.pem", key)
	savePublicPEMKey("public.pem", publicKey)
}

//保存 PEM 格式的私钥
func savePEMKey(filename string, key *rsa.PrivateKey) {
	//打开输出流
	outfile, err := os.Create(filename)
	checkError(err)
	defer outfile.Close()

	//PKCS#1 私钥结构体
	privateKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	//打印密钥
	fmt.Println("私钥： private.pem")
	fmt.Println(string(pem.EncodeToMemory(privateKey)))
	//转换格式并写入密钥
	err = pem.Encode(outfile, privateKey)
	checkError(err)
}

//保存 PEM 格式的公钥
func savePublicPEMKey(filename string, key rsa.PublicKey) {
	//打开输出流
	outfile, err := os.Create(filename)
	checkError(err)
	defer outfile.Close()

	//PKCS#1 公钥结构体
	publicKey := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&key),
	}
	//打印密钥
	fmt.Println("公钥： public.pem")
	fmt.Println(string(pem.EncodeToMemory(publicKey)))
	//转换格式并写入密钥
	err = pem.Encode(outfile, publicKey)
	checkError(err)
}
