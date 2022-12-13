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

package gles

import (
	"strings"

	"github.com/google/gapid/core/stream"
	"github.com/google/gapid/core/stream/fmts"
	"github.com/google/gapid/gapis/vertex"
)

var semanticPatterns = []struct {
	pattern  string
	semantic vertex.Semantic_Type
}{
	// Ordered from highest priority to lowest
	{"position", vertex.Semantic_Position},
	{"normal", vertex.Semantic_Normal},
	{"tangent", vertex.Semantic_Tangent},
	{"bitangent", vertex.Semantic_Bitangent},
	{"binormal", vertex.Semantic_Bitangent},
	{"texcoord", vertex.Semantic_Texcoord},
	{"pos", vertex.Semantic_Position},
	{"uv", vertex.Semantic_Texcoord},
	{"vertex", vertex.Semantic_Position},
}

var semanticFormats = []struct {
	format   *stream.Format
	semantic vertex.Semantic_Type
}{
	// Ordered from highest priority to lowest
	{fmts.XYZ_F32, vertex.Semantic_Position},
	{fmts.XYZ_S8_NORM, vertex.Semantic_Normal},
}

// guessSemantics uses string and format matching to try and guess the semantic
// usage of a vertex stream.
// This is a big fat hack. See: https://github.com/google/gapid/issues/960
func guessSemantics(vb *vertex.Buffer, hints map[string]vertex.Semantic_Type) {
	taken := map[vertex.Semantic_Type]bool{}
	if len(hints) > 0 {
		for _, s := range vb.Streams {
			if semantic, ok := hints[s.Name]; ok && !taken[semantic] {
				taken[semantic] = true
				s.Semantic.Type = semantic
			}
		}
	}

	for _, p := range semanticPatterns {
		if taken[p.semantic] {
			continue
		}
		for _, s := range vb.Streams {
			if needsGuess(s) && strings.Contains(strings.ToLower(s.Name), p.pattern) {
				s.Semantic.Type = p.semantic
				taken[p.semantic] = true
				break
			}
		}
	}
	for _, p := range semanticFormats {
		if taken[p.semantic] {
			continue
		}
		for _, s := range vb.Streams {
			if needsGuess(s) && s.Format.String() == p.format.String() {
				s.Semantic.Type = p.semantic
				taken[p.semantic] = true
				break
			}
		}
	}
}

func needsGuess(s *vertex.Stream) bool {
	return s.Semantic.Type == vertex.Semantic_Unknown
}
