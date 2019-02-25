package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"reflect"

	"golang.org/x/crypto/hkdf"
)

var defaultHash = sha256.New

func AesEncrypt(key, plain []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	text := make([]byte, aes.BlockSize+len(plain))
	iv := text[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(text[aes.BlockSize:], plain)

	return text, nil
}

func AesDecrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(text) < aes.BlockSize {
		return nil, errors.New("cipher text too short")
	}

	cfb := cipher.NewCFBDecrypter(block, text[:aes.BlockSize])
	plain := make([]byte, len(text)-aes.BlockSize)
	cfb.XORKeyStream(plain, text[aes.BlockSize:])

	return plain, nil
}

func EciesEncrypt(raw []byte, plain []byte) ([]byte, error) {
	cert, err := PEMtoCertificate(raw)
	if err != nil {
		return nil, err
	}

	pub, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("the publickey is not ecdsa %v\n", reflect.TypeOf(cert.PublicKey))
	}

	params := pub.Curve

	// Select an ephemeral elliptic curve key pair associated with
	// elliptic curve domain parameters params
	priv, Rx, Ry, err := elliptic.GenerateKey(pub.Curve, rand.Reader)

	// Convert R=(Rx,Ry) to an octed string R bar
	// This is uncompressed
	Rb := elliptic.Marshal(pub.Curve, Rx, Ry)

	// Derive a shared secret field element z from the ephemeral secret key k
	// and convert z to an octet string Z
	z, _ := params.ScalarMult(pub.X, pub.Y, priv)
	Z := z.Bytes()

	// generate keying data K of length ecnKeyLen + macKeyLen octects from Z
	// ans s1
	kE := make([]byte, 32)
	kM := make([]byte, 32)
	hkdfo := hkdf.New(defaultHash, Z, nil, nil)
	_, err = hkdfo.Read(kE)
	if err != nil {
		return nil, err
	}
	_, err = hkdfo.Read(kM)
	if err != nil {
		return nil, err
	}

	// Use the encryption operation of the symmetric encryption scheme
	// to encrypt m under EK as ciphertext EM
	EM, err := AesEncrypt(kE, plain)

	// Use the tagging operation of the MAC scheme to compute
	// the tag D on EM || s2
	mac := hmac.New(defaultHash, kM)
	mac.Write(EM)

	D := mac.Sum(nil)

	// Output R,EM,D
	ciphertext := make([]byte, len(Rb)+len(EM)+len(D))
	copy(ciphertext, Rb)
	copy(ciphertext[len(Rb):], EM)
	copy(ciphertext[len(Rb)+len(EM):], D)

	return ciphertext, nil
}

func EciesDecrypt(raw []byte, ciphertext []byte) ([]byte, error) {
	privInterface, err := PEMtoPrivateKey(raw, nil)
	if err != nil {
		return nil, err
	}

	priv, ok := privInterface.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("the private is not ecdsa %v\n", reflect.TypeOf(privInterface))
	}

	params := priv.Curve
	var (
		rLen   int
		hLen   = defaultHash().Size()
		mStart int
		mEnd   int
	)

	switch ciphertext[0] {
	case 2, 3:
		rLen = ((priv.PublicKey.Curve.Params().BitSize + 7) / 8) + 1
		if len(ciphertext) < (rLen + hLen + 1) {
			return nil, fmt.Errorf("Invalid ciphertext len [First byte = %d]", ciphertext[0])
		}
		break
	case 4:
		rLen = 2*((priv.PublicKey.Curve.Params().BitSize+7)/8) + 1
		if len(ciphertext) < (rLen + hLen + 1) {
			return nil, fmt.Errorf("Invalid ciphertext len [First byte = %d]", ciphertext[0])
		}
		break

	default:
		return nil, fmt.Errorf("Invalid ciphertext. Invalid first byte. [%d]", ciphertext[0])
	}

	mStart = rLen
	mEnd = len(ciphertext) - hLen

	Rx, Ry := elliptic.Unmarshal(priv.Curve, ciphertext[:rLen])
	if Rx == nil {
		return nil, errors.New("Invalid ephemeral PK")
	}
	if !priv.Curve.IsOnCurve(Rx, Ry) {
		return nil, errors.New("Invalid point on curve")
	}

	// Derive a shared secret field element z from the ephemeral secret key k
	// and convert z to an octet string Z
	z, _ := params.ScalarMult(Rx, Ry, priv.D.Bytes())
	Z := z.Bytes()

	// generate keying data K of length ecnKeyLen + macKeyLen octects from Z
	// ans s1
	kE := make([]byte, 32)
	kM := make([]byte, 32)
	hkdfo := hkdf.New(defaultHash, Z, nil, nil)
	_, err = hkdfo.Read(kE)
	if err != nil {
		return nil, err
	}
	_, err = hkdfo.Read(kM)
	if err != nil {
		return nil, err
	}

	// Use the tagging operation of the MAC scheme to compute
	// the tag D on EM || s2 and then compare
	mac := hmac.New(defaultHash, kM)
	mac.Write(ciphertext[mStart:mEnd])

	D := mac.Sum(nil)

	if subtle.ConstantTimeCompare(ciphertext[mEnd:], D) != 1 {
		return nil, errors.New("Tag check failed")
	}

	// Use the decryption operation of the symmetric encryption scheme
	// to decryptr EM under EK as plaintext

	plaintext, err := AesDecrypt(kE, ciphertext[mStart:mEnd])

	return plaintext, err
}

// PEMtoCertificate converts pem to x509
func PEMtoCertificate(raw []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("No PEM block available")
	}

	if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
		return nil, errors.New("Not a valid CERTIFICATE PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// PEMtoPrivateKey unmarshals a pem to private key
func PEMtoPrivateKey(raw []byte, pwd []byte) (interface{}, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("raw is nil")
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("Failed decoding [% x]", raw)
	}

	cert, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, err
}

func GenerateKey(len int) []byte {
	key := make([]byte, len)

	_, err := rand.Read(key)
	if err != nil {
		return nil
	}

	return key
}
