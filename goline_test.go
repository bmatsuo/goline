// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goline
/*
 *  Filename:    goline_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Aug 13 02:28:54 PDT 2011
 *  Description: 
 *  Usage:       gotest
 */
import (
    "testing"
    "bytes"
    "io"
    "os"
)

//  Returns true if f raises a runtime panic, false otherwise.
func CausesPanic(f func()) (paniced bool) {
    defer func() {
        if e := recover(); e != nil {
            paniced = true
        }
    }()
    f()
    return
}

func FuncOutput(fn func(wr io.Writer)) string {
    pout := new(bytes.Buffer)
    fn(pout)
    return string(pout.Bytes())
}

func StringsShouldBeEqual(a, b string) os.Error {
    if a != b {
        return os.NewError("Not Equal")
    }
    return nil
}

func OutputEqualityTest(name string, fn func(io.Writer), expect string, T *testing.T) {
    defer func() {
        if e := recover(); e != nil {
            T.Errorf("%s: Function paniced %#v", name, e)
        }
    }()
    output := FuncOutput(fn)
    if output != expect {
        T.Errorf("%s: Unexpected output\n\noutput '%s'\n\nexpected '%s'\n\n", name, output, expect)
    } else {
        T.Logf("%s: PASS", name)
    }
}

func TestGoline(T *testing.T) {
}

func TestSay(T *testing.T) {
    OutputEqualityTest("It adds a newline when none is present",
        func(wr io.Writer) {
            fSay(wr, "Hello, World!", false)
        }, "Hello, World!\n", T)
    OutputEqualityTest("It does't add a newline when there is trailing space ",
        func(wr io.Writer) {
            fSay(wr, "Hello, World! ", false)
        }, "Hello, World! ", T)
    OutputEqualityTest("It can trim trailing space, leaving one trailing newline",
        func(wr io.Writer) {
            fSay(wr, "Hello, World!   \t\n", true)
        }, "Hello, World!\n", T)
}

func TestList(T *testing.T) {
    // goline.List([]string{"cat", "dog", "go fish"}, goline.ColumnsAcross, 15)
    // /* Outputs:
    // *  cat     dog
    // *  go fish
    // */
    OutputEqualityTest("It prints columns alphabetically left to right",
        func(wr io.Writer) {
            fList(wr, []string{"cat", "dog", "go fish"}, ColumnsAcross, 15)
        }, "cat     dog\ngo fish\n", T)

    //      goline.List([]string{"cat", "dog", "go fish"}, goline.ColumnsDown, 15)
    //      /* Outputs:
    //       *  cat     go fish
    //       *  dog
    //       */
    OutputEqualityTest("It prints columns alphabetically up to down",
        func(wr io.Writer) {
            fList(wr, []string{"cat", "dog", "go fish"}, ColumnsDown, 15)
        }, "cat     go fish\ndog\n", T)

    //      goline.List([]string{"cat", "dog", "go fish"}, goline.Inline, "and ")
    //      /* Outputs:
    //       *  cat, dog, and go fish
    //       */
    OutputEqualityTest("It prints inline lists with custom separators",
        func(wr io.Writer) {
            fList(wr, []string{"cat", "dog", "go fish"}, Inline, "and ")
        }, "cat, dog, and go fish\n", T)

    //      goline.List([]string{"cat", "dog", "go fish"}, goline.Rows, nil)
    //      /* Outputs:
    //       *  cat
    //       *  dog
    //       *  go fish
    //       */
    OutputEqualityTest("It prints everything on separate rows",
        func(wr io.Writer) {
            fList(wr, []string{"cat", "dog", "go fish"}, Rows, nil)
        }, "cat\ndog\ngo fish\n", T)
}
