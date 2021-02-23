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
	. "github.com/dimchat/core-go/core"
	. "github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
	"strings"
)

/**
 *  Common Processor
 *  ~~~~~~~~~~~~~~~~
 */
type CommonProcessor struct {
	MessengerProcessor
}

func (processor *CommonProcessor) Init(transceiver ICommonMessenger) *CommonProcessor {
	if processor.MessengerProcessor.Init(transceiver) != nil {
	}
	return processor
}

func (processor *CommonProcessor) Messenger() ICommonMessenger {
	return processor.Transceiver().(ICommonMessenger)
}

func (processor *CommonProcessor) Facebook() ICommonFacebook {
	return processor.MessengerProcessor.Facebook().(ICommonFacebook)
}

// check whether group info empty
func (processor *CommonProcessor) isEmptyGroup(group ID) bool {
	facebook := processor.Facebook()
	members := facebook.GetMembers(group)
	if members == nil || len(members) == 0 {
		return true
	} else {
		return facebook.GetOwner(group) == nil
	}
}

// check whether need to update group
func (processor *CommonProcessor) isWaitingGroup(content Content, sender ID) bool {
	// Check if it is a group message, and whether the group members info needs update
	group := content.Group()
	if group == nil || group.IsBroadcast() {
		// 1. personal message
		// 2. broadcast message
		return false
	}
	// check meta for new group ID
	messenger := processor.Messenger()
	facebook := processor.Facebook()
	meta := facebook.GetMeta(group)
	if meta == nil {
		// NOTICE: if meta for group not found,
		//         facebook should query it from DIM network automatically
		// TODO: insert the message to a temporary queue to wait meta
		//throw new NullPointerException("group meta not found: " + group);
		return true
	}
	// query group info
	if processor.isEmptyGroup(group) {
		// NOTICE: if the group info not found, and this is not an 'invite' command
		//         query group info from the sender
		cmd, ok := content.(Command)
		if ok {
			name := cmd.CommandName()
			if name == INVITE || name == RESET {
				// FIXME: can we trust this stranger?
				//        may be we should keep this members list temporary,
				//        and send 'query' to the owner immediately.
				// TODO: check whether the members list is a full list,
				//       it should contain the group owner(owner)
				return false
			}
		}
		return messenger.QueryGroupInfo(group, []ID{sender})
	} else if facebook.ContainMember(sender, group) ||
		facebook.ContainAssistant(sender, group) ||
		facebook.IsOwner(sender, group) {
		// normal membership
		return false
	} else {
		var ok1, ok2 bool
		// if assistants exist, query them
		bots := facebook.GetAssistants(group)
		if bots == nil || len(bots) == 0 {
			ok1 = false
		} else {
			ok1 = messenger.QueryGroupInfo(group, bots)
		}
		// if owner found, query it too
		owner := facebook.GetOwner(group)
		if owner == nil {
			ok2 = false
		} else {
			ok2 = messenger.QueryGroupInfo(group, []ID{owner})
		}
		return ok1 && ok2
	}
}

func (processor *CommonProcessor) ProcessContent(content Content, rMsg ReliableMessage) Content {
	sender := rMsg.Sender()
	if processor.isWaitingGroup(content, sender) {
		// save this message in a queue to wait group meta response
		group := content.Group()
		rMsg.Set("waiting", group.String())
		processor.Messenger().SuspendReliableMessage(rMsg)
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			text, ok := r.(string)
			if ok && strings.Contains(text, "failed to get meta for") {
				pos := strings.Index(text, ": ")
				if pos > 0 {
					waiting := IDParse(text[pos+2:])
					if waiting == nil {
						panic("failed to get ID: " + text)
					} else {
						rMsg.Set("waiting", waiting.String())
						processor.Messenger().SuspendReliableMessage(rMsg)
					}
				}
			}
		}
	}()
	return processor.MessengerProcessor.ProcessContent(content, rMsg)
}
