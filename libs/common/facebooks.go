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
	. "github.com/dimchat/core-go/core"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Data Source for Facebook
 *  ~~~~~~~~~~~~~~~~~~~~~
 */
type CommonFacebookDataSource struct {
	BarrackSource
}

func (shadow *CommonFacebookDataSource) Init(facebook ICommonFacebook) *CommonFacebookDataSource {
	if shadow.BarrackSource.Init(facebook) != nil {
	}
	return shadow
}

func (shadow *CommonFacebookDataSource) Facebook() ICommonFacebook {
	return shadow.Barrack().(ICommonFacebook)
}

func (shadow *CommonFacebookDataSource) DB() IFacebookDatabase {
	return shadow.Facebook().DB()
}

//-------- EntityDataSource

func (shadow *CommonFacebookDataSource) GetMeta(identifier ID) Meta {
	if identifier.IsBroadcast() {
		// broadcast ID has no meta
		return nil
	}
	return shadow.DB().GetMeta(identifier)
}

func (shadow *CommonFacebookDataSource) GetDocument(identifier ID, docType string) Document {
	return shadow.DB().GetDocument(identifier, docType)
}

//-------- UserDataSource

func (shadow *CommonFacebookDataSource) GetContacts(user ID) []ID {
	return shadow.DB().GetContacts(user)
}

func (shadow *CommonFacebookDataSource) GetPrivateKeysForDecryption(user ID) []DecryptKey {
	return shadow.DB().GetPrivateKeysForDecryption(user)
}

func (shadow *CommonFacebookDataSource) GetPrivateKeyForSignature(user ID) SignKey {
	return shadow.DB().GetPrivateKeyForSignature(user)
}

func (shadow *CommonFacebookDataSource) GetPrivateKeyForVisaSignature(user ID) SignKey {
	return shadow.DB().GetPrivateKeyForVisaSignature(user)
}

//-------- GroupDataSource

func (shadow *CommonFacebookDataSource) GetFounder(group ID) ID {
	founder := shadow.DB().GetFounder(group)
	if founder == nil {
		founder = shadow.BarrackSource.GetFounder(group)
	}
	return founder
}

func (shadow *CommonFacebookDataSource) GetOwner(group ID) ID {
	owner := shadow.DB().GetOwner(group)
	if owner == nil {
		owner = shadow.BarrackSource.GetOwner(group)
	}
	return owner
}

func (shadow *CommonFacebookDataSource) GetMembers(group ID) []ID {
	members := shadow.DB().GetMembers(group)
	if members == nil || len(members) == 0 {
		members = shadow.BarrackSource.GetMembers(group)
	}
	return members
}

func (shadow *CommonFacebookDataSource) GetAssistants(group ID) []ID {
	bots := shadow.BarrackSource.GetAssistants(group)
	if bots != nil && len(bots) > 0 {
		return bots
	}
	// TODO: try ANS record

	return shadow.DB().GetAssistants(group)
}
