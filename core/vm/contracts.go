// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	bwcrypto "github.com/immesys/bw2/crypto"
	"github.com/immesys/bw2/objects"
	"github.com/immesys/bw2bc/common"
	"github.com/immesys/bw2bc/crypto"
	"github.com/immesys/bw2bc/logger"
	"github.com/immesys/bw2bc/logger/glog"
	"github.com/immesys/bw2bc/params"
)

const (
	BWGas = 3000
)

// import sha3; f = lambda x : sha3.sha3_256(x).hexdigest()[:8]
var SigVerifyEd25519 = []byte{0x3a, 0x05, 0x40, 0xba} // VerifyEd25519(bytes) -> bool
//flags is [struct/sig valid], [numrevokers], 000000
// UnpackDOT(bytes) (bytes32 flags, uint64 expiry, bytes32 srcvk, bytes32 dstvk, hash dothash)
var SigUnpackDOT = []byte{0x3a, 0xab, 0x47, 0xb2} //UnpackDOT(bytes)
// GetDOTDelegatedRevoker(bytes, index) (bytes32)
var SigGetDOTDelegatedRevoker = []byte{0x5a, 0x26, 0x93, 0xb8} //GetDOTDelegatedRevoker(bytes,uint8)
// UnpackEntity(bytes) (bytes32 flags, uint64 expiry, bytes32 vk)
var SigUnpackEntity = []byte{0xe7, 0xb6, 0x86, 0xa7} //UnpackEntity(bytes)
// GetEntityDelegatedRevoker(bytes, index) (bytes32)
var SigGetEntityDelegatedRevoker = []byte{0xa3, 0xb3, 0xf9, 0x25} //GetEntityDelegatedRevoker(bytes,uint8)
// flags isflags is [struct/sig valid], [numdots], 000000
// UnpackDChain(bytes) (bytes32 flags, hash chainhash)
var SigUnpackDChain = []byte{0xb3, 0x8b, 0xbf, 0xbf} //UnpackDChain(bytes)
// GetChainDOTHash(bytes, index) (hash dothash)
var SigGetChainDOTHash = []byte{0xdc, 0xc4, 0xaa, 0x85} //GetChainDOTHash(bytes,uint8)
// SliceByte32(bytes, offset) (byte32)
var SigSliceByte32 = []byte{0xce, 0x7a, 0x94, 0xeb} //SliceByte32(bytes,uint32)
// flags is [sigvalid] 000000
// UnpackRevocation(bytes) (bytes32 flags, bytes32 target)
var SigUnpackRevocation = []byte{0xe5, 0x73, 0x1b, 0x77} //UnpackRevocation(bytes)
// PrecompiledAccount represents a native ethereum contract
type PrecompiledAccount struct {
	Gas func(l int) *big.Int
	fn  func(in []byte) []byte
}

// Call calls the native function
func (self PrecompiledAccount) Call(in []byte) []byte {
	return self.fn(in)
}

// Precompiled contains the default set of ethereum contracts
var Precompiled = PrecompiledContracts()

// PrecompiledContracts returns the default set of precompiled ethereum
// contracts defined by the ethereum yellow paper.
func PrecompiledContracts() map[string]*PrecompiledAccount {
	return map[string]*PrecompiledAccount{
		// ECRECOVER
		string(common.LeftPadBytes([]byte{1}, 20)): &PrecompiledAccount{func(l int) *big.Int {
			return params.EcrecoverGas
		}, ecrecoverFunc},

		// SHA256
		string(common.LeftPadBytes([]byte{2}, 20)): &PrecompiledAccount{func(l int) *big.Int {
			n := big.NewInt(int64(l+31) / 32)
			n.Mul(n, params.Sha256WordGas)
			return n.Add(n, params.Sha256Gas)
		}, sha256Func},

		// RIPEMD160
		string(common.LeftPadBytes([]byte{3}, 20)): &PrecompiledAccount{func(l int) *big.Int {
			n := big.NewInt(int64(l+31) / 32)
			n.Mul(n, params.Ripemd160WordGas)
			return n.Add(n, params.Ripemd160Gas)
		}, ripemd160Func},

		string(common.LeftPadBytes([]byte{4}, 20)): &PrecompiledAccount{func(l int) *big.Int {
			n := big.NewInt(int64(l+31) / 32)
			n.Mul(n, params.IdentityWordGas)

			return n.Add(n, params.IdentityGas)
		}, memCpy},

		// BWFunctions
		string(common.LeftPadBytes([]byte{0x2, 0x85, 0x89}, 20)): &PrecompiledAccount{func(l int) *big.Int {
			return big.NewInt(BWGas)
		}, bosswave},
	}
}

