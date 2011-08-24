// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 *  Filename:    goline.go
 *  Package:     goline
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Aug 13 02:28:54 PDT 2011
 *  Description: 
 */

// Package goline manages prompts in the spirit of highline.
package goline
import (
    "reflect"
    "strings"
    "unicode"
    "utf8"
    "fmt"
    "os"
)

func Say(msg string) (int, os.Error) {
    if c, _ := utf8.DecodeLastRuneInString(msg); unicode.IsSpace(c) {
        return fmt.Print(msg)
    }
    return fmt.Println(msg)
}

type Stringer interface {
    String() string
}

var (
    stringerZero Stringer
    stringerType = reflect.TypeOf(stringerZero)
)

type ListMode uint

const (
    ColumnsAcross ListMode = iota
    ColumnsDown
    Inline
    Rows
)

func List(items interface{}, mode ListMode, option interface{}) {
    ival := reflect.ValueOf(items)
    itype := ival.Type()
    if k := itype.Kind(); k != reflect.Slice {
        panic(os.NewError("List given non-Slice types."))
    }
    strs := make([]string, ival.Len())
    for i := range strs {
        v := ival.Index(i).Interface()
        switch v.(type) {
        case Stringer:
            strs[i] = v.(Stringer).String()
        case string:
            strs[i] = v.(string)
        default:
            panic(os.NewError("List items contain non-string, non-Stringer item"))
        }
    }
    switch mode {
    case ColumnsAcross:
        fallthrough
    case ColumnsDown:
        wrap := 80
        switch option.(type) {
        case nil:
        case int:
            wrap = option.(int)
        default:
            panic(os.NewError("List option of unacceptable type"))
        }

        var width int
        for i := range strs {
            if n := len(strs[i]); n > width {
                width = n
            }
        }

        n := len(strs)
        ncols := (wrap + 1) / (width + 1)

        if ncols <= 1 {
            // Just print rows if no more than 1 column fits.
            for i := range strs {
                Say(strings.TrimRightFunc(strs[i], unicode.IsSpace))
            }
            break
        }

        nrows := (n + ncols - 1) / ncols

        sfmt := fmt.Sprintf("%%-%ds", width)
        for i := range strs {
            strs[i] = fmt.Sprintf(sfmt, strs[i])
        }

        switch mode {
        case ColumnsAcross:
            for i := 0; i < n; i += ncols {
                end := i+ncols
                if end > n {
                    end = n
                }
                row := strs[i:end]
                Say(strings.TrimRightFunc(strings.Join(row, " "), unicode.IsSpace))
            }
        case ColumnsDown:
            for i := 0; i < nrows; i++ {
                var row []string
                for j := 0; j < ncols; j++ {
                    index := j*nrows + i
                    if index >= n {
                        break
                    }
                    row = append(row, strs[index])
                }
                Say(strings.TrimRightFunc(strings.Join(row, " "), unicode.IsSpace))
            }
        }
    case Inline:
        n := len(strs)
        if n == 1 {
            Say(strings.TrimRightFunc(strs[0], unicode.IsSpace))
            break
        }
        join := " or "
        switch option.(type) {
        case nil:
        case string:
            join = option.(string)
        default:
            panic(os.NewError("List option of unacceptable type"))
        }
        if n == 2 {
            Say(strings.Join([]string{strs[n-2], join, strs[n-2], "\n"}, ""))
            break
        }
        strs[n-1] = join + strings.TrimRightFunc(strs[n-1], unicode.IsSpace)
        Say(strings.Join(strs, ", ") + "\n")
    case Rows:
        for i := range strs {
            Say(strings.TrimRightFunc(strs[i], unicode.IsSpace))
        }
    default:
        panic(os.NewError("Unknown mode"))
    }
}
