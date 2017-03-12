// Copyright 2016 Michael Andersen
// This file is part of the bw2bc library.
//
// The bw2bc library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The bw2bc library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the bw2bc library. If not, see <http://www.gnu.org/licenses/>.

package vm

import "sync"

// The ScratchDatabase is used to store per-run data in the EVM
type ScratchDatabase struct {
	db   map[[32]byte]interface{}
	lock sync.Mutex
}

// NewScratchDatabase makes a new empty database
func NewScratchDatabase() *ScratchDatabase {
	return &ScratchDatabase{
		db: make(map[[32]byte]interface{}),
	}
}

// Lookup returns the stored value, or nil
func (sdb *ScratchDatabase) Lookup(key [32]byte) interface{} {
	sdb.lock.Lock()
	defer sdb.lock.Unlock()
	rv, ok := sdb.db[key]
	if !ok {
		return nil
	}
	return rv
}

// LookupSlice will convert the slice to a 32 byte array, truncating or
// zero extending as necessary, then call Lookup()
func (sdb *ScratchDatabase) LookupSlice(key []byte) interface{} {
	var arr [32]byte
	e := len(key)
	if e > 32 {
		e = 32
	}
	for i := 0; i < e; i++ {
		arr[i] = key[i]
	}
	return sdb.Lookup(arr)
}

// Clear will empty the database
func (sdb *ScratchDatabase) Clear() {
	sdb.lock.Lock()
	sdb.db = make(map[[32]byte]interface{})
	sdb.lock.Unlock()
}
func (sdb *ScratchDatabase) Insert(key [32]byte, val interface{}) {
	sdb.lock.Lock()
	sdb.db[key] = val
	sdb.lock.Unlock()
}
func (sdb *ScratchDatabase) InsertSlice(key []byte, val interface{}) {
	var arr [32]byte
	e := len(key)
	if e > 32 {
		e = 32
	}
	for i := 0; i < e; i++ {
		arr[i] = key[i]
	}
	sdb.Insert(arr, val)
}
