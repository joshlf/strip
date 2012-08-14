// Copyright 2012 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The strip package provides a wrapper for an io.Reader which ignores custom byte sequences.
package strip

import "io"

type stripper struct {
	r io.Reader
	// Character sequence to ignore
	i []byte
	// Buffer of back-write characters
	b    []byte
	len  int
	ipos int
	// Buffer start pos
	bspos int
	// Buffer end pos
	bepos int
}

func NewReader(rdr io.Reader, b []byte) io.Reader {
	return &stripper{rdr, b, make([]byte, len(b)), len(b), 0, 0, 0}
}

func (s *stripper) Read(p []byte) (int, error) {
	n, err := s.r.Read(p)
	p = p[:n]
	plen := len(p)
	var wc, rc int
	for rc < plen && wc < plen {
		// If possible, write one 
		// character from the buffer
		if s.bspos != s.bepos {
			p[wc] = s.b[s.bspos]
			s.bspos = (s.bspos + 1) % s.len
			wc++
		}

		if p[rc] == s.i[s.ipos] {
			s.ipos++
			if s.ipos == s.len {
				s.ipos = 0
			}
		} else {
			if s.ipos != 0 {
				for i := 0; i < s.ipos; i++ {
					s.b[(s.bepos+i)%s.len] = s.i[i]
				}
				s.bepos = (s.bepos + s.ipos) % s.len
				s.ipos = 0	
			}
			// Either write another character
			// from the buffer and write one
			// from the input to the buffer or
			// write directly from input.
			if s.bspos != s.bepos {
				p[wc] = s.b[s.bspos]
				s.b[s.bepos] = p[rc]
				s.bspos = (s.bspos + 1) % s.len
				s.bepos = (s.bepos + 1) % s.len
			} else {
				p[wc] = p[rc]
			}
			wc++
		}
		rc++
	}
	
	// Flush the buffer
	for s.bspos != s.bepos && wc < plen {
		p[wc] = s.b[s.bspos]
		s.bspos = (s.bspos + 1) % s.len
		wc++
	}
	return wc, err
}
