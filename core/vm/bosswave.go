package vm

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	bwcrypto "github.com/immesys/bw2/crypto"
	"github.com/immesys/bw2/objects"
	bwutil "github.com/immesys/bw2/util"
	"github.com/immesys/bw2bc/common"
)

// The gas used by any bosswave function
const (
	BWGas = 3000
)

func getUInt64Param(args []byte, paramnum int) uint64 {
	return common.Bytes2Big(args[32*paramnum : (32*paramnum + 1)]).Uint64()
}
func getIntParam(args []byte, paramnum int) int {
	return int(common.Bytes2Big(args[32*paramnum : (32*paramnum + 1)]).Uint64())
}
func getBytes32Param(args []byte, paramnum int) []byte {
	return args[32*paramnum : (32*paramnum + 1)]
}

// This is a python script that can be used to work out the
// sha3 sig for a function
// import sha3; f = lambda x : sha3.sha3_256(x).hexdigest()[:8]

// VerifyEd25519Packed(bytes object)
// sig: VerifyEd25519Packed(bytes) (bool)
// returns true if valid, false otherwise
var SigVerifyEd25519Packed = []byte{0x70, 0xd6, 0x95, 0xf7}

func bwVerifyEd25519Packed(args []byte, vm *Vm) []byte {
	payload := getBytesParam(args, 0)
	if len(payload) < 96 {
		return nil
	}
	vk := payload[:32]
	bodyEnd := len(payload) - 64
	sig := payload[bodyEnd:]
	body := payload[:bodyEnd]
	fmt.Printf("Doing VFB:\nvk =%s\nsig=%s\nbdy=%s\n",
		hex.EncodeToString(vk), hex.EncodeToString(sig), hex.EncodeToString(body))
	if bwcrypto.VerifyBlob(vk, sig, body) {
		fmt.Println("BWVFBLOB WAS OK\n")
		return common.LeftPadBytes([]byte{1}, 32)
	} else {
		fmt.Println("BWVFBLOB WAS BAD\n")
		return common.LeftPadBytes([]byte{0}, 32)
	}
}

// VerifyEd25519(bytes32 vk, bytes sig, bytes body)
// sig: VerifyEd25519(bytes32,bytes,bytes) (bool)
var SigVerifyEd25519 = []byte{0x0b, 0x35, 0xfe, 0x44}

func bwVerifyEd25519(args []byte, vm *Vm) []byte {
	vk := args[0:32]
	sig := getBytesParam(args, 1)
	body := getBytesParam(args, 2)
	if len(sig) != 64 {
		return nil
	}
	if bwcrypto.VerifyBlob(vk, sig, body) {
		return common.LeftPadBytes([]byte{1}, 32)
	} else {
		return common.LeftPadBytes([]byte{0}, 32)
	}
}

// UnpackDOT(bytes dot)
// sig: UnpackDOT(bytes) (bool valid, uint8 numrevokers, bool ispermission,
//												uint64 expiry, bytes32 srcvk, bytes32 dstvk, bytes32 dothash)
var SigUnpackDOT = []byte{0x3a, 0xab, 0x47, 0xb2}

func bwUnpackDOT(args []byte, vm *Vm) []byte {
	// Bit of a hack, we read the DOT type from the actual DOT itself
	blob := getBytesParam(args, 0)
	res := make([]byte, 32*7)
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
		res[0*32] = 1
	} else {
		return res
	}
	res[1*32] = byte(len(dot.GetRevokers()))
	if ronum == objects.ROPermissionDOT {
		res[2*32] = 1
	}
	expiry := dot.GetExpiry()
	if expiry != nil {
		copy(res[3*32:4*32], common.BigToBytes(big.NewInt(expiry.Unix()), 256))
	}
	copy(res[4*32:5*32], dot.GetGiverVK())
	copy(res[5*32:6*32], dot.GetReceiverVK())
	copy(res[6*32:7*32], dot.GetHash())
	// We can now refer to the DOT by its hash
	vm.Scratch().InsertSlice(dot.GetHash(), dot)
	return res
}

// GetDOTDelegatedRevoker(bytes32 dothash, uint8 index)
// sig: GetDOTDelegatedRevoker(bytes32,uint8) (bytes32)
// The DOT must be in scratch
var SigGetDOTDelegatedRevoker = []byte{0xe0, 0x03, 0x1b, 0x1d}

func bwGetDOTDelegatedRevoker(args []byte, vm *Vm) []byte {
	dothash := getBytes32Param(args, 0)
	indx := getIntParam(args, 1)
	dot := vm.Scratch().LookupSlice(dothash).(*objects.DOT)
	return dot.GetRevokers()[indx]
}

// UnpackEntity(bytes entity)
// sig: UnpackEntity(bytes) (bool valid, uint8 numrevokers, uint64 expiry, bytes32 vk)
var SigUnpackEntity = []byte{0xe7, 0xb6, 0x86, 0xa7}

