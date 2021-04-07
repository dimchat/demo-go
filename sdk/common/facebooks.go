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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

//-------- EntityDataSource

func (facebook *CommonFacebook) GetMeta(identifier ID) Meta {
	if identifier.IsBroadcast() {
		// broadcast ID has no meta
		return nil
	}
	return facebook.DB().GetMeta(identifier)
}

func (facebook *CommonFacebook) GetDocument(identifier ID, docType string) Document {
	return facebook.DB().GetDocument(identifier, docType)
}

//-------- UserDataSource

func (facebook *CommonFacebook) GetContacts(user ID) []ID {
	return facebook.DB().GetContacts(user)
}

func (facebook *CommonFacebook) GetPrivateKeysForDecryption(user ID) []DecryptKey {
	return facebook.DB().GetPrivateKeysForDecryption(user)
}

func (facebook *CommonFacebook) GetPrivateKeyForSignature(user ID) SignKey {
	return facebook.DB().GetPrivateKeyForSignature(user)
}

func (facebook *CommonFacebook) GetPrivateKeyForVisaSignature(user ID) SignKey {
	return facebook.DB().GetPrivateKeyForVisaSignature(user)
}

//-------- GroupDataSource

func (facebook *CommonFacebook) GetFounder(group ID) ID {
	founder := facebook.DB().GetFounder(group)
	if founder == nil {
		founder = facebook.Facebook.GetFounder(group)
	}
	return founder
}

func (facebook *CommonFacebook) GetOwner(group ID) ID {
	owner := facebook.DB().GetOwner(group)
	if owner == nil {
		owner = facebook.Facebook.GetOwner(group)
	}
	return owner
}

func (facebook *CommonFacebook) GetMembers(group ID) []ID {
	members := facebook.DB().GetMembers(group)
	if members == nil || len(members) == 0 {
		members = facebook.Facebook.GetMembers(group)
	}
	return members
}

func (facebook *CommonFacebook) GetAssistants(group ID) []ID {
	bots := facebook.DB().GetAssistants(group)
	if bots == nil || len(bots) == 0 {
		bots = facebook.Facebook.GetAssistants(group)
		if bots == nil || len(bots) == 0 {
			assistant := IDParse("assistant")
			if assistant != nil {
				bots = []ID{assistant}
			}
		}
	}
	return bots
}
