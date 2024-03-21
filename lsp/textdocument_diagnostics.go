package lsp

type TextDocumentDiagnosticNotification struct {
	Notification
	Params TextDocumentDiagnosticParams `json:"params"`
}

type TextDocumentDiagnosticParams struct {
	URI         string       `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity"`
	Source   string `json:"source"`
	Message  string `json:"message"`
}
