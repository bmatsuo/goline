package goline
/*
 *  Filename:    errors.go
 *  Package:     goline
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Tue Aug 23 19:42:49 PDT 2011
 *  Description: 
 */
import (
    "reflect"
    "os"
    "fmt"
)

type RecoverableError interface {
    os.Error
    IsRecoverable() bool
}

func panicUnrecoverable(err os.Error) {
    if err != nil {
        switch err.(type) {
        case RecoverableError:
            break
        default:
            panic(err)
        }
    }
}

type ErrorParse struct {
    string
    os.Error
}

func (e ErrorParse) String() string      { return fmt.Sprintf("Parsing %#v", e.string) }
func (e ErrorParse) IsRecoverable() bool { return true }

type ErrorPrecision struct {
    Wide, Thin interface{}
}

func (e ErrorPrecision) String() string {
    return fmt.Sprintf("Input out of destination range (%v -> %v)", e.Wide, e.Thin)
}
func (e ErrorPrecision) IsRecoverable() bool { return true }


type ErrorNotInSet struct{ os.Error }

func (err ErrorNotInSet) IsRecoverable() bool { return true }

func (a *Answer) makeErrorNotInSet(val interface{}) ErrorNotInSet {
    return ErrorNotInSet{
        fmt.Errorf("%s %s (%#v)", a.Responses[NotInSet], a.set.String(), val)}
}

func (a *Answer) makeTypeError(expect, recv interface{}) os.Error {
    return fmt.Errorf("%s (%s != %s)",
        a.Responses[InvalidType],
        reflect.ValueOf(recv).Kind().String(),
        reflect.ValueOf(expect).Kind().String())
}

/*
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
*/

type ErrorEmptyInput uint

func (oor ErrorEmptyInput) String() string      { return "Can not use empty value" }
func (oor ErrorEmptyInput) IsRecoverable() bool { return true }


func errorEmptyRange(min, max interface{}) os.Error {
    return fmt.Errorf("Range max is less than min (%v < %v)", min, max)
}

func errorSetMemberType(set, member interface{}) os.Error {
    return fmt.Errorf("Set type %v cannot contain type %v",
        reflect.TypeOf(set).String(),
        reflect.TypeOf(member).String())
}

var errorNoChoices = os.NewError("No Menu choices given")
