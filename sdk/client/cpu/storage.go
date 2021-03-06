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
package cpu

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/demo-go/sdk/extensions"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
	. "github.com/dimchat/sdk-go/protocol"
)

type StorageCommandProcessor struct {
	BaseCommandProcessor
}

func (cpu *StorageCommandProcessor) Init() *StorageCommandProcessor {
	if cpu.BaseCommandProcessor.Init() != nil {
	}
	return cpu
}

func (cpu *StorageCommandProcessor) decryptWithPassword(sCmd StorageCommand, pwd SymmetricKey) interface{} {
	// 1. get encrypted data
	data := sCmd.Data()
	if data == nil {
		// data not found
		panic(sCmd)
	}
	// 2. decrypt data
	data = pwd.Decrypt(data)
	if data == nil {
		// failed to decrypt data
		panic(sCmd)
	}
	// 3. decode data
	return JSONDecode(data)
}

func (cpu *StorageCommandProcessor) decryptData(sCmd StorageCommand) interface{} {
	// 1. get encrypt key
	key := sCmd.Key()
	if key == nil {
		// key not found
		panic(sCmd)
	}
	// 2. get user ID
	identifier := sCmd.ID()
	if identifier == nil {
		// ID not found
		panic(sCmd)
	}
	// 3. decrypt key
	user := cpu.Facebook().GetUser(identifier)
	key = user.Decrypt(key)
	if key == nil {
		// failed to decrypt key
		panic(sCmd)
	}
	// 4. decrypt key
	dict := JSONDecode(key)
	password := SymmetricKeyParse(dict)
	// 5. decrypt data
	return cpu.decryptWithPassword(sCmd, password)
}

//---- Contacts

func (cpu *StorageCommandProcessor) saveContacts(contacts []string, user ID) Content {
	// TODO: save contacts when import your account in a new app
	return nil
}

func (cpu *StorageCommandProcessor) processContacts(sCmd StorageCommand) Content {
	contacts, ok := sCmd.Get("contacts").([]string)
	if !ok || contacts == nil {
		contacts, _ = cpu.decryptData(sCmd).([]string)
		if contacts == nil {
			panic(sCmd)
		}
	}
	identifier := sCmd.ID()
	return cpu.saveContacts(contacts, identifier)
}

//---- Private Key

func (cpu *StorageCommandProcessor) savePrivateKey(key PrivateKey, user ID) Content {
	// TODO: save private key when import your accounts from network
	return nil
}

func (cpu *StorageCommandProcessor) processPrivateKey(sCmd StorageCommand) Content {
	string := "<TODO: input your password>"
	password := GeneratePassword(string)
	dict := cpu.decryptWithPassword(sCmd, password)
	key := PrivateKeyParse(dict)
	if key == nil {
		// failed to decrypt private key
		panic(sCmd)
	}
	identifier := sCmd.ID()
	return cpu.savePrivateKey(key, identifier)
}

func (cpu *StorageCommandProcessor) Execute(cmd Command, _ ReliableMessage) Content {
	sCmd, _ := cmd.(StorageCommand)
	title := sCmd.Title()
	if title == CONTACTS {
		return cpu.processContacts(sCmd)
	} else if title == PRIVATE_KEY {
		return cpu.processPrivateKey(sCmd)
	}
	panic(sCmd)
}
