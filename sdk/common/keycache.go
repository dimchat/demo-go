/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package dimp

import (
	. "github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/demo-go/sdk/common/db"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/plugins/crypto"
)

type KeyCache struct {
	CipherKeyDelegate

	_keyMap map[ID]map[ID]SymmetricKey

	_database MsgKeyTable
}

func (cache *KeyCache) Init() *KeyCache {
	cache._keyMap = make(map[ID]map[ID]SymmetricKey)
	cache._database = nil
	return cache
}

func (cache *KeyCache) SetKeyTable(table MsgKeyTable) {
	cache._database = table
}

func (cache *KeyCache) GetCipherKey(sender, receiver ID, generate bool) SymmetricKey {
	if receiver.IsBroadcast() {
		// broadcast message has no key
		return GetPlainKey()
	}
	var key SymmetricKey
	// try from memory cache
	table := cache._keyMap[sender]
	if table == nil {
		table = make(map[ID]SymmetricKey)
		cache._keyMap[sender] = table
	} else {
		key = table[receiver]
		if key != nil {
			return key
		}
	}
	// try from database
	key = cache._database.GetKey(sender, receiver)
	if key != nil {
		// cache it
		table[receiver] = key
	} else if generate {
		// generate new key and store it
		key = SymmetricKeyGenerate(AES)
		cache._database.SaveKey(sender, receiver, key)
		// cache it
		table[receiver] = key
	}
	return key
}

func (cache *KeyCache) CacheCipherKey(sender, receiver ID, key SymmetricKey) {
	if receiver.IsBroadcast() {
		// broadcast message has no key
		return
	}
	// save into database
	if cache._database.SaveKey(sender, receiver, key) {
		// cache it
		table := cache._keyMap[sender]
		if table == nil {
			table = make(map[ID]SymmetricKey)
			cache._keyMap[sender] = table
		}
		table[receiver] = key
	}
}
