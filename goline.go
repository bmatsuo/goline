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

/*
Package goline is a command line interfacing (prompting) library inspired
by Ruby's HighLine.

Differences for HighLine users:

    - To be more Go-ish, where HighLine uses the term "strip", the package
      uses "trim".

    - Instead of an `Agree(question, config) bool` function, the package
      provides a function `Confirm(question, yesorno, config) bool`. This is
      because the author things the term "agree" implies the desire of a
      positive response to the question ("yes").
*/
package goline

import (
    "reflect"
    "strings"
    "unicode"
    "bufio"
    "utf8"
    "fmt"
    "os"
)

//  A simple function for printing (single-line) messages and prompts to
//  os.Stdout. If trailing whitespace is present in the given message, it
//  will be printed as given. Otherwise, a trailing newline '\n' will be
//  printed after the message.
//      goline.Say("Hello, World!") // Prints "Hello, World!\n"
//      goline.Say("Hello, World! ") // Prints "Hello, World! "
//  See also, SayTrimmed.
func Say(msg string) (int, os.Error) {
    if c, _ := utf8.DecodeLastRuneInString(msg); unicode.IsSpace(c) {
        return fmt.Print(msg)
    }
    return fmt.Println(msg)
}

//  Like Say, but trailing whitespace is removed from the message before
//  an internal call to Say is made.
//      goline.SayTrimmed("Hello, World! \n\t\t") // Prints "Hello, World!\n"
func SayTrimmed(msg string) (int, os.Error) {
    return Say(strings.TrimRightFunc(msg, unicode.IsSpace))
}

type ListMode uint

const (
    ColumnsAcross ListMode = iota
    ColumnsDown
    Inline
    Rows
)

//  Print a list of items to os.Stdout. The list can be formatted into rows
//  or into a matrix using the ListMode argument. The third argument has
//  different meaning (and type) depending on the mode.
//      MODE        OPTION  DEFAULT     MEANING
//      Rows        n/a
//      Inline      string  " or "      Join terminal element (e.g. "a, b, or c")
//      Columns*    int     80          Maximum line width
//  If the default option is desired, it should be passed as nil.
//      goline.List([]string{"cat", "dog", "go fish"}, goline.ColumnsAcross, nil)
//      /* Outputs:
//       *  cat     dog     go fish
//       */
//      goline.List([]string{"cat", "dog", "go fish"}, goline.ColumnsDown, 15)
//      /* Outputs:
//       *  cat     go fish
//       *  dog
//       */
//      goline.List([]string{"cat", "dog", "go fish"}, goline.Inline, " and ")
//      /* Outputs:
//       *  cat, dog, and go fish
//       */
//      goline.List([]string{"cat", "dog", "go fish"}, goline.Rows, nil)
//      /* Outputs:
//       *  cat
//       *  dog
//       *  go fish
//       */
//  See subdirectory examples/goline-lists.
func List(items interface{}, mode ListMode, option interface{}) {
    ival := reflect.ValueOf(items)

    // Ensure that items is a slice type.
    if k := ival.Type().Kind(); k != reflect.Slice {
        panic(os.NewError("List given non-Slice types."))
    }

    // Stringify each entry in items.
    n := len(ival.Len())
    strs := make([]string, n)
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

    // Print the list.
    switch mode {
    case ColumnsAcross:
        // There is another switch statement in the ColumnsDown action.
        fallthrough
    case ColumnsDown:
        wrap := 80

        // Try to interpret the option variable as a wrap width.
        switch option.(type) {
        case nil:
        case int:
            wrap = option.(int)
        default:
            panic(os.NewError("List option of unacceptable type"))
        }

        // Determine the width of each column, and number of columns to print.
        var width int
        for i := range strs {
            if m := len(strs[i]); m > width {
                width = m
            }
        }
        ncols := (wrap + 1) / (width + 1)

        // Treat the special case of 1 column.
        if ncols <= 1 {
            // Just print rows if no more than 1 column fits.
            for i := range strs {
                SayTrimmed(strs[i])
            }
            break
        }

        nrows := (n + ncols - 1) / ncols

        // Pad each string so that it is one column wide.
        sfmt := fmt.Sprintf("%%-%ds", width) // e.g. "%20"
        for i := range strs {
            strs[i] = fmt.Sprintf(sfmt, strs[i])
        }

        // Print the list according to the mode.
        switch mode {
        case ColumnsAcross:
            for i := 0; i < n; i += ncols {
                // Determine the number of elements in the row.
                end := i + ncols
                if end > n {
                    end = n
                }

                // Print the row (trimming excess padding).
                row := strs[i:end]
                SayTrimmed(strings.Join(row, " "))
            }
        case ColumnsDown:
            for i := 0; i < nrows; i++ {
                // Select the items in the row using column-major order.
                var row []string
                for j := 0; j < ncols; j++ {
                    index := j*nrows + i
                    if index >= n {
                        break
                    }
                    row = append(row, strs[index])
                }

                // Print the row with no excess padding.
                SayTrimmed(strings.Join(row, " "))
            }
        }
    case Inline:
        n := len(strs)

        // Handle the special zero-/single-element case.
        if n == 0 {
            Say("")
            break
        } else if n == 1 {
            SayTrimmed(strs[0])
            break
        }

        // Try to interpret the option argument as a joining string.
        join := "or "
        switch option.(type) {
        case nil:
        case string:
            join = option.(string)
        default:
            panic(os.NewError("List option of unacceptable type"))
        }

        // Handle the special two-element case.
        if n == 2 {
            SayTrimmed(strings.Join([]string{strs[0], " ", join, strs[1], "\n"}, ""))
            break
        }

        // Create and print the inline list.
        strs[n-1] = join + strs[n-1]
        SayTrimmed(strings.Join(strs, ", "))
    case Rows:
        // Print each item on its own row.
        for i := range strs {
            SayTrimmed(strs[i])
        }
    default:
        panic(os.NewError("Unknown mode"))
    }
}

