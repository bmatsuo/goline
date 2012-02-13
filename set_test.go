package goline
/*
 *  Filename:    set_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Fri Aug 19 03:15:13 PDT 2011
 *  Description: 
 *  Usage:       gotest
 */
import (
    "testing"
)

func TestShellSet(T *testing.T) {
    set := shellCommandSet{"echo", "ls", "which", "exit"}
    if !set.Has("exit") {
        T.Errorf("simplest good case failure")
    }
    if !set.Has("  echo") {
        T.Errorf("preceding space error")
    }
    if !set.Has("  echo\t") {
        T.Errorf("preceding & trailing space error")
    }
    if !set.Has("which 6g") {
        T.Errorf("simple argument error")
    }
}
