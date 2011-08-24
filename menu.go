package goline
/*
 *  Filename:    menu.go
 *  Package:     goline
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Aug 13 23:31:27 PDT 2011
 *  Description: 
 */
import (
    "strconv"
    "reflect"
    "fmt"
    "os"
)

var errorNoChoices = os.NewError("No Menu choices given")

func Choose(config func(*Menu)) (i int, v interface{}) {
    i = -1
    m := newMenu()
    config(m)
    if m.Len() == 0 {
        if m.Panic != nil {
            m.Panic(errorNoChoices)
            return
        }
    }

    if len(m.Header) > 0 {
        Say(m.Header)
    }

    raw, selections, tr := m.Selections()
    List(raw, m.ListMode, nil)
    ok := true
    var resp string
    Ask(&resp, m.Question, func(a *Answer) {
        a.In(StringSet(selections))
        a.Panic = func(err os.Error) {
            ok = false
            m.Panic(err)
        }
    })
    if !ok {
        return
    }

    i = tr[resp]
    v = m.Choices[i]

    return
}

//  Construct an IndexMode by combining index options and suffix options.
//      mode1 := Literal | DefaultSuffix // Items like "- Hello"
//      mode2 := Number | LiteralSuffix // Items like "1::Hello"
//      mode3 := Letter | DefaultSuffix // Items like "a. Hello"
//  Do not combine multiple index options, or multiple suffix options. There
//  will likely be unintended consequences.
type IndexMode uint

const (
    Literal IndexMode = iota
    Number
    Letter
)

const (
    DefaultSuffix IndexMode = iota << 8
    LiteralSuffix
)

func (imode IndexMode) UseLiteral() bool       { return imode&0xFF == Literal }
func (imode IndexMode) UseNumber() bool        { return imode&0xFF == Number }
func (imode IndexMode) UseLetter() bool        { return imode&0xFF == Letter }
func (imode IndexMode) UseDefaultSuffix() bool { return imode&(0xFF00) == DefaultSuffix }
func (imode IndexMode) UseLiteralSuffix() bool { return imode&(0xFF00) == LiteralSuffix }

var alpha = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
var alphaLen = len(alpha)

func getLetterIndex(i int) string {
    p := make([]byte, 0, 1)
    for i > 0 {
        p = append(p, 0)
        if len(p) > 1 {
            copy(p[1:], p)
        }
        p[0] = alpha[i%alphaLen]
        i /= alphaLen
    }
    return string(p)
}

func (m *Menu) getIndexNoSuffix(i int) string {
    switch {
    case m.UseLiteral():
        return m.Index
    case m.UseNumber():
        return strconv.Itoa(i)
    case m.UseLetter():
        return getLetterIndex(i)
    }
    panic("Index option error.")
}

func (m *Menu) getIndex(i int) string {
    ind := m.getIndexNoSuffix(i)

    var s string
    switch {
    case m.UseLiteralSuffix():
        s = m.IndexSuffix
    case m.UseDefaultSuffix():
        switch {
        case m.UseLiteral():
            // This case might produce a warning...
            s = " "
        case m.UseNumber():
            fallthrough
        case m.UseLetter():
            s = ". "
        }
    default:
        panic("Suffix option error.")
    }

    return ind + s
}

type SelectMode uint

const (
    IndexSelect SelectMode = 1 << iota
    NameSelect
)

func (smode SelectMode) SelectIndices() bool { return smode&IndexSelect != 0 }
func (smode SelectMode) SelectNames() bool   { return smode&NameSelect != 0 }

type Menu struct {
    // A list of Menu choices. See Menu.Choice and Menu.SetChoices
    Choices []Stringer
    // A header text (describing the Menu).
    Header string
    // The text to prompt the user with after displaying the Menu.
    Question string
    // This mode is passed directly to the function List().
    ListMode
    // The index string mode.
    IndexMode
    // The selection mode for the choice prompt.
    SelectMode
    // The index and suffix used for all choices if IndexMode is Literal.
    Index       string
    IndexSuffix string
    // A handler function for any errors encountered.
    Panic func(os.Error)
}

func newMenu() *Menu {
    m := new(Menu)
    m.ListMode = Rows
    m.IndexMode = Number
    m.SelectMode = IndexSelect | NameSelect
    return m
}

//  The number of choices currently in the Menu.
func (m *Menu) Len() int { return len(m.Choices) }

//  Create a list of menu items (with indices) and a translation table that
//  maps menu selections (possibly name and index) to an integer index into
//  m.Choices.
func (m *Menu) Selections() (choices []string, selections []string, tr map[string]int) {
    selectIndices, selectNames := m.SelectIndices(), m.SelectNames()
    if m.UseLiteral() {
        // Can't select indices if all choices have the same index.
        selectIndices = false
    }
    // Compute the necessary size of the structures and allocate them.
    n := m.Len()
    trSize := n
    if selectIndices && selectNames {
        trSize += n
    }
    choices = make([]string, n)
    selections = make([]string, 0, trSize)
    tr = make(map[string]int, trSize)

    addSelection := func(i int, s string) {
        if _, present := tr[s]; present {
            panic(fmt.Errorf("Selection conflict %s", s))
        }
        tr[s] = i
        selections = append(selections, s)
    }

    for i := range m.Choices {
        choices[i] = m.getIndex(i) + m.Choices[i].String()
        if selectIndices {
            addSelection(i, m.getIndexNoSuffix(i))
        }
        if selectNames {
            addSelection(i, m.Choices[i].String())
        }
    }
    return
}

// Make a []Stringer with objects from a slice of arbitrary (interface) type.
// This should be called before calling m.Choice() to add single choices.
func (m *Menu) SetChoices(cs interface{}) {
    // Zero out the old choice list (even if there is an error)
    var zero []Stringer
    m.Choices = zero

    // Make sure cs is a slice containing Stringer objects.
    csval := reflect.ValueOf(cs)
    cstyp := csval.Type()
    if k := cstyp.Kind(); k != reflect.Slice {
        panic("Argument of Menu.ChoicesSlice must be a slice.")
    }

    n := csval.Len()
    m.Choices = make([]Stringer, n)
    for i := 0; i < n; i++ {
        m.Choices[i] = makeStringer(csval.Index(i).Interface())
    }
}

// Append a choice (either string or Stringer) to m.Choices.
func (m *Menu) Choice(s interface{}) { m.Choices = append(m.Choices, makeStringer(s)) }
// Prepend a choice (either string or Stringer) to the front (top) of m.Choices.
func (m *Menu) ChoicePre(s interface{}) {
    m.Choices = append(m.Choices, zeroStringer)
    if m.Len() > 1 {
        copy(m.Choices[1:], m.Choices)
    }
    m.Choices[0] = makeStringer(s)
}
