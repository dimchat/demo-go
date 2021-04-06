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
package db

import (
	. "github.com/dimchat/demo-go/sdk/common/db"
	. "github.com/dimchat/demo-go/sdk/utils"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
)

//-------- PrivateKeyTable

func (db *Storage) SavePrivateKey(user ID, key PrivateKey, keyType string, sign bool, decrypt bool) bool {
	if keyType == META_KEY {
		if cacheIdentityKey(db, user, key) {
			return saveIdentityKey(db, user, key)
		} else {
			return false
		}
	} else {
		if cacheCommunicationKey(db, user, key) {
			keys := getCommunicationKeys(db, user)
			return saveCommunicationKeys(db, user, keys)
		} else {
			return false
		}
	}
}

func (db *Storage) GetPrivateKeysForDecryption(user ID) []DecryptKey {
	return getDecryptionKeys(db, user)
}

func (db *Storage) GetPrivateKeyForSignature(user ID) PrivateKey {
	keys := getCommunicationKeys(db, user)
	if len(keys) > 0 {
		// sign message with communication key
		return keys[0]
	} else {
		// if communication keys not exists, use identity key to sign message
		return getIdentityKey(db, user)
	}
}

func (db *Storage) GetPrivateKeyForVisaSignature(user ID) PrivateKey {
	return getIdentityKey(db, user)
}

/**
 *  Private Key file for Local Users
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 *  1. Identity Key      - paired to meta.key, CONSTANT
 *     file path: '.dim/private/{ADDRESS}/secret'
 *
 *  2. Communication Key - paired to visa.key, VOLATILE
 *     file path: '.dim/private/{ADDRESS}/secret_keys'
 */

func identityKeyPath(db *Storage, identifier ID) string {
	return PathJoin(db.Root(), "private", identifier.Address().String(), "secret")
}

func communicationKeysPath(db *Storage, identifier ID) string {
	return PathJoin(db.Root(), "private", identifier.Address().String(), "secret_keys")
}

func loadIdentityKey(db *Storage, identifier ID) PrivateKey {
	path := identityKeyPath(db, identifier)
	db.log("Loading identity key: " + path)
	data := db.readSecret(path)
	if data == nil {
		return nil
	} else {
		return PrivateKeyParse(JSONDecodeMap(data))
	}
}
func loadCommunicationKeys(db *Storage, identifier ID) []PrivateKey {
	keys := make([]PrivateKey, 0, 1)
	path := communicationKeysPath(db, identifier)
	db.log("Loading communication keys: " + path)
	data := db.readSecret(path)
	if data != nil {
		arr := JSONDecodeList(data)
		for _, item := range arr {
			k := PrivateKeyParse(item)
			if k == nil {
				panic(item)
			} else {
				keys = append(keys, k)
			}
		}
	}
	return keys
}

func saveIdentityKey(db *Storage, identifier ID, key PrivateKey) bool {
	info := key.GetMap(false)
	path := identityKeyPath(db, identifier)
	db.log("Saving identity key: " + path)
	return db.writeSecret(path, JSONEncodeMap(info))
}
func saveCommunicationKeys(db *Storage, identifier ID, keys []PrivateKey) bool {
	arr := make([]interface{}, 0, len(keys))
	for _, item := range keys {
		arr = append(arr, item.GetMap(false))
	}
	path := communicationKeysPath(db, identifier)
	db.log("Saving communication keys: " + path)
	return db.writeSecret(path, JSONEncodeList(arr))
}

// place holder
var emptyPrivateKey = PrivateKeyGenerate(ECC)

func getIdentityKey(db *Storage, identifier ID) PrivateKey {
	// 1. try from memory cache
	key := db._identityKeys[identifier]
	if key == nil {
		// 2. try from local storage
		key = loadIdentityKey(db, identifier)
		if key == nil {
			// place an empty key for cache
			db._identityKeys[identifier] = emptyPrivateKey
		} else {
			// cache it
			db._identityKeys[identifier] = key
		}
	} else if key == emptyPrivateKey {
		db.error("Private key not found: " + identifier.String())
		key = nil
	}
	return key
}

func getCommunicationKeys(db *Storage, identifier ID) []PrivateKey {
	// 1. try from memory cache
	keys := db._communicationKeys[identifier]
	if keys == nil {
		// 2. try from local storage
		keys = loadCommunicationKeys(db, identifier)
		// 3. cache them
		db._communicationKeys[identifier] = keys
	}
	return keys
}
func getDecryptionKeys(db *Storage, identifier ID) []DecryptKey {
	// 1. try from memory cache
	keys := db._decryptionKeys[identifier]
	if keys == nil || len(keys) == 0 {
		var decKey DecryptKey
		var ok bool
		// 2. get communication keys
		msgKeys := getCommunicationKeys(db, identifier)
		keys = make([]DecryptKey, 0, len(msgKeys) + 1)
		for _, item := range msgKeys {
			decKey, ok = item.(DecryptKey)
			if ok && decKey != nil {
				keys = append(keys, decKey)
			}
		}
		// 3. check identity key
		idKey := getIdentityKey(db, identifier)
		decKey, ok = idKey.(DecryptKey)
		if ok && decKey != nil && findKey(msgKeys, idKey) < 0 {
			keys = append(keys, decKey)
		}
		// 4. cache them
		db._decryptionKeys[identifier] = keys
	}
	return keys
}

func cacheIdentityKey(db *Storage, identifier ID, key PrivateKey) bool {
	old := getIdentityKey(db, identifier)
	if old == nil {
		db._identityKeys[identifier] = key
		return true
	} else {
		// identity key won't change
		return false
	}
}

func cacheCommunicationKey(db *Storage, identifier ID, key PrivateKey) bool {
	keys := getCommunicationKeys(db, identifier)
	index := findKey(keys, key)
	if index == 0 {
		return false                   // nothing changed
	} else if index > 0 {
		keys = removeKey(keys, index)  // move to the front
	} else if len(keys) > 2 {
		keys = keys[:2]                // keep only last three records
	}
	keys = insertKey(keys, key)
	db._communicationKeys[identifier] = keys
	// reset decryption keys
	delete(db._decryptionKeys, identifier)
	return true
}

func findKey(keys []PrivateKey, key PrivateKey) int {
	for index, item := range keys {
		if key.Equal(item) {
			return index
		}
	}
	return -1
}
func removeKey(keys []PrivateKey, index int) []PrivateKey {
	return append(keys[:index], keys[index+1:]...)
}
func insertKey(keys []PrivateKey, key PrivateKey) []PrivateKey {
	arr := make([]PrivateKey, 0, len(keys) + 1)
	arr = append(arr, key)
	return append(arr, keys...)
}
