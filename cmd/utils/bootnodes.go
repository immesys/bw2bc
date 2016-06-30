// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package utils

import "github.com/immesys/bw2bc/p2p/discover"

// FrontierBootNodes are the enode URLs of the P2P bootstrap nodes running on
// the Frontier network.
var FrontierBootNodes = []*discover.Node{
	// BOSSWAVE boot nodes
	// Castle
	discover.MustParseNode("enode://b2304f29230f9ceddb5e64e24ce5681f869d331a1dc41328eb4a7c26fedc92e24f34b87e775d7ce1793df376d63ae47ca00792ae7ecc01080aeebec14548e93b@128.32.37.201:30303"),
	// Asylum
	discover.MustParseNode("enode://686f709677c4d0f2cd58cf651ea8ce1375bef22dcf29065994e34c1c4fd6f86691698321460f43059cc6cea536cd66ef534208869cd27765c4455f577a42a107@128.32.37.241:30303"),
	// BW2.io
	discover.MustParseNode("enode://9cda6d7d65c465b92c413e6befd69e47588bb24806782a7bab0663de303e73f0cd2416e0d7c68cc9cba398d160ec34853252d382566bf6f706423b7bd7a712ef@54.183.221.12:30303"),
}

// TestNetBootNodes are the enode URLs of the P2P bootstrap nodes running on the
// Morden test network.
var TestNetBootNodes = []*discover.Node{}
