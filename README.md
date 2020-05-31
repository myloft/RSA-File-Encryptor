# RSA-File-Encryptor
基于 RSA 非对称加密算法的文件分组加密工具

## 说明
### 编译
```bash
make
```
### 运行单元测试
```bash
make unittest
```
### 运行性能对比测试
```bash
make bench
```

## 功能
- 生成 RSA 密钥

- 加密文件

- 解密文件

## 使用示例
### 查看帮助
```
rfe -h
```
### 生成密钥对
```
rfe -g
```
### 加密文件
```
rfe -e -k [公钥文件路径] -i [待加密文件路径] -o [加密文件输出路径]
```
### 解密文件
```
rfe -d -k [私钥文件路径] -i [待解密文件路径] -o [解密文件输出路径]
```

## 设计文档

### myrsa

prime.go中实现了两个函数，分别是：

- isPrime：用Miller-Rabin算法检测一个数是否有极大的可能是质数（以前15个质数为底来检测）
- Prime：随机生成一个指定二进制位数的质数(极大的可能)
  - 实现思路是：随机生成一个大奇数，然后利用isPrime判定它是否是质数，如果不是，就继续向后找，由于质数分布规律（分布密度与ln(n)成反比），一般来说可以在O(log(n))的时间里找到一个质数

rsa.go中主要实现了三个函数，分别是：

- GenerateKey：生成rsa公私秘钥对
- Encrypt：使用公钥加密一串字节buffer
  - 在明文最前面补一个0x01，防止出现第一个字节正好是零时，解密时会漏掉这个字节的情况
- Decrypt：使用私钥解密一串字节buffer
  - 对应的，解密后返回的明文也会去掉开头的第一个字节

### 主程序

主程序主要包含三大功能模块，分别是：

- RSA 密钥对生成

- 文本密钥文件解码

- 文件加解密


RSA 密钥对生成由generateKeyPair.go实现，包含以下三个函数：

- generateKey：调用myrsa中实现的rsa公私秘钥对函数生成密钥。

- savePEMKey：依据PKCS#1标准将生成的原始密钥导出为PEM格式的私钥。

- savePublicPEMKey：依据PKCS#1标准将生成的原始密钥导出为PEM格式的公钥。


文本密钥文件解码由decryptKeyPair.go实现，包含以下两个函数：

- decryptPEMKey：依据PKCS#1标准将PEM格式的私钥转换成原始密钥。

- decryptPublicPEMKey：依据PKCS#1标准将PEM格式的公钥转换成原始密钥。


文件加解密由fileEncryptDecrypt.go实现，包含以下两个函数：

- encryptFile：使用使用myrsa中实现的Encrypt函数进行加密。
  - 由于2048bit长的密钥最大加密长度为245byte，因此将需要加密的文件分为长度为245byte的分组后分别加密。

- decryptFile：使用myrsa中实现的Decrypt函数进行解密。
  - 由于2048bit长的密钥加密后的密文长度为256byte，因此将需要解密的文件分为长度为256byte的分组后分别解密。
