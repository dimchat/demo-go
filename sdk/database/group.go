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
	. "github.com/dimchat/mkm-go/protocol"
	"strings"
)

//-------- GroupTable

func (db *Storage) GetFounder(group ID) ID {
	// TODO: get group founder
	return nil
}

func (db *Storage) GetOwner(group ID) ID {
	// TODO: get group owner
	return nil
}

func (db *Storage) GetMembers(group ID) []ID {
	arr := db._members[group]
	if arr == nil {
		arr = loadMembers(db, group)
		db._members[group] = arr
	}
	return arr
}

func (db *Storage) GetAssistants(group ID) []ID {
	// TODO: get group assistants
	return nil
}

func (db *Storage) AddMember(member ID, group ID) bool {
	arr := db.GetMembers(group)
	for _, item := range arr {
		if member.Equal(item) {
			// duplicated
			return false
		}
	}
	arr = append(arr, member)
	return db.SaveMembers(arr, group)
}

func (db *Storage) RemoveMember(member ID, group ID) bool {
	arr := db.GetMembers(group)
	var pos = -1
	for index, id := range arr {
		if member.Equal(id) {
			pos = index
			break
		}
	}
	if pos == -1 {
		// contact ID not found
		return false
	} else {
		arr = append(arr[:pos], arr[pos+1:]...)
		return db.SaveMembers(arr, group)
	}
}

func (db *Storage) SaveMembers(members []ID, group ID) bool {
	db._members[group] = members
	return saveMembers(db, group, members)
}

func (db *Storage) RemoveGroup(group ID) bool {
	// TODO: remove group info
	return false
}

/**
 *  Members file for Group
 *  ~~~~~~~~~~~~~~~~~~~~~~
 *
 *  file path: '.dim/mkm/{zzz}/{ADDRESS}/members.txt'
 */

func membersPath(db *Storage, group ID) string {
	return PathJoin(db.mkmDir(group), "members.txt")
}

func loadMembers(db *Storage, group ID) []ID {
	path := membersPath(db, group)
	db.log("Loading members for group: " + group.String())
	text := db.readText(path)
	lines := strings.Split(text, "\n")
	members := make([]ID, 0, len(lines))
	for _, rec := range lines {
		id := IDParse(rec)
		if id != nil {
			members = append(members, id)
		}
	}
	return members
}

func saveMembers(db *Storage, group ID, members []ID) bool {
	text := ""
	lines := IDRevert(members)
	for _, rec := range lines {
		text = text + rec + "\n"
	}
	path := membersPath(db, group)
	db.log("Saving members for group: " + group.String())
	return db.writeText(path, text)
}
