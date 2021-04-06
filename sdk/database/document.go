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

//-------- DocumentTable

func (db *Storage) SaveDocument(doc Document) bool {
	if cacheDocument(db, doc) {
		return saveDocument(db, doc)
	} else {
		return false
	}
}

func (db *Storage) GetDocument(entity ID, docType string) Document {
	docType = documentType(docType, entity)
	return getDocument(db, entity, docType)
}

/**
 *  Document for Entities (User/Group)
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 *  file path: '.dim/mkm/{zzz}/{ADDRESS}/doc_{type}.js'
 */

func documentPath(db *Storage, identifier ID, docType string) string {
	return PathJoin(db.mkmDir(identifier), documentFile(docType))
}
func documentFile(docType string) string {
	if docType == VISA {
		return "visa.js"
	} else {
		// TODO: other types?
		return "doc.js"
	}
}
func documentType(docType string, identifier ID) string {
	if docType != "" && docType != "*" {
		return docType
	} else if identifier.IsUser() {
		return VISA
	} else if identifier.IsGroup() {
		return BULLETIN
	} else {
		return PROFILE
	}
}

func loadDocument(db *Storage, identifier ID, docType string) Document {
	path := documentPath(db, identifier, docType)
	db.log("Loading document: " + path)
	return DocumentParse(db.readMap(path))
}

func saveDocument(db *Storage, doc Document) bool {
	info := doc.GetMap(false)
	path := documentPath(db, doc.ID(), doc.Type())
	db.log("Saving document: " + path)
	return db.writeMap(path, info)
}

// place holder
var emptyProfile = DocumentCreate(PROFILE, ANYONE, nil, nil)

func getDocument(db *Storage, identifier ID, docType string) Document {
	// 1. try from memory cache
	var doc Document
	table := db._docs[docType]
	if table == nil {
		// FIXME: document type not support?
		table = make(map[ID]Document)
		db._docs[docType] = table
		doc = nil
	} else {
		doc = table[identifier]
	}
	if doc == nil {
		// 2. try from local storage
		doc = loadDocument(db, identifier, docType)
		if doc == nil {
			// place an empty doc for cache
			table[identifier] = emptyProfile
		} else {
			// cache it
			table[identifier] = doc
		}
	} else if doc == emptyProfile {
		doc = nil
	}
	return doc
}

func cacheDocument(db *Storage, doc Document) bool {
	// 1. check valid
	if doc.IsValid() == false {
		return false
	}
	// 2. prepare table with document type
	identifier := doc.ID()
	docType := documentType(doc.Type(), identifier)
	table := db._docs[docType]
	if table == nil {
		// FIXME: document type not support?
		table = make(map[ID]Document)
		db._docs[docType] = table
	}
	// 3. cache it
	table[identifier] = doc
	return true
}
