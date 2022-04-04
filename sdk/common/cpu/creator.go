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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
	. "github.com/dimchat/sdk-go/dimp/cpu"
	. "github.com/dimchat/sdk-go/dimp/protocol"
)

/**
 *  CPU Creator
 *  ~~~~~~~~~~~
 *
 *  Delegate for CPU factory
 */
type CommonProcessorCreator struct {
	BaseCreator
}

//-------- IProcessorCreator

func (factory *CommonProcessorCreator) CreateContentProcessor(msgType ContentType) ContentProcessor {
	switch msgType {
	// file contents
	case FILE:
		return NewFileContentProcessor(factory.Facebook(), factory.Messenger())
	case IMAGE:
		return NewFileContentProcessor(factory.Facebook(), factory.Messenger())
	case AUDIO:
		return NewFileContentProcessor(factory.Facebook(), factory.Messenger())
	case VIDEO:
		return NewFileContentProcessor(factory.Facebook(), factory.Messenger())
	default:
	}
	// others
	cpu := factory.BaseCreator.CreateContentProcessor(msgType)
	if cpu != nil {
		return cpu
	}
	// unknown
	return NewBaseContentProcessor(factory.Facebook(), factory.Messenger())
}

func (factory *CommonProcessorCreator) CreateCommandProcessor(msgType ContentType, cmdName string) ContentProcessor {
	switch cmdName {
	// receipt command
	case RECEIPT:
		return NewReceiptCommandProcessor(factory.Facebook(), factory.Messenger())
	// mute command
	case MUTE:
		return NewMuteCommandProcessor(factory.Facebook(), factory.Messenger())
	// block command
	case BLOCK:
		return NewBlockCommandProcessor(factory.Facebook(), factory.Messenger())
	// storage
	default:
	}
	// others
	return factory.BaseCreator.CreateCommandProcessor(msgType, cmdName)
}
