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
package utils

import (
	"fmt"
	"time"
)

const (
	debugFlag   = 0x01
	infoFlag    = 0x02
	warningFlag = 0x04
	errorFlag   = 0x08
)

const (
	debug   = 0xFF
	develop = 0xFE
	release = 0xFC
)

var LogLevel = develop

func now() string {
	current := time.Now()
	return current.Format("2006-01-02 15:04:05")
}

func LogDebug(msg string) {
	if LogLevel & debugFlag == 0 {
		return
	}
	fmt.Printf("[%s] DEBUG - %s\n", now(), msg)
}

func LogInfo(msg string) {
	if LogLevel & infoFlag == 0 {
		return
	}
	fmt.Printf("[%s] %s\n", now(), msg)
}

func LogWarning(msg string) {
	if LogLevel & warningFlag == 0 {
		return
	}
	fmt.Printf("[%s] %s\n", now(), msg)
}

func LogError(msg string) {
	if LogLevel & errorFlag == 0 {
		return
	}
	fmt.Printf("[%s] ERROR - %s\n", now(), msg)
}

type Logging struct {

}

func (loggable *Logging) Debug(msg string) {
	msg = fmt.Sprintf("%T > %s", loggable, msg)
	LogDebug(msg)
}

func (loggable *Logging) Info(msg string) {
	msg = fmt.Sprintf("%T > %s", loggable, msg)
	LogInfo(msg)
}

func (loggable *Logging) Warning(msg string) {
	msg = fmt.Sprintf("%T > %s", loggable, msg)
	LogWarning(msg)
}

func (loggable *Logging) Error(msg string) {
	msg = fmt.Sprintf("%T > %s", loggable, msg)
	LogError(msg)
}
