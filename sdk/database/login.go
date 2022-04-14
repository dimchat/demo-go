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
package db

import (
	. "github.com/dimchat/demo-go/sdk/utils"
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp/protocol"
)

//-------- LoginTable

func (db *Storage) GetLoginCommand(user ID) LoginCommand {
	cmd, _ := getLoginInfo(db, user)
	return cmd
}

func (db *Storage) GetLoginMessage(user ID) ReliableMessage {
	_, msg := getLoginInfo(db, user)
	return msg
}

func (db *Storage) SaveLoginCommandMessage(cmd LoginCommand, msg ReliableMessage) bool {
	if cacheLoginInfo(db, cmd, msg) {
		return saveLoginInfo(db, cmd, msg)
	} else {
		return false
	}
}

/**
 *  Login info for Users
 *  ~~~~~~~~~~~~~~~~~~~~
 *
 *  file path: '.dim/mkm/{zzz}/{ADDRESS}/login.js'
 */

func loginInfoPath(db *Storage, identifier ID) string {
	return PathJoin(db.mkmDir(identifier), "login.js")
}

func loadLoginInfo(db *Storage, identifier ID) (cmd LoginCommand, msg ReliableMessage) {
	path := loginInfoPath(db, identifier)
	db.log("Loading login info: " + path)
	info := db.readMap(path)
	if info == nil {
		return nil, nil
	}
	cmd, _ = ContentParse(info["cmd"]).(LoginCommand)
	msg = ReliableMessageParse(info["msg"])
	return cmd, msg
}

func saveLoginInfo(db *Storage, cmd LoginCommand, msg ReliableMessage) bool {
	info := make(map[string]interface{})
	info["cmd"] = cmd.Map()
	info["msg"] = msg.Map()
	identifier := cmd.ID()
	path := loginInfoPath(db, identifier)
	db.log("Saving login info: " + path)
	return db.writeMap(path, info)
}

// place holder
var emptyMessage = NewReliableMessage(nil)

func getLoginInfo(db *Storage, identifier ID) (cmd LoginCommand, msg ReliableMessage) {
	// 1. try from memory cache
	msg = db._loginMessages[identifier]
	if msg == nil {
		// 2. try from local storage
		cmd, msg = loadLoginInfo(db, identifier)
		if msg == nil {
			// place an empty message for cache
			db._loginMessages[identifier] = emptyMessage
		} else {
			// cache them
			db._loginCommands[identifier] = cmd
			db._loginMessages[identifier] = msg
		}
	} else if msg == emptyMessage {
		cmd = nil
		msg = nil
	} else {
		cmd = db._loginCommands[identifier]
	}
	return cmd, msg
}

func cacheLoginInfo(db *Storage, cmd LoginCommand, msg ReliableMessage) bool {
	// 1. verify sender ID
	identifier := cmd.ID()
	if msg.Sender().Equal(identifier) == false {
		return false
	}
	// 2. check last login time
	old, _ := getLoginInfo(db, identifier)
	if old != nil {
		oldTime := old.Time().Unix()
		newTime := cmd.Time().Unix()
		if newTime <= oldTime {
			// expired command, drop it
			return false
		}
	}
	// 3. cache them
	db._loginCommands[identifier] = cmd
	db._loginMessages[identifier] = msg
	return true
}
