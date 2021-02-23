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
	. "github.com/dimchat/core-go/core"
	. "github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/demo-go/libs/common/cpu"
	. "github.com/dimchat/demo-go/libs/common/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
	. "github.com/dimchat/sdk-go/protocol"
)

type IMessengerExtension interface {

	/**
	 *  Query meta from network with ID
	 */
	QueryMeta(identifier ID) bool

	/**
	 *  Query document from network with ID & type
	 */
	QueryDocument(identifier ID, docType string) bool

	/**
	 *  Query group info from members
	 */
	QueryGroupInfo(group ID, members []ID) bool
}

type ICommonMessenger interface {
	IMessenger
	IMessengerExtension
}

/**
 *  Common Messenger
 *  ~~~~~~~~~~~~~~~~
 */
type CommonMessenger struct {
	Messenger
	IMessengerExtension
}

func (messenger *CommonMessenger) Init() *CommonMessenger {
	if messenger.Messenger.Init() != nil {
	}
	return messenger
}

func (messenger *CommonMessenger) Facebook() ICommonFacebook {
	return messenger.Messenger.Facebook().(ICommonFacebook)
}

func (messenger *CommonMessenger) SetTransmitter(transmitter ICommonTransmitter) {
	messenger.Messenger.SetTransmitter(transmitter)
}
func (messenger *CommonMessenger) Transmitter() ICommonTransmitter {
	return messenger.Messenger.Transmitter().(ICommonTransmitter)
}

//-------- IMessengerExtension

func (messenger *CommonMessenger) QueryMeta(identifier ID) bool {
	return messenger.Transmitter().QueryMeta(identifier)
}

func (messenger *CommonMessenger) QueryDocument(identifier ID, docType string) bool {
	return messenger.Transmitter().QueryDocument(identifier, docType)
}

func (messenger *CommonMessenger) QueryGroupInfo(group ID, members []ID) bool {
	return messenger.Transmitter().QueryGroupInfo(group, members)
}

/**
 *  Register common parsers
 */
func RegisterCommonFactories() {
	// register command parsers
	CommandRegister(SEARCH, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(SearchCommand).Init(dict)
	}))
	CommandRegister(ONLINE_USERS, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(SearchCommand).Init(dict)
	}))

	CommandRegister(REPORT, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(ReportCommand).Init(dict)
	}))
	CommandRegister(ONLINE, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(ReportCommand).Init(dict)
	}))
	CommandRegister(OFFLINE, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(ReportCommand).Init(dict)
	}))

	// register content processors
	ContentProcessorRegister(0, new(AnyContentProcessor).Init())

	// register command processors
	CommandProcessorRegister(RECEIPT, new(ReceiptCommandProcessor).Init())
	CommandProcessorRegister(MUTE, new(MuteCommandProcessor).Init())
	CommandProcessorRegister(BLOCK, new(BlockCommandProcessor).Init())
}

func init() {
	RegisterCommonFactories()
}
