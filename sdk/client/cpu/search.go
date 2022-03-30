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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/demo-go/sdk/common/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
)

type SearchCommandProcessor struct {
	BaseCommandProcessor
}

func (cpu *SearchCommandProcessor) parse(cmd *SearchCommand) {
	results := cmd.Results()
	if results == nil {
		return
	}
	facebook := cpu.Facebook()
	var identifier ID
	var meta Meta
	for key, value := range results {
		identifier = IDParse(key)
		meta = MetaParse(value)
		if identifier == nil || meta == nil || !MetaMatchID(meta, identifier) {
			// TODO: meta error
			continue
		}
		facebook.SaveMeta(meta, identifier)
	}
}

func (cpu *SearchCommandProcessor) Execute(cmd Command, _ ReliableMessage) Content {
	sCmd, _ := cmd.(*SearchCommand)

	cpu.parse(sCmd)

	// message
	message := sCmd.Get("message")
	fmt.Println("search respond:", message)
	// users
	users := sCmd.Get("users")
	if users != nil {
		fmt.Println("	users:", UTF8Decode(JSONEncode(users)))
	}
	// results
	results := sCmd.Get("results")
	if results != nil {
		fmt.Println("	results:", UTF8Decode(JSONEncode(results)))
	}
	return nil
}
