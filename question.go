package goline
/*
 *  Filename:    question.go
 *  Package:     goline
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Aug 13 02:30:29 PDT 2011
 *  Description: 
 */
import (
    "reflect"
    "strings"
    "unicode"
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
    NotInSet:    "Answer not contained in",
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

func TypeOf(v interface{}) (typ Type, err os.Error) {
    switch v.(type) {
    case uint:
        typ = Uint
    case uint8:
        typ = Uint
    case uint16:
        typ = Uint
    case uint32:
        typ = Uint
    case uint64:
        typ = Uint
    case int:
        typ = Int
    case int8:
        typ = Int
    case int16:
        typ = Int
    case int32:
        typ = Int
    case int64:
        typ = Int
    case float32:
        typ = Float
    case float64:
        typ = Float
    case string:
        typ = String
    default:
        err = fmt.Errorf("Unrecognizable type %s", reflect.TypeOf(v).Name())
    }
    return
}

type Question struct {
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
    set   AnswerSet
    typ   Type
    val   interface{}
    def   interface{}
}

func newQuestion(t Type) *Question {
    q := new(Question)
    q.typ = t
    q.Responses = makeResponses()
    q.Default = nil
    q.FirstAnswer = nil
    switch q.typ {
    case String:
        q.Whitespace = Trim
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
        q.Whitespace = Trim | Collapse
    }
    q.Sep = " "
    q.set = nil
    return q
}

func (q *Question) setHas(x interface{}) bool {
    if q.set != nil {
        return q.set.Has(x)
    }
    return true
}

func (q *Question) typeCast(v interface{}) (val interface{}, err os.Error) {
    switch q.typ {
    case String:
        switch v.(type) {
        case string:
            val = v
        default:
            err = q.makeTypeError("", v)
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
            err = q.makeTypeError(int64(1), v)
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
            err = q.makeTypeError(uint64(1), v)
        }
    case Float:
        switch v.(type) {
        case float32:
            val = float64(v.(float32))
        case float64:
            val = v.(float64)
        default:
            err = q.makeTypeError(float64(1), v)
        }
    case StringSlice:
        fallthrough
    case IntSlice:
        fallthrough
    case UintSlice:
        fallthrough
    case FloatSlice:
        err = fmt.Errorf("%s unimplemented", q.typ.String())
    }
    return
}

func (q *Question) tryFirstAnswer() os.Error {
    if q.FirstAnswer != nil {
        if val, err := q.typeCast(q.FirstAnswer); err != nil {
            return err
        } else {
            q.val = val
        }
    }
    return nil
}

func (q *Question) defaultString(suffix string) string {
    if q.Default != nil {
        return fmt.Sprintf("|%v|%s", q.Default, suffix)
    }
    return ""
}
func (q *Question) tryDefault() (val interface{}, err os.Error) {
    val = nil
    if q.Default != nil {
        return q.typeCast(q.Default)
    }
    return
}

//  Specify a set of answers in which the response much be contained.
func (q *Question) In(s AnswerSet) { q.set = s }

//  Returns the Type which is enforced by the Answer.
func (q *Question) Type() Type { return q.typ }

var spaceRE = regexp.MustCompile("[ \t]+")

func (q *Question) parse(in string) os.Error {
    q.val = nil // Clear the parse value for good measure.
    // Perform all pre-processing on the input.
    if q.Whitespace&Remove > 0 {
        in = strings.Join(strings.FieldsFunc(in, unicode.IsSpace), "")
    } else {
        if q.Whitespace&Chomp > 0 {
            if in[len(in)-1] == '\r' {
                in = in[:len(in)]
            }
        }
        if q.Whitespace&Trim > 0 {
            in = strings.TrimSpace(in)
        }
        if q.Whitespace&Collapse > 0 {
            in = spaceRE.ReplaceAllString(in, " ")
        }
    }
    switch q.Case {
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
    useDefault := noInput && q.Default != nil
    if useDefault {
        if def, err = q.tryDefault(); err != nil {
            return err
        }
    }

    // Parse the user's input.
    var val interface{}
    switch q.typ {
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
            return ErrorEmptyInput
        } else if x, err = strconv.Atoi64(in); err != nil {
            return q.makeTypeError(x, in)
        }
        val = x
    case Uint:
        var x uint64
        if useDefault {
            x = def.(uint64)
        } else if noInput {
            return ErrorEmptyInput
        } else if x, err = strconv.Atoui64(in); err != nil {
            return q.makeTypeError(x, in)
        }
        val = x
    case Float:
        var x float64
        if useDefault {
            x = def.(float64)
        } else if noInput {
            return ErrorEmptyInput
        } else if x, err = strconv.Atof64(in); err != nil {
            return q.makeTypeError(x, in)
        }
        val = x
    case StringSlice:
        fallthrough
    case IntSlice:
        fallthrough
    case UintSlice:
        fallthrough
    case FloatSlice:
        err = fmt.Errorf("%s unimplemented", q.typ.String())
    }

    // Check set membership
    switch q.typ {
    case String:
        fallthrough
    case Int:
        fallthrough
    case Uint:
        fallthrough
    case Float:
        if !q.setHas(val) {
            return q.makeErrorNotInSet(val)
        }
    case StringSlice:
        fallthrough
    case IntSlice:
        fallthrough
    case UintSlice:
        fallthrough
    case FloatSlice:
        err = fmt.Errorf("%s unimplemented", q.typ.String())
    }

    // Set the parsed value if there was no error.
    if err == nil {
        q.val = val
    }

    return err
}

// Cast a value result from a wide (e.g. 64bit) type to the desired type.
// This should not fail under any normal circumstances, so failure
// should break the loop.
func (q *Question) setDest(dest interface{}) os.Error {
    switch q.Type() {
    case Uint:
        switch dest.(type) {
        case *uint:
            d := dest.(*uint)
            *(d) = uint(q.val.(uint64))
            if x := uint64(*(d)); x != q.val.(uint64) {
                return ErrorPrecision{q.val.(uint64), x}
            }
        case *uint8:
            d := dest.(*uint8)
            *(d) = uint8(q.val.(uint64))
            if x := uint64(*(d)); x != q.val.(uint64) {
                return ErrorPrecision{q.val.(uint64), x}
            }
        case *uint16:
            d := dest.(*uint16)
            *(d) = uint16(q.val.(uint64))
            if x := uint64(*(d)); x != q.val.(uint64) {
                return ErrorPrecision{q.val.(uint64), x}
            }
        case *uint32:
            d := dest.(*uint32)
            *(d) = uint32(q.val.(uint64))
            if x := uint64(*(d)); x != q.val.(uint64) {
                return ErrorPrecision{q.val.(uint64), x}
            }
        case *uint64:
            *(dest.(*uint64)) = q.val.(uint64)
        default:
            return fmt.Errorf("Unexpected cast type")
        }
    case Int:
        switch dest.(type) {
        case *int:
            d := dest.(*int)
            *(d) = int(q.val.(int64))
            if x := int64(*(d)); x != q.val.(int64) {
                return ErrorPrecision{q.val.(int64), x}
            }
        case *int8:
            d := dest.(*int8)
            *(d) = int8(q.val.(int64))
            if x := int64(*(d)); x != q.val.(int64) {
                return ErrorPrecision{q.val.(int64), x}
            }
        case *int16:
            d := dest.(*int16)
            *(d) = int16(q.val.(int64))
            if x := int64(*(d)); x != q.val.(int64) {
                return ErrorPrecision{q.val.(int64), x}
            }
        case *int32:
            d := dest.(*int32)
            *(d) = int32(q.val.(int64))
            if x := int64(*(d)); x != q.val.(int64) {
                return ErrorPrecision{q.val.(int64), x}
            }
        case *int64:
            *(dest.(*int64)) = q.val.(int64)
        default:
            return fmt.Errorf("Unexpected cast type")
        }
    case Float:
        switch dest.(type) {
        case *float32:
            d := dest.(*float32)
            *(d) = float32(q.val.(float64))
            if x := float64(*(d)); x != q.val.(float64) {
                return ErrorPrecision{q.val.(float64), x}
            }
        case *float64:
            *(dest.(*float64)) = q.val.(float64)
        default:
            return fmt.Errorf("Unexpected cast type")
        }
    case String:
        switch dest.(type) {
        case *string:
            *(dest.(*string)) = q.val.(string)
        default:
            return fmt.Errorf("Unexpected cast type")
        }
    }
    return nil
}
