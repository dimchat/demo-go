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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/demo-go/sdk/common"
	. "github.com/dimchat/demo-go/sdk/utils"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp/protocol"
)

type ClientProcessor struct {
	CommonProcessor
}

func (processor *ClientProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	responses := processor.CommonProcessor.ProcessContent(content, rMsg)
	if responses == nil || len(responses) == 0 {
		// respond nothing
		return nil
	}
	if _, ok := responses[0].(HandshakeCommand); ok {
		// urgent command
		return responses
	}

	// check receiver
	receiver := rMsg.Receiver()
	user := processor.Facebook().SelectLocalUser(receiver)
	if user == nil {
		panic(receiver)
	}
	sender := rMsg.Sender()
	//messenger := processor.Messenger()
	// check responses
	for _, res := range responses {
		if res == nil {
			// should not happen
			continue
		} else if _, ok := res.(ReceiptCommand); ok {
			if sender.Type() == STATION {
				// no need to respond receipt to station
				LogInfo("drop receipt responding to station: " + sender.String())
				continue
			}
		} else if _, ok := res.(TextContent); ok {
			if sender.Type() == STATION {
				// no need to respond text message to station
				LogInfo("drop text msg responding to station: " + sender.String())
				continue
			}
		}
		// normal response
		//messenger.SendContent()
	}
	// DON'T respond to station directly
	return nil
}
