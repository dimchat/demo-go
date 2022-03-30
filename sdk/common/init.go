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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/demo-go/sdk/common/protocol"
	_ "github.com/dimchat/sdk-go/dimp/cpu"
	_ "github.com/dimchat/sdk-go/plugins"
)

/**
 *  Register common parsers
 */
func RegisterCommonFactories() {
	// register command parsers
	CommandSetFactory(SEARCH, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(SearchCommand).Init(dict)
	}))
	CommandSetFactory(ONLINE_USERS, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(SearchCommand).Init(dict)
	}))

	CommandSetFactory(REPORT, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(ReportCommand).Init(dict)
	}))
	CommandSetFactory("broadcast", NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(ReportCommand).Init(dict)
	}))
	CommandSetFactory(ONLINE, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(ReportCommand).Init(dict)
	}))
	CommandSetFactory(OFFLINE, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(ReportCommand).Init(dict)
	}))

	//// register content processors
	//ContentProcessorRegister(0, new(AnyContentProcessor).Init())
	//
	//// register command processors
	//CommandProcessorRegister(RECEIPT, new(ReceiptCommandProcessor).Init())
	//CommandProcessorRegister(MUTE, new(MuteCommandProcessor).Init())
	//CommandProcessorRegister(BLOCK, new(BlockCommandProcessor).Init())
}

func init() {
	UpgradeIDFactory()
	RegisterCommonFactories()
}
