package myrsa

import (
	"crypto/rand"
	"testing"
)

func TestPrime(t *testing.T) {
	for i := 0; i < 100; i++ {
		p, err := Prime(rand.Reader, 512)
		if err != nil {
			t.FailNow()
		}
		if p.ProbablyPrime(20) != isPrime(p) || p.BitLen() != 512 {
			t.FailNow()
		}
	}
}

func BenchmarkPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p, err := Prime(rand.Reader, 1024)
		if err != nil {
			b.FailNow()
		}
		_ = p.BitLen()
	}
}

func BenchmarkStandardPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p, err := rand.Prime(rand.Reader, 1024)
		if err != nil {
			b.FailNow()
		}
		_ = p.BitLen()
	}
}
