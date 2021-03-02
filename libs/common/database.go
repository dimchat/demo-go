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
	. "github.com/dimchat/demo-go/libs/common/db"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

type IFacebookDatabase interface {
	PrivateKeyTable
	MetaTable
	DocumentTable

	UserTable
	GroupTable
	ContactTable
}

type FacebookDatabase struct {
	IFacebookDatabase

	// databases
	_privateTable PrivateKeyTable
	_metaTable MetaTable
	_docTable DocumentTable

	_userTable UserTable
	_groupTable GroupTable
	_contactTable ContactTable
}

func (db *FacebookDatabase) Init() *FacebookDatabase {
	db._privateTable = nil
	db._metaTable = nil
	db._docTable = nil
	db._userTable = nil
	db._groupTable = nil
	db._contactTable = nil
	return db
}

func (db *FacebookDatabase) SetPrivateKeyTable(table PrivateKeyTable) {
	db._privateTable = table
}
func (db *FacebookDatabase) SetMetaTable(table MetaTable) {
	db._metaTable = table
}
func (db *FacebookDatabase) SetDocumentTable(table DocumentTable) {
	db._docTable = table
}

func (db *FacebookDatabase) SetUserTable(table UserTable) {
	db._userTable = table
}
func (db *FacebookDatabase) SetGroupTable(table GroupTable) {
	db._groupTable = table
}
func (db *FacebookDatabase) SetContactTable(table ContactTable) {
	db._contactTable = table
}

//-------- PrivateKeyTable

func (db *FacebookDatabase) SavePrivateKey(user ID, key PrivateKey, keyType string, sign bool, decrypt bool) bool {
	return db._privateTable.SavePrivateKey(user, key, keyType, sign, decrypt)
}

func (db *FacebookDatabase) GetPrivateKeysForDecryption(user ID) []DecryptKey {
	return db._privateTable.GetPrivateKeysForDecryption(user)
}

func (db *FacebookDatabase) GetPrivateKeyForSignature(user ID) PrivateKey {
	return db._privateTable.GetPrivateKeyForSignature(user)
}

func (db *FacebookDatabase) GetPrivateKeyForVisaSignature(user ID) PrivateKey {
	return db._privateTable.GetPrivateKeyForVisaSignature(user)
}

//-------- MetaTable

func (db *FacebookDatabase) SaveMeta(meta Meta, entity ID) bool {
	return db._metaTable.SaveMeta(meta, entity)
}

func (db *FacebookDatabase) GetMeta(entity ID) Meta {
	return db._metaTable.GetMeta(entity)
}

//-------- DocumentTable

func (db *FacebookDatabase) SaveDocument(doc Document) bool {
	return db._docTable.SaveDocument(doc)
}

func (db *FacebookDatabase) GetDocument(entity ID, docType string) Document {
	return db._docTable.GetDocument(entity, docType)
}

//-------- UserTable

func (db *FacebookDatabase) AllUsers() []ID {
	return db._userTable.AllUsers()
}

func (db *FacebookDatabase) AddUser(user ID) bool {
	return db._userTable.AddUser(user)
}

func (db *FacebookDatabase) RemoveUser(user ID) bool {
	return db._userTable.RemoveUser(user)
}

func (db *FacebookDatabase) SetCurrentUser(user ID) {
	db._userTable.SetCurrentUser(user)
}

func (db *FacebookDatabase) GetCurrentUser() ID {
	return db._userTable.GetCurrentUser()
}

//-------- GroupTable

func (db *FacebookDatabase) GetFounder(group ID) ID {
	return db._groupTable.GetFounder(group)
}

func (db *FacebookDatabase) GetOwner(group ID) ID {
	return db._groupTable.GetFounder(group)
}

func (db *FacebookDatabase) GetMembers(group ID) []ID {
	return db._groupTable.GetMembers(group)
}

func (db *FacebookDatabase) GetAssistants(group ID) []ID {
	return db._groupTable.GetAssistants(group)
}

func (db *FacebookDatabase) AddMember(member ID, group ID) bool {
	return db._groupTable.AddMember(member, group)
}

func (db *FacebookDatabase) RemoveMember(member ID, group ID) bool {
	return db._groupTable.RemoveMember(member, group)
}

func (db *FacebookDatabase) SaveMembers(members []ID, group ID) bool {
	return db._groupTable.SaveMembers(members, group)
}

func (db *FacebookDatabase) RemoveGroup(group ID) bool {
	return db._groupTable.RemoveGroup(group)
}

//-------- ContactTable

func (db *FacebookDatabase) GetContacts(user ID) []ID {
	return db._contactTable.GetContacts(user)
}

func (db *FacebookDatabase) AddContact(contact ID, user ID) bool {
	return db._contactTable.AddContact(contact, user)
}

func (db *FacebookDatabase) RemoveContact(contact ID, user ID) bool {
	return db._contactTable.RemoveContact(contact, user)
}

func (db *FacebookDatabase) SaveContacts(contacts []ID, user ID) bool {
	return db._contactTable.SaveContacts(contacts, user)
}
