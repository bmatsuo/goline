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

// An interface for sets of values.
type AnswerSet interface {
    Has(x interface{}) bool
    String() string
}

type Direction uint

const (
    Above Direction = iota
    Below
)

var infty = []string{Above: "Infinity", Below: "-Infinity"}

func (d Direction) Infinity() string { return infty[d] }

// A range of uint64 values [Min, Max]
type UintRange struct {
    Min, Max uint64
}
// A range of int64 values [Min, Max]
type IntRange struct {
    Min, Max int64
}
// A range of float64 values [Min, Max]
type FloatRange struct {
    Min, Max float64
}
// A range of string values [Min, Max]
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
    length := 4*n + 4
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

type UintBounded struct {
    Direction
    X   uint64
}
type IntBounded struct {
    Direction
    X   int64
}
type FloatBounded struct {
    Direction
    X   float64
}
type StringBounded struct {
    Direction
    X   string
}

func (r UintBounded) String() string   {
    if r.Direction == Above {
        return fmt.Sprintf("range [%v, %s)", r.X, r.Infinity())
    }
    return fmt.Sprintf("range (%s, %v]", r.Infinity(), r.X)
}
func (r IntBounded) String() string    {
    if r.Direction == Above {
        return fmt.Sprintf("range [%v, %s)", r.X, r.Infinity())
    }
    return fmt.Sprintf("range (%s, %v]", r.Infinity(), r.X)
}
func (r FloatBounded) String() string  {
    if r.Direction == Above {
        return fmt.Sprintf("range [%v, %s)", r.X, r.Infinity())
    }
    return fmt.Sprintf("range (%s, %v]", r.Infinity(), r.X)
}
func (r StringBounded) String() string {
    if r.Direction == Above {
        return fmt.Sprintf("range [%v, %s)", r.X, r.Infinity())
    }
    return fmt.Sprintf("range (%#s, %v]", r.Infinity(), r.X)
}

func (r UintBounded) Has(x interface{}) bool {
    switch x.(type) {
    case uint64:
        y := x.(uint64)
        switch r.Direction {
        case Above:
            return y >= r.X
        case Below:
            return y <= r.X
        }
    }
    panic(errorSetMemberType(r, x))
}
func (r IntBounded) Has(x interface{}) bool {
    switch x.(type) {
    case int64:
        y := x.(int64)
        switch r.Direction {
        case Above:
            return y >= r.X
        case Below:
            return y <= r.X
        }
    }
    panic(errorSetMemberType(r, x))
}
func (r FloatBounded) Has(x interface{}) bool {
    switch x.(type) {
    case float64:
        y := x.(float64)
        switch r.Direction {
        case Above:
            return y >= r.X
        case Below:
            return y <= r.X
        }
    }
    panic(errorSetMemberType(r, x))
}
func (r StringBounded) Has(x interface{}) bool {
    switch x.(type) {
    case string:
        y := x.(string)
        switch r.Direction {
        case Above:
            return y >= r.X
        case Below:
            return y <= r.X
        }
    }
    panic(errorSetMemberType(r, x))
}

type UintBoundedStrictly UintBounded
type IntBoundedStrictly IntBounded
type FloatBoundedStrictly FloatBounded
type StringBoundedStrictly StringBounded

// TODO Fix the String method so that it returns open intervals

func (r UintBoundedStrictly) String() string   { return UintBounded(r).String() }
func (r IntBoundedStrictly) String() string    { return IntBounded(r).String() }
func (r FloatBoundedStrictly) String() string  { return FloatBounded(r).String() }
func (r StringBoundedStrictly) String() string { return StringBounded(r).String() }

func (r UintBoundedStrictly) Has(x interface{}) bool {
    switch x.(type) {
    case uint64:
        y := x.(uint64)
        switch r.Direction {
        case Above:
            return y > r.X
        case Below:
            return y < r.X
        }
    }
    panic(errorSetMemberType(r, x))
}
func (r IntBoundedStrictly) Has(x interface{}) bool {
    switch x.(type) {
    case int64:
        y := x.(int64)
        switch r.Direction {
        case Above:
            return y > r.X
        case Below:
            return y < r.X
        }
    }
    panic(errorSetMemberType(r, x))
}
func (r FloatBoundedStrictly) Has(x interface{}) bool {
    switch x.(type) {
    case float64:
        y := x.(float64)
        switch r.Direction {
        case Above:
            return y > r.X
        case Below:
            return y < r.X
        }
    }
    panic(errorSetMemberType(r, x))
}
func (r StringBoundedStrictly) Has(x interface{}) bool {
    switch x.(type) {
    case string:
        y := x.(string)
        switch r.Direction {
        case Above:
            return y > r.X
        case Below:
            return y < r.X
        }
    }
    panic(errorSetMemberType(r, x))
}
