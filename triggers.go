package kovacs

import (
	"fmt"
	"time"
)

type ExpressionScope struct{ string }
type CaseSensitivity struct{ bool }

var (
	Wholename = ExpressionScope{"wholename"}
	Basename  = ExpressionScope{"basename"}

	CaseSensitive   = CaseSensitivity{true}
	CaseInsensitive = CaseSensitivity{false}
)

var (
	StdinDevNull     StdinType = stdinString("/dev/null")
	StdinNamePerLine StdinType = stdinString("NAME_PER_LINE")
)

type StdinType interface {
	stdinNoop()
}

type stdinString string

func (s stdinString) stdinNoop() {}

type StdinArray []string

func (s StdinArray) stdinNoop() {}

type Trigger struct {
	Name          string      `json:"name"`
	Command       []string    `json:"command"`
	AppendFiles   bool        `json:"append_files,omitempty"`
	Expression    interface{} `json:"expression"`
	Stdin         StdinType   `json:"stdin"`
	Stdout        string      `json:"stdout"`
	Stderr        string      `json:"stderr"`
	MaxFilesStdin int         `json:"max_files_stdin"`
	Chdir         string      `json:"chdir"`
}

type Expression interface {
	noopExpr()
}

type exprSlice []Expression

func (sl exprSlice) noopExpr() {}

func fromStringSlice(strs []string) exprSlice {
	sl := make(exprSlice, len(strs))
	for i := range strs {
		sl[i] = exprString(strs[i])
	}

	return sl
}

type exprString string

func (s exprString) noopExpr() {}

func caseName(name string, cs CaseSensitivity) exprString {
	if cs == CaseSensitive {
		return exprString(name)
	}

	return exprString("i" + name)
}

// https://facebook.github.io/watchman/docs/expr/allof.html
func AllOf(expr ...Expression) Expression {
	return append(exprSlice{exprString("allof")}, expr...)
}

// https://facebook.github.io/watchman/docs/expr/anyof.html
func AnyOf(expr ...Expression) Expression {
	return append(exprSlice{exprString("anyof")}, expr...)
}

// https://facebook.github.io/watchman/docs/expr/dirname.html
func Dirname(cs CaseSensitivity, dir string) Expression {
	return nil
}

// https://facebook.github.io/watchman/docs/expr/empty.html
func Empty() Expression {
	return exprSlice{exprString("empty")}
}

// https://facebook.github.io/watchman/docs/expr/exists.html
func Exists() Expression {
	return exprSlice{exprString("exists")}
}

// https://facebook.github.io/watchman/docs/expr/match.html
func Match(cs CaseSensitivity, scope ExpressionScope, pattern string) Expression {
	return exprSlice{caseName("match", cs), exprString(pattern), exprString(scope.string)}
}

// https://facebook.github.io/watchman/docs/expr/name.html
func Name(cs CaseSensitivity, scope ExpressionScope, names ...string) Expression {
	return exprSlice{caseName("name", cs), fromStringSlice(names), exprString("basename")}
}

// https://facebook.github.io/watchman/docs/expr/not.html
func Not(expr Expression) Expression {
	return exprSlice{exprString("not"), expr}
}

// https://facebook.github.io/watchman/docs/expr/pcre.html
func Pcre(cs CaseSensitivity, scope ExpressionScope, pattern string) Expression {
	return exprSlice{caseName("pcre", cs), exprString(pattern), exprString("basename")}
}

// https://facebook.github.io/watchman/docs/expr/since.html
func SinceClock(clock string) Expression {
	return exprSlice{exprString("since"), exprString(clock), exprString("oclock")}
}

const (
	TimeFieldModified = "mtime"
	TimeFieldCreated  = "ctime"
)

func SinceTime(t time.Time, tf string) Expression {
	val := fmt.Sprintf("%d", t.Unix())
	return exprSlice{exprString("since"), exprString(val), exprString(tf)}
}

// https://facebook.github.io/watchman/docs/expr/size.html
func Size() Expression {
	return nil
}

// https://facebook.github.io/watchman/docs/expr/suffix.html
func Suffix(suffix string) Expression {
	return exprSlice{exprString("suffix"), exprString(suffix)}
}

const (
	TypeBlockSpecialFile     = "b"
	TypeCharacterSpecialFile = "c"
	TypeDirectory            = "d"
	TypeRegularFile          = "f"
	TypeNamedPipe            = "p"
	TypeSymbolicLink         = "l"
	TypeSocket               = "s"
	TypeSolarisDoor          = "D"
)

// https://facebook.github.io/watchman/docs/expr/type.html
func Type(typ string) Expression {
	return exprSlice{exprString("type"), exprString(typ)}
}