func getBytesParam(in []byte, paramnumber int) []byte {
	offset := common.BytesToBig(in[(paramnumber * 32) : (paramnumber+1)*32]).Uint64()
	length := common.BytesToBig(in[offset : offset+32]).Uint64()
	return in[offset+32 : offset+32+length]
}
func bosswave(in []byte) []byte {
	fmt.Printf("GOT: %d bytes: %v\n", len(in), hex.EncodeToString(in))
	if len(in) < 4 {
		return nil
	}
	sig := in[:4]
	args := in[4:]
	switch {
	case bytes.Equal(sig, SigVerifyEd25519):
		fmt.Printf("args len is %v\n", len(args))
		startOfPayload := common.BytesToBig(args[:32]).Uint64()
		lengthOfPayload := common.BytesToBig(args[startOfPayload : startOfPayload+32]).Uint64()
		actualArgs := args[startOfPayload+32 : startOfPayload+32+lengthOfPayload]
		if lengthOfPayload < 96 {
			return nil
		}
		vk := actualArgs[:32]
		bodyEnd := lengthOfPayload - 64
		sig := actualArgs[bodyEnd:]
		body := actualArgs[:bodyEnd]
		fmt.Printf("Doing VFB:\nvk =%s\nsig=%s\nbdy=%s\n",
			hex.EncodeToString(vk), hex.EncodeToString(sig), hex.EncodeToString(body))
		if bwcrypto.VerifyBlob(vk, sig, body) {
			fmt.Println("BWVFBLOB WAS OK\n")
			return common.LeftPadBytes([]byte{1}, 32)
		} else {
			fmt.Println("BWVFBLOB WAS BAD\n")
			return common.LeftPadBytes([]byte{0}, 32)
		}
	case bytes.Equal(sig, SigUnpackDOT):
		// UnpackDOT(bytes) -> (bytes32 flags, uint64 expiry, bytes32 srcvk, bytes32 dstvk, hash dothash)
		// Bit of a hack, we read the DOT type from the actual DOT itself
		blob := getBytesParam(args, 0)
		res := make([]byte, 32*5)
		if len(blob) < 96 {
			return res
		}
		ronum := objects.ROAccessDOT
		if blob[65] == 0x02 {
			ronum = objects.ROPermissionDOT
		}
		ro, err := objects.NewDOT(ronum, blob)
		if err != nil {
			//return flags zero
			return res
		}
		dot := ro.(*objects.DOT)
		// sigok / structok
		if dot.SigValid() {
			res[0] = 1
		} else {
			return res
		}
		res[1] = byte(len(dot.GetRevokers()))
		expiry := dot.GetExpiry()
		if expiry != nil {
			copy(res[32:64], common.BigToBytes(big.NewInt(expiry.Unix()), 256))
		}
		copy(res[64:96], dot.GetGiverVK())
		copy(res[96:128], dot.GetReceiverVK())
		copy(res[128:160], dot.GetHash())
		return res
	case bytes.Equal(sig, SigGetDOTDelegatedRevoker):
		// GetDOTDelegatedRevoker(bytes,uint8) -> (bytes32)
		blob := getBytesParam(args, 0)
		idx := int(common.Bytes2Big(args[32:64]).Uint64())
		res := make([]byte, 32)
		if len(blob) < 96 {
			return res
		}
		ronum := objects.ROAccessDOT
		if blob[65] == 0x02 {
			ronum = objects.ROPermissionDOT
		}
		ro, err := objects.NewDOT(ronum, blob)
		if err != nil {
			//return zero as error
			return res
		}
		dot := ro.(*objects.DOT)
		if !dot.SigValid() {
			return res
		}
		if idx < 0 || len(dot.GetRevokers()) <= idx {
			return res
		}
		rvk := dot.GetRevokers()[idx]
		return rvk
	case bytes.Equal(sig, SigUnpackEntity):
		// UnpackEntity(bytes) -> (bytes32 flags, uint64 expiry, bytes32 vk)
		blob := getBytesParam(args, 0)
		res := make([]byte, 32*3)
		if len(blob) < 96 {
			fmt.Printf("XXX a %v\n", len(blob))
			return res
		}
		ro, err := objects.NewEntity(objects.ROEntity, blob)
		if err != nil {
			fmt.Printf("XXX b: %v\n", err)
			//return flags zero
			return res
		}
		e := ro.(*objects.Entity)
		// sigok / structok
		if e.SigValid() {
			res[0] = 1
		} else {
			fmt.Printf("XXX c\n")
			return res
		}
		res[1] = byte(len(e.GetRevokers()))
		expiry := e.GetExpiry()
		if expiry != nil {
			copy(res[32:64], common.BigToBytes(big.NewInt(expiry.Unix()), 256))
		}
		copy(res[64:96], e.GetVK())
		return res
	case bytes.Equal(sig, SigGetEntityDelegatedRevoker):
		// GetEntityDelegatedRevoker(bytes,uint8) -> (bytes32)
		blob := getBytesParam(args, 0)
		idx := int(common.Bytes2Big(args[32:64]).Uint64())
		res := make([]byte, 32)
		if len(blob) < 96 {
			return res
		}
		ro, err := objects.NewEntity(objects.ROEntity, blob)
		if err != nil {
			//return flags zero
			return res
		}
		e := ro.(*objects.Entity)
		// sigok / structok
		if !e.SigValid() {
			return res
		}
		if len(e.GetRevokers()) <= idx {
			return res
		}
		return e.GetRevokers()[idx]
	case bytes.Equal(sig, SigUnpackDChain):
		// UnpackDChain(bytes) -> (bytes32 flags, hash chainhash)
		// flags is [struct/sig valid], [numdots], 000000
		blob := getBytesParam(args, 0)
		rv := make([]byte, 32*2)
		if len(blob)%32 != 0 {
			return rv
		}
		rv[0] = 1
		rv[1] = byte(len(blob) / 32)
		chainhash := sha256.Sum256(blob)
		copy(rv[32:64], chainhash[:])
		return rv
	case bytes.Equal(sig, SigGetChainDOTHash):
		// GetChainDOTHash(bytes,uint8) -> (hash dothash)
		blob := getBytesParam(args, 0)
		idx := int(common.Bytes2Big(args[32:64]).Uint64())
		rv := make([]byte, 32)
		if len(blob)%32 != 0 {
			return rv
		}
		if idx < 0 || idx >= len(blob)/32 {
			return rv
		}
		return rv[idx*32 : (idx+1)*32]
	case bytes.Equal(sig, SigSliceByte32):
		// SliceByte32(bytes,uint64) -> (byte32)
		blob := getBytesParam(args, 0)
		idx := int(common.Bytes2Big(args[32:64]).Uint64())
		rv := make([]byte, 32)
		if idx < 0 || (idx+32) >= len(blob) {
			return rv
		}
		return rv[idx : idx+32]
	case bytes.Equal(sig, SigUnpackRevocation):
		// UnpackRevocation(bytes) -> (bytes32 flags, bytes32 target)
		blob := getBytesParam(args, 0)
		res := make([]byte, 32*2)
		ro, err := objects.NewRevocation(objects.RORevocation, blob)
		if err != nil {
			//return flags zero
			return res
		}
		rvk := ro.(*objects.Revocation)
		// sigok / structok
		if rvk.SigValid() {
			res[0] = 1
		} else {
			return res
		}
		copy(res[64:96], rvk.GetTarget())
		return res
	default:
		return nil
	}
}
func bwtest(in []byte) []byte {
	fmt.Printf("GOT: %d bytes: %v\n", len(in), hex.EncodeToString(in))
	return common.LeftPadBytes([]byte{1, 3, 3, 7}, 32)
}

