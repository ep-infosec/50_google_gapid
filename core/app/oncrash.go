// Copyright (C) 2017 Google Inc.
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

package app

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"

	"github.com/google/gapid/core/fault/stacktrace"
	"github.com/google/gapid/core/fault/stacktrace/crunch"
)

func decodeCrashCode(s string) stacktrace.Callstack {
	d := base64.NewDecoder(base64.RawStdEncoding, bytes.NewReader([]byte(s)))
	data, err := ioutil.ReadAll(d)
	if err != nil {
		return stacktrace.Callstack{}
	}
	return crunch.Uncrunch(data)
}
