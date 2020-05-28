package myrsa

import (
	"crypto/rsa"
	"errors"
	"io"
	"math/big"
)

// PublicKey 公钥
type PublicKey struct {
	N *big.Int
	E int
}

// PrivateKey 私钥，模仿golang标准库里的设计，让私钥继承公钥
type PrivateKey struct {
	PublicKey
	D  *big.Int
	P1 *big.Int
	P2 *big.Int
}

// stdKeyWrapper 将自己定义的key类型转化成golang标准库的key类型
func stdKeyWrapper(key *PrivateKey) (stdKey *rsa.PrivateKey) {
	stdKey = new(rsa.PrivateKey)
	stdKey.D = key.D
	stdKey.E = key.E
	stdKey.N = key.N
	stdKey.Primes = []*big.Int{key.P1, key.P2}
	stdKey.Precompute()
	return stdKey
}

// GenerateKey 生成公私秘钥对
func GenerateKey(random io.Reader, bits int) (*rsa.PrivateKey, error) {
	for {
		twoPrimesOK := false
		var p1 *big.Int
		var p2 *big.Int
		var err error
		for !twoPrimesOK {
			// 生成第一个大质数, 长度为bits/2
			p1, err = Prime(random, bits/2)
			if err != nil {
				return nil, err
			}
			// 生成第二个大质数, 长度为bits-bits/2
			p2, err = Prime(random, bits-bits/2)
			if err != nil {
				return nil, err
			}
			// 如果两个质数不一样，则继续，否则再生成
			if p1.Cmp(p2) != 0 {
				twoPrimesOK = true
			}
		}
		// 秘钥初始化
		key := new(PrivateKey)
		key.P1 = p1
		key.P2 = p2
		key.N = new(big.Int)
		key.N.Mul(p1, p2)
		key.E = 65537 // 按照惯例把E设为65537

		key.D = new(big.Int)

		// phi = (p1-1)*(p2-1)
		phi := new(big.Int).Mul(new(big.Int).Sub(p1, oneBig), new(big.Int).Sub(p2, oneBig))
		// 求满足(e*d) mod phi == 1的d，即求e,phi的模逆元，如果有解就返回，无解就再重复生成
		if key.D = key.D.ModInverse(big.NewInt(int64(key.E)), phi); key.D != nil {
			return stdKeyWrapper(key), nil
		}
	}
}

// Encrypt 加密一串明文，在明文最前面补一个0x01，防止出现第一个字节正好是零时，解密时会漏掉这个字节的情况
func Encrypt(key *rsa.PublicKey, msg []byte) ([]byte, error) {
	if 1+len(msg) >= key.Size() {
		return nil, errors.New("msg too long")
	}
	wrappedMsg := make([]byte, 1+len(msg))
	wrappedMsg[0] = 1
	copy(wrappedMsg[1:], msg)
	en := make([]byte, key.Size())
	// 计算m^e mod N
	cBig := new(big.Int).Exp(new(big.Int).SetBytes(wrappedMsg), big.NewInt(int64(key.E)), key.N)
	c := cBig.Bytes()

	// 如果密文长度小于N，则在左侧补零
	if len(c) < key.Size() {
		d := key.Size() - len(c)
		copy(en[d:], c)
		for i := range en[0:d] {
			en[i] = 0
		}
	} else {
		copy(en, c)
	}

	return en, nil
}

// Decrypt 解密一串密文
func Decrypt(key *rsa.PrivateKey, c []byte) ([]byte, error) {
	m := new(big.Int)
	// 计算c^d mod N
	m.Exp(new(big.Int).SetBytes(c), key.D, key.N)
	// 解密的明文需要去掉第一个字节
	return m.Bytes()[1:], nil
}
