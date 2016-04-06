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

package core

import (
	"compress/gzip"
	"encoding/base64"
	"io"
	"strings"
)

func NewDefaultGenesisReader() (io.Reader, error) {
	return gzip.NewReader(base64.NewDecoder(base64.StdEncoding, strings.NewReader(defaultGenesisBlock)))
}

// {
//     "nonce": "0xb055deadbeefdeadbeef",
//     "timestamp": "0x0",
//     "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
//     "extraData": "0x0",
//     "gasLimit": "0x8000000",
//     "difficulty": "200000",
//     "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
//     "coinbase": "0x3333333333333333333333333333333333333333",
//     "alloc": {
//     }
// }

const defaultGenesisBlock = "H4sICKu9+VYAA2Zvb2dlbi5qc29uAKvm4lTKy89LTlWyUlAyqEgyMDVNSU1MSUpNTYPRSjpANSWZuanFJYm5BRB1BmDBgsSi1LwSj8TiDKgohQBsaGpFSVGiS2JJIrJN6YnFPpm5mSUQMQsk5SmZaWmZyaU5JZUgOSOERG5mRQZ1XZacn5mXlFgMDSpjIgFYa2JOTn4yUF81F2ctVy0XAF+q+gd2AQAA"
