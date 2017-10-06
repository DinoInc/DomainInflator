package Thrift

var reservedSet map[string]bool

var resolvedList = []string{"BEGIN", "END", "__CLASS__", "__DIR__", "__FILE__", "__FUNCTION__",
	"__LINE__", "__METHOD__", "__NAMESPACE__", "abstract", "alias", "and", "args", "as", "assert",
	"begin", "break", "case", "catch", "class", "clone", "continue", "declare", "def", "default",
	"del", "delete", "do", "dynamic", "elif", "else", "elseif", "elsif", "end", "enddeclare",
	"endfor", "endforeach", "endif", "endswitch", "endwhile", "ensure", "except", "exec", "finally",
	"float", "for", "foreach", "from", "function", "global", "goto", "if", "implements", "import",
	"in", "inline", "instanceof", "interface", "is", "lambda", "module", "native", "new", "next",
	"nil", "not", "or", "package", "pass", "public", "print", "private", "protected", "raise", "redo",
	"rescue", "retry", "register", "return", "self", "sizeof", "static", "super", "switch", "synchronized",
	"then", "this", "throw", "transient", "try", "undef", "unless", "unsigned", "until", "use", "var",
	"virtual", "volatile", "when", "while", "with", "xor", "yield"}

func IsReservedWord(word string) bool {

	if reservedSet == nil {
		reservedSet = make(map[string]bool)
		for _, reservedWord := range resolvedList {
			reservedSet[reservedWord] = true
		}
	}

	_, isInReservedWord := reservedSet[word]
	return isInReservedWord
}
