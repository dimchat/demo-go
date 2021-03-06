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
	. "github.com/dimchat/demo-go/sdk/common"
	. "github.com/dimchat/demo-go/sdk/utils"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/protocol"
	"time"
)

type ClientProcessor struct {
	CommonProcessor
}

func (processor *ClientProcessor) Init(transceiver ICommonMessenger) *ClientProcessor {
	if processor.CommonProcessor.Init(transceiver) != nil {
	}
	return processor
}

func (processor *ClientProcessor) ProcessContent(content Content, rMsg ReliableMessage) Content {
	res := processor.CommonProcessor.ProcessContent(content, rMsg)
	if res == nil {
		// respond nothing
		return nil
	}
	if _, ok := res.(HandshakeCommand); ok {
		// urgent command
		return res
	}

	sender := rMsg.Sender()
	if _, ok := res.(ReceiptCommand); ok {
		if sender.Type() == STATION {
			// no need to respond receipt to station
			return nil
		}
		LogInfo("receipt to sender: " + sender.String())
	}

	// check receiver
	receiver := rMsg.Receiver()
	user := processor.Facebook().SelectLocalUser(receiver)
	if user == nil {
		panic(receiver)
	}
	// pack message
	env := EnvelopeCreate(user.ID(), sender, time.Now())
	iMsg := InstantMessageCreate(env, res)
	// normal response
	processor.Messenger().SendInstantMessage(iMsg, nil, 1)
	// DON'T respond to station directly
	return nil
}
