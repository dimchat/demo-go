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
	"math/rand"
)

type UserInfo struct {

	ID ID

	Meta Meta
	Visa Document

	IdentityKey SignKey
	CommunicationKey DecryptKey
}

type GroupInfo struct {

	ID ID

	Meta Meta
	Bulletin Document
}

/**
 *  Generate user account info
 *
 * @param name - nickname
 * @param avatar - photo URL
 * @return user info
 */
func GenerateUserInfo(name string, avatar string) *UserInfo {
	//
	//  Step 1. generate private key (with asymmetric algorithm)
	//
	identityKey := PrivateKeyGenerate(ECC)
	//
	//  Step 2. generate meta with private key
	//
	meta := MetaGenerate(ETH, identityKey, "")
	//
	//  Step 3. generate ID with meta
	//
	identifier := IDGenerate(meta, MAIN, "")
	//
	//  Step 4. generate visa document and sign with private key
	//
	communicationKey := PrivateKeyGenerate(RSA)
	visaKey := communicationKey.PublicKey().(EncryptKey)
	visa := DocumentCreate(VISA, identifier, nil, nil).(Visa)
	visa.SetName(name)
	visa.SetAvatar(avatar)
	visa.SetKey(visaKey)
	visa.Sign(identityKey)
	//
	//  OK
	//
	return &UserInfo{
		ID: identifier,
		Meta: meta,
		Visa: visa,
		IdentityKey: identityKey,
		CommunicationKey: communicationKey.(DecryptKey),
	}
}

/**
 *  Generate robot account info
 *
 * @param seed - meta seed for ID.name
 * @param name - robot name
 * @param avatar - photo URL
 * @return robot info
 */
func GenerateRobotInfo(seed string, name string, avatar string) *UserInfo {
	if seed == "" {
		seed = "robot"
	}
	//
	//  Step 1. generate private key (with asymmetric algorithm)
	//
	identityKey := PrivateKeyGenerate(RSA)
	//
	//  Step 2. generate meta with private key
	//
	meta := MetaGenerate(MKM, identityKey, seed)
	//
	//  Step 3. generate ID with meta
	//
	identifier := IDGenerate(meta, ROBOT, "")
	//
	//  Step 4. generate visa document and sign with private key
	//
	visa := DocumentCreate(VISA, identifier, nil, nil).(Visa)
	visa.SetName(name)
	visa.SetAvatar(avatar)
	visa.Sign(identityKey)
	//
	//  OK
	//
	return &UserInfo{
		ID: identifier,
		Meta: meta,
		Visa: visa,
		IdentityKey: identityKey,
		CommunicationKey: identityKey.(DecryptKey),
	}
}

/**
 *  Generate station account info
 *
 * @param seed - meta seed for ID.name
 * @param name - station name
 * @param logo - service provider logo URL
 * @param host - station IP
 * @param port - station port
 * @return station info
 */
func GenerateStationInfo(seed string, name string, logo string, host string, port uint16) *UserInfo {
	if seed == "" {
		seed = "station"
	}
	//
	//  Step 1. generate private key (with asymmetric algorithm)
	//
	identityKey := PrivateKeyGenerate(RSA)
	//
	//  Step 2. generate meta with private key
	//
	meta := MetaGenerate(MKM, identityKey, seed)
	//
	//  Step 3. generate ID with meta
	//
	identifier := IDGenerate(meta, STATION, "")
	//
	//  Step 4. generate visa document and sign with private key
	//
	profile := DocumentCreate(PROFILE, identifier, nil, nil)
	profile.SetName(name)
	profile.Set("logo", logo)
	profile.Set("host", host)
	profile.Set("port", port)
	profile.Sign(identityKey)
	//
	//  OK
	//
	return &UserInfo{
		ID: identifier,
		Meta: meta,
		Visa: profile,
		IdentityKey: identityKey,
		CommunicationKey: identityKey.(DecryptKey),
	}
}

/**
 *  Generate group account info
 *
 * @param founder - group founder
 * @param name - group name
 * @return Group object
 */
func GenerateGroupInfo(founder *UserInfo, name string, seed string) *GroupInfo {
	//
	//  Step 1. prepare seed for group ID
	//
	if seed == "" {
		r := rand.Int31n(999990000) + 10000  // 10,000 ~ 999,999,999
		seed = "Group-" + string(r)
	}
	//
	//  Step 2. generate meta with founder's private key
	//
	meta := MetaGenerate(MKM, founder.IdentityKey, seed)
	//
	//  Step 3. generate ID with meta
	//
	identifier := IDGenerate(meta, POLYLOGUE, "")
	//
	//  Step 4. generate bulletin document and sign with founder's private key
	//
	bulletin := DocumentCreate(BULLETIN, identifier, nil, nil)
	bulletin.SetName(name)
	bulletin.Sign(founder.IdentityKey)
	//
	//  OK
	//
	return &GroupInfo{
		ID: identifier,
		Meta: meta,
		Bulletin: bulletin,
	}
}
