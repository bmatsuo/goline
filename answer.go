package goline
/*
 *  Filename:    answer.go
 *  Package:     goline
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Aug 13 02:33:00 PDT 2011
 *  Description: 
 */
import (
    "strings"
    "strconv"
    "reflect"
    "regexp"
    "fmt"
    "os"
)

type Options uint

const (
    Trim Options = 1 << iota
    Collapse
)

//  64-bit types are always used to read data. It is then cast to a thinner
//  type dynamically with the "reflect" package.
type Type uint

const (
    Int Type = iota
    Uint
    Float
    String
    // Slice types are not implemented yet
    StringSlice
    UintSlice
    IntSlice
    FloatSlice
)

var tstring = []string{
    Int:         "Int",
    Uint:        "Uint",
    Float:       "Float",
    String:      "String",
    StringSlice: "StringSlice",
    UintSlice:   "UintSlice",
    IntSlice:    "IntSlice",
    FloatSlice:  "FloatSlice",
}

func (t Type) String() string { return tstring[t] }

func (t Type) IsSliceType() bool {
    return t >= StringSlice
}

type Answer struct {
    Opt     Options
    Set     AnswerSet
    t       Type
    v       interface{}
    d       interface{}
    inRange bool
    umin    uint64
    umax    uint64
    fmin    float64
    fmax    float64
    imin    int64
    imax    int64
    smin    string
    smax    string
    sep     string
}

func newAnswer(t Type) *Answer {
    a := new(Answer)
    a.t = t
    a.d = nil
    switch a.t {
    //case String:
    case Int:
        fallthrough
    case Uint:
        fallthrough
    case Float:
        fallthrough
    case StringSlice:
        fallthrough
    case IntSlice:
        fallthrough
    case UintSlice:
        fallthrough
    case FloatSlice:
        a.Opt = Trim | Collapse
    }
    a.sep = " "
    a.Set = nil
    return a
}

func errorEmptyRange(min, max interface{}) os.Error {
    return fmt.Errorf("Range max is less than min (%v < %v)", min, max)
}
func errorTypeError(expect, recv interface{}) os.Error {
    return fmt.Errorf("Received type %s not equal to expected type %s",
        reflect.ValueOf(recv).Kind().String(), reflect.ValueOf(expect).Kind().String())
}

func (a *Answer) SetHas(x interface{}) bool {
    if a.Set != nil {
        return a.Set.Has(x)
    }
    return true
}

func (a *Answer) Default() interface{} { return a.d }
func (a *Answer) DefaultString() string {
    if a.d != nil {
        return fmt.Sprintf("|%v|  ", a.d)
    }
    return ""
}
func (a *Answer) SetDefault(v interface{}) os.Error {
    var err os.Error
    switch a.t {
    case String:
        switch v.(type) {
        case string:
            a.d = v
        default:
            errorTypeError("", v)
        }
    case Int:
        switch v.(type) {
        case int:
            a.d = int64(v.(int))
        case int8:
            a.d = int64(v.(int8))
        case int16:
            a.d = int64(v.(int16))
        case int32:
            a.d = int64(v.(int32))
        case int64:
            a.d = int64(v.(int64))
        default:
            errorTypeError(int64(1), v)
        }
    case Uint:
        switch v.(type) {
        case uint:
            a.d = uint64(v.(uint))
        case uint8:
            a.d = uint64(v.(uint8))
        case uint16:
            a.d = uint64(v.(uint16))
        case uint32:
            a.d = uint64(v.(uint32))
        case uint64:
            a.d = v.(uint64)
        default:
            errorTypeError(uint64(1), v)
        }
    case Float:
        switch v.(type) {
        case float32:
            a.d = float64(v.(float32))
        case float64:
            a.d = v.(float64)
        default:
            errorTypeError(float64(1), v)
        }
    case StringSlice:
        fallthrough
    case IntSlice:
        fallthrough
    case UintSlice:
        fallthrough
    case FloatSlice:
        err = fmt.Errorf("%s unimplemented", a.t.String())
    }
    return err
}

