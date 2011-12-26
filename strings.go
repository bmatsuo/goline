package goline
/*
 *  Filename:    strings.go
 *  Package:     goline
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Wed Aug 24 22:56:19 PDT 2011
 *  Description: 
 */
import (
    "reflect"
    "strings"
)

//  Returns the index i of the longest terminal substring s[i:] such that f
//  returns true for all runes in s[i:]. Returns -1 if there is no such i.
func stringSuffixIndexFunc(s string, f func(c rune) bool) (i int) {
    var hasSuffix bool
    i = strings.LastIndexFunc(s, func(c rune) (done bool) {
        if done = !f(c); !hasSuffix {
            hasSuffix = !done
        }
        return
    })
    if i++; !hasSuffix {
        i = -1
    }
    return
}

//  Return the suffix string corresponding to the same call to
//  stringSuffixIndexFunc.
func stringSuffixFunc(s string, f func(c rune) bool) (suff string) {
    if i := stringSuffixIndexFunc(s, f); i >= 0 {
        suff = s[i:]
    }
    return
}

//  A string type that implements Stringer.
type simpleString string

//  The empty string.
var zeroSimpleString simpleString

//  Simply recast the simpleString.
func (s simpleString) String() string { return string(s) }

//  An interface that Menu items must implement..
type Stringer interface {
    String() string
}

var (
    //  The zero Stringer value.
    zeroStringer Stringer
    //  The Stringer reflect type.
    typeStringer = reflect.TypeOf(Stringer(zeroSimpleString))
)

//  Return a Stringer object given either a string or an object that implements
//  Stringer.
func makeStringer(s interface{}) Stringer {
    switch s.(type) {
    case string:
        return simpleString(s.(string))
    case Stringer:
        return s.(Stringer)
    default:
        panic("Value must be type 'string' or 'Stringer'")
    }
    return zeroStringer
}
