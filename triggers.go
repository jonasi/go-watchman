package watchman

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

func AllOf(expr ...Expression) Expression {
	return append(exprSlice{exprString("allof")}, expr...)
}

func AnyOf(expr ...Expression) Expression {
	return append(exprSlice{exprString("anyof")}, expr...)
}

func Not(expr Expression) Expression {
	return exprSlice{exprString("not"), expr}
}

func True() Expression {
	return exprString("true")
}

func False() Expression {
	return exprString("false")
}

func Suffix(suffix string) Expression {
	return exprSlice{exprString("suffix"), exprString(suffix)}
}

func Match(cs CaseSensitivity, scope ExpressionScope, pattern string) Expression {
	return exprSlice{caseName("match", cs), exprString(pattern), exprString(scope.string)}
}

func Pcre(cs CaseSensitivity, scope ExpressionScope, pattern string) Expression {
	return exprSlice{caseName("pcre", cs), exprString(pattern), exprString("basename")}
}

func Name(cs CaseSensitivity, scope ExpressionScope, names ...string) Expression {
	return exprSlice{caseName("name", cs), fromStringSlice(names), exprString("basename")}
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

func Type(typ string) Expression {
	return exprSlice{exprString("type"), exprString(typ)}
}

func Empty() Expression {
	return exprSlice{exprString("empty")}
}

func Exists() Expression {
	return exprSlice{exprString("exists")}
}

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
