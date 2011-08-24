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
    "regexp"
    "fmt"
    "os"
)

type WhitespaceOption uint

const (
    NilWhitespace WhitespaceOption = 0
    Strip         WhitespaceOption = 1 << iota
    Chomp
    Trim
    Collapse
    Remove
)

type CaseOption uint

const (
    NilCase CaseOption = iota
    Upper
    Lower
    // TODO implement Capitalize
    Capitalize
)

type Response uint

const (
    AskOnError Response = iota
    InvalidType
    //NoCompletion
    //AmbiguousCompletion
    NotInSet
    //NotValid
)

type Responses [3]string

var defaultResponses = Responses{
    AskOnError:  "Please retry:  ",
    InvalidType: "Type mismatch",
    //NoCompletion: "No auto-completion",
    //AmbiguousCompletion: "Ambiguous auto-completion",
    NotInSet: "Answer not contained in",
    //NotValid: "Answer did not pass validity test.",
}

func makeResponses() Responses {
    var r Responses
    copy(r[:], defaultResponses[:])
    return r
}

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
    // The "prompt" message for the user.
    Question string
    // Options for pre-processing scanned input whitespace
    Whitespace WhitespaceOption
    // Options for pre-processing scanned input symbol case
    Case CaseOption
    // A list of responses to various errors.
    Responses
    // Field separator for slice (list) inputs.
    FirstAnswer interface{}
    // The default value used when the user inputs an empty string.
    Default interface{}
    Sep     string
    set     AnswerSet
    typ     Type
    val     interface{}
    def     interface{}
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
    a.Responses = makeResponses()
    a.Default = nil
    a.FirstAnswer = nil
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

func (a *Answer) SetHas(x interface{}) bool {
    if a.set != nil {
        return a.set.Has(x)
    }
    return true
}

func (a *Answer) typeCast(v interface{}) (interface{}, os.Error) {
    var zero interface{}
    switch a.typ {
    case String:
        switch v.(type) {
        case string:
            return v, nil
        default:
            return zero, errorTypeError(a.Responses, "", v)
        }
    case Int:
        switch v.(type) {
        case int:
            return int64(v.(int)), nil
        case int8:
            return int64(v.(int8)), nil
        case int16:
            return int64(v.(int16)), nil
        case int32:
            return int64(v.(int32)), nil
        case int64:
            return int64(v.(int64)), nil
        default:
            return zero, errorTypeError(a.Responses, int64(1), v)
        }
    case Uint:
        switch v.(type) {
        case uint:
            return uint64(v.(uint)), nil
        case uint8:
            return uint64(v.(uint8)), nil
        case uint16:
            return uint64(v.(uint16)), nil
        case uint32:
            return uint64(v.(uint32)), nil
        case uint64:
            return v.(uint64), nil
        default:
            return zero, errorTypeError(a.Responses, uint64(1), v)
        }
    case Float:
        switch v.(type) {
        case float32:
            return float64(v.(float32)), nil
        case float64:
            return v.(float64), nil
        default:
            return zero, errorTypeError(a.Responses, float64(1), v)
        }
    case StringSlice:
        fallthrough
    case IntSlice:
        fallthrough
    case UintSlice:
        fallthrough
    case FloatSlice:
        return zero, fmt.Errorf("%s unimplemented", a.typ.String())
    }
    return zero, fmt.Errorf("%s unimplemented", a.typ.String())
}

func (a *Answer) tryFirstAnswer() os.Error {
    if a.FirstAnswer != nil {
        if val, err := a.typeCast(a.FirstAnswer); err != nil {
            return err
        } else {
            a.FirstAnswer = val
        }
    }
    return nil
}

func (a *Answer) DefaultString() string {
    if a.Default != nil {
        return fmt.Sprintf("|%v|  ", a.Default)
    }
    return ""
}
func (a *Answer) tryDefault() (val interface{}, err os.Error) {
    if a.Default != nil {
        if val, err = a.typeCast(a.Default); err != nil {
            return
        } else {
            return
        }
    }
    val = nil
    return
}
/*
func (a *Answer) Default() interface{} { return a.def }
func (a *Answer) SetDefault(v interface{}) os.Error {
    if def, err := a.typeCast(v); err != nil {
        return err
    } else {
        a.def = def
    }
    return nil
}
*/

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

var spaceRE = regexp.MustCompile("[ \t]+")

func (a *Answer) parse(in string) os.Error {
    a.val = nil // Clear the parse value for good measure.
    // Perform all pre-processing on the input.
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
    switch a.Case {
    case Upper:
        in = strings.ToUpper(in)
    case Lower:
        in = strings.ToLower(in)
    case Capitalize:
        in = strings.ToLower(in)
    }

    // Handle the default value if necessary.
    var (
        def interface{}
        err os.Error
    )
    noInput := len(in) == 0
    useDefault := noInput && a.Default != nil
    if useDefault {
        if def, err = a.tryDefault(); err != nil {
            return err
        }
    }

    // Parse the user's input.
    var val interface{}
    switch a.typ {
    case String:
        if useDefault {
            val = def.(string)
        } else {
            val = in
        }
    case Int:
        var x int64
        if useDefault {
            x = def.(int64)
        } else if noInput {
            return ErrorEmptyInput(0)
        } else if x, err = strconv.Atoi64(in); err != nil {
            return err
        }
        val = x
    case Uint:
        var x uint64
        if useDefault {
            x = def.(uint64)
        } else if noInput {
            return ErrorEmptyInput(0)
        } else if x, err = strconv.Atoui64(in); err != nil {
            return err
        }
        val = x
    case Float:
        var x float64
        if useDefault {
            x = def.(float64)
        } else if noInput {
            return ErrorEmptyInput(0)
        } else if x, err = strconv.Atof64(in); err != nil {
            return err
        }
        val = x
    case StringSlice:
        fallthrough
    case IntSlice:
        fallthrough
    case UintSlice:
        fallthrough
    case FloatSlice:
        err = fmt.Errorf("%s unimplemented", a.typ.String())
    }

    // Check set membership
    switch a.typ {
    case String:
        fallthrough
    case Int:
        fallthrough
    case Uint:
        fallthrough
    case Float:
        if !a.SetHas(val) {
            return a.makeErrorNotInSet(a.Responses, val)
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

    // Set the parsed value if there was no error.
    if err == nil {
        a.val = val
    }

    return err
}
