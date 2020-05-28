package myrsa

import (
	"errors"
	"io"
	"math/big"
)

// 前15个质数
var smallPrimes = []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53}

// 前15个质数的乘积
var smallPrimesMult = new(big.Int).SetUint64(16294579238595022365)

var smallPrimesBig = []*big.Int{new(big.Int).SetUint64(2), new(big.Int).SetUint64(3), new(big.Int).SetUint64(5), new(big.Int).SetUint64(7),
	new(big.Int).SetUint64(11), new(big.Int).SetUint64(13), new(big.Int).SetUint64(17), new(big.Int).SetUint64(19), new(big.Int).SetUint64(23), new(big.Int).SetUint64(29),
	new(big.Int).SetUint64(31), new(big.Int).SetUint64(37), new(big.Int).SetUint64(41), new(big.Int).SetUint64(43), new(big.Int).SetUint64(47), new(big.Int).SetUint64(53)}

var oneBig = new(big.Int).SetUint64(1)
var twoBig = new(big.Int).SetUint64(2)

// isPrime 用Miller-Rabin算法检测n是否有极大的可能是质数
func isPrime(n *big.Int) bool {
	subOne := new(big.Int).Sub(n, oneBig)
	d := new(big.Int).Set(subOne) // 这里令d为n-1
	tmp := new(big.Int)
	// 先对于前15个小质数，进行费马素性测试，如果不通过直接就返回
	for _, a := range smallPrimesBig {
		if tmp.Exp(a, d, n).Cmp(oneBig) != 0 {
			return false
		}
	}
	// 下面进行Miller—Rabin检测
	d.Div(d, twoBig)
	for d.Bit(0) == 0 { // 当仍然d为偶数时，继续循环
		for _, a := range smallPrimesBig { // 选择15个小质数作为底
			tmp.Exp(a, d, n)
			if !(tmp.Cmp(oneBig) == 0 || tmp.Cmp(subOne) == 0) { // Miller-Rabin的判定条件
				return false
			}
		}
		d.Div(d, twoBig) // d除以2
	}
	// 此时d已经是奇数
	for _, a := range smallPrimesBig {
		tmp.Exp(a, d, n)
		if !(tmp.Cmp(oneBig) == 0 || tmp.Cmp(subOne) == 0) {
			return false
		}
	}
	// 通过了Miller-Rabin检测
	return true
}

// Prime 生成bits个比特位的质数,bits应为8的整数倍
func Prime(random io.Reader, bits int) (*big.Int, error) {
	if bits < 2 {
		return nil, errors.New("bits too low")
	}
	if bits%8 != 0 {
		return nil, errors.New("bits mod eight is not zero")
	}
	buf := make([]byte, bits/8)
	res := new(big.Int)

	bigDelta := new(big.Int)

	for {
		// 从随机流里读入bits个比特
		if _, err := io.ReadFull(random, buf); err != nil {
			return nil, err
		}
		// 将第一个比特位设为1，让这个数尽量的大
		buf[0] |= (1 << 7)
		// 将最后一个比特位设为1，让这个数为奇数
		buf[len(buf)-1] |= 1

		// 随机生成了一个大奇数，存入res，不妨记为N
		res.SetBytes(buf)

		// 令 bigDelta = N mod (p1*p2*p3*...*p15) ，即对前15个小质数的乘积取模，然后转成uint64
		// 这样做是为了加速运算，因为big.Int的计算过于慢
		bigDelta.Mod(res, smallPrimesMult)
		bigDeltaU64 := bigDelta.Uint64()

		// delta是偏移量，从0开始枚举delta，当发现N+delta通过了前15个小质数的初步检测后，再用Miller-Rabin算法检测N+delta
		for delta := uint64(0); delta < smallPrimesMult.Uint64(); delta++ {
			isDeltaOK := true
			// 用bigDelta+delta模前15个小质数来初步验证N+delta是否是合数，如果是合数就继续枚举delta
			for _, p := range smallPrimes {
				if (bigDeltaU64+delta)%p == 0 {
					isDeltaOK = false
					break
				}
			}
			if isDeltaOK { // delta通过了前15个小质数的初步检测，之后用Miller-Rabin再检测
				bigDeltaU64 = delta
				break
			}
		}
		bigDelta.SetUint64(bigDeltaU64)
		res.Add(res, bigDelta)
		if isPrime(res) && res.BitLen() == bits { // 通过Miller-Rabin检测，并且比特数也的确是bits（因为不断累加delta的过程中，有可能N进位了，长度就变成bits+1了）
			return res, nil
		}
	}
}
