package lsp

type InitializeRequest struct {
	Request
	Params InitializeParams `json:"params"`
}

type InitializeParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
	// ... more fields
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo"`
}

type ServerCapabilities struct {
	// ... more fields
	TextDocumentSync   int            `json:"textDocumentSync"`
	HoverProvider      bool           `json:"hoverProvider"`
	DefinitionProvider bool           `json:"definitionProvider"`
	CodeActionProvider bool           `json:"codeActionProvider"`
	CompletionProvider map[string]any `json:"completionProvider"`
	DiagnosticProvider bool           `json:"diagnosticProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) *InitializeResponse {
	return &InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   1,
				HoverProvider:      true,
				DefinitionProvider: true,
				CodeActionProvider: true,
				CompletionProvider: map[string]any{},
				DiagnosticProvider: true,
			},
			ServerInfo: &ServerInfo{
				Name:    "educationallsp",
				Version: "0.0.0.0.0.0-beta1.final",
			},
		},
	}
}
