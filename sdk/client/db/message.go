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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
)

type ConversationTable interface {

	/**
	 *  Get how many chat boxes
	 *
	 * @return conversations count
	 */
	NumberOfConversations() int

	/**
	 *  Get chat box info
	 *
	 * @param index - sorted index
	 * @return conversation ID
	 */
	ConversationAtIndex(index int) ID

	/**
	 *  Remove one chat box
	 *
	 * @param index - chat box index
	 * @return true on row(s) affected
	 */
	RemoveConversationAtIndex(index int) bool

	/**
	 *  Remove the chat box
	 *
	 * @param entity - conversation ID
	 * @return true on row(s) affected
	 */
	RemoveConversation(entity ID) bool
}

type MessageTable interface {

	/**
	 *  Get message count in this conversation for an entity
	 *
	 * @param entity - conversation ID
	 * @return total count
	 */
	NumberOfMessages(entity ID) int

	/**
	 *  Get unread message count in this conversation for an entity
	 *
	 * @param entity - conversation ID
	 * @return unread count
	 */
	NumberOfUnreadMessages(entity ID) int

	/**
	 *  Clear unread flag in this conversation for an entity
	 *
	 * @param entity - conversation ID
	 * @return true on row(s) affected
	 */
	ClearUnreadMessages(entity ID) bool

	/**
	 *  Get last message of this conversation
	 *
	 * @param entity - conversation ID
	 * @return instant message
	 */
	LastMessage(entity ID) InstantMessage

	/**
	 *  Get last received message from all conversations
	 *
	 * @param user - current user ID
	 * @return instant message
	 */
	LastReceivedMessage(user ID) InstantMessage

	/**
	 *  Get message at index of this conversation
	 *
	 * @param index - start from 0, latest first
	 * @param entity - conversation ID
	 * @return instant message
	 */
	MessageAtIndex(index int, entity ID) InstantMessage

	/**
	 *  Save the new message to local storage
	 *
	 * @param iMsg - instant message
	 * @param entity - conversation ID
	 * @return true on success
	 */
	InsertMessage(iMsg InstantMessage, entity ID) bool

	/**
	 *  Delete the message
	 *
	 * @param iMsg - instant message
	 * @param entity - conversation ID
	 * @return true on row(s) affected
	 */
	RemoveMessage(iMsg InstantMessage, entity ID) bool

	/**
	 *  Try to withdraw the message, maybe won't success
	 *
	 * @param iMsg - instant message
	 * @param entity - conversation ID
	 * @return true on success
	 */
	WithdrawMessage(iMsg InstantMessage, entity ID) bool

	/**
	 *  Update message state with receipt
	 *
	 * @param iMsg - message with receipt content
	 * @param entity - conversation ID
	 * @return true while target message found
	 */
	SaveReceipt(iMsg InstantMessage, entity ID) bool
}
