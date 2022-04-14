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
)

//-------- MetaTable

func (db *Storage) SaveMeta(meta Meta, entity ID) bool {
	if cacheMeta(db, meta, entity) {
		return saveMeta(db, meta, entity)
	} else {
		return false
	}
}

func (db *Storage) GetMeta(entity ID) Meta {
	return getMeta(db, entity)
}

/**
 *  Meta file for Entities (User/Group)
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 *  file path: '.dim/mkm/{zzz}/{ADDRESS}/meta.js'
 */

func metaPath(db *Storage, identifier ID) string {
	return PathJoin(db.mkmDir(identifier), "meta.js")
}

func loadMeta(db *Storage, identifier ID) Meta {
	path := metaPath(db, identifier)
	db.log("Loading meta: " + path)
	return MetaParse(db.readMap(path))
}

func saveMeta(db *Storage, meta Meta, identifier ID) bool {
	info := meta.Map()
	path := metaPath(db, identifier)
	db.log("Saving meta: " + path)
	return db.writeMap(path, info)
}

// place holder
var emptyMeta = MetaGenerate(MKM, emptyPrivateKey, "empty")

func getMeta(db *Storage, identifier ID) Meta {
	// 1. try from memory cache
	meta := db._metas[identifier]
	if meta == nil {
		// 2. try from local storage
		meta = loadMeta(db, identifier)
		if meta == nil {
			// place an empty meta for cache
			db._metas[identifier] = emptyMeta
		} else {
			// cache it
			db._metas[identifier] = meta
		}
	} else if meta == emptyMeta {
		meta = nil
	}
	return meta
}

func cacheMeta(db *Storage, meta Meta, identifier ID) bool {
	// 1. verify meta with ID
	if MetaMatchID(meta, identifier) {
		// 2. cache it
		db._metas[identifier] = meta
		return true
	} else {
		return false
	}
}
