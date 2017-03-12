// Copyright 2015 The go-ethereum Authors
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

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{
	// BOSSWAVE boot nodes
	//boota ipv4
	"enode://6ae73d0621c9c9a6bdac4a332900f1f57ea927f1a03aef5c2ffffa70fca0fada636da3ceac45ee4a2addbdb2bdbe9cb129b3a098d57fa09ff451712ac9c80fc9@54.215.189.111:30301",
	//boota ipv6
	"enode://6ae73d0621c9c9a6bdac4a332900f1f57ea927f1a03aef5c2ffffa70fca0fada636da3ceac45ee4a2addbdb2bdbe9cb129b3a098d57fa09ff451712ac9c80fc9@[2600:1f1c:c2f:a400:2f8f:3b34:1f55:3f7a]:30301",
	//bootb ipv4
	"enode://832c5a520a1079190e9fb57827306ee3882231077a3c543c8cae4c3a386703b3a4e0fd3ca9cb6b00b0d5482efc3e4dd8aafdb7fedb061d74a9d500f230e45873@54.183.54.213:30301",
	//bootb ipv6
	"enode://832c5a520a1079190e9fb57827306ee3882231077a3c543c8cae4c3a386703b3a4e0fd3ca9cb6b00b0d5482efc3e4dd8aafdb7fedb061d74a9d500f230e45873@[2600:1f1c:c2f:a400:5c38:c2f5:7e26:841c]:30301",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Morden test network.
var TestnetBootnodes = []string{
	"enode://e4533109cc9bd7604e4ff6c095f7a1d807e15b38e9bfeb05d3b7c423ba86af0a9e89abbf40bd9dde4250fef114cd09270fa4e224cbeef8b7bf05a51e8260d6b8@94.242.229.4:40404",
	"enode://8c336ee6f03e99613ad21274f269479bf4413fb294d697ef15ab897598afb931f56beb8e97af530aee20ce2bcba5776f4a312bc168545de4d43736992c814592@94.242.229.203:30303",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	//boota ipv4
	"enode://6ae73d0621c9c9a6bdac4a332900f1f57ea927f1a03aef5c2ffffa70fca0fada636da3ceac45ee4a2addbdb2bdbe9cb129b3a098d57fa09ff451712ac9c80fc9@54.215.189.111:30304",
	//boota ipv6
	"enode://6ae73d0621c9c9a6bdac4a332900f1f57ea927f1a03aef5c2ffffa70fca0fada636da3ceac45ee4a2addbdb2bdbe9cb129b3a098d57fa09ff451712ac9c80fc9@[2600:1f1c:c2f:a400:2f8f:3b34:1f55:3f7a]:30304",
	//bootb ipv4
	"enode://832c5a520a1079190e9fb57827306ee3882231077a3c543c8cae4c3a386703b3a4e0fd3ca9cb6b00b0d5482efc3e4dd8aafdb7fedb061d74a9d500f230e45873@54.183.54.213:30304",
	//bootb ipv6
	"enode://832c5a520a1079190e9fb57827306ee3882231077a3c543c8cae4c3a386703b3a4e0fd3ca9cb6b00b0d5482efc3e4dd8aafdb7fedb061d74a9d500f230e45873@[2600:1f1c:c2f:a400:5c38:c2f5:7e26:841c]:30304",
	// Asylum
	"enode://686f709677c4d0f2cd58cf651ea8ce1375bef22dcf29065994e34c1c4fd6f86691698321460f43059cc6cea536cd66ef534208869cd27765c4455f577a42a107@128.32.37.241:30304",
}
