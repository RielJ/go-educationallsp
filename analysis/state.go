package analysis

import (
	"fmt"
	"strings"

	"github.com/rielj/go-educationallsp/lsp"
)

type State struct {
	// Documents is a map of URIs to documents.
	Documents map[string]Document
}

type Document struct {
	Text string
}

func NewState() State {
	return State{
		Documents: make(map[string]Document),
	}
}

func getDiagnostics(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	for row, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "VS Code") {
			idx := strings.Index(line, "VS Code")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range: LineRange(
					row,
					idx,
					idx+len("VS Code"),
				),
				Severity: 1,
				Source:   "Common Sense",
				Message:  "VS Code is not allowed",
			})
		}
	}
	return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = Document{
		Text: text,
	}
	return getDiagnostics(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	doc, ok := s.Documents[uri]
	if !ok {
		return []lsp.Diagnostic{}
	}

	doc.Text = text
	s.Documents[uri] = doc
	return getDiagnostics(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	documents := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File %s, Characters %d", uri, len(documents.Text)),
		},
	}
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}

func (s *State) TextDocumentCodeAction(
	id int,
	uri string,
) lsp.TextDocumentCodeActionResponse {
	text := s.Documents[uri].Text

	actions := []lsp.CodeAction{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "Neovim",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS Code with Neovim",
				Edit: &lsp.WorkspaceEdit{
					Changes: replaceChange,
				},
			})

			censorChange := map[string][]lsp.TextEdit{}
			censorChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: strings.Repeat("*", len("VS Code")),
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Censor VS Code",
				Edit: &lsp.WorkspaceEdit{
					Changes: censorChange,
				},
			})
		}
	}

	return lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}
}

func (s *State) TextDocumentCompletion(
	id int,
	uri string,
) lsp.TextDocumentCompletionResponse {
	items := []lsp.CompletionItem{
		{
			Label:         "Neovim (BTW)",
			Detail:        "Very Cool Editor",
			Documentation: "For people who like to type `:wq`",
		},
	}

	response := lsp.TextDocumentCompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}
	return response
}

func LineRange(row, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      row,
			Character: start,
		},
		End: lsp.Position{
			Line:      row,
			Character: end,
		},
	}
}
