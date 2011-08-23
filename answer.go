package goline
/*
 *  Filename:    answer.go
 *  Package:     goline
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Aug 13 02:33:00 PDT 2011
 *  Description: 
 */
import (
    "unicode"
    "strings"
    "strconv"
    "reflect"
    "regexp"
    "fmt"
    "os"
)

type WhitespaceOptions uint

const (
    NilWhitespace WhitespaceOptions = 1 << iota
    Strip
    Chomp
    Trim
    Collapse
    Remove
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
    // Options for pre-processing scanned input
    Whitespace WhitespaceOptions
    // Field separator for slice (list) inputs.
    Sep string
    set AnswerSet
    typ Type
    val interface{}
    def interface{}
    /*
       inRange bool
       umin    uint64
       umax    uint64
       fmin    float64
       fmax    float64
       imin    int64
       imax    int64
       smin    string
       smax    string
    */
}

func newAnswer(t Type) *Answer {
    a := new(Answer)
    a.typ = t
    a.def = nil
    switch a.typ {
    case String:
        a.Whitespace = Trim
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
        a.Whitespace = Trim | Collapse
    }
    a.Sep = " "
    a.set = nil
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
    if a.set != nil {
        return a.set.Has(x)
    }
    return true
}

func (a *Answer) Default() interface{} { return a.def }
func (a *Answer) DefaultString() string {
    if a.def != nil {
        return fmt.Sprintf("|%v|  ", a.def)
    }
    return ""
}
func (a *Answer) SetDefault(v interface{}) os.Error {
    var err os.Error
    switch a.typ {
    case String:
        switch v.(type) {
        case string:
            a.def = v
        default:
            errorTypeError("", v)
        }
    case Int:
        switch v.(type) {
        case int:
            a.def = int64(v.(int))
        case int8:
            a.def = int64(v.(int8))
        case int16:
            a.def = int64(v.(int16))
        case int32:
            a.def = int64(v.(int32))
        case int64:
            a.def = int64(v.(int64))
        default:
            errorTypeError(int64(1), v)
        }
    case Uint:
        switch v.(type) {
        case uint:
            a.def = uint64(v.(uint))
        case uint8:
            a.def = uint64(v.(uint8))
        case uint16:
            a.def = uint64(v.(uint16))
        case uint32:
            a.def = uint64(v.(uint32))
        case uint64:
            a.def = v.(uint64)
        default:
            errorTypeError(uint64(1), v)
        }
    case Float:
        switch v.(type) {
        case float32:
            a.def = float64(v.(float32))
        case float64:
            a.def = v.(float64)
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
        err = fmt.Errorf("%s unimplemented", a.typ.String())
    }
    return err
}

func (a *Answer) In(s AnswerSet) { a.set = s }

/*
func (a *Answer) InRange(min, max interface{}) {
    switch a.typ {
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
*/

func (a *Answer) Type() Type { return a.typ }

type RecoverableError interface {
    os.Error
    IsRecoverable() bool
}

type ErrorNotInSet struct{ os.Error }

func (err ErrorNotInSet) IsRecoverable() bool { return true }

func (a *Answer) makeErrorNotInSet(val interface{}) ErrorNotInSet {
    return ErrorNotInSet{
        fmt.Errorf("Value %v not in set %s", val, a.set.String())}
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
    if a.Whitespace&Remove > 0 {
        in = strings.Join(strings.FieldsFunc(in, unicode.IsSpace), "")
    } else {
        if a.Whitespace&Chomp > 0 {
            if in[len(in)-1] == '\r' {
                in = in[:len(in)]
            }
        }
        if a.Whitespace&Trim > 0 {
            in = strings.TrimSpace(in)
        }
        if a.Whitespace&Collapse > 0 {
            in = spaceRE.ReplaceAllString(in, " ")
        }
    }
    a.val = a.def
    noInput := len(in) == 0
    useDefault := noInput && a.def != nil
    var err os.Error
    switch a.typ {
    case String:
        if useDefault {
            in = a.def.(string)
        }
        if !a.SetHas(in) {
            return a.makeErrorNotInSet(in)
        }
        a.val = in
    case Int:
        var x int64
        if useDefault {
            x = a.def.(int64)
        } else if noInput {
            return ErrorEmptyInput(0)
        } else if x, err = strconv.Atoi64(in); err != nil {
            return err
        }
        if !a.SetHas(x) {
            return a.makeErrorNotInSet(x)
        }
        a.val = x
    case Uint:
        var x uint64
        if useDefault {
            x = a.def.(uint64)
        } else if noInput {
            return ErrorEmptyInput(0)
        } else if x, err = strconv.Atoui64(in); err != nil {
            return err
        }
        if !a.SetHas(x) {
            return a.makeErrorNotInSet(x)
        }
        a.val = x
    case Float:
        var x float64
        if useDefault {
            x = a.def.(float64)
        } else if noInput {
            return ErrorEmptyInput(0)
        } else if x, err = strconv.Atof64(in); err != nil {
            return err
        }
        if !a.SetHas(x) {
            return a.makeErrorNotInSet(x)
        }
        a.val = x
    case StringSlice:
        fallthrough
    case IntSlice:
        fallthrough
    case UintSlice:
        fallthrough
    case FloatSlice:
        err = fmt.Errorf("%s unimplemented", a.typ.String())
    }
    return err
}
