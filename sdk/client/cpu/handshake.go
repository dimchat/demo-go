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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
	. "github.com/dimchat/sdk-go/protocol"
)

type HandshakeCommandProcessor struct {
	BaseCommandProcessor
}

func (cpu *HandshakeCommandProcessor) Init() *HandshakeCommandProcessor {
	if cpu.BaseCommandProcessor.Init() != nil {
	}
	return cpu
}

func (cpu *HandshakeCommandProcessor) Execute(cmd Command, _ ReliableMessage) Content {
	hsCmd, _ := cmd.(HandshakeCommand)
	message := hsCmd.Message()
	if message == "DIM?" {
		// station ask client to handshake again
		return HandshakeCommandRestart(hsCmd.Session())
	} else if message == "DIM!" {
		// handshake accepted by station
		//server := cpu.Messenger().Server()
		//server.OnHandshakeAccepted(hsCmd)
		return nil
	} else {
		panic(cmd)
	}
}
