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
	. "github.com/dimchat/demo-go/sdk/client/cpu"
	. "github.com/dimchat/demo-go/sdk/common"
	. "github.com/dimchat/sdk-go/dimp"
	. "github.com/dimchat/sdk-go/dimp/cpu"
)

func createKeyCache() CipherKeyDelegate {
	cache := new(KeyCache)
	cache.Init()
	return cache
}
func createProcessor(facebook IClientFacebook, messenger IClientMessenger) Processor {
	// CPU creator
	creator := new(ClientProcessorCreator)
	creator.Init(facebook, messenger)
	// CPU factory
	factory := new(CPFactory)
	factory.Init(facebook, messenger)
	factory.SetCreator(creator)
	// message processor
	processor := new(ClientProcessor)
	processor.Init(facebook, messenger)
	processor.SetFactory(factory)
	return processor
}
func createPacker(facebook IClientFacebook, messenger IClientMessenger) Packer {
	packer := new(CommonPacker)
	packer.Init(facebook, messenger)
	return packer
}
//func createTransmitter(messenger IClientMessenger) ICommonTransmitter {
//	return new(CommonTransmitter).Init(messenger)
//}

type IClientMessenger interface {
	ICommonMessenger
}

type ClientMessenger struct {
	CommonMessenger
}

func (messenger *ClientMessenger) Init(facebook IClientFacebook) *ClientMessenger {
	if messenger.CommonMessenger.Init() != nil {
		// initialize delegates for Transceiver
		messenger.SetCipherKeyDelegate(createKeyCache())
		messenger.SetEntityDelegate(facebook)
		messenger.SetPacker(createPacker(facebook, messenger))
		messenger.SetProcessor(createProcessor(facebook, messenger))
		//messenger.SetTransformer(createTransformer(messenger))
		// initialize delegates for Messenger
		//messenger.SetTransmitter(createTransmitter(messenger))
	}
	return messenger
}

//
//  Singleton
//
var sharedMessenger *ClientMessenger

func SharedMessenger() IClientMessenger {
	return sharedMessenger
}

func init() {
	sharedMessenger = new(ClientMessenger)
	sharedMessenger.Init(SharedFacebook())
}
