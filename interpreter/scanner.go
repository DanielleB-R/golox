package interpreter

type SourceScanner struct {
	source string
}

func NewSourceScanner(source string) SourceScanner {
	return SourceScanner{source}
}

func (*SourceScanner) ScanTokens() []string {
	return []string{}
}
