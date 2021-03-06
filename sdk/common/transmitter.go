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
	. "github.com/dimchat/sdk-go/dimp"
)

type ICommonTransmitter interface {
	Transmitter
	IMessengerExtension
}

/**
 *  Common Processor
 *  ~~~~~~~~~~~~~~~~
 *
 *  Abstract Methods:
 *      // IMessengerExtension
 *      QueryMeta(identifier ID) bool
 *      QueryDocument(identifier ID, docType string) bool
 *      QueryGroupInfo(group ID, members []ID) bool
 */
type CommonTransmitter struct {
	MessengerTransmitter
	IMessengerExtension
}

func (transmitter *CommonTransmitter) Init(transceiver ICommonMessenger) *CommonTransmitter {
	if transmitter.MessengerTransmitter.Init(transceiver) != nil {
	}
	return transmitter
}

func (transmitter *CommonTransmitter) SendInstantMessage(iMsg InstantMessage, callback MessengerCallback, priority int) bool {
	go func() {
		messenger := transmitter.Messenger()
		sMsg := messenger.EncryptMessage(iMsg)
		if sMsg == nil {
			// public key not found?
			//panic(iMsg)
			return
		}
		rMsg := messenger.SignMessage(sMsg)
		if rMsg == nil {
			// TODO: set iMsg.state = error
			panic(sMsg)
		}
		transmitter.SendReliableMessage(rMsg, callback, priority)
		// TODO: if OK, set iMsg.state = sending; else set iMsg.state = waiting

		// save signature for receipt
		iMsg.Set("signature", rMsg.Get("signature"))

		messenger.SaveMessage(iMsg)
	}()
	return true
}

//-------- IMessengerExtension
