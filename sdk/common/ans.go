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
	. "github.com/dimchat/demo-go/sdk/common/db"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/dimp"
)

/**
 *  Data Source for ANS
 *  ~~~~~~~~~~~~~~~~~~~
 */
type AddressNameDataSource struct {
	AddressNameServer

	_ansTable AddressNameTable
}

func (ans *AddressNameDataSource) Init() *AddressNameDataSource {
	if ans.AddressNameServer.Init() != nil {
		ans._ansTable = nil
	}
	return ans
}

func (ans *AddressNameDataSource) GetID(alias string) ID {
	identifier := ans.AddressNameServer.GetID(alias)
	if identifier == nil && ans._ansTable != nil {
		identifier = ans._ansTable.GetIdentifier(alias)
		if identifier != nil {
			// FIXME: is reserved name?
			ans.Cache(alias, identifier)
		}
	}
	return identifier
}

func (ans *AddressNameDataSource) Save(alias string, identifier ID) bool {
	if ans.AddressNameServer.Save(alias, identifier) == false {
		return false
	} else if ValueIsNil(identifier) {
		return ans._ansTable.RemoveRecord(alias)
	} else {
		return ans._ansTable.AddRecord(identifier, alias)
	}
}

/**
 *  ID Factory
 *  ~~~~~~~~~~
 */
type CommonIDFactory struct {
	IDFactory

	_ans AddressNameService
	_origin IDFactory
}

func (factory *CommonIDFactory) Init(ans AddressNameService, origin IDFactory) *CommonIDFactory {
	factory._ans = ans
	factory._origin = origin
	return factory
}

func (factory *CommonIDFactory) GenerateID(meta Meta, network NetworkType, terminal string) ID {
	return factory._origin.GenerateID(meta, network, terminal)
}

func (factory *CommonIDFactory) CreateID(name string, address Address, terminal string) ID {
	return factory._origin.CreateID(name, address, terminal)
}

func (factory *CommonIDFactory) ParseID(identifier string) ID {
	// try ANS record
	id := factory._ans.GetID(identifier)
	if id == nil {
		// parse by original factory
		id = factory._origin.ParseID(identifier)
	}
	return id
}

func UpgradeIDFactory() {
	// ANS
	ans := new(AddressNameDataSource)
	ans.Init()

	// origin ID factory
	origin := IDGetFactory()

	// wrap
	factory := new(CommonIDFactory)
	factory.Init(ans, origin)
	IDSetFactory(factory)
}
