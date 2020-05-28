package myrsa

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)

func TestMyRSA(t *testing.T) {
	key, err := GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 512; i++ {
		randomText := make([]byte, 64)
		io.ReadFull(rand.Reader, randomText)
		if randomText[0]&1 == 0 {
			randomText[0] = 0 // 测试第一个字节为零的情况
		}
		c, err := Encrypt(&key.PublicKey, randomText)
		if err != nil {
			t.Fatal(err)
		}
		text, err := Decrypt(key, c)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(randomText, text) != 0 {
			fmt.Printf("origin text is %v\n", randomText)
			fmt.Printf("encrypted text is %v\n", c)
			fmt.Printf("decrypted text is %v\n", text)
			t.Fatalf("uncorrect!")
		}
	}
}