func bwUnpackEntity(args []byte, vm *Vm) []byte {
	// UnpackEntity(bytes) -> (bytes32 flags, uint64 expiry, bytes32 vk)
	blob := getBytesParam(args, 0)
	res := make([]byte, 32*4)
	ro, err := objects.NewEntity(objects.ROEntity, blob)
	if err != nil {
		//return flags zero
		return res
	}
	e := ro.(*objects.Entity)
	// sigok / structok
	if e.SigValid() {
		res[0*32] = 1
	} else {
		return res
	}
	res[1*32] = byte(len(e.GetRevokers()))
	expiry := e.GetExpiry()
	if expiry != nil {
		copy(res[2*32:3*32], common.BigToBytes(big.NewInt(expiry.Unix()), 256))
	}
	copy(res[3*32:4*32], e.GetVK())
	vm.Scratch().InsertSlice(e.GetVK(), e)
	return res
}

// GetEntityDelegatedRevoker(bytes32 vk, uint8 index)
// sig: GetEntityDelegatedRevoker(bytes32,index) (bytes32)
// Returns a delegated revoker for an entity.
// Entity must be in scratch
var SigGetEntityDelegatedRevoker = []byte{0x3a, 0xfe, 0x3a, 0x8a}

func bwGetEntityDelegatedRevoker(args []byte, vm *Vm) []byte {
	vk := getBytes32Param(args, 0)
	indx := getIntParam(args, 1)
	e := vm.Scratch().LookupSlice(vk).(*objects.Entity)
	return e.GetRevokers()[indx]
}

// UnpackAccessDChain(bytes obj)
// sig: UnpackAccessDChain(bytes) (bool valid, uint8 numdots, bytes32 chainhash)
// obj len must be a multiple of 32
// Also puts the dchain in scratch
var SigUnpackAccessDChain = []byte{0x22, 0xaf, 0x1b, 0x27}

func bwUnpackAccessDChain(args []byte, vm *Vm) []byte {
	// UnpackDChain(bytes) -> (bytes32 flags, hash chainhash)
	// flags is [structvalid], [numdots], 000000
	blob := getBytesParam(args, 0)
	dci, err := objects.LoadRoutingObject(objects.ROAccessDChain, blob)
	rv := make([]byte, 32*3)
	if err != nil {
		return rv
	}
	dc := dci.(*objects.DChain)
	rv[0*32] = 1
	rv[1*32] = byte(len(blob) / 32)
	chainhash := dc.GetChainHash()
	copy(rv[2*32:3*32], chainhash)
	// We might be augmenting chains, don't overwrite it if it is there
	if vm.Scratch().LookupSlice(chainhash) == nil {
		vm.Scratch().InsertSlice(chainhash, dc)
	}
	return rv
}

// GetDChainDOTHash(bytes32 chainhash, index) (bytes32 dothash)
// sig: GetDChainDOTHash(bytes32,uint8) (bytes32 dothash)
// chain must be in scratch
var SigGetDChainDOTHash = []byte{0xda, 0x3c, 0xd6, 0x74}

func bwGetDChainDOTHash(args []byte, vm *Vm) []byte {
	chainhash := getBytes32Param(args, 0)
	indx := getIntParam(args, 1)
	dc := vm.Scratch().LookupSlice(chainhash).(*objects.DChain)
	return dc.GetDotHash(indx)
}

// SliceByte32(bytes, offset) (bytes32)
// sig: SliceByte32(bytes,uint32) (bytes32)
var SigSliceByte32 = []byte{0xce, 0x7a, 0x94, 0xeb}

func bwSliceByte32(args []byte, vm *Vm) []byte {
	blob := getBytesParam(args, 0)
	idx := getIntParam(args, 1)
	return blob[idx : idx+32]
}

// UnpackRevocation(bytes) (bool valid, bytes32 vk, bytes32 target)
// sig: UnpackRevocation(bytes) (bool,bytes32,bytes32)
var SigUnpackRevocation = []byte{0xe5, 0x73, 0x1b, 0x77} //UnpackRevocation(bytes)
func bwUnpackRevocation(args []byte, vm *Vm) []byte {
	blob := getBytesParam(args, 0)
	res := make([]byte, 32*3)
	ro, err := objects.NewRevocation(objects.RORevocation, blob)
	if err != nil {
		//return flags zero
		return res
	}
	rvk := ro.(*objects.Revocation)
	// sigok / structok
	if rvk.SigValid() {
		res[0*32] = 1
	} else {
		return res
	}
	copy(res[1*32:2*32], rvk.GetTarget())
	copy(res[2*32:3*32], rvk.GetVK())
	key := make([]byte, 32)
	copy(key, rvk.GetTarget())
	key[0] = ^key[0]
	//Check if a slice of revocations for the target exists
	eslice, _ := vm.Scratch().LookupSlice(key).([]*objects.Revocation)
	if eslice == nil {
		eslice = make([]*objects.Revocation, 0, 1)
	}
	//Check if this revocation is in the slice. Compare on hash
	found := false
	for _, ervk := range eslice {
		if bytes.Equal(ervk.GetHash(), rvk.GetHash()) {
			found = true
			break
		}
	}
	if !found {
		eslice = append(eslice, rvk)
	}
	vm.Scratch().InsertSlice(key, eslice)
	return res
}