func (a *Answer) InRange(min, max interface{}) {
    switch a.t {
    case String:
        fallthrough
    case StringSlice:
        a.smin = min.(string)
        a.smax = max.(string)
        if a.smax < a.smin {
            panic(errorEmptyRange(min, max))
        }
    case Int:
        fallthrough
    case IntSlice:
        a.imin = min.(int64)
        a.imax = max.(int64)
        if a.imax < a.imin {
            panic(errorEmptyRange(min, max))
        }
    case Uint:
        fallthrough
    case UintSlice:
        a.umin = min.(uint64)
        a.umax = max.(uint64)
        if a.umax < a.umin {
            panic(errorEmptyRange(min, max))
        }
    case Float:
        fallthrough
    case FloatSlice:
        a.fmin = min.(float64)
        a.fmax = max.(float64)
        if a.fmax < a.fmin {
            panic(errorEmptyRange(min, max))
        }
    }
    a.inRange = true
}

func (a *Answer) Type() Type { return a.t }

type RecoverableError interface {
    os.Error
    IsRecoverable() bool
}

type ErrorNotInSet struct{ os.Error }

func (err ErrorNotInSet) IsRecoverable() bool { return true }

func (a *Answer) makeErrorNotInSet(val interface{}) ErrorNotInSet {
    return ErrorNotInSet{
        fmt.Errorf("Value %v not in set %s", val, a.Set.String())}
}

type ErrorOutOfRange struct {
    value, min, max interface{}
    err             os.Error
}

func (oor ErrorOutOfRange) String() string      { return oor.err.String() }
func (oor ErrorOutOfRange) IsRecoverable() bool { return true }

func errorOOR(val, min, max interface{}) ErrorOutOfRange {
    return ErrorOutOfRange{min, max, val,
        fmt.Errorf("Value %v out of range [%v, %v]", val, min, max),
    }
}

type ErrorEmptyInput uint

func (oor ErrorEmptyInput) String() string      { return "Can not use empty value" }
func (oor ErrorEmptyInput) IsRecoverable() bool { return true }

var spaceRE = regexp.MustCompile("[ \t]+")

func (a *Answer) parse(in string) os.Error {
    if a.Opt&Trim > 0 {
        in = strings.TrimSpace(in)
    }
    if a.Opt&Collapse > 0 {
        in = spaceRE.ReplaceAllString(in, " ")
    }
    a.v = a.d
    noInput := len(in) == 0
    useDefault := noInput && a.d != nil
    var err os.Error
    switch a.t {
    case String:
        if useDefault {
            in = a.d.(string)
        }
        if !a.SetHas(in) {
            return a.makeErrorNotInSet(in)
        }
        a.v = in
    case Int:
        var x int64
        if useDefault {
            x = a.d.(int64)
        } else if noInput {
            return ErrorEmptyInput(0)
        } else if x, err = strconv.Atoi64(in); err != nil {
            return err
        }
        if !a.SetHas(x) {
            return a.makeErrorNotInSet(x)
        }
        a.v = x
    case Uint:
        var x uint64
        if useDefault {
            x = a.d.(uint64)
        } else if noInput {
            return ErrorEmptyInput(0)
        } else if x, err = strconv.Atoui64(in); err != nil {
            return err
        }
        if !a.SetHas(x) {
            return a.makeErrorNotInSet(x)
        }
        a.v = x
    case Float:
        var x float64
        if useDefault {
            x = a.d.(float64)
        } else if noInput {
            return ErrorEmptyInput(0)
        } else if x, err = strconv.Atof64(in); err != nil {
            return err
        }
        if !a.SetHas(x) {
            return a.makeErrorNotInSet(x)
        }
        a.v = x
    case StringSlice:
        fallthrough
    case IntSlice:
        fallthrough
    case UintSlice:
        fallthrough
    case FloatSlice:
        err = fmt.Errorf("%s unimplemented", a.t.String())
    }
    return err
}
