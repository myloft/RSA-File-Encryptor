package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var generate, help, encrypt, decrypt bool
	var keyPath, input, output string
	flag.BoolVar(&help, "h", false, "打印帮助")
	flag.BoolVar(&generate, "g", false, "生成密钥对")
	flag.BoolVar(&encrypt, "e", false, "加密文件")
	flag.BoolVar(&decrypt, "d", false, "解密文件")
	flag.StringVar(&keyPath, "k", "", "指定密钥")
	flag.StringVar(&input, "i", "", "输入文件")
	flag.StringVar(&output, "o", "", "输出文件")

	flag.Parse()

	if help {
		printHelp()
		os.Exit(0)
	}
	if generate {
		fmt.Println("密钥对生成成功")
		fmt.Printf("密钥长度： 2048\n\n")
		generateKey()
		os.Exit(0)
	}

	if encrypt && keyPath != "" && input != "" && output != "" && !decrypt {
		encryptFile(input, output, decryptPublicPEMKey(keyPath))
		os.Exit(0)
	}

	if decrypt && keyPath != "" && input != "" && output != "" && !encrypt {
		decryptFile(input, output, decryptPEMKey(keyPath))
		os.Exit(0)
	}
	printHelp()
}

func printHelp() {
	fmt.Printf("基于 RSA 非对称加密算法的文件分组加密工具 Copyright 信息安全小组\n\n")

	fmt.Printf("参数列表：\n")
	fmt.Printf("-h\t打印帮助\n")
	fmt.Printf("-g\t生成密钥对\n")
	fmt.Printf("-e\t加密文件\n")
	fmt.Printf("-d\t解密文件\n")
	fmt.Printf("-k\t指定密钥位置\n")
	fmt.Printf("-i\t输入文件位置\n")
	fmt.Printf("-o\t输出文件位置\n\n")

	fmt.Printf("使用示例：\n")
	fmt.Printf("生成密钥对：\trfe -g\n")
	fmt.Printf("加密文件：\trfe -e -k [公钥文件路径] -i [待加密文件路径] -o [加密文件输出路径]\n")
	fmt.Printf("解密文件：\trfe -d -k [私钥文件路径] -i [待解密文件路径] -o [解密文件输出路径]\n")
}

//异常处理
func checkError(err error) {
	if err != nil {
		fmt.Println("致命错误：", err.Error())
	}
}
