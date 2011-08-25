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


func TestGoline(T *testing.T) {
}
