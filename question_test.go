package goline
/*
 *  Filename:    question_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Aug 13 02:30:29 PDT 2011
 *  Description: 
 *  Usage:       gotest
 */
import (
    "testing"
    "os"
)

func testOptions(T *testing.T) {
    a := newAnswer(String)
    if a.Case != NilCase {
        T.Errorf("New String answer's case is not nil")
    }
}

func testGood(T *testing.T, a *Answer, name, in string, v interface{}) os.Error {
    t := a.Type()
    err := a.parse(in)
    if err != nil {
        T.Errorf("%s %s parse failed %s", name, t.String(), err.String())
    } else if a.val != v {
        T.Errorf("Parsed value != expected value (%#v != %#v)", a.val, v)
    }
    return err
}
func testBad(T *testing.T, a *Answer, name, in string) os.Error {
    t := a.Type()
    err := a.parse(in)
    if err == nil {
        T.Errorf("%s improper %s parse succeded!", name, t.String())
    }
    if a.val != nil {
        T.Errorf("Parsed %s value != nil (%#v != nil)", t.String(), a.val)
    }
    return err
}

func TestAnswerInt(T *testing.T) {
    a := newAnswer(Int)
    // Standard parse tests
    testGood(T, a, "Simple", "1234", int64(1234))
    testBad(T, a, "Simple", "123JumpStreet")
    // Ranged parse tests
    a.In(IntRange{int64(-3), int64(10)})
    testGood(T, a, "In-range", "5", int64(5))
    testGood(T, a, "Edge-low", "-3", int64(-3))
    testGood(T, a, "Edge-high", "10", int64(10))
    testBad(T, a, "Low", "-4")
    testBad(T, a, "High", "11")
}

func TestAnswerUint(T *testing.T) {
    a := newAnswer(Uint)
    // Standard parse tests
    testGood(T, a, "Simple", "1234", uint64(1234))
    testGood(T, a, "Trimmed", " 4321  \n", uint64(4321))
    testBad(T, a, "Simple", "123JumpStreet")
    // Ranged parse tests
    a.In(UintRange{uint64(1), uint64(10)})
    testGood(T, a, "In-range", "5", uint64(5))
    testGood(T, a, "Edge-low", "1", uint64(1))
    testGood(T, a, "Edge-high", "10", uint64(10))
    testBad(T, a, "Low", "0")
    testBad(T, a, "High", "11")
}

func TestAnswerFloat(T *testing.T) {
    a := newAnswer(Float)
    // NOTE: Only use floating point number with a lossless 64bit representation.
    // Standard parse tests
    testGood(T, a, "Int", "1234", float64(1234))
    testGood(T, a, "Simple", "-12.75", float64(-12.75))
    testGood(T, a, "Scientific", "123.5e+1", float64(1235.0))
    testBad(T, a, "Complex", "123.5+i5")
    // Ranged parse tests
    a.In(FloatRange{float64(1), float64(10)})
    testGood(T, a, "In-range", "5", float64(5))
    testGood(T, a, "Edge-low", "1", float64(1))
    testGood(T, a, "Edge-high", "10", float64(10))
    testBad(T, a, "Low", "0")
    testBad(T, a, "High", "11")
}


func TestAnswerString(T *testing.T) {
    a := newAnswer(String)
    // Standard parse tests
    testGood(T, a, "Simple", "1234", "1234")
    a.Whitespace &= ^Trim
    testGood(T, a, "No trim default", " 4321  \n", " 4321  \n")
    a.Whitespace |= Collapse
    testGood(T, a, "Collapse", " 4321  \tabc  \n", " 4321 abc \n")
    // Ranged parse tests
    a.In(StringRange{"aaa", "zzz"})
    testGood(T, a, "In-range", "tidoeids", "tidoeids")
    testGood(T, a, "Edge-low", "aaa", "aaa")
    testGood(T, a, "Edge-high", "zzz", "zzz")
    testBad(T, a, "Low", "ZZZ")
    testBad(T, a, "High", "{")
    a.In(StringSet([]string{"abc", "def", "ghi"}))
    testGood(T, a, "In set", "def", "def")
    testBad(T, a, "Not in set", "blah")
}
