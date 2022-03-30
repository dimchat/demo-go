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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/digest"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/sdk-go/dimp"
	. "github.com/dimchat/sdk-go/plugins/crypto"
)

/**
 *  Common Packer
 *  ~~~~~~~~~~~~~
 */
type CommonPacker struct {
	MessagePacker
}

func (packer *CommonPacker) Init(facebook IFacebook, messenger IMessenger) *CommonPacker {
	if packer.MessagePacker.Init(facebook, messenger) != nil {
	}
	return packer
}

func (packer *CommonPacker) CipherKeyDelegate() CipherKeyDelegate {
	return packer.Messenger().CipherKeyDelegate()
}

func (packer *CommonPacker) attachKeyDigest(rMsg ReliableMessage) {
	if rMsg.Delegate() == nil {
		rMsg.SetDelegate(packer.Messenger())
	}
	if rMsg.EncryptedKey() != nil {
		// 'key' exists
		return
	}
	keys := rMsg.EncryptedKeys()
	if keys == nil {
		keys = make(map[string]string)
	} else if keys["digest"] != "" {
		// key digest already exists
		return
	}
	// get key with direction
	var key SymmetricKey
	sender := rMsg.Sender()
	group := rMsg.Group()
	if group == nil {
		receiver := rMsg.Receiver()
		key = packer.CipherKeyDelegate().GetCipherKey(sender, receiver, false)
	} else {
		key = packer.CipherKeyDelegate().GetCipherKey(sender, group, false)
	}
	if key == nil {
		// broadcast message has no key
		return
	}
	// get key data
	data := key.Data()
	if data == nil || len(data) < 6 {
		if key.Algorithm() == PLAIN {
			// broadcast message has no key
			return
		}
		panic(key)
	}
	// get digest
	part := data[len(data)-6:]
	b64 := Base64Encode(SHA256(part))
	digest := b64[len(b64)-8:]
	keys["digest"] = digest
	rMsg.Set("keys", keys)
}

func (packer *CommonPacker) EncryptMessage(iMsg InstantMessage) SecureMessage {
	sMsg := packer.MessagePacker.EncryptMessage(iMsg)

	receiver := iMsg.Receiver()
	if receiver.IsGroup() {
		// reuse group message keys
		sender := iMsg.Sender()
		key := packer.CipherKeyDelegate().GetCipherKey(sender, receiver, false)
		key.Set("reused", true)
	}
	// TODO: reuse personal message key?

	return sMsg
}
