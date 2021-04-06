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
	"fmt"
	. "github.com/dimchat/demo-go/sdk/common/db"
	. "github.com/dimchat/demo-go/sdk/extensions"
	. "github.com/dimchat/demo-go/sdk/utils"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/plugins/crypto"
	. "github.com/dimchat/sdk-go/protocol"
)

type Database interface {

	PrivateKeyTable
	MetaTable
	DocumentTable

	AddressNameTable
	LoginTable

	UserTable
	ContactTable
	GroupTable

	// root directory for database
	SetRoot(root string)
}

/**
 *  Local Storage
 *  ~~~~~~~~~~~~~
 */
type Storage struct {
	Database

	_root string

	_password SymmetricKey

	//
	//  memory caches
	//

	_identityKeys map[ID]PrivateKey         // meta keys: ID -> SK
	_communicationKeys map[ID][]PrivateKey  // visa keys: ID -> []SK
	_decryptionKeys map[ID][]DecryptKey     // visa keys: ID -> []SK

	_metas map[ID]Meta                // meta: ID -> meta

	_docs map[string]map[ID]Document  // document: type -> ID -> doc

	_ans map[string]ID                // ANS: string -> ID

	_loginCommands map[ID]LoginCommand     // ID -> Login Command
	_loginMessages map[ID]ReliableMessage  // ID -> Login Message

	_users []ID
	_contacts map[ID][]ID             // user contacts: ID -> []ID

	_members map[ID][]ID              // group members: ID -> []ID
}

func (db *Storage) Init() *Storage {

	db._root = "/tmp/.dim"

	db._password = GetPlainKey()

	// private keys
	db._identityKeys = make(map[ID]PrivateKey)
	db._communicationKeys = make(map[ID][]PrivateKey)
	db._decryptionKeys = make(map[ID][]DecryptKey)

	// meta
	db._metas = make(map[ID]Meta)

	// documents
	docs := make(map[string]map[ID]Document)
	docs[VISA] = make(map[ID]Document)
	docs[PROFILE] = make(map[ID]Document)
	docs[BULLETIN] = make(map[ID]Document)
	db._docs = docs

	// ANS
	db._ans = loadANS(db)  // make(map[string]ID)

	// login info
	db._loginCommands = make(map[ID]LoginCommand)
	db._loginMessages = make(map[ID]ReliableMessage)

	// local users
	db._users = make([]ID, 0, 1)
	db._contacts = make(map[ID][]ID)

	// group info
	db._members = make(map[ID][]ID)

	return db
}

/**
 *  Password for private key encryption
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 */
func (db *Storage) Password() SymmetricKey {
	return db._password
}
func (db *Storage) SetPassword(password string) {
	db._password = GeneratePassword(password)
}

/**
 *  Root Directory
 *  ~~~~~~~~~~~~~~
 *
 *  File directory for database
 */
func (db *Storage) Root() string {
	return db._root
}
func (db *Storage) SetRoot(root string) {
	if PathIsExist(root) {
		db._root = root
	} else {
		panic(root)
	}
}

// Directory for MKM entity: '.dim/mkm/{zzz}/{ADDRESS}'
func (db *Storage) mkmDir(identifier ID) string {
	address := identifier.Address().String()
	pos := len(address)
	z := string(address[pos-1])
	y := string(address[pos-2])
	x := string(address[pos-3])
	w := string(address[pos-4])
	return PathJoin(db.Root(), "mkm", z, y, x, w, address)
}

func (db *Storage) prepareDir(filepath string) bool {
	dir := PathDir(filepath)
	return MakeDirs(dir)
}

//
//  DOS
//

//func (db *Storage) isExist(path string) bool {
//	return PathIsExist(path)
//}
//func (db *Storage) remove(path string) bool {
//	return PathRemove(path)
//}

func (db *Storage) readText(path string) string {
	return ReadTextFile(path)
}
func (db *Storage) readMap(path string) map[string]interface{} {
	info := ReadJSONFile(path)
	if info == nil {
		return nil
	} else {
		return info.(map[string]interface{})
	}
}

func (db *Storage) writeText(path string, text string) bool {
	if db.prepareDir(path) {
		return WriteTextFile(path, text)
	} else {
		panic(path)
	}
}
func (db *Storage) writeMap(path string, container interface{}) bool {
	if db.prepareDir(path) {
		return WriteJSONFile(path, container)
	} else {
		panic(path)
	}
}

//
//  Log
//

func (db *Storage) debug(msg string) {
	msg = fmt.Sprintf("Storage > %s", msg)
	LogDebug(msg)
}

func (db *Storage) log(msg string) {
	msg = fmt.Sprintf("Storage > %s", msg)
	LogInfo(msg)
}

func (db *Storage) warning(msg string) {
	msg = fmt.Sprintf("Storage > %s", msg)
	LogWarning(msg)
}

func (db *Storage) error(msg string) {
	msg = fmt.Sprintf("Storage > %s", msg)
	LogError(msg)
}

//
//  Singleton
//
var sharedDatabase Database

func SharedDatabase() Database {
	return sharedDatabase
}

func init() {
	sharedDatabase= new(Storage).Init()
}
