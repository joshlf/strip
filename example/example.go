// Copyright 2012 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"fmt"
	"io"
	"github.com/joshlf13/strip"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <filename>; falling back to stdin\n", os.Args[0])
		// This reader ignores c-style comment delimeters
		io.Copy(os.Stdout, strip.NewReader(os.Stdin, []byte{'/', '/'}))
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	// You can nest them for multiple sets of byte sequences
	inner := strip.NewReader(file, []byte{'/', '/'})
	outer := strip.NewReader(inner, []byte{'/', '\n', '/'})
	io.Copy(os.Stdout, outer)
}