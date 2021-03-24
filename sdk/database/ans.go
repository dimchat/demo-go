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

//-------- AddressNameTable

func (db *Storage) GetIdentifier(alias string) ID {
	return db._ans[alias]
}

func (db *Storage) AddRecord(identifier ID, alias string) bool {
	if len(alias) == 0 || identifier == nil {
		return false
	}
	if len(db._ans) == 0 {
		panic("ANS not initialized")
	}
	// cache it
	db._ans[alias] = identifier
	// save them
	return saveANS(db, db._ans)
}

func (db *Storage) RemoveRecord(alias string) bool {
	if len(alias) == 0 || db._ans[alias] == nil {
		return false
	}
	// remove it
	delete(db._ans, alias)
	// save them
	return saveANS(db, db._ans)
}

/**
 *  Address Name Service
 *  ~~~~~~~~~~~~~~~~~~~~
 *
 *  file path: '.dim/ans.txt'
 */

func ansPath(db *Storage) string {
	return PathJoin(db.Root(), "ans.txt")
}

func loadANS(db *Storage) map[string]ID {
	table := make(map[string]ID)

	path := ansPath(db)
	db.log("Loading ANS records from: " + path)
	text := db.readText(path)
	lines := strings.Split(text, "\n")
	for _, rec := range lines {
		pair := strings.Split(rec, "\t")
		if len(pair) != 2 {
			db.error("Invalid ANS record: " + rec)
			continue
		}
		table[strings.TrimSpace(pair[0])] = IDParse(strings.TrimSpace(pair[1]))
	}
	//
	//  Reserved names
	//
	table["all"] = EVERYONE
	table[EVERYONE.Name()] = EVERYONE
	table[ANYONE.Name()] = ANYONE
	table["owner"] = ANYONE
	table["founder"] = FOUNDER

	return table
}

func saveANS(db *Storage, records map[string]ID) bool {
	text := ""
	for key, value := range records {
		if value != nil {
			text += key + "\t" + value.String() + "\n"
		}
	}
	path := ansPath(db)
	db.log("Saving ANS records into: " + path)
	return db.writeText(path, text)
}
