// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
 *  Filename:    goline-lists.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Wed Aug 24 02:07:52 PDT 2011
 *  Description: 
 *  Usage:       goline-lists [options] ARGUMENT ...
 */
import (
	"fmt"
	"goline"
)

var opt = parseFlags()

var items = []string{
	"Hello, World!",
	"0xDEADBEEF",
	"NDQ +3.35%",
	"T-t-t-to the top of the world!",
	"Ain't nobody gonna take my H away.",
	"I never hexed a man I didn't like",
}

func main() {
	action := func(choice string, sharg string) {
		fmt.Printf(`"%v" is my final answer.`+"\n", choice)
	}
	i, _ := goline.Choose(func(m *goline.Menu) {
		m.Header = "Choose an item"
		m.Question = "Which item do you want?"
		for i := range items {
			m.Choice(items[i], action)
		}
		m.Choice("Do nothing.", nil)
	})
	fmt.Println("Selected", i)
	fmt.Println()
	for cont := true; cont; {
		goline.Choose(func(m *goline.Menu) {
			m.Header = "Enter a command: "
			m.Question = "?> "
			m.Shell = true
			m.ListMode = goline.Inline
			m.IndexMode = goline.NoIndex
			m.Choice("echo", func(s string, args string) { fmt.Println(s, "|", args) })
			m.Choice("quit", func(s string, args string) { cont = false })
			m.Choice("exit", func(s string, args string) { cont = false })
		})
	}
}
