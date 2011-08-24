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
    // Perform no whitespace preprocessing.
    NilWhitespace WhitespaceOption = 0
    // Remove the trailing newline from input.
    Chomp WhitespaceOption = 1 << iota
    // Remove all leading and trailing whitespace from the input. See unicode.IsSpace.
    Trim
    // Replace all runs of whitespace with a single space character ' '.
    Collapse
    // Remove all whitespace from string input.
    Remove
)

type CaseOption uint

const (
    // Perform no case adjustment.
    NilCase CaseOption = iota
    // Convert string input to its upper case equivalent.
    Upper
    // Convert string input to its lower case equivalent.
    Lower
    // Capitalize the first character of input (after whitespace processing) (TODO)
    Capitalize
)

type Response uint

const (
    // A question that prompts the user when an error was encountered.
    AskOnError Response = iota
    // The error message printed when an Answer's FirstAnswer or Default has
    // a type incompatible with its Type.
    InvalidType
    // The error message printed when the parsed input was not in the Answer's
    // AnswerSet.
    NotInSet
    // The error message printed when the answer did not pass any validity
    // test.
    //NotValid
    // The error message printed when there are no auto-completion results.
    //NoCompletion
    // The error message printed when auto-completion is ambiguous.
    //AmbiguousCompletion
)

type Responses [3]string

