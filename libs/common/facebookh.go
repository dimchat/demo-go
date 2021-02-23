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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
	"time"
)

const (
	EXPIRES_KEY = "expires"
	DOCUMENT_EXPIRES = 30 * 60
)

type ICommonFacebookHandler interface {
	IFacebookHandler
	IFacebookExtension
}

/**
 *  Delegate for Facebook
 *  ~~~~~~~~~~~~~~~~~~~~~
 */
type CommonFacebookHandler struct {
	FacebookHandler

	_users []User
}

func (shadow *CommonFacebookHandler) Init(facebook ICommonFacebook) *CommonFacebookHandler {
	if shadow.FacebookHandler.Init(facebook) != nil {
		shadow._users = nil
	}
	return shadow
}

func (shadow *CommonFacebookHandler) Facebook() ICommonFacebook {
	return shadow.Barrack().(ICommonFacebook)
}

func (shadow *CommonFacebookHandler) DB() IFacebookDatabase {
	return shadow.Facebook().DB()
}

//-------- EntityHandler

func (shadow *CommonFacebookHandler) isWaitingMeta(entity ID) bool {
	if entity.IsBroadcast() {
		return false
	}
	return shadow.Facebook().GetMeta(entity) == nil
}

func (shadow *CommonFacebookHandler) CreateUser(identifier ID) User {
	if shadow.isWaitingMeta(identifier) {
		return nil
	}
	return shadow.FacebookHandler.CreateUser(identifier)
}

func (shadow *CommonFacebookHandler) CreateGroup(identifier ID) Group {
	if shadow.isWaitingMeta(identifier) {
		return nil
	}
	return shadow.FacebookHandler.CreateGroup(identifier)
}

func (shadow *CommonFacebookHandler) GetLocalUsers() []User {
	if shadow._users == nil {
		shadow._users = make([]User, 0)
		users := shadow.DB().AllUsers()
		if users != nil {
			barrack := shadow.Barrack()
			var usr User
			for _, id := range users {
				usr = barrack.GetUser(id)
				if usr == nil {
					panic(id)
				} else {
					shadow._users = append(shadow._users, usr)
				}
			}
		}
	}
	return shadow._users
}

//-------- EntityManager

func (shadow *CommonFacebookHandler) SaveMeta(meta Meta, identifier ID) bool {
	return shadow.DB().SaveMeta(meta, identifier)
}

func (shadow *CommonFacebookHandler) SaveDocument(doc Document) bool {
	if shadow.Facebook().CheckDocument(doc) == false {
		return false
	}
	doc.Set(EXPIRES_KEY, nil)
	return shadow.DB().SaveDocument(doc)
}

func (shadow *CommonFacebookHandler) SaveMembers(members []ID, group ID) bool {
	return shadow.DB().SaveMembers(members, group)
}

//-------- Private Key

func (shadow *CommonFacebookHandler) SavePrivateKey(key PrivateKey, keyType string, user User) bool {
	_, ok := key.(DecryptKey)
	return shadow.DB().SavePrivateKey(user.ID(), key, keyType, true, ok)
}

//-------- Local Users

func (shadow *CommonFacebookHandler) SetCurrentUser(user User) {
	shadow._users = nil
	shadow.DB().SetCurrentUser(user.ID())
}

func (shadow *CommonFacebookHandler) AddUser(user User) bool {
	shadow._users = nil
	return shadow.DB().AddUser(user.ID())
}

func (shadow *CommonFacebookHandler) RemoveUser(user User) bool {
	shadow._users = nil
	return shadow.DB().RemoveUser(user.ID())
}

//-------- Contacts

func (shadow *CommonFacebookHandler) AddContact(contact ID, user ID) bool {
	return shadow.DB().AddContact(contact, user)
}

func (shadow *CommonFacebookHandler) RemoveContact(contact ID, user ID) bool {
	return shadow.DB().RemoveContact(contact, user)
}

//-------- Relationship

func (shadow *CommonFacebookHandler) AddMember(member ID, group ID) bool {
	return shadow.DB().AddMember(member, group)
}
func (shadow *CommonFacebookHandler) RemoveMember(member ID, group ID) bool {
	return shadow.DB().RemoveMember(member, group)
}
func (shadow *CommonFacebookHandler) ContainMember(member ID, group ID) bool {
	members := shadow.Facebook().GetMembers(group)
	if members != nil {
		for _, item := range members {
			if member.Equal(item) {
				return true
			}
		}
	}
	owner := shadow.Facebook().GetOwner(group)
	return owner != nil && owner.Equal(members)
}
func (shadow *CommonFacebookHandler) ContainAssistant(bot ID, group ID) bool {
	assistants := shadow.Facebook().GetAssistants(group)
	if assistants != nil {
		for _, item := range assistants {
			if bot.Equal(item) {
				return true
			}
		}
	}
	return false
}
func (shadow *CommonFacebookHandler) RemoveGroup(group ID) bool {
	return shadow.DB().RemoveGroup(group)
}

//-------- profiles

func (shadow *CommonFacebookHandler) GetName(entity ID) string {
	// get name from document
	doc := shadow.Facebook().GetDocument(entity, "*")
	if doc != nil {
		name := doc.Name()
		if name != "" {
			return name
		}
	}
	// get name fro ID
	return AnonymousGetName(entity)
}

func (shadow *CommonFacebookHandler) IsExpiredDocument(doc Document, reset bool) bool {
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
