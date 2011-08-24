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
    "unicode"
    //"strings"
    "bufio"
    "fmt"
    "os"
)

//  Prompt the user for text input. The result is stored in dest, which must
//  be a pointer to a native Go type (int, uint16, string, float32, ...).
//  Slice types are not currently supported. List input must be done with a
//  *string destination and post-processing.
func Ask(dest interface{}, msg string, config func(*Answer)) (e os.Error) {
    var a *Answer
    defer func() {
        if err := recover(); err != nil {
            switch err.(type) {
            case os.Error:
                // Call a panic method...
                if a.Panic != nil {
                    a.Panic(err.(os.Error))
                }
            default:
                panic(err)
            }
        }
    }()
    if k := reflect.TypeOf(dest).Kind(); k != reflect.Ptr && k != reflect.Slice {
        panicUnrecoverable(fmt.Errorf("Ask(...) requires a Ptr type, not %s", k.String()))
        return
    } else if k == reflect.Slice {
        panicUnrecoverable(fmt.Errorf("Ask(...) can not currently assign to slices."))
        return
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
    a = newAnswer(t)
    a.Question = msg
    if config != nil {
        config(a)
    }

    if err := a.tryFirstAnswer(); err == nil && a.val != nil {
        if err := a.setDest(dest); err != nil {
            panicUnrecoverable(err)
            a.val = nil
        }
        return
    }

    prompt := msg
    contFunc := func(err os.Error) {
        Say(fmt.Sprintf("Error: %s\n", err.String()))
        prompt = a.Responses[AskOnError]
    }
    r := bufio.NewReader(os.Stdin)
    for {
        tail := stringSuffixFunc(prompt, unicode.IsSpace)
        Say(prompt + a.defaultString(tail))
        var resp []byte
        for cont := true; cont; {
            s, isPrefix, err := r.ReadLine()
            if err != nil {
                panicUnrecoverable(err)
                return
            }
            resp = append(resp, s...)
            cont = isPrefix
        }
        if err := a.parse(string(resp)); err != nil {
            panicUnrecoverable(err)
            contFunc(err)
            continue
        }

        // Cast the result from a wide (e.g. 64bit) type to the desired type.
        // This should not fail under any normal circumstances, so failure
        // should break the loop.
        if err := a.setDest(dest); err != nil {
            panicUnrecoverable(err)
            contFunc(err)
            continue
        }
        break
    }
    return
}
