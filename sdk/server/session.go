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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/plugins/types"
)

// format "(IP, Port)"
type SessionAddress string

type SessionHandler interface {

	PushMessage(msg ReliableMessage) bool
}

type Session interface {

	// user ID
	ID() ID
	SetID(identifier ID)

	// session key
	Key() string

	// connection target "(IP, port)"
	ClientAddress() SessionAddress

	// when the client entered background, it should be set to False
	IsActive() bool
	SetActive(active bool)

	// Push message when session active
	PushMessage(msg ReliableMessage) bool
}

func generateSessionKey() string {
	return HexEncode(RandomBytes(32))
}

type BaseSession struct {
	Session

	_identifier ID
	_key string
	_address SessionAddress
	_active bool
	_handler SessionHandler
}

func NewSession(address SessionAddress, handler SessionHandler) Session {
	return new(BaseSession).Init(address, handler)
}

func (session *BaseSession) Init(address SessionAddress, handler SessionHandler) *BaseSession {
	session._identifier = nil
	session._key = generateSessionKey()
	session._address = address
	session._active = true
	session._handler = handler
	return session
}

//-------- Session

func (session *BaseSession) ID() ID {
	return session._identifier
}
func (session *BaseSession) SetID(identifier ID) {
	session._identifier = identifier
}

func (session *BaseSession) Key() string {
	return session._key
}

func (session *BaseSession) ClientAddress() SessionAddress {
	return session._address
}

func (session *BaseSession) IsActive() bool {
	return session._active
}
func (session *BaseSession) SetActive(active bool) {
	session._active = active
}

func (session *BaseSession) PushMessage(msg ReliableMessage) bool {
	if session._active {
		return session._handler.PushMessage(msg)
	} else {
		return false
	}
}

/**
 *  Session Server
 *  ~~~~~~~~~~~~~~
 */
type SessionServer struct {

	_clientAddresses map[ID][]SessionAddress
	_sessions map[SessionAddress]Session
}

func (server *SessionServer) Init() *SessionServer {
	server._clientAddresses = make(map[ID][]SessionAddress)
	server._sessions = make(map[SessionAddress]Session)
	return server
}

// Session factory
func (server *SessionServer) GetSession(address SessionAddress, handler SessionHandler) Session {
	session := server._sessions[address]
	if session == nil && handler != nil {
		// create a new session and cache it
		session := NewSession(address, handler)
		server._sessions[address] = session
	}
	return session
}

func (server *SessionServer) insert(address SessionAddress, identifier ID) {
	array := server._clientAddresses[identifier]
	if array == nil {
		array = make([]SessionAddress, 0, 1)
	//} else {
	//	for _, item := range array {
	//		if item == address {
	//			// already exists
	//			return
	//		}
	//	}
	}
	server._clientAddresses[identifier] = append(array, address)
}

func (server *SessionServer) remove(address SessionAddress, identifier ID) {
	array := server._clientAddresses[identifier]
	if array == nil {
		// not exists
		return
	}
	pos := len(array)
	for pos > 0 {
		pos--
		if array[pos] == address {
			array = append(array[:pos], array[pos+1:]...)
		}
	}
	if len(array) == 0 {
		// all sessions removed
		delete(server._clientAddresses, identifier)
	}
}

// Insert a session with ID into memory cache
func (server *SessionServer) UpdateSession(session Session, identifier ID) {
	address := session.ClientAddress()
	old := session.ID()
	if old != nil {
		// 0. remove client_address from old ID
		server.remove(address, old)
	}
	// 1. insert client_address for new ID
	server.insert(address, identifier)
	// 2. update session ID
	session.SetID(identifier)
}

// Remove the session from memory cache
func (server *SessionServer) RemoveSession(session Session) {
	identifier := session.ID()
	address := session.ClientAddress()
	if identifier != nil {
		// 1. remove client_address with ID
		server.remove(address, identifier)
	}
	// 2. remove session with client_address
	session.SetActive(false)
	delete(server._sessions, address)
}

// Get all sessions of this user
func (server *SessionServer) AllSessions(identifier ID) []Session {
	results := make([]Session, 0, 1)
	// 1. get all client_address with ID
	array := server._clientAddresses[identifier]
	if array != nil {
		// 2. get session by each client_address
		var session Session
		for _, item := range array {
			session = server._sessions[item]
			if session != nil {
				results = append(results, session)
			}
		}
	}
	return results
}
func (server *SessionServer) ActiveSessions(identifier ID) []Session {
	results := make([]Session, 0, 1)
	// 1. get all sessions
	all := server.AllSessions(identifier)
	for _, item := range all {
		// 2. check session active
		if item.IsActive() {
			results = append(results, item)
		}
	}
	return results
}

//
//  Users
//

func (server *SessionServer) AllUsers() []ID {
	users := make([]ID, 0, 8)
	for key := range server._clientAddresses {
		users = append(users, key)
	}
	return users
}

func (server *SessionServer) IsActive(identifier ID) bool {
	sessions := server.AllSessions(identifier)
	for _, item := range sessions {
		if item.IsActive() {
			return true
		}
	}
	return false
}

func (server *SessionServer) ActiveUsers() []ID {
	users := make([]ID, 0, 8)
	all := server.AllUsers()
	for _, item := range all {
		if server.IsActive(item) {
			users = append(users, item)
		}
	}
	return users
}
