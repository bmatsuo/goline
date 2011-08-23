package goline
/*
 *  Filename:    menu.go
 *  Package:     goline
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Aug 13 23:31:27 PDT 2011
 *  Description: 
 */
import (
    "reflect"
    "fmt"
    "os"
)

type Stringer interface {
    String() string
}

var (
    stringerZero Stringer
    stringerType = reflect.TypeOf(stringerZero)
)

func Choose(dest, choices interface{}, label, msg string) os.Error {
    dval := reflect.ValueOf(dest)
    cval := reflect.ValueOf(choices)
    if k := cval.Kind(); k != reflect.Slice {
        return fmt.Errorf("Choices must be a Slice type")
    }
    if k := dval.Kind(); k != reflect.Ptr {
        return fmt.Errorf("Destination must be a Ptr type")
    }
    if cval.Type().Elem().Name() != dval.Type().Elem().Name() {
        return fmt.Errorf("Type mismatch %s != %s",
            cval.Type().Elem().Name(),
            dval.Type().Elem().Name())
    }
    if !cval.Type().Elem().Implements(stringerType) {
        return fmt.Errorf("Choices do not implement Stringer")
    }
    toset := reflect.Indirect(dval)
    if !toset.CanSet() {
        return fmt.Errorf("Can not set the destination")
    }
    var schoices []string
    switch choices.(type) {
    case []string:
        schoices = choices.([]string)
    default:
        schoices = make([]string, cval.Len())
        for i := range schoices {
            elval := cval.Index(i)
            stringer := elval.MethodByName("String")
            sval := stringer.Call([]reflect.Value{})
            schoices[i] = sval[0].Interface().(string)
        }
    }
    fmt.Println(label)
    for i, s := range schoices {
        fmt.Println(i+1, s)
    }
    var ichosen int
    err := Ask(&ichosen, msg, func(a *Answer) {
        // TODO Make a string set.
        //a.In(int64(1), int64(len(schoices)))
    })
    if err != nil {
        return err
    }
    chosen := cval.Index(ichosen)
    toset.Set(chosen)
    return nil
}

type Menu struct {
}
