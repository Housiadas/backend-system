// Copyright 2024-2025 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris

package protoplugin

import (
	"os"
	"syscall"
)

// extraInterruptSignals are signals beyond os.Interrupt that we want to be handled
// as interrupts.
//
// For unix-like platforms, this adds syscall.SIGTERM.
var extraInterruptSignals = []os.Signal{syscall.SIGTERM}
