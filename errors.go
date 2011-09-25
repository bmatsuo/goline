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

//  Some simple errors with no dynamic String() elements.
var (
    ErrorEmptyInput = NewErrorRecoverable("Can not use empty string as value")
    ErrorNoChoices  = NewError("No Menu choices given")
)

//  An interface for errors which prompts can recover from.
type RecoverableError interface {
    os.Error
    IsRecoverable() bool
}

type RespondableError interface {
    RecoverableError
    Response() Response
}

func NewError(msg string) os.Error { return os.NewError(msg) }
func NewErrorRecoverable(msg string) SimpleRecoverableError {
    return SimpleRecoverableError(msg)
}

//  A simple wrapper turning strings into RecoverableErrors
type SimpleRecoverableError string

func (err SimpleRecoverableError) IsRecoverable() bool { return true }
func (err SimpleRecoverableError) String() string      { return string(err) }

//  Returns true if error e implements RecoverableError.
func ErrorIsRecoverable(e os.Error) (ok bool) {
    switch e.(type) {
    case RecoverableError:
        ok = e.(RecoverableError).IsRecoverable()
    }
    return
}

//  Returns true if error e implements RespondableError.
func ErrorHasResponse(e os.Error) (ok bool) {
    switch e.(type) {
    case RespondableError:
        ok = true
    }
    return
}

//  Raises a run-time panic if the error err is not a RecoverableError.
func panicUnrecoverable(err os.Error) {
    if err != nil && !ErrorIsRecoverable(err) {
        panic(err)
    }
}

//  Errors returned when the input is too large to fit in a standard data type.
type ErrorPrecision struct{ Wide, Thin interface{} }

var errPrecisionMsg = "Input out of destination range (%v -> %v)"

func (e ErrorPrecision) String() string      { return fmt.Sprintf(errPrecisionMsg, e.Wide, e.Thin) }
func (e ErrorPrecision) IsRecoverable() bool { return true }
func (e ErrorPrecision) Response() Response { return Precision }


//  Errors returned when the input provided was not in a Question's AnswerSet.
type ErrorNotInSet struct{ os.Error }

func (a *Question) makeErrorNotInSet(val interface{}) ErrorNotInSet {
    return ErrorNotInSet{
        fmt.Errorf("%s %s (%#v)", a.Responses[NotInSet], a.set.String(), val)}
}
func (err ErrorNotInSet) IsRecoverable() bool { return true }
func (err ErrorNotInSet) Response() Response { return NotInSet }

//  Errors raised when the input (or default, or first-answer) are not of the
//  prompting Question's type.
type ErrorType struct {
    msg      string
    exp, rec reflect.Value
}

func makeTypeError(msg string, exp, rec interface{}) os.Error {
    return ErrorType{msg, reflect.ValueOf(exp), reflect.ValueOf(rec)}
}
func (a *Question) makeTypeError(expect, recv interface{}) os.Error {
    return makeTypeError(a.Responses[InvalidType], expect, recv)
}
func (e ErrorType) IsRecoverable() bool { return true }
func (e ErrorType) ExpKind() string     { return e.exp.Kind().String() }
func (e ErrorType) RecKind() string     { return e.rec.Kind().String() }
func (e ErrorType) String() string {
    return fmt.Sprintf("%s (%s != %s)", e.msg, e.RecKind(), e.ExpKind())
}
func (e ErrorType) Response() Response { return InvalidType }

//  Errors raised when an AnswerSet of improper type was given to the Question.
type ErrorMemberType struct{ Set, Member reflect.Type }

func makeErrorMemberType(s AnswerSet, member interface{}) os.Error {
    return ErrorMemberType{reflect.TypeOf(s), reflect.TypeOf(member)}
}
func (err ErrorMemberType) Type() string    { return err.Set.String() }
func (err ErrorMemberType) MemberType() string { return err.Member.String() }
func (err ErrorMemberType) String() string {
    return fmt.Sprintf("%s can't contain %s", err.Type(), err.MemberType())
}
