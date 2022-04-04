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
package cpu

import (
	"fmt"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
	. "github.com/dimchat/sdk-go/dimp/cpu"
	. "github.com/dimchat/sdk-go/dimp/dkd"
)

type AnyContentProcessor struct {
	BaseContentProcessor
}

func NewAnyContentProcessor(facebook IFacebook, messenger IMessenger) *AnyContentProcessor {
	cpu := new(AnyContentProcessor)
	cpu.Init(facebook, messenger)
	return cpu
}

func (cpu *AnyContentProcessor) Process(content Content, rMsg ReliableMessage) []Content {
	var text string

	// File: Image, Audio, Video
	msgType := content.Type()
	switch msgType {
	case FILE:
		text = "File Received"
	case IMAGE:
		text = "Image received"
	case AUDIO:
		text = "Audio Received"
	case VIDEO:
		text = "Video received"
	case PAGE:
		text = "Web page received"
	default:
		text := fmt.Sprintf(FmtContentNotSupport, content.Type())
		return cpu.RespondText(text, content.Group())
	}

	group := content.Group()
	if group != nil {
		// DON'T response group message for disturb reason
		return nil
	}

	// response
	sn := content.SN()
	env := rMsg.Envelope()
	signature := rMsg.Get("signature")
	receipt := NewReceiptCommand(text, env, sn, nil)
	receipt.Set("signature", signature)
	return cpu.RespondContent(receipt)
}