//  Prompt the user for text input. The result is stored in dest, which must
//  be a pointer to a native Go type (int, uint16, string, float32, ...).
//  Slice types are not currently supported. List input must be done with a
//  *string destination and post-processing.
//      package main
//      import (
//          "goline"
//          "os"
//      )
//      func main() {
//          timeout := 5e3
//          goline.Ask(&timeout, "Timeout (ms)? ", func(q *goline.Question) {
//              q.Default = timeout
//              q.In(goline.IntBoundedStrictly(goline.Above, 0))
//              q.Panic = func(e os.Error) { panic(e) }
//          })
//      }
func Ask(dest interface{}, msg string, config func(*Question)) (e os.Error) {
    var q *Question
    defer func() {
        // Recover from thrown errors and use the question's Panic method instead.
        if err := recover(); err != nil {
            switch err.(type) {
            case os.Error:
                // Call a panic method...
                if q.Panic != nil {
                    q.Panic(err.(os.Error))
                }
            default:
                panic(err)
            }
        }
    }()

    // Check the type of the dest variable.
    if k := reflect.TypeOf(dest).Kind(); k != reflect.Ptr && k != reflect.Slice {
        panicUnrecoverable(fmt.Errorf("Ask(...) requires a Ptr type, not %s", k.String()))
        return
    } else if k == reflect.Slice {
        panicUnrecoverable(fmt.Errorf("Ask(...) can not currently assign to slices."))
        return
    }

    // Determine the question type from the destination.
    t, err := TypeOf(reflect.Indirect(reflect.ValueOf(dest)).Interface())
    if err != nil {
        panic(err)
    }

    // Create a new Question and configure it.
    q = newQuestion(t)
    q.Question = msg
    if config != nil {
        config(q)
    }

    // Attempt to short circuit if a first-answer was configured.
    if err := q.tryFirstAnswer(); err == nil && q.val != nil {
        if err := q.setDest(dest); err != nil {
            panicUnrecoverable(err)
            q.val = nil
        }
        return
    }

    // Prepare the prompt and an error function to reset it.
    prompt := msg
    contFunc := func(err os.Error) {
        Say(fmt.Sprintf("Error: %s\n", err.String()))
        prompt = q.Responses[AskOnError]
    }

    // Prompt the user and interpret the next line of input.
    r := bufio.NewReader(os.Stdin)
    for {
        // Say the prompt, preserving the trailing space.
        tail := stringSuffixFunc(prompt, unicode.IsSpace)
        Say(prompt + q.defaultString(tail))

        // Read a line of input.
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

        // Parse the input and check for errors.
        if err := q.parse(string(resp)); err != nil {
            panicUnrecoverable(err)
            contFunc(err)
            continue
        }

        // Cast the result from a wide (e.g. 64bit) type to the desired type.
        // This should not fail under any normal circumstances, so failure
        // should break the loop.
        if err := q.setDest(dest); err != nil {
            panicUnrecoverable(err)
            contFunc(err)
            continue
        }
        break
    }
    return
}

//  Corresponds to HighLine's `agree` method. A simple wrapper around Ask for
//  yes/no questions. Confirm is given a string to prompt the user with, and a
//  default (or expected) value (yes=true, no=false). Returns the value of the
//  input.
//      if Confirm("Fetch data from the server? ", true, nil) {
//          var server string
//          Ask(&server, "Server (host:port)? ", nil)
//          // Fetch some data...
//      }
func Confirm(question string, yes bool, config func(*Question)) bool {
    def := "no"
    if yes {
        def = "yes"
    }

    // Ask a yes/no question.
    var okstr string
    var err os.Error
    Ask(&okstr, question, func(q *Question) {
        q.Default = def
        q.In(StringSet{"yes", "y", "no", "n"})
        if config != nil {
            config(q)
        }
        if q.Panic != nil {
            f := q.Panic
            q.Panic = func(e os.Error) {
                err = e
                f(e)
            }
        }
    })

    // Interpret the result.
    if err != nil {
        return panic(err)
    }
    if okstr[0] == 'y' {
        return true
    }
    return false
}

//  Prompt the user to choose from a list of choices. Return the index
//  of the chosen item, and the item itself in an empty interface. See
//  Menu for more information about configuring the prompt.
func Choose(config func(*Menu)) (i int, v interface{}) {
    i = -1 // The chosen index.

    // Create a Menu and run the config function.
    m := newMenu()
    config(m)
    if m.Len() == 0 {
        if m.Panic != nil {
            m.Panic(ErrorNoChoices)
            return
        }
    }

    // Print the menu header if there is one.
    if len(m.Header) > 0 {
        Say(m.Header)
    }

    // Display the list.
    raw, selections, tr := m.Selections()
    List(raw, m.ListMode, nil)
    ok := true

    // Ask for a selection.
    var resp string
    Ask(&resp, m.Question, func(q *Question) {
        q.In(StringSet(selections))
        q.Panic = func(err os.Error) {
            ok = false
            m.Panic(err)
        }
    })
    if !ok {
        return
    }

    // Translate the response into an index and choice.
    i = tr[resp]
    v = m.Choices[i]

    // Return the index and choice selected.
    return
}
