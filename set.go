package goline
/*
 *  Filename:    set.go
 *  Package:     goline
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Fri Aug 19 03:15:13 PDT 2011
 *  Description: 
 */
import (
    "strings"
    "fmt"
)

// An interface for sets of values.
type AnswerSet interface {
    Has(x interface{}) bool
    String() string
}

//  Composite answer sets (AnswerSetUnion and AnswerSetIntersection objects)
//  are tools available to create more complex AnswerSet objects not provided
//  by goline. For example open, or half-open, intervals can be constructed by
//  taking the intersection of two bounded answer sets.
type CompositeAnswerSet interface {
    AnswerSet
    Size() int
    Set(i int) AnswerSet
}

func compositeString(composite CompositeAnswerSet, name string) string {
    strs := make([]string, composite.Size())
    for i := range strs {
        strs[i] = composite.Set(i).String()
    }
    return fmt.Sprintf("%s of %s", strings.Join(strs, ", and "))
}

//  When making set intersections it is much easier to unknowingly create
//  empty AnswerSets that cannot be detected effectively. Exercise extreme
//  caution if creating set intersetions dynamically.
type AnswerSetIntersection []AnswerSet
type AnswerSetUnion []AnswerSet

var (
    Universe = AnswerSetIntersection{}
    EmptySet = AnswerSetUnion{}
)

func (set AnswerSetIntersection) Size() int { return len(set) }
func (set AnswerSetUnion) Size() int        { return len(set) }

func (set AnswerSetIntersection) Set(i int) AnswerSet { return set[i] }
func (set AnswerSetUnion) Set(i int) AnswerSet        { return set[i] }

func (set AnswerSetIntersection) String() string { return compositeString(set, "intersection") }
func (set AnswerSetUnion) String() string        { return compositeString(set, "union") }

func (set AnswerSetIntersection) Has(x interface{}) bool {
    for i := range set {
        if !set[i].Has(x) {
            return false
        }
    }
    return true
}
func (set AnswerSetUnion) Has(x interface{}) bool {
    for i := range set {
        if set[i].Has(x) {
            return true
        }
    }
    return false
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
    panic(makeErrorMemberType(ur, x))
}
func (ur IntRange) Has(x interface{}) bool {
    switch x.(type) {
    case int64:
        y := x.(int64)
        return y >= ur.Min && y <= ur.Max
    }
    panic(makeErrorMemberType(ur, x))
}
func (ur FloatRange) Has(x interface{}) bool {
    switch x.(type) {
    case float64:
        y := x.(float64)
        return y >= ur.Min && y <= ur.Max
    }
    panic(makeErrorMemberType(ur, x))
}
func (ur StringRange) Has(x interface{}) bool {
    switch x.(type) {
    case string:
        y := x.(string)
        return y >= ur.Min && y <= ur.Max
    }
    panic(makeErrorMemberType(ur, x))
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
    panic(makeErrorMemberType(set, x))
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

//  An interval with only one bound, X.
type UintBounded struct {
    Direction
    X   uint64
}
//  An interval with only one bound, X.
type IntBounded struct {
    Direction
    X   int64
}
//  An interval with only one bound, X.
type FloatBounded struct {
    Direction
    X   float64
}
//  An interval with only one bound, X.
type StringBounded struct {
    Direction
    X   string
}

func (r UintBounded) String() string {
    if r.Direction == Above {
        return fmt.Sprintf("range [%v, %s)", r.X, r.Infinity())
    }
    return fmt.Sprintf("range (%s, %v]", r.Infinity(), r.X)
}
func (r IntBounded) String() string {
    if r.Direction == Above {
        return fmt.Sprintf("range [%v, %s)", r.X, r.Infinity())
    }
    return fmt.Sprintf("range (%s, %v]", r.Infinity(), r.X)
}
func (r FloatBounded) String() string {
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
    panic(makeErrorMemberType(r, x))
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
    panic(makeErrorMemberType(r, x))
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
    panic(makeErrorMemberType(r, x))
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
    panic(makeErrorMemberType(r, x))
}

//  In interval strictly bounded by a single number.
type UintBoundedStrictly UintBounded
//  In interval strictly bounded by a single number.
type IntBoundedStrictly IntBounded
//  In interval strictly bounded by a single number.
type FloatBoundedStrictly FloatBounded
//  In interval strictly bounded by a single number.
type StringBoundedStrictly StringBounded

// TODO: Fix the String method so that it returns open intervals
func (r UintBoundedStrictly) String() string   { return UintBounded(r).String() }
// TODO: Fix the String method so that it returns open intervals
func (r IntBoundedStrictly) String() string    { return IntBounded(r).String() }
// TODO: Fix the String method so that it returns open intervals
func (r FloatBoundedStrictly) String() string  { return FloatBounded(r).String() }
// TODO: Fix the String method so that it returns open intervals
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
    panic(makeErrorMemberType(r, x))
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
    panic(makeErrorMemberType(r, x))
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
    panic(makeErrorMemberType(r, x))
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
    panic(makeErrorMemberType(r, x))
}
