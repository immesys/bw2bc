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

const defaultGenesisBlock = "H4sICPiewVYAA2dlbmVzaXMuanNvbgCr5lKAAqW8/LzkVCUrBSWDipTUxJSkVCNjkxRTI1Pj5MQUJR2EupLM3NTiksTcAohaA2S5gsSi1LwSj8TiDKgkhQDZ7NSKkqJEl8SSRCz2picW+2TmZpZApCwwNadkpqVlJpfmlFRClJhgqMjNrMigibuT8zPzkhKLU0kzGdmExJyc/GSg9mq4UC1XLRcXAL/4Mm+8AQAA"