// ADChainGrants(bytes32 chainhash, bytes8 adps, bytes32 mvk, bytes urisuffix)
// sig: ADChainGrants(bytes32,bytes8,bytes32,bytes) (uint16)
// rv = 200 if chain is valid, and all dots are valid and unexpired and
//          it grants a superset of the passed adps, mvk and suffix
//          and all the entities are known to be unexpired
// rv = 201 same as above, but some entities were not present in Scratch
// rv else  a BWStatus code that something went wrong
var SigADChainGrants = []byte{0x8c, 0x75, 0x65, 0xdc}

func wrappedBWChainGrants(dc *objects.DChain, adpspacked []byte, mvk []byte, suffix []byte, vm *Vm) int {
	sSuffix := string(suffix)
	valid, _, _, _, _ := bwutil.AnalyzeSuffix(sSuffix)
	if !valid {
		return bwutil.BWStatusBadURI
	}
	ADPS := objects.DecodeADPS(adpspacked)
	now := time.Unix(vm.Env().Time().Int64(), 0)
	getDOT := func(k []byte) *objects.DOT {
		rv, ok := vm.Scratch().LookupSlice(k).(*objects.DOT)
		if !ok {
			return nil
		}
		return rv
	}
	getEntity := func(k []byte) *objects.Entity {
		rv, ok := vm.Scratch().LookupSlice(k).(*objects.Entity)
		if !ok {
			return nil
		}
		return rv
	}
	//We actually only store one revocation. We know its valid
	//so there is no point having more revocations
	getRevocation := func(k []byte) []*objects.Revocation {
		nk := make([]byte, len(k))
		copy(nk, k)
		nk[0] = ^nk[0]
		r, ok := vm.Scratch().LookupSlice(nk).(*objects.Revocation)
		if !ok {
			return []*objects.Revocation{}
		}
		return []*objects.Revocation{r}
	}
	// Down the rabbit hole
	return dc.CheckAccessGrants(&now, ADPS, mvk, sSuffix, getDOT,
		getEntity, getRevocation)
}
func bwADChainGrants(args []byte, vm *Vm) []byte {
	//remember to check all dots are access
	chainhash := getBytes32Param(args, 0)
	adpspacked := getBytes32Param(args, 1)
	mvk := getBytes32Param(args, 2)
	suffix := getBytesParam(args, 3)

	dc := vm.Scratch().LookupSlice(chainhash).(*objects.DChain)
	// Wow, such abstraction. This is like pages of code lol:
	result := wrappedBWChainGrants(dc, adpspacked, mvk, suffix, vm)
	return common.BigToBytes(big.NewInt(int64(result)), 256)
}

// GetDOTNumRevokableHashes(bytes32 dothash)
// sig: GetDOTNumRevokableHashes(bytes32) (uint32)
// Gets the total number of vulnerable hashes for the given dot
var SigGetDOTNumRevokableHashes = []byte{0x84, 0xea, 0x2e, 0x31}

func bwHelperDOTGetRevokableHashes(dot *objects.DOT) [][]byte {
	//DOT just has hash, src, dst
	rv := make([][]byte, 3)
	rv[0] = dot.GetHash()
	rv[1] = dot.GetGiverVK()
	rv[2] = dot.GetReceiverVK()
	return rv
}

func bwGetDOTNumRevokableHashes(args []byte, vm *Vm) []byte {
	dhash := getBytes32Param(args, 0)
	dot := vm.Scratch().LookupSlice(dhash).(*objects.DOT)
	nrh := len(bwHelperDOTGetRevokableHashes(dot))
	return common.BigToBytes(big.NewInt(int64(nrh)), 256)
}

// GetDOTRevokableHash(bytes32 dothash, uint32 index)
// sig: GetDOTRevokableHash(bytes32,uint32) (bytes32)
var SigGetDOTRevokableHash = []byte{0x24, 0xf6, 0x18, 0xb6}

func bwGetDOTRevokableHash(args []byte, vm *Vm) []byte {
	dhash := getBytes32Param(args, 0)
	indx := getIntParam(args, 1)
	dot := vm.Scratch().LookupSlice(dhash).(*objects.DOT)
	return bwHelperDOTGetRevokableHashes(dot)[indx]
}

