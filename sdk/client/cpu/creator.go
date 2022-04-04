/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2022 Albert Moky
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
	. "github.com/dimchat/demo-go/sdk/common/cpu"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
)

/**
 *  CPU Creator
 *  ~~~~~~~~~~~
 *
 *  Delegate for CPU factory
 */
type ClientProcessorCreator struct {
	CommonProcessorCreator
}

//-------- IProcessorCreator

func (factory *ClientProcessorCreator) CreateContentProcessor(msgType ContentType) ContentProcessor {
	// text

	// others
	return factory.CommonProcessorCreator.CreateContentProcessor(msgType)
}

func (factory *ClientProcessorCreator) CreateCommandProcessor(msgType ContentType, cmdName string) ContentProcessor {
	// handshake
	if cmdName == HANDSHAKE {
		return NewHandshakeCommandProcessor(factory.Facebook(), factory.Messenger())
	}

	// login

	// others
	return factory.CommonProcessorCreator.CreateCommandProcessor(msgType, cmdName)
}
