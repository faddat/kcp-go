package kcp

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"hash/crc32"
	"io"
	"testing"
)

func TestSM4(t *testing.T) {
	bc, err := NewSM4BlockCrypt(pass[:16])
	if err != nil {
		t.Fatal(err)
	}
	cryptTest(t, bc)
}

func TestAES(t *testing.T) {
	bc, err := NewAESBlockCrypt(pass[:32])
	if err != nil {
		t.Fatal(err)
	}
	cryptTest(t, bc)
}

func TestTEA(t *testing.T) {
	bc, err := NewTEABlockCrypt(pass[:16])
	if err != nil {
		t.Fatal(err)
	}
	cryptTest(t, bc)
}

func TestXOR(t *testing.T) {
	bc, err := NewSimpleXORBlockCrypt(pass[:32])
	if err != nil {
		t.Fatal(err)
	}
	cryptTest(t, bc)
}

func TestNone(t *testing.T) {
	bc, err := NewNoneBlockCrypt(pass[:32])
	if err != nil {
		t.Fatal(err)
	}
	cryptTest(t, bc)
}

func TestCast5(t *testing.T) {
	bc, err := NewCast5BlockCrypt(pass[:16])
	if err != nil {
		t.Fatal(err)
	}
	cryptTest(t, bc)
}

func TestSalsa20(t *testing.T) {
	bc, err := NewSalsa20BlockCrypt(pass[:32])
	if err != nil {
		t.Fatal(err)
	}
	cryptTest(t, bc)
}

func cryptTest(t *testing.T, bc BlockCrypt) {
	t.Helper()
	data := make([]byte, mtuLimit)
	io.ReadFull(rand.Reader, data)
	dec := make([]byte, mtuLimit)
	enc := make([]byte, mtuLimit)
	bc.Encrypt(enc, data)
	bc.Decrypt(dec, enc)
	if !bytes.Equal(data, dec) {
		t.Fail()
	}
}

func BenchmarkSM4(b *testing.B) {
	bc, err := NewSM4BlockCrypt(pass[:16])
	if err != nil {
		b.Fatal(err)
	}
	benchCrypt(b, bc)
}

func BenchmarkAES128(b *testing.B) {
	bc, err := NewAESBlockCrypt(pass[:16])
	if err != nil {
		b.Fatal(err)
	}

	benchCrypt(b, bc)
}

func BenchmarkAES192(b *testing.B) {
	bc, err := NewAESBlockCrypt(pass[:24])
	if err != nil {
		b.Fatal(err)
	}

	benchCrypt(b, bc)
}

func BenchmarkAES256(b *testing.B) {
	bc, err := NewAESBlockCrypt(pass[:32])
	if err != nil {
		b.Fatal(err)
	}

	benchCrypt(b, bc)
}

func BenchmarkTEA(b *testing.B) {
	bc, err := NewTEABlockCrypt(pass[:16])
	if err != nil {
		b.Fatal(err)
	}
	benchCrypt(b, bc)
}

func BenchmarkXOR(b *testing.B) {
	bc, err := NewSimpleXORBlockCrypt(pass[:32])
	if err != nil {
		b.Fatal(err)
	}
	benchCrypt(b, bc)
}

func BenchmarkNone(b *testing.B) {
	bc, err := NewNoneBlockCrypt(pass[:32])
	if err != nil {
		b.Fatal(err)
	}
	benchCrypt(b, bc)
}

func BenchmarkCast5(b *testing.B) {
	bc, err := NewCast5BlockCrypt(pass[:16])
	if err != nil {
		b.Fatal(err)
	}
	benchCrypt(b, bc)
}

func BenchmarkSalsa20(b *testing.B) {
	bc, err := NewSalsa20BlockCrypt(pass[:32])
	if err != nil {
		b.Fatal(err)
	}
	benchCrypt(b, bc)
}

func benchCrypt(b *testing.B, bc BlockCrypt) {
	b.Helper()
	data := make([]byte, mtuLimit)
	io.ReadFull(rand.Reader, data)
	dec := make([]byte, mtuLimit)
	enc := make([]byte, mtuLimit)

	b.ReportAllocs()
	b.SetBytes(int64(len(enc) * 2))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bc.Encrypt(enc, data)
		bc.Decrypt(dec, enc)
	}
}

func BenchmarkCRC32(b *testing.B) {
	content := make([]byte, 1024)
	b.SetBytes(int64(len(content)))
	for i := 0; i < b.N; i++ {
		crc32.ChecksumIEEE(content)
	}
}

func BenchmarkCsprngSystem(b *testing.B) {
	data := make([]byte, sha256.Size)
	b.SetBytes(int64(len(data)))

	for i := 0; i < b.N; i++ {
		io.ReadFull(rand.Reader, data)
	}
}

func BenchmarkCsprngSHA256(b *testing.B) {
	var data [sha256.Size]byte
	b.SetBytes(sha256.Size)

	for i := 0; i < b.N; i++ {
		data = sha256.Sum256(data[:])
	}
}

func BenchmarkCsprngNonceMD5(b *testing.B) {
	var ng nonceSHA256
	ng.Init()
	b.SetBytes(sha256.Size)
	data := make([]byte, sha256.Size)
	for i := 0; i < b.N; i++ {
		ng.Fill(data)
	}
}

func BenchmarkCsprngNonceAES128(b *testing.B) {
	var ng nonceAES128
	ng.Init()

	b.SetBytes(aes.BlockSize)
	data := make([]byte, aes.BlockSize)
	for i := 0; i < b.N; i++ {
		ng.Fill(data)
	}
}
