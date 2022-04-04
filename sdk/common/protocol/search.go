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

import (
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/mkm-go/protocol"
)

const (
	SEARCH = "search"
	ONLINE_USERS = "users"  // search online users
)

/**
 *  Command message: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command  : "search",        // or "users"
 *
 *      keywords : "keywords",      // keyword string
 *      users    : ["ID"],          // user ID list
 *      results  : {"ID": {meta}, } // user's meta map
 *  }
 */
type SearchCommand struct {
	BaseCommand
}

func (cmd *SearchCommand) Init(dict map[string]interface{}) *SearchCommand {
	if cmd.BaseCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *SearchCommand) InitWithKeywords(keywords string) *SearchCommand {
	var command string
	if keywords == ONLINE_USERS {
		command = ONLINE_USERS
		keywords = ""
	} else {
		command = SEARCH
	}
	if cmd.BaseCommand.InitWithCommand(command) != nil {
		if keywords != "" {
			cmd.Set("keywords", keywords)
		}
	}
	return cmd
}

/**
 *  Get user ID list
 *
 * @return ID string list
 */
func (cmd *SearchCommand) Users() []ID {
	users := cmd.Get("users")
	if users == nil {
		return nil
	}
	return IDConvert(users)
}

/**
 *  Get user metas mapping to ID strings
 *
 * @return meta dictionary
 */
func (cmd *SearchCommand) Results() map[string]interface{} {
	results := cmd.Get("results")
	if results == nil {
		return nil
	}
	return results.(map[string]interface{})
}
