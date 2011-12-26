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
)

func testOptions(T *testing.T) {
    q := newQuestion(String)
    if q.Case != NilCase {
        T.Errorf("New String answer's case is not nil")
    }
}

func testGood(T *testing.T, q *Question, name, in string, v interface{}) error {
    t := q.Type()
    err := q.parse(in)
    if err != nil {
        T.Errorf("%s %s parse failed %s", name, t.String(), err.Error())
    } else if q.val != v {
        T.Errorf("Parsed value != expected value (%#v != %#v)", q.val, v)
    }
    return err
}
func testBad(T *testing.T, q *Question, name, in string) error {
    t := q.Type()
    err := q.parse(in)
    if err == nil {
        T.Errorf("%s improper %s parse succeded!", name, t.String())
    }
    if q.val != nil {
        T.Errorf("Parsed %s value != nil (%#v != nil)", t.String(), q.val)
    }
    return err
}

func TestQuestionInt(T *testing.T) {
    q := newQuestion(Int)
    // Standard parse tests
    testGood(T, q, "Simple", "1234", int64(1234))
    testBad(T, q, "Simple", "123JumpStreet")
    // Ranged parse tests
    q.In(IntRange{int64(-3), int64(10)})
    testGood(T, q, "In-range", "5", int64(5))
    testGood(T, q, "Edge-low", "-3", int64(-3))
    testGood(T, q, "Edge-high", "10", int64(10))
    testBad(T, q, "Low", "-4")
    testBad(T, q, "High", "11")
}

func TestQuestionUint(T *testing.T) {
    q := newQuestion(Uint)
    // Standard parse tests
    testGood(T, q, "Simple", "1234", uint64(1234))
    testGood(T, q, "Trimmed", " 4321  \n", uint64(4321))
    testBad(T, q, "Simple", "123JumpStreet")
    // Ranged parse tests
    q.In(UintRange{uint64(1), uint64(10)})
    testGood(T, q, "In-range", "5", uint64(5))
    testGood(T, q, "Edge-low", "1", uint64(1))
    testGood(T, q, "Edge-high", "10", uint64(10))
    testBad(T, q, "Low", "0")
    testBad(T, q, "High", "11")
}

func TestQuestionFloat(T *testing.T) {
    q := newQuestion(Float)
    // NOTE: Only use floating point number with a lossless 64bit representation.
    // Standard parse tests
    testGood(T, q, "Int", "1234", float64(1234))
    testGood(T, q, "Simple", "-12.75", float64(-12.75))
    testGood(T, q, "Scientific", "123.5e+1", float64(1235.0))
    testBad(T, q, "Complex", "123.5+i5")
    // Ranged parse tests
    q.In(FloatRange{float64(1), float64(10)})
    testGood(T, q, "In-range", "5", float64(5))
    testGood(T, q, "Edge-low", "1", float64(1))
    testGood(T, q, "Edge-high", "10", float64(10))
    testBad(T, q, "Low", "0")
    testBad(T, q, "High", "11")
}


func TestQuestionString(T *testing.T) {
    q := newQuestion(String)
    // Standard parse tests
    testGood(T, q, "Simple", "1234", "1234")
    q.Whitespace &= ^Trim
    testGood(T, q, "No trim default", " 4321  \n", " 4321  \n")
    q.Whitespace |= Collapse
    testGood(T, q, "Collapse", " 4321  \tabc  \n", " 4321 abc \n")
    // Ranged parse tests
    q.In(StringRange{"aaa", "zzz"})
    testGood(T, q, "In-range", "tidoeids", "tidoeids")
    testGood(T, q, "Edge-low", "aaa", "aaa")
    testGood(T, q, "Edge-high", "zzz", "zzz")
    testBad(T, q, "Low", "ZZZ")
    testBad(T, q, "High", "{")
    q.In(StringSet([]string{"abc", "def", "ghi"}))
    testGood(T, q, "In set", "def", "def")
    testBad(T, q, "Not in set", "blah")
}
