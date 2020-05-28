package main

import (
	"crypto/rsa"
	"main/myrsa"
	"os"
)

//加密文件
func encryptFile(inputPath string, outputPath string, key *rsa.PublicKey) {
	//打开输入流
	file, err := os.Open(inputPath)
	checkError(err)
	defer file.Close()

	//打开输出流
	outfile, err := os.Create(outputPath)
	checkError(err)
	defer outfile.Close()

	//读取文件描述符
	info, err := file.Stat()
	checkError(err)
	//开辟空间
	b := make([]byte, info.Size())
	//读取文件
	file.Read(b)

	//分组加密文件
	len := len(b)
	size := len / 245
	//循环加密
	for i := 0; i < size; i++ {
		block := b[i*245 : (i+1)*245]
		res, err := myrsa.Encrypt(key, block)
		checkError(err)
		outfile.Write(res)
	}
	//处理剩余部分
	block := b[size*245:]
	res, err := myrsa.Encrypt(key, block)
	checkError(err)
	outfile.Write(res)
}

func decryptFile(inputPath string, outputPath string, key *rsa.PrivateKey) {
	//打开输入流
	file, err := os.Open(inputPath)
	checkError(err)
	defer file.Close()

	//打开输出流
	outfile, err := os.Create(outputPath)
	checkError(err)
	defer outfile.Close()

	//读取文件描述符
	info, err := file.Stat()
	checkError(err)
	//开辟空间
	b := make([]byte, info.Size())
	//读取文件
	file.Read(b)

	//分组解密文件
	len := len(b)
	size := len / 256
	//循环解密
	for i := 0; i < size; i++ {
		block := b[i*256 : (i+1)*256]
		res, err := myrsa.Decrypt(key, block)
		checkError(err)
		outfile.Write(res)
	}
}
