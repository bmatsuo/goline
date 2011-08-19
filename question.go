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
    "bufio"
    "fmt"
    "os"
)

type ErrorPrecision struct {
    Wide, Thin interface{}
}

func (e ErrorPrecision) String() string {
    return fmt.Sprintf("Input out of destination range (%v -> %v)", e.Wide, e.Thin)
}

//  Prompt the user for text input. The result is stored in dest, which must
//  be a pointer to a native Go type (int, uint16, string, float32, ...).
//  Slice types are not currently supported. List input must be done with a
//  *string destination and post-processing.
func Ask(dest interface{}, msg string, config func(*Answer)) os.Error {
    if k := reflect.TypeOf(dest).Kind(); k != reflect.Ptr && k != reflect.Slice {
        return fmt.Errorf("Ask(...) requires a Ptr type, not %s", k.String())
    } else if k == reflect.Slice {
        return fmt.Errorf("Ask(...) can not currently assign to slices.")
    }

    var t Type
    switch dest.(type) {
    case *uint:
        t = Uint
    case *uint8:
        t = Uint
    case *uint16:
        t = Uint
    case *uint32:
        t = Uint
    case *uint64:
        t = Uint
    case *int:
        t = Int
    case *int8:
        t = Int
    case *int16:
        t = Int
    case *int32:
        t = Int
    case *int64:
        t = Int
    case *float32:
        t = Float
    case *float64:
        t = Float
    case *string:
        t = String
    default:
        fmt.Errorf("Unusable destination")
    }
    a := newAnswer(t)
    config(a)

    pstr := strings.Trim(msg, "\n") + a.DefaultString()
    r := bufio.NewReader(os.Stdin)
    for {
        fmt.Print(pstr)
        var resp []byte
        for cont := true; cont; {
            s, isPrefix, err := r.ReadLine()
            cont = isPrefix
            if err != nil {
                return err
            }
            resp = append(resp, s...)
        }
        if err := a.parse(string(resp)); err != nil {
            switch err.(type) {
            case RecoverableError:
                fmt.Println(err.String())
                continue
            default:
                return err
            }
        }
        var errCast os.Error
        switch dest.(type) {
        case *uint:
            d := dest.(*uint)
            *(d) = uint(a.v.(uint64))
            if x := uint64(*(d)); x != a.v.(uint64) {
                errCast = ErrorPrecision{a.v.(uint64), x}
            }
        case *uint8:
            d := dest.(*uint8)
            *(d) = uint8(a.v.(uint64))
            if x := uint64(*(d)); x != a.v.(uint64) {
                errCast = ErrorPrecision{a.v.(uint64), x}
            }
        case *uint16:
            d := dest.(*uint16)
            *(d) = uint16(a.v.(uint64))
            if x := uint64(*(d)); x != a.v.(uint64) {
                errCast = ErrorPrecision{a.v.(uint64), x}
            }
        case *uint32:
            d := dest.(*uint32)
            *(d) = uint32(a.v.(uint64))
            if x := uint64(*(d)); x != a.v.(uint64) {
                errCast = ErrorPrecision{a.v.(uint64), x}
            }
        case *uint64:
            *(dest.(*uint64)) = a.v.(uint64)
        case *int:
            d := dest.(*int)
            *(d) = int(a.v.(int64))
            if x := int64(*(d)); x != a.v.(int64) {
                errCast = ErrorPrecision{a.v.(int64), x}
            }
        case *int8:
            d := dest.(*int8)
            *(d) = int8(a.v.(int64))
            if x := int64(*(d)); x != a.v.(int64) {
                errCast = ErrorPrecision{a.v.(int64), x}
            }
        case *int16:
            d := dest.(*int16)
            *(d) = int16(a.v.(int64))
            if x := int64(*(d)); x != a.v.(int64) {
                errCast = ErrorPrecision{a.v.(int64), x}
            }
        case *int32:
            d := dest.(*int32)
            *(d) = int32(a.v.(int64))
            if x := int64(*(d)); x != a.v.(int64) {
                errCast = ErrorPrecision{a.v.(int64), x}
            }
        case *int64:
            *(dest.(*int64)) = a.v.(int64)
        case *float32:
            d := dest.(*float32)
            *(d) = float32(a.v.(float64))
            if x := float64(*(d)); x != a.v.(float64) {
                errCast = ErrorPrecision{a.v.(float64), x}
            }
        case *float64:
            *(dest.(*float64)) = a.v.(float64)
        case *string:
            *(dest.(*string)) = a.v.(string)
        default:
            errCast = fmt.Errorf("Unexpected cast type")
        }
        if errCast != nil {
            return errCast
        }
        break
    }
    return nil
}
