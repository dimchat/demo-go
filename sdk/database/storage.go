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
	. "github.com/dimchat/demo-go/sdk/utils"
)

/**
 *  Local Storage
 *  ~~~~~~~~~~~~~
 */
type Storage struct {

	PrivateKeyTable
	MetaTable
	DocumentTable

	AddressNameTable
	LoginTable

	UserTable
	ContactTable
	GroupTable

	ProviderTable
	StationTable

	ConversationTable
	MessageTable

	MsgKeyTable

	_root string
}

func (db *Storage) Init() *Storage {
	db._root = "/tmp/.dim"
	return db
}

func (db *Storage) SetRoot(root string) {
	if PathIsExist(root) {
		db._root = root
	} else {
		panic(root)
	}
}

//
//  DOS
//

func (db *Storage) IsExist(path string) bool {
	return PathIsExist(path)
}
func (db *Storage) Remove(path string) bool {
	return PathRemove(path)
}

//
//  Log
//

func (db *Storage) Debug(msg string) {
	msg = fmt.Sprintf("Storage > %s", msg)
	LogDebug(msg)
}

func (db *Storage) Info(msg string) {
	msg = fmt.Sprintf("Storage > %s", msg)
	LogInfo(msg)
}

func (db *Storage) Warning(msg string) {
	msg = fmt.Sprintf("Storage > %s", msg)
	LogWarning(msg)
}

func (db *Storage) Error(msg string) {
	msg = fmt.Sprintf("Storage > %s", msg)
	LogError(msg)
}

//
//  Singleton
//
var sharedDatabase = new(Storage).Init()

func Database() *Storage {
	return sharedDatabase
}