var defaultResponses = Responses{
    AskOnError:  "Please retry:  ",
    InvalidType: "Type mismatch",
    NotInSet: "Answer not contained in",
    //NotValid: "Answer did not pass validity test.",
    //NoCompletion: "No auto-completion",
    //AmbiguousCompletion: "Ambiguous auto-completion",
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

func (t Type) String() string    { return tstring[t] }
func (t Type) IsSliceType() bool { return t >= StringSlice }

type Answer struct {
    // The "prompt" message for the user.
    Question string
    // Pre-processing options for whitespace. See WhitespaceOption.
    Whitespace WhitespaceOption
    // Pre-processing options string case. See CaseOption.
    Case CaseOption
    // A set of responses to various errors. See Response.
    Responses
    // If this field is not nil then then Ask() will assign it to the
    // destination variable without prompting the user; provided that it
    // satisfies all the validity and set membership tests of the Answer.
    FirstAnswer interface{}
    // The default value used when the user inputs an empty string.
    Default interface{}
    // Separator for list (slice) input (TODO)
    Sep string
    // Called when an error forces the prompt to halt without a value.
    Panic func(os.Error)
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

func (a *Answer) setHas(x interface{}) bool {
    if a.set != nil {
        return a.set.Has(x)
    }
    return true
}

func (a *Answer) typeCast(v interface{}) (val interface{}, err os.Error) {
    switch a.typ {
    case String:
        switch v.(type) {
        case string:
            val = v
        default:
            err = a.makeTypeError("", v)
        }
    case Int:
        switch v.(type) {
        case int:
            val = int64(v.(int))
        case int8:
            val = int64(v.(int8))
        case int16:
            val = int64(v.(int16))
        case int32:
            val = int64(v.(int32))
        case int64:
            val = int64(v.(int64))
        default:
            err = a.makeTypeError(int64(1), v)
        }
    case Uint:
        switch v.(type) {
        case uint:
            val = uint64(v.(uint))
        case uint8:
            val = uint64(v.(uint8))
        case uint16:
            val = uint64(v.(uint16))
        case uint32:
            val = uint64(v.(uint32))
        case uint64:
            val = v.(uint64)
        default:
            err = a.makeTypeError(uint64(1), v)
        }
    case Float:
        switch v.(type) {
        case float32:
            val = float64(v.(float32))
        case float64:
            val = v.(float64)
        default:
            err = a.makeTypeError(float64(1), v)
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
    return
}

func (a *Answer) tryFirstAnswer() os.Error {
    if a.FirstAnswer != nil {
        if val, err := a.typeCast(a.FirstAnswer); err != nil {
            return err
        } else {
            a.val = val
        }
    }
    return nil
}

func (a *Answer) defaultString(suffix string) string {
    if a.Default != nil {
        return fmt.Sprintf("|%v|%s", a.Default, suffix)
    }
    return ""
}
func (a *Answer) tryDefault() (val interface{}, err os.Error) {
    val = nil
    if a.Default != nil {
        return a.typeCast(a.Default)
    }
    return
}

//  Specify a set of answers in which the response much be contained.
func (a *Answer) In(s AnswerSet) { a.set = s }

//  Returns the Type which is enforced by the Answer.
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
            return ErrorParse{in,err}
        }
        val = x
    case Uint:
        var x uint64
        if useDefault {
            x = def.(uint64)
        } else if noInput {
            return ErrorEmptyInput(0)
        } else if x, err = strconv.Atoui64(in); err != nil {
            return ErrorParse{in,err}
        }
        val = x
    case Float:
        var x float64
        if useDefault {
            x = def.(float64)
        } else if noInput {
            return ErrorEmptyInput(0)
        } else if x, err = strconv.Atof64(in); err != nil {
            return ErrorParse{in,err}
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
        if !a.setHas(val) {
            return a.makeErrorNotInSet(val)
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

// Cast a value result from a wide (e.g. 64bit) type to the desired type.
// This should not fail under any normal circumstances, so failure
// should break the loop.
func (a *Answer) setDest(dest interface{}) os.Error {
    switch a.Type() {
    case Uint:
        switch dest.(type) {
        case *uint:
            d := dest.(*uint)
            *(d) = uint(a.val.(uint64))
            if x := uint64(*(d)); x != a.val.(uint64) {
                return ErrorPrecision{a.val.(uint64), x}
            }
        case *uint8:
            d := dest.(*uint8)
            *(d) = uint8(a.val.(uint64))
            if x := uint64(*(d)); x != a.val.(uint64) {
                return ErrorPrecision{a.val.(uint64), x}
            }
        case *uint16:
            d := dest.(*uint16)
            *(d) = uint16(a.val.(uint64))
            if x := uint64(*(d)); x != a.val.(uint64) {
                return ErrorPrecision{a.val.(uint64), x}
            }
        case *uint32:
            d := dest.(*uint32)
            *(d) = uint32(a.val.(uint64))
            if x := uint64(*(d)); x != a.val.(uint64) {
                return ErrorPrecision{a.val.(uint64), x}
            }
        case *uint64:
            *(dest.(*uint64)) = a.val.(uint64)
        default:
            return fmt.Errorf("Unexpected cast type")
        }
    case Int:
        switch dest.(type) {
        case *int:
            d := dest.(*int)
            *(d) = int(a.val.(int64))
            if x := int64(*(d)); x != a.val.(int64) {
                return ErrorPrecision{a.val.(int64), x}
            }
        case *int8:
            d := dest.(*int8)
            *(d) = int8(a.val.(int64))
            if x := int64(*(d)); x != a.val.(int64) {
                return ErrorPrecision{a.val.(int64), x}
            }
        case *int16:
            d := dest.(*int16)
            *(d) = int16(a.val.(int64))
            if x := int64(*(d)); x != a.val.(int64) {
                return ErrorPrecision{a.val.(int64), x}
            }
        case *int32:
            d := dest.(*int32)
            *(d) = int32(a.val.(int64))
            if x := int64(*(d)); x != a.val.(int64) {
                return ErrorPrecision{a.val.(int64), x}
            }
        case *int64:
            *(dest.(*int64)) = a.val.(int64)
        default:
            return fmt.Errorf("Unexpected cast type")
        }
    case Float:
        switch dest.(type) {
        case *float32:
            d := dest.(*float32)
            *(d) = float32(a.val.(float64))
            if x := float64(*(d)); x != a.val.(float64) {
                return ErrorPrecision{a.val.(float64), x}
            }
        case *float64:
            *(dest.(*float64)) = a.val.(float64)
        default:
            return fmt.Errorf("Unexpected cast type")
        }
    case String:
        switch dest.(type) {
        case *string:
            *(dest.(*string)) = a.val.(string)
        default:
            return fmt.Errorf("Unexpected cast type")
        }
    }
    return nil
}