func sha256Func(in []byte) []byte {
	return crypto.Sha256(in)
}

func ripemd160Func(in []byte) []byte {
	return common.LeftPadBytes(crypto.Ripemd160(in), 32)
}

const ecRecoverInputLength = 128

func ecrecoverFunc(in []byte) []byte {
	in = common.RightPadBytes(in, 128)
	// "in" is (hash, v, r, s), each 32 bytes
	// but for ecrecover we want (r, s, v)

	r := common.BytesToBig(in[64:96])
	s := common.BytesToBig(in[96:128])
	// Treat V as a 256bit integer
	vbig := common.Bytes2Big(in[32:64])
	v := byte(vbig.Uint64())

	// tighter sig s values in homestead only apply to tx sigs
	if !crypto.ValidateSignatureValues(v, r, s, false) {
		glog.V(logger.Debug).Infof("EC RECOVER FAIL: v, r or s value invalid")
		return nil
	}

	// v needs to be at the end and normalized for libsecp256k1
	vbignormal := new(big.Int).Sub(vbig, big.NewInt(27))
	vnormal := byte(vbignormal.Uint64())
	rsv := append(in[64:128], vnormal)
	pubKey, err := crypto.Ecrecover(in[:32], rsv)
	// make sure the public key is a valid one
	if err != nil {
		glog.V(logger.Error).Infof("EC RECOVER FAIL: ", err)
		return nil
	}

	// the first byte of pubkey is bitcoin heritage
	return common.LeftPadBytes(crypto.Sha3(pubKey[1:])[12:], 32)
}

func memCpy(in []byte) []byte {
	return in
}
