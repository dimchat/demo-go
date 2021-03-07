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
	. "github.com/dimchat/demo-go/sdk/common"
)

func createKeyCache() CipherKeyDelegate {
	return new(KeyCache).Init()
}
func createTransformer(messenger IClientMessenger) Transformer {
	return new(CommonTransformer).Init(messenger)
}
func createProcessor(messenger IClientMessenger) Processor {
	return new(ClientProcessor).Init(messenger)
}
func createPacker(messenger IClientMessenger) Packer {
	return new(CommonPacker).Init(messenger)
}
func createTransmitter(messenger IClientMessenger) ICommonTransmitter {
	return new(CommonTransmitter).Init(messenger)
}

type IClientMessenger interface {
	ICommonMessenger
}

type ClientMessenger struct {
	CommonMessenger
}

func (messenger *ClientMessenger) Init() *ClientMessenger {
	if messenger.CommonMessenger.Init() != nil {
		// initialize delegates for Transceiver
		messenger.SetCipherKeyDelegate(createKeyCache())
		messenger.SetTransformer(createTransformer(messenger))
		messenger.SetProcessor(createProcessor(messenger))
		messenger.SetPacker(createPacker(messenger))
		// initialize delegates for Messenger
		messenger.SetTransmitter(createTransmitter(messenger))
	}
	return messenger
}

//
//  Singleton
//
var sharedMessenger IClientMessenger

func SharedMessenger() IClientMessenger {
	return sharedMessenger
}

func init() {
	sharedMessenger = new(ClientMessenger).Init()
}
