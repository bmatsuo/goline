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
    fmt.Println("ColumnsAcross")
    goline.List(items, goline.ColumnsAcross, nil)
    fmt.Println()

    fmt.Println("ColumnsDown")
    goline.List(items, goline.ColumnsDown, nil)
    fmt.Println()

    fmt.Println("Inline")
    goline.List(items, goline.Inline, nil)
    fmt.Println()

    fmt.Println("Rows")
    goline.List(items, goline.Rows, nil)
    fmt.Println()

    goline.List([]string{"cat", "dog", "go fish"}, goline.ColumnsAcross, nil)
    /* Outputs:
     *  cat     dog     go fish
     */
    goline.Say("")

    goline.List([]string{"cat", "dog", "go fish"}, goline.ColumnsDown, 15)
    /* Outputs:
     *  cat     go fish
     *  dog
     */
    goline.Say("")

    goline.List([]string{"cat", "dog", "go fish"}, goline.Inline, "and ")
    /* Outputs:
     *  cat, dog, and go fish
     */
    goline.Say("")

    goline.List([]string{"cat", "dog", "go fish"}, goline.Rows, nil)
    /* Outputs:
     *  cat
     *  dog
     *  go fish
     */
    goline.Say("")

}
