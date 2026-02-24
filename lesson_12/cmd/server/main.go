package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/Nick2603/golang/lesson_12/internal/documentstore"
	"github.com/Nick2603/golang/lesson_12/internal/protocol"
)

type Server struct {
	store  *documentstore.Store
	logger *slog.Logger
}

func NewServer(logger *slog.Logger) *Server {
	return &Server{
		store:  documentstore.NewStoreWithLogger(logger),
		logger: logger,
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	s.logger.Info("client connected", "address", clientAddr)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		req, err := protocol.ReceiveRequest(reader)
		if err != nil {
			s.logger.Info("client disconnected", "address", clientAddr, "error", err)
			return
		}

		s.logger.Info("received command", "command", req.Command, "client", clientAddr)

		resp := s.handleRequest(req)

		if err := protocol.SendResponse(writer, resp); err != nil {
			s.logger.Error("failed to send response", "error", err)
			return
		}

		if req.Command == protocol.QuitCommand {
			s.logger.Info("client requested disconnect", "address", clientAddr)
			return
		}
	}
}

func (s *Server) handleRequest(req *protocol.Request) protocol.Response {
	switch req.Command {
	case protocol.PingCommand:
		return s.handlePing()
	case protocol.CreateCollectionCommand:
		return s.handleCreateCollection(req.Payload)
	case protocol.DeleteCollectionCommand:
		return s.handleDeleteCollection(req.Payload)
	case protocol.GetCollectionCommand:
		return s.handleGetCollection(req.Payload)
	case protocol.ListCollectionsCommand:
		return s.handleListCollections()
	case protocol.PutDocumentCommand:
		return s.handlePutDocument(req.Payload)
	case protocol.GetDocumentCommand:
		return s.handleGetDocument(req.Payload)
	case protocol.DeleteDocumentCommand:
		return s.handleDeleteDocument(req.Payload)
	case protocol.ListDocumentsCommand:
		return s.handleListDocuments(req.Payload)
	case protocol.CreateIndexCommand:
		return s.handleCreateIndex(req.Payload)
	case protocol.DeleteIndexCommand:
		return s.handleDeleteIndex(req.Payload)
	case protocol.QueryCommand:
		return s.handleQuery(req.Payload)
	case protocol.QuitCommand:
		return protocol.Response{Success: true, Data: "Goodbye!"}
	default:
		return protocol.Response{
			Success: false,
			Error:   fmt.Sprintf("unknown command: %s", req.Command),
		}
	}
}

func (s *Server) handlePing() protocol.Response {
	resp := protocol.PingResponse{Message: "pong"}
	data, _ := json.Marshal(resp)
	return protocol.Response{Success: true, Data: string(data)}
}

func (s *Server) handleCreateCollection(payload string) protocol.Response {
	var req protocol.CreateCollectionRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return protocol.Response{Success: false, Error: "invalid payload"}
	}

	cfg := &documentstore.CollectionConfig{PrimaryKey: req.PrimaryKey}
	_, err := s.store.CreateCollection(req.Name, cfg)
	if err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	return protocol.Response{Success: true, Data: "collection created"}
}

func (s *Server) handleDeleteCollection(payload string) protocol.Response {
	var req protocol.DeleteCollectionRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return protocol.Response{Success: false, Error: "invalid payload"}
	}

	if err := s.store.DeleteCollection(req.Name); err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	return protocol.Response{Success: true, Data: "collection deleted"}
}

func (s *Server) handleGetCollection(payload string) protocol.Response {
	var req protocol.GetCollectionRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return protocol.Response{Success: false, Error: "invalid payload"}
	}

	coll, err := s.store.GetCollection(req.Name)
	if err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	info := protocol.CollectionInfo{
		Name:       req.Name,
		PrimaryKey: coll.GetConfig().PrimaryKey,
	}
	data, _ := json.Marshal(info)
	return protocol.Response{Success: true, Data: string(data)}
}

func (s *Server) handleListCollections() protocol.Response {
	resp := protocol.ListCollectionsResponse{Collections: []protocol.CollectionInfo{}}
	data, _ := json.Marshal(resp)
	return protocol.Response{Success: true, Data: string(data)}
}

func (s *Server) handlePutDocument(payload string) protocol.Response {
	var req protocol.PutDocumentRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return protocol.Response{Success: false, Error: "invalid payload"}
	}

	coll, err := s.store.GetCollection(req.Collection)
	if err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	if err := coll.Put(req.Document); err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	return protocol.Response{Success: true, Data: "document added"}
}

func (s *Server) handleGetDocument(payload string) protocol.Response {
	var req protocol.GetDocumentRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return protocol.Response{Success: false, Error: "invalid payload"}
	}

	coll, err := s.store.GetCollection(req.Collection)
	if err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	doc, err := coll.Get(req.Key)
	if err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	resp := protocol.GetDocumentResponse{Document: *doc}
	data, _ := json.Marshal(resp)
	return protocol.Response{Success: true, Data: string(data)}
}

func (s *Server) handleDeleteDocument(payload string) protocol.Response {
	var req protocol.DeleteDocumentRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return protocol.Response{Success: false, Error: "invalid payload"}
	}

	coll, err := s.store.GetCollection(req.Collection)
	if err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	if err := coll.Delete(req.Key); err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	return protocol.Response{Success: true, Data: "document deleted"}
}

func (s *Server) handleListDocuments(payload string) protocol.Response {
	var req protocol.ListDocumentsRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return protocol.Response{Success: false, Error: "invalid payload"}
	}

	coll, err := s.store.GetCollection(req.Collection)
	if err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	docs := coll.List()
	resp := protocol.ListDocumentsResponse{Documents: docs}
	data, _ := json.Marshal(resp)
	return protocol.Response{Success: true, Data: string(data)}
}

func (s *Server) handleCreateIndex(payload string) protocol.Response {
	var req protocol.CreateIndexRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return protocol.Response{Success: false, Error: "invalid payload"}
	}

	coll, err := s.store.GetCollection(req.Collection)
	if err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	if err := coll.CreateIndex(req.FieldName); err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	return protocol.Response{Success: true, Data: "index created"}
}

func (s *Server) handleDeleteIndex(payload string) protocol.Response {
	var req protocol.DeleteIndexRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return protocol.Response{Success: false, Error: "invalid payload"}
	}

	coll, err := s.store.GetCollection(req.Collection)
	if err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	if err := coll.DeleteIndex(req.FieldName); err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	return protocol.Response{Success: true, Data: "index deleted"}
}

func (s *Server) handleQuery(payload string) protocol.Response {
	var req protocol.QueryRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return protocol.Response{Success: false, Error: "invalid payload"}
	}

	coll, err := s.store.GetCollection(req.Collection)
	if err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	docs, err := coll.Query(req.FieldName, req.Params)
	if err != nil {
		return protocol.Response{Success: false, Error: err.Error()}
	}

	resp := protocol.QueryResponse{Documents: docs}
	data, _ := json.Marshal(resp)
	return protocol.Response{Success: true, Data: string(data)}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	server := NewServer(logger)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
	defer listener.Close()

	logger.Info("server started", "address", ":8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("failed to accept connection", "error", err)
			continue
		}

		go server.handleConnection(conn)
	}
}
