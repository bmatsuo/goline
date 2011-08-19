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
        a.SetDefault(13)
        a.InRange(int64(26), int64(62))
    })
    if err != nil {
        fmt.Printf("Error: %s\n", err.String())
    }
    fmt.Printf("Integer %d\n", b)
}
