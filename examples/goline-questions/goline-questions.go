// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main
/*
 *  Filename:    goline-questions.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Aug 13 21:39:57 PDT 2011
 *  Description: 
 *  Usage:       goline-questions [options] ARGUMENT ...
 */
import (
    "goline"
    "fmt"
)

var opt = parseFlags()

func main() {
    var b int32
    err := goline.Ask(&b, "Enter an int:  ", func(a *goline.Answer) {
        a.Responses[goline.AskOnError] = a.Question
        a.Default = 13
        a.In(goline.IntRange{26, 62})
    })
    if err != nil {
        fmt.Printf("Error: %s\n", err.String())
    }
    fmt.Printf("Integer %d\n", b)
    var s string
    err = goline.Ask(&s, "Exit?  ", func(a *goline.Answer) {
        a.Default = "yes"
        a.In(goline.StringSet([]string{"yes", "y", "no", "n"}))
    })
    if err != nil {
        fmt.Printf("Error: %s\n", err.String())
    }
    fmt.Printf("String %s\n", s)
}
