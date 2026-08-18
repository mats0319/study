package test

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/mats0319/secure_transfer/internal"
	"github.com/mats0319/secure_transfer/internal/lib"
	"github.com/mats0319/secure_transfer/utils"
	mlog "github.com/mats0319/secure_transfer/utils/log"
)

const (
	publicKeyFileName  = "PUB.KEY"
	originFileName     = "message.txt"
	encFileName_Once   = "CIPHER_ONCE.TXT"
	encFileName_Stream = "CIPHER_STREAM.TXT"
)

func BenchmarkEncryptSameFile(b *testing.B) {
	err := generateRandomFile("message.txt", 200)
	if err != nil {
		b.Fatal(err)
	}

	err = internal.GenerateKeyPair(true)
	if err != nil {
		b.Fatal(err)
	}

	pubKey, e := deserializePublicKey()
	if e != nil {
		b.Fatal(err)
	}

	b.Run("Encrypt Once", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			enc := lib.NewEncryptorOnce(pubKey)
			err := enc.Encrypt(originFileName, encFileName_Once)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Encrypt Stream", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			enc := lib.NewEncryptorStream(pubKey)
			err := enc.Encrypt(originFileName, encFileName_Stream)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func deserializePublicKey() (pubKey *ecdh.PublicKey, e *utils.Error) {
	pubKeyBytes, err := os.ReadFile(publicKeyFileName)
	if err != nil {
		e = utils.ErrReadPublicKey().WithCause(err)
		mlog.Error(e.String())
		return
	}

	block, _ := pem.Decode(pubKeyBytes)
	if block == nil {
		e = utils.ErrDecodePublicKey()
		mlog.Error(e.String())
		return
	}

	pubKeyI, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		e = utils.ErrParsePublicKey().WithCause(err)
		mlog.Error(e.String())
		return
	}

	pubKey, ok := pubKeyI.(*ecdh.PublicKey)
	if !ok {
		ecdsaPubKey, ok := pubKeyI.(*ecdsa.PublicKey) // x509 parser usually return this type
		if !ok {
			e = utils.ErrInvalidPublicKey()
			mlog.Error(e.String())
			return
		}

		pubKey, err = ecdsaPubKey.ECDH()
		if err != nil {
			e = utils.ErrInvalidPublicKey().WithCause(err)
			mlog.Error(e.String())
			return
		}
	}

	return
}
