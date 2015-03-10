package watchman

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

type exprSl []Expression

func (sl exprSl) noopExpr() {}

type exprString string

func (s exprString) noopExpr() {}

func caseName(name string, cs CaseSensitivity) exprString {
	if cs == CaseSensitive {
		return exprString(name)
	}

	return exprString("i" + name)
}

func AllOf(expr ...Expression) Expression {
	return append(exprSl{exprString("allof")}, expr...)
}

func AnyOf(expr ...Expression) Expression {
	return append(exprSl{exprString("anyof")}, expr...)
}

func Not(expr Expression) Expression {
	return exprSl{exprString("not"), expr}
}

func True() Expression {
	return exprString("true")
}

func False() Expression {
	return exprString("false")
}

func Suffix(suffix string) Expression {
	return exprSl{exprString("suffix"), exprString(suffix)}
}

func Match(pattern string, cs CaseSensitivity, scope ExpressionScope) Expression {
	return exprSl{caseName("match", cs), exprString(pattern), exprString(scope.string)}
}

func Pcre(pattern string, cs CaseSensitivity, scope ExpressionScope) Expression {
	return exprSl{caseName("pcre", cs), exprString(pattern), exprString("basename")}
}

func Name(scope ExpressionScope, cs CaseSensitivity, names ...string) Expression {
	return nil
}

func Type() {

}

func Empty() {

}

func Exists() {

}

func Since() {

}
