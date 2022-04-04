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
package protocol

import . "github.com/dimchat/core-go/dkd"

const (
	REPORT = "report"
	ONLINE = "online"
	OFFLINE = "offline"
)

/**
 *  Command message: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command  : "report",
 *      title    : "online",      // or "offline"
 *      //---- extra info
 *      time     : 1234567890,    // timestamp?
 *  }
 */
type ReportCommand struct {
	BaseCommand
}

func (cmd *ReportCommand) Init(dict map[string]interface{}) *ReportCommand {
	if cmd.BaseCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *ReportCommand) InitWithTitle(title string) *ReportCommand {
	if cmd.BaseCommand.InitWithCommand(REPORT) != nil {
		cmd.SetTitle(title)
	}
	return cmd
}

func (cmd *ReportCommand) Title() string {
	text := cmd.Get("title")
	if text == nil {
		return ""
	}
	return text.(string)
}

func (cmd *ReportCommand) SetTitle(title string) {
	cmd.Set("title", title)
}
