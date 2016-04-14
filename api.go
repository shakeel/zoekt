// Copyright 2016 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zoekt

import (
	"time"
)

// FileMatch contains all the matches within a file.
type FileMatch struct {
	// Ranking; the higher, the better.
	Score   float64
	Name    string
	Matches []Match
}

// Match is a match for a single atomic query within a file.
type Match struct {
	Line    string
	LineNum int

	// Offset within the line.
	LineOff int

	// Offset from file start
	Offset      uint32
	MatchLength int

	// If set, this was a match on the filename.
	FileName bool

	// The higher the better. Only ranks the quality of the match
	// within the file, does not take rank of file into account
	Score float64
}

// Stats contains interesting numbers on the search
type Stats struct {
	NgramMatches    int
	FilesConsidered int
	FilesLoaded     int
	FileCount       int
	MatchCount      int
	Duration        time.Duration
}

func (s *Stats) Add(o Stats) {
	s.NgramMatches += o.NgramMatches
	s.FilesLoaded += o.FilesLoaded
	s.MatchCount += o.MatchCount
	s.FileCount += o.FileCount
	s.FilesConsidered += o.FilesConsidered
}

// SearchResult contains search matches and extra data
type SearchResult struct {
	Stats
	Files []FileMatch
}

type Searcher interface {
	Search(query Query) (*SearchResult, error)
	Close() error
}
