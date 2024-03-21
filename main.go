package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/rielj/go-educationallsp/analysis"
	"github.com/rielj/go-educationallsp/lsp"
	"github.com/rielj/go-educationallsp/rpc"
)

func main() {
	logger := getLogger("/home/rielj/engineering/projects/go-educationallsp/educationallsp.log")
	logger.Println("Starting up")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Println("Failed to decode message", err)
			continue
		}

		handleMessage(logger, writer, state, method, contents)
	}
}

func handleMessage(
	logger *log.Logger,
	writer io.Writer,
	state analysis.State,
	method string,
	contents []byte,
) {
	logger.Printf("Received message: %s\n", method)

	switch method {
	case "initialize":
		var req lsp.InitializeRequest
		if err := json.Unmarshal(contents, &req); err != nil {
			logger.Println("Failed to decode initialize request", err)
		}
		logger.Printf(
			"Connected to: %s %s",
			req.Params.ClientInfo.Name,
			req.Params.ClientInfo.Version,
		)

		msg := lsp.NewInitializeResponse(req.ID)

		writeResponse(writer, msg)

		logger.Println("Sent initialize response")
	case "textDocument/didOpen":
		var notif lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &notif); err != nil {
			logger.Println("Failed to decode didOpen notification", err)
			return
		}

		logger.Printf(
			"Opened document: %s",
			notif.Params.TextDocument.URI,
		)

		diagnostics := state.OpenDocument(
			notif.Params.TextDocument.URI,
			notif.Params.TextDocument.Text,
		)
		writeResponse(writer, lsp.TextDocumentDiagnosticNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.TextDocumentDiagnosticParams{
				URI:         notif.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		})

	case "textDocument/didChange":
		var notif lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &notif); err != nil {
			logger.Println("Failed to decode didChange notification", err)
			return
		}

		logger.Printf(
			"Changed document: %s",
			notif.Params.TextDocument.URI,
		)

		for _, change := range notif.Params.ContentChanges {
			diagnostics := state.UpdateDocument(
				notif.Params.TextDocument.URI,
				change.Text,
			)

			writeResponse(writer, lsp.TextDocumentDiagnosticNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.TextDocumentDiagnosticParams{
					URI:         notif.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})

		}
	case "textDocument/hover":
		var req lsp.HoverRequest
		if err := json.Unmarshal(contents, &req); err != nil {
			logger.Println("Failed to decode hover request", err)
			return
		}

		logger.Printf(
			"Hovering over: %s",
			req.Params.TextDocumentPositionParams.TextDocument.URI,
		)

		response := state.Hover(
			req.ID,
			req.Params.TextDocumentPositionParams.TextDocument.URI,
			req.Params.TextDocumentPositionParams.Position,
		)

		writeResponse(writer, response)
	case "textDocument/definition":
		var req lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &req); err != nil {
			logger.Println("Failed to decode definition request", err)
			return
		}

		logger.Printf(
			"Defining: %s",
			req.Params.TextDocumentPositionParams.TextDocument.URI,
		)

		response := state.Definition(
			req.ID,
			req.Params.TextDocumentPositionParams.TextDocument.URI,
			req.Params.TextDocumentPositionParams.Position,
		)

		writeResponse(writer, response)

	case "textDocument/codeAction":
		var req lsp.CodeActionRequest
		if err := json.Unmarshal(contents, &req); err != nil {
			logger.Println("Failed to decode code action request", err)
			return
		}

		logger.Printf(
			"Code action: %s",
			req.Params.TextDocument.URI,
		)

		response := state.TextDocumentCodeAction(
			req.ID,
			req.Params.TextDocument.URI,
		)

		writeResponse(writer, response)
	case "textDocument/completion":
		var req lsp.TextDocumentCompletionRequest
		if err := json.Unmarshal(contents, &req); err != nil {
			logger.Println("Failed to decode code action request", err)
			return
		}

		logger.Printf(
			"Code action: %s",
			req.Params.TextDocument.URI,
		)

		response := state.TextDocumentCompletion(
			req.ID,
			req.Params.TextDocument.URI,
		)

		writeResponse(writer, response)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to open log file", err)
	}

	return log.New(logfile, "[educationallsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
