package goline
/*
 *  Filename:    set.go
 *  Package:     goline
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Fri Aug 19 03:15:13 PDT 2011
 *  Description: 
 */
import (
    "reflect"
    "fmt"
    "os"
)

func errorSetMemberType(set, member interface{}) os.Error {
    return fmt.Errorf("Set type %v cannot contain type %v",
        reflect.TypeOf(set).String(),
        reflect.TypeOf(member).String())
}

type AnswerSet interface {
    Has(x interface{}) bool
    String() string
}

type UintRange struct {
    Min, Max uint64
}

type IntRange struct {
    Min, Max int64
}

type FloatRange struct {
    Min, Max float64
}

type StringRange struct {
    Min, Max string
}

func (r UintRange) String() string   { return fmt.Sprintf("range [%v, %v]", r.Min, r.Max) }
func (r IntRange) String() string    { return fmt.Sprintf("range [%v, %v]", r.Min, r.Max) }
func (r FloatRange) String() string  { return fmt.Sprintf("range [%v, %v]", r.Min, r.Max) }
func (r StringRange) String() string { return fmt.Sprintf("range [%#v, %#v]", r.Min, r.Max) }

func (ur UintRange) Has(x interface{}) bool {
    switch x.(type) {
    case uint64:
        y := x.(uint64)
        return y >= ur.Min && y <= ur.Max
    }
    panic(errorSetMemberType(ur, x))
}

func (ur IntRange) Has(x interface{}) bool {
    switch x.(type) {
    case int64:
        y := x.(int64)
        return y >= ur.Min && y <= ur.Max
    }
    panic(errorSetMemberType(ur, x))
}

func (ur FloatRange) Has(x interface{}) bool {
    switch x.(type) {
    case float64:
        y := x.(float64)
        return y >= ur.Min && y <= ur.Max
    }
    panic(errorSetMemberType(ur, x))
}

func (ur StringRange) Has(x interface{}) bool {
    switch x.(type) {
    case string:
        y := x.(string)
        return y >= ur.Min && y <= ur.Max
    }
    panic(errorSetMemberType(ur, x))
}

type StringSet []string

func (set StringSet) Has(x interface{}) bool {
    switch x.(type) {
    case string:
        y := x.(string)
        for _, s := range set {
            if s == y {
                return true
            }
        }
        return false
    }
    panic(errorSetMemberType(set, x))
}

func (set StringSet) String() string {
    n := len(set)
    if n == 0 {
        return "{}"
    }
    length := 4 * n + 4
    for _, s := range set {
        length += len(s)
    }
    var j int
    p := make([]byte, length)
    j += copy(p, "set {")
    for i, s := range set {
        j += copy(p[j:], []byte{'"'})
        j += copy(p[j:], s)
        j += copy(p[j:], []byte{'"'})
        if i < n-1 {
            j += copy(p[j:], ", ")
        }
    }
    j += copy(p[j:], "}")
    return string(p)
}
