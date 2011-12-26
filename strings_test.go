package goline
/*
 *  Filename:    strings_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Wed Aug 24 22:56:19 PDT 2011
 *  Description: 
 *  Usage:       gotest
 */
import (
    "testing"
    "unicode"
)

func TestSuffix(T *testing.T) {
    s := "Blah   "
    spaces := "   "
    if i := stringSuffixIndexFunc(s, unicode.IsSpace); i != 4 {
        T.Errorf("Did not correctly identify a real suffix (%d != %d)", i, 4)
    }
    if suff := stringSuffixFunc(s, unicode.IsSpace); suff != spaces {
        T.Errorf("Did not correctly return the suffix (%#v != %#v)", suff, spaces)
    }

    all := func(c rune) bool { return true }
    if i := stringSuffixIndexFunc(s, all); i != 0 {
        T.Errorf("Did not correctly identify a the whole string as a suffix")
    }
    if suff:= stringSuffixFunc(s, all); suff != s {
        T.Errorf("Did not correctly return the whole string as a suffix")
    }

    none := func(c rune) bool { return false }
    if i := stringSuffixIndexFunc(s, none); i != -1 {
        T.Errorf("Did not correctly identify an empty suffix (%d != %d)", i, -1)
    }
    if suff := stringSuffixFunc(s, none); suff != "" {
        T.Errorf("Did not correctly return the empty suffix")
    }

}

func TestStringer(T *testing.T) {
    var s interface{} = simpleString("abc")
    switch s.(type) {
    case Stringer:
        T.Logf("Object %#v implements Stringer as expected.", s)
    default:
        T.Errorf("Object %#v does not implement Stringer.", s)
    }
}

func TestMakeStringer(T *testing.T) {
    if s := simpleString("abc"); CausesPanic(func() { makeStringer(s) }) {
        T.Errorf("Error making Stringer out of %#v", s)
    }
    if s := "abc"; CausesPanic(func() { makeStringer(s) }) {
        T.Errorf("Error making Stringer out o %#v", s)
    }
    if s := 123; !CausesPanic(func() { makeStringer(s) }) {
        T.Errorf("Error Stringer created from %#v", s)
    }
}
