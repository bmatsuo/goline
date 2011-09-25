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
    "goline"
    "fmt"
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
    i, _ := goline.Choose(func(m *goline.Menu) {
        m.Header = "Choose an item"
        m.Question = "Which item do you want?"
        m.SetChoices(items)
    })
    fmt.Println("Selected", i)
}
