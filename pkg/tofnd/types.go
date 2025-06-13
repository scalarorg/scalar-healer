package tofnd

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"

	ec "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signature [crypto.SignatureLength]byte

// NewSignature is the constructor of Signature
func NewSignature(bz []byte) (sig Signature, err error) {
	if len(bz) != crypto.SignatureLength {
		return Signature{}, errors.New("invalid signature length")
	}

	copy(sig[:], bz)

	return sig, nil
}

// Hex returns the hex-encoding of the given Signature
func (s Signature) Hex() string {
	return hex.EncodeToString(s[:])
}

// ToHomesteadSig converts signature to openzeppelin compatible
func (s Signature) ToHomesteadSig() []byte {
	/* TODO: We have to make v 27 or 28 due to openzeppelin's implementation at https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/cryptography/ECDSA.sol
	requiring that. Consider copying and modifying it to require v to be just 0 or 1
	instead.
	*/
	bz := s[:]
	if bz[crypto.SignatureLength-1] == 0 || bz[crypto.SignatureLength-1] == 1 {
		bz[crypto.SignatureLength-1] += 27
	}

	return bz
}

// ToSignature transforms an Scalar generated signature into a recoverable signature
func ToSignature(sig ec.Signature, hash common.Hash) (*Signature, *ecdsa.PublicKey, error) {
	s := Signature{}
	encSig := sig.Serialize()

	// read R length
	encSig = encSig[3:]
	rLen := int(encSig[0])
	encSig = encSig[1:]

	// extract R
	encR := encSig[:rLen]
	if encR[0] == 0 {
		encR = encR[1:]
	}
	copy(s[:32], common.LeftPadBytes(encR, 32))
	encSig = encSig[rLen:]

	// read S length
	encSig = encSig[1:]
	sLen := int(encSig[0])
	encSig = encSig[1:]

	// extract S
	encS := encSig[:sLen]
	if encS[0] == 0 {
		encS = encS[1:]
	}
	copy(s[32:], common.LeftPadBytes(encS, 32))

	// s[64] = 0 implicit

	derivedPK, err := crypto.SigToPub(hash.Bytes(), s[:])
	if err != nil {
		return nil, nil, err
	}

	s[64] = 1

	return &s, derivedPK, nil
}
