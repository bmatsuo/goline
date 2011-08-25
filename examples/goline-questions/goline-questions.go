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
    "os"
)

var opt = parseFlags()

func main() {
    cont := true
    for cont {
        // "Read a byte" but short-circuit the prompt.
        var a uint8
        goline.Ask(&a, "This should not appear:  ", func(a *goline.Question) {
            a.FirstAnswer = uint(0xFF)
            fmt.Println(a.FirstAnswer)
            a.In(goline.UintRange{200, 255})
            a.Panic = func(err os.Error) {
                fmt.Printf("Error: %s\n", err.String())
            }
        })
        fmt.Printf("byte 0x%X\n", a)

        // Read a bounded integer.
        var b int32
        goline.Ask(&b, "Enter an int:  ", func(a *goline.Question) {
            a.Responses[goline.AskOnError] = a.Question
            a.Default = 13
            a.In(goline.IntRange{26, 62})
            a.Panic = func(err os.Error) {
                fmt.Printf("Error: %s\n", err.String())
            }
        })
        fmt.Printf("Integer %d\n", b)

        // Read a string contained in a set of possible values.
        var broken bool
        cont = !goline.Confirm("Exit?  ", true, func(a *goline.Question) {
            a.Panic = func(err os.Error) {
                if err == os.EOF {
                    broken = true
                }
                fmt.Printf("Error: %s\n", err.String())
            }
        })
        cont = cont && !broken
    }
}
