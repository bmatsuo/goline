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

See github.com/bmatsuo/goline/examples for examples using goline.
*/
package goline

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

func fSay(wr io.Writer, msg string, trim bool) (int, error) {
	if trim {
		msg = strings.TrimRightFunc(msg, unicode.IsSpace)
	}
	if c, _ := utf8.DecodeLastRuneInString(msg); unicode.IsSpace(c) {
		return fmt.Fprint(wr, msg)
	}
	return fmt.Fprintln(wr, msg)
}

//  A simple function for printing (single-line) messages and prompts to
//  os.Stdout. If trailing whitespace is present in the given message, it
//  will be printed as given. Otherwise, a trailing newline '\n' will be
//  printed after the message.
//      goline.Say("Hello, World!") // Prints "Hello, World!\n"
//      goline.Say("Hello, World! ") // Prints "Hello, World! "
//  See also, SayTrimmed.
func Say(msg string) (int, error) {
	if c, _ := utf8.DecodeLastRuneInString(msg); unicode.IsSpace(c) {
		return fmt.Print(msg)
	}
	return fmt.Println(msg)
}

//  Like Say, but trailing whitespace is removed from the message before
//  an internal call to Say is made.
//      goline.SayTrimmed("Hello, World! \n\t\t") // Prints "Hello, World!\n"
func SayTrimmed(msg string) (int, error) {
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
	itype := ival.Type()
	if k := itype.Kind(); k != reflect.Slice {
		panic(errors.New("List given non-Slice types."))
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
			panic(errors.New("List items contain non-string, non-Stringer item"))
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
			panic(errors.New("List option of unacceptable type"))
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
				SayTrimmed(strs[i])
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
				end := i + ncols
				if end > n {
					end = n
				}
				row := strs[i:end]
				SayTrimmed(strings.Join(row, " "))
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
				SayTrimmed(strings.Join(row, " "))
			}
		}
	case Inline:
		n := len(strs)
		if n == 1 {
			SayTrimmed(strs[0])
			break
		}
		join := "or "
		switch option.(type) {
		case nil:
		case string:
			join = option.(string)
		default:
			panic(errors.New("List option of unacceptable type"))
		}
		if n == 2 {
			Say(strings.Join([]string{strs[n-2], join, strs[n-2], "\n"}, ""))
			break
		}
		strs[n-1] = join + strs[n-1]
		SayTrimmed(strings.Join(strs, ", "))
	case Rows:
		for i := range strs {
			SayTrimmed(strs[i])
		}
	default:
		panic(errors.New("Unknown mode"))
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
func Ask(dest interface{}, msg string, config func(*Question)) (e error) {
	var q *Question
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case error:
				// Call a panic method...
				if q.Panic != nil {
					q.Panic(err.(error))
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
	q = newQuestion(t)
	q.Question = msg
	if config != nil {
		config(q)
	}

	if err := q.tryFirstAnswer(); err == nil && q.val != nil {
		if err := q.setDest(dest); err != nil {
			panicUnrecoverable(err)
			q.val = nil
		}
		return
	}

	prompt := msg
	contFunc := func(err error) {
		Say(fmt.Sprintf("Error: %s\n", err.Error()))
		prompt = q.Responses[AskOnError]
	}
	r := bufio.NewReader(os.Stdin)
	for {
		tail := stringSuffixFunc(prompt, unicode.IsSpace)
		Say(prompt + q.defaultString(tail))
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

	var okstr string
	var err error
	Ask(&okstr, question, func(q *Question) {
		q.Default = def
		q.In(StringSet{"yes", "y", "no", "n"})
		if config != nil {
			config(q)
		}
		if q.Panic != nil {
			f := q.Panic
			q.Panic = func(e error) {
				err = e
				f(e)
			}
		}
	})
	if err != nil {
		return false
	}
	if okstr[0] == 'y' {
		return true
	}
	return false
}

func splitShellCmd(cmd string) (name, args string) {
	cmd = strings.TrimLeftFunc(cmd, unicode.IsSpace)
	pre := strings.IndexFunc(cmd, unicode.IsSpace)
	if pre > 0 {
		name = cmd[0:pre]
		args = strings.TrimLeftFunc(cmd[pre:], unicode.IsSpace)
	} else if pre == -1 {
		name = cmd
		args = ""
	} else {
		panic("unexpected case (untrimmed)")
	}
	return
}

//  Prompt the user to choose from a list of choices. Return the index
//  of the chosen item, and the item itself in an empty interface. See
//  Menu for more information about configuring the prompt.
func Choose(config func(*Menu)) (i int, v interface{}) {
	i = -1
	m := newMenu()
	config(m)
	if m.Len() == 0 {
		if m.Panic != nil {
			m.Panic(ErrorNoChoices)
			return
		}
		panic(ErrorNoChoices)
	}

	if len(m.Header) > 0 {
		Say(m.Header)
	}

	raw, selections, tr := m.Selections()
	List(raw, m.ListMode, nil)
	ok := true
	var resp, args string
	Ask(&resp, m.Question, func(q *Question) {
		var set AnswerSet = StringSet(selections)
		if m.Shell {
			set = shellCommandSet(set.(StringSet))
		}
		q.In(set)
		q.Panic = func(err error) {
			ok = false
			if m.Panic != nil {
				m.Panic(err)
			} else {
				panic(err)
			}
		}
	})
	if !ok {
		return
	}
	if m.Shell {
		resp, args = splitShellCmd(resp)
	}

	i = tr[resp]
	v = m.Choices[i]

	if m.Actions[i] != nil {
		m.Actions[i](resp, args)
	}

	return
}
