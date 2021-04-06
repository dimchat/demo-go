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
	. "github.com/dimchat/demo-go/sdk/extensions"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
	"time"
)

type IFacebookExtension interface {

	SavePrivateKey(key PrivateKey, keyType string, user ID) bool

	SetCurrentUser(user User)
	AddUser(user User) bool
	RemoveUser(user User) bool

	AddContact(contact ID, user ID) bool
	RemoveContact(contact ID, user ID) bool

	AddMember(member ID, group ID) bool
	RemoveMember(member ID, group ID) bool
	ContainMember(member ID, group ID) bool
	ContainAssistant(bot ID, group ID) bool
	RemoveGroup(group ID) bool

	GetName(entity ID) string

	IsExpiredDocument(doc Document, reset bool) bool
}

type ICommonFacebook interface {
	IFacebook
	IFacebookExtension

	DB() IFacebookDatabase
	SetDB(db IFacebookDatabase)
}

const (
	EXPIRES_KEY = "expires"
	DOCUMENT_EXPIRES = 30 * 60
)

/**
 *  Common Facebook
 *  ~~~~~~~~~~~~~~~
 *  Barrack for Server/Client
 */
type CommonFacebook struct {
	Facebook
	IFacebookExtension

	_db IFacebookDatabase

	_users []User
}

func (facebook *CommonFacebook) Init() *CommonFacebook {
	if facebook.Facebook.Init() != nil {
		facebook._db = nil
		facebook._users = nil
	}
	return facebook
}

func (facebook *CommonFacebook) self() ICommonFacebook {
	return facebook.Facebook.Self().(ICommonFacebook)
}

func (facebook *CommonFacebook) SetDB(db IFacebookDatabase) {
	facebook._db = db
}
func (facebook *CommonFacebook) DB() IFacebookDatabase {
	return facebook._db
}

//-------- EntityCreator

func (facebook *CommonFacebook) isWaitingMeta(entity ID) bool {
	if entity.IsBroadcast() {
		return false
	}
	return facebook.self().GetMeta(entity) == nil
}

func (facebook *CommonFacebook) CreateUser(identifier ID) User {
	if facebook.isWaitingMeta(identifier) {
		return nil
	}
	return facebook.Facebook.CreateUser(identifier)
}

func (facebook *CommonFacebook) CreateGroup(identifier ID) Group {
	if facebook.isWaitingMeta(identifier) {
		return nil
	}
	return facebook.Facebook.CreateGroup(identifier)
}

func (facebook *CommonFacebook) GetLocalUsers() []User {
	if facebook._users == nil {
		facebook._users = make([]User, 0, 1)
		users := facebook.DB().AllUsers()
		if users != nil {
			self := facebook.self()
			var usr User
			for _, id := range users {
				usr = self.GetUser(id)
				if usr == nil {
					panic(id)
				} else {
					facebook._users = append(facebook._users, usr)
				}
			}
		}
	}
	return facebook._users
}

//-------- EntityManager

func (facebook *CommonFacebook) SaveMeta(meta Meta, identifier ID) bool {
	return facebook.DB().SaveMeta(meta, identifier)
}

func (facebook *CommonFacebook) SaveDocument(doc Document) bool {
	if facebook.self().CheckDocument(doc) == false {
		return false
	}
	doc.Set(EXPIRES_KEY, nil)
	return facebook.DB().SaveDocument(doc)
}

func (facebook *CommonFacebook) SaveMembers(members []ID, group ID) bool {
	return facebook.DB().SaveMembers(members, group)
}


//-------- IFacebookExtSource

func (facebook *CommonFacebook) SavePrivateKey(key PrivateKey, keyType string, user ID) bool {
	_, ok := key.(DecryptKey)
	return facebook.DB().SavePrivateKey(user, key, keyType, true, ok)
}

//-------- Local Users

func (facebook *CommonFacebook) SetCurrentUser(user User) {
	facebook._users = nil
	facebook.DB().SetCurrentUser(user.ID())
}

func (facebook *CommonFacebook) AddUser(user User) bool {
	facebook._users = nil
	return facebook.DB().AddUser(user.ID())
}

func (facebook *CommonFacebook) RemoveUser(user User) bool {
	facebook._users = nil
	return facebook.DB().RemoveUser(user.ID())
}

//-------- Contacts

func (facebook *CommonFacebook) AddContact(contact ID, user ID) bool {
	return facebook.DB().AddContact(contact, user)
}

func (facebook *CommonFacebook) RemoveContact(contact ID, user ID) bool {
	return facebook.DB().RemoveContact(contact, user)
}

//-------- Relationship

func (facebook *CommonFacebook) AddMember(member ID, group ID) bool {
	return facebook.DB().AddMember(member, group)
}
func (facebook *CommonFacebook) RemoveMember(member ID, group ID) bool {
	return facebook.DB().RemoveMember(member, group)
}
func (facebook *CommonFacebook) ContainMember(member ID, group ID) bool {
	members := facebook.self().GetMembers(group)
	if members != nil {
		for _, item := range members {
			if member.Equal(item) {
				return true
			}
		}
	}
	owner := facebook.self().GetOwner(group)
	return owner != nil && owner.Equal(members)
}
func (facebook *CommonFacebook) ContainAssistant(bot ID, group ID) bool {
	assistants := facebook.self().GetAssistants(group)
	if assistants != nil {
		for _, item := range assistants {
			if bot.Equal(item) {
				return true
			}
		}
	}
	return false
}
func (facebook *CommonFacebook) RemoveGroup(group ID) bool {
	return facebook.DB().RemoveGroup(group)
}

//-------- profiles

func (facebook *CommonFacebook) GetName(entity ID) string {
	// get name from document
	doc := facebook.self().GetDocument(entity, "*")
	if doc != nil {
		name := doc.Name()
		if name != "" {
			return name
		}
	}
	// get name fro ID
	return AnonymousGetName(entity)
}

func (facebook *CommonFacebook) IsExpiredDocument(doc Document, reset bool) bool {
	now := time.Now().Unix()
	expires := doc.Get(EXPIRES_KEY)
	if expires == nil {
		// set expired time
		doc.Set(EXPIRES_KEY, now + DOCUMENT_EXPIRES)
		return false
	}
	if now > expires.(int64) {
		if reset {
			// update expired time
			doc.Set(EXPIRES_KEY, now + DOCUMENT_EXPIRES)
		}
		return true
	} else {
		return false
	}
}
