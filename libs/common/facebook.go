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
)

type IFacebookExtension interface {

	SavePrivateKey(key PrivateKey, keyType string, user User) bool

	SetCurrentUser(user User)
	AddUser(user User) bool
	RemoveUser(user User) bool

	AddContact(contact ID, user ID) bool
	RemoveContact(contact ID, user ID) bool

	AddMember(member ID, group ID) bool
	RemoveMember(member ID, group ID) bool
	SaveMembers(members []ID, group ID) bool
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
}

/**
 *  Common Facebook
 *  ~~~~~~~~~~~~~~~
 *  Barrack for Server/Client
 */
type CommonFacebook struct {
	Facebook
	IFacebookExtension

	_db IFacebookDatabase
}

func (facebook *CommonFacebook) Init() *CommonFacebook {
	if facebook.Facebook.Init() != nil {
		facebook._db = nil
	}
	return facebook
}

func (facebook *CommonFacebook) SetDB(db IFacebookDatabase) {
	facebook._db = db
}
func (facebook *CommonFacebook) DB() IFacebookDatabase {
	return facebook._db
}

func (facebook *CommonFacebook) SetHandler(handler ICommonFacebookHandler) {
	facebook.Facebook.SetHandler(handler)
	facebook.Facebook.SetManager(handler)
}
func (facebook *CommonFacebook) Handler() ICommonFacebookHandler {
	return facebook.Facebook.Handler().(ICommonFacebookHandler)
}

//-------- IFacebookExtSource

func (facebook *CommonFacebook) SavePrivateKey(key PrivateKey, keyType string, user User) bool {
	return facebook.Handler().SavePrivateKey(key, keyType, user)
}

func (facebook *CommonFacebook) SetCurrentUser(user User) {
	facebook.Handler().SetCurrentUser(user)
}
func (facebook *CommonFacebook) AddUser(user User) bool {
	return facebook.Handler().AddUser(user)
}
func (facebook *CommonFacebook) RemoveUser(user User) bool {
	return facebook.Handler().RemoveUser(user)
}

func (facebook *CommonFacebook) AddContact(contact ID, user ID) bool {
	return facebook.Handler().AddContact(contact, user)
}
func (facebook *CommonFacebook) RemoveContact(contact ID, user ID) bool {
	return facebook.Handler().RemoveContact(contact, user)
}

func (facebook *CommonFacebook) AddMember(member ID, group ID) bool {
	return facebook.Handler().AddMember(member, group)
}
func (facebook *CommonFacebook) RemoveMember(member ID, group ID) bool {
	return facebook.Handler().RemoveMember(member, group)
}
func (facebook *CommonFacebook) SaveMembers(members []ID, group ID) bool {
	return facebook.Handler().SaveMembers(members, group)
}
func (facebook *CommonFacebook) ContainMember(member ID, group ID) bool {
	return facebook.Handler().ContainMember(member, group)
}
func (facebook *CommonFacebook) ContainAssistant(bot ID, group ID) bool {
	return facebook.Handler().ContainAssistant(bot, group)
}
func (facebook *CommonFacebook) RemoveGroup(group ID) bool {
	return facebook.Handler().RemoveGroup(group)
}

func (facebook *CommonFacebook) GetName(entity ID) string {
	return facebook.Handler().GetName(entity)
}

func (facebook *CommonFacebook) IsExpiredDocument(doc Document, reset bool) bool {
	return facebook.Handler().IsExpiredDocument(doc, reset)
}