// GetDChainNumRevokableHashes(bytes32 chainhash)
// sig: GetDChainNumRevokableHashes(bytes32) (uint32)
var SigGetDChainNumRevokableHashes = []byte{0x05, 0xd5, 0x6a, 0x4e}

func bwHelperDChainGetRevokableHashes(dc *objects.DChain) [][]byte {

	rv := make([][]byte, 2*dc.NumHashes()+1)
	rv[0] = dc.GetGiverVK()
	for i := 0; i < dc.NumHashes(); i++ {
		rv[i*2+1] = dc.GetDOT(i).GetHash()
		rv[i*2+2] = dc.GetDOT(i).GetReceiverVK()
	}
	return rv
}
func bwHelperAugmentDC(dc *objects.DChain, s *ScratchDatabase) {
	for i := 0; i < dc.NumHashes(); i++ {
		dh := dc.GetDotHash(i)
		dt := s.LookupSlice(dh).(*objects.DOT)
		dc.AugmentBy(dt)
	}
}
func bwGetDChainNumRevokableHashes(args []byte, vm *Vm) []byte {
	chainhash := getBytes32Param(args, 0)
	dc := vm.Scratch().LookupSlice(chainhash).(*objects.DChain)
	bwHelperAugmentDC(dc, vm.Scratch())
	return common.BigToBytes(big.NewInt(int64(len(bwHelperDChainGetRevokableHashes(dc)))), 256)
}

// GetDChainRevokableHash(bytes32 chainhash, uint32 index)
// sig: GetDChainRevokableHash(bytes32,uint32) (bytes32)
var SigGetDChainRevokableHash = []byte{0xee, 0xf9, 0x36, 0x11}

func bwGetDChainRevokableHash(args []byte, vm *Vm) []byte {
	chainhash := getBytes32Param(args, 0)
	indx := getIntParam(args, 1)
	dc := vm.Scratch().LookupSlice(chainhash).(*objects.DChain)
	bwHelperAugmentDC(dc, vm.Scratch())
	return bwHelperDChainGetRevokableHashes(dc)[indx]
}

func getBytesParam(in []byte, paramnumber int) []byte {
	offset := common.BytesToBig(in[(paramnumber * 32) : (paramnumber+1)*32]).Uint64()
	length := common.BytesToBig(in[offset : offset+32]).Uint64()
	return in[offset+32 : offset+32+length]
}

func bosswave(in []byte, vm *Vm) (rv []byte) {
	defer func() {
		if r := recover(); r != nil {
			rv = nil
			return
		}
	}()
	fmt.Printf("GOT: %d bytes: %v\n", len(in), hex.EncodeToString(in))
	if len(in) < 4 {
		return nil
	}
	sig := in[:4]
	args := in[4:]
	switch {
	//Low level
	case bytes.Equal(sig, SigVerifyEd25519):
		return bwVerifyEd25519(args, vm)
	case bytes.Equal(sig, SigVerifyEd25519Packed):
		return bwVerifyEd25519Packed(args, vm)
	case bytes.Equal(sig, SigSliceByte32):
		return bwSliceByte32(args, vm)
		//Entities
	case bytes.Equal(sig, SigUnpackEntity):
		return bwUnpackEntity(args, vm)
	case bytes.Equal(sig, SigGetEntityDelegatedRevoker):
		return bwGetEntityDelegatedRevoker(args, vm)
		//DOTs
	case bytes.Equal(sig, SigUnpackDOT):
		return bwUnpackDOT(args, vm)
	case bytes.Equal(sig, SigGetDOTDelegatedRevoker):
		return bwGetDOTDelegatedRevoker(args, vm)
	case bytes.Equal(sig, SigGetDOTNumRevokableHashes):
		return bwGetDOTNumRevokableHashes(args, vm)
	case bytes.Equal(sig, SigGetDOTRevokableHash):
		return bwGetDOTRevokableHash(args, vm)
	//Chains
	case bytes.Equal(sig, SigUnpackAccessDChain):
		return bwUnpackAccessDChain(args, vm)
	case bytes.Equal(sig, SigGetDChainDOTHash):
		return bwGetDChainDOTHash(args, vm)
	case bytes.Equal(sig, SigGetDChainNumRevokableHashes):
		return bwGetDChainNumRevokableHashes(args, vm)
	case bytes.Equal(sig, SigGetDChainRevokableHash):
		return bwGetDChainRevokableHash(args, vm)
	case bytes.Equal(sig, SigADChainGrants):
		return bwADChainGrants(args, vm)
	//Revocations
	case bytes.Equal(sig, SigUnpackRevocation):
		return bwUnpackRevocation(args, vm)

	default:
		fmt.Println("Hit default sig comparison: sig:", sig)
		return nil
	}
}
