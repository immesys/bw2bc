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
// Ropsten test network.
var TestnetBootnodes = []string{
	"enode://6ce05930c72abc632c58e2e4324f7c7ea478cec0ed4fa2528982cf34483094e9cbc9216e7aa349691242576d552a2a56aaeae426c5303ded677ce455ba1acd9d@13.84.180.240:30303", // US-TX
	"enode://20c9ad97c081d63397d7b685a412227a40e23c8bdc6688c6f37e97cfbc22d2b4d1db1510d8f61e6a8866ad7f0e17c02b14182d37ea7c3c8b9c2683aeb6b733a1@52.169.14.227:30303", // IE
}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{
	"enode://a24ac7c5484ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab56516d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303", // IE
}

// RinkebyV5Bootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network for the experimental RLPx v5 topic-discovery network.
var RinkebyV5Bootnodes = []string{
	"enode://a24ac7c5484ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab56516d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303?discport=30304", // IE
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
