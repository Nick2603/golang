package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/Nick2603/golang/lesson_13/internal/documentstore"
	"github.com/Nick2603/golang/lesson_13/internal/protocol"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) sendRequest(req protocol.Request) (*protocol.Response, error) {
	if err := protocol.SendRequest(c.writer, req); err != nil {
		return nil, err
	}

	resp, err := protocol.ReceiveResponse(c.reader)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) printHelp() {
	help := `
Available commands:

Collection operations:
  create_collection <name> <primary_key>     - Create a new collection
  delete_collection <name>                   - Delete a collection
  get_collection <name>                      - Get collection info
  list_collections                           - List all collections

Document operations:
  put <collection> <json>                    - Add/update document
  get <collection> <key>                     - Get document by key
  delete <collection> <key>                  - Delete document
  list <collection>                          - List all documents in collection

Index operations:
  create_index <collection> <field>          - Create an index on a field
  delete_index <collection> <field>          - Delete an index
  query <collection> <field> [params_json]   - Query using an index

Other:
  ping                                       - Test connection
  help                                       - Show this help
  quit                                       - Exit

Examples:
  create_collection users id
  put users {"id":"user:1","name":"Alice"}
  get users user:1
  create_index users name
  query users name {"desc":true}
`
	fmt.Println(help)
}

func (c *Client) handleCommand(line string) error {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "help":
		c.printHelp()
		return nil

	case "ping":
		return c.handlePing()

	case "create_collection":
		if len(args) != 2 {
			fmt.Println("Usage: create_collection <name> <primary_key>")
			return nil
		}
		return c.handleCreateCollection(args[0], args[1])

	case "delete_collection":
		if len(args) != 1 {
			fmt.Println("Usage: delete_collection <name>")
			return nil
		}
		return c.handleDeleteCollection(args[0])

	case "get_collection":
		if len(args) != 1 {
			fmt.Println("Usage: get_collection <name>")
			return nil
		}
		return c.handleGetCollection(args[0])

	case "list_collections":
		return c.handleListCollections()

	case "put":
		if len(args) < 2 {
			fmt.Println("Usage: put <collection> <json>")
			return nil
		}
		jsonStr := strings.Join(args[1:], " ")
		return c.handlePut(args[0], jsonStr)

	case "get":
		if len(args) != 2 {
			fmt.Println("Usage: get <collection> <key>")
			return nil
		}
		return c.handleGet(args[0], args[1])

	case "delete":
		if len(args) != 2 {
			fmt.Println("Usage: delete <collection> <key>")
			return nil
		}
		return c.handleDelete(args[0], args[1])

	case "list":
		if len(args) != 1 {
			fmt.Println("Usage: list <collection>")
			return nil
		}
		return c.handleList(args[0])

	case "create_index":
		if len(args) != 2 {
			fmt.Println("Usage: create_index <collection> <field>")
			return nil
		}
		return c.handleCreateIndex(args[0], args[1])

	case "delete_index":
		if len(args) != 2 {
			fmt.Println("Usage: delete_index <collection> <field>")
			return nil
		}
		return c.handleDeleteIndex(args[0], args[1])

	case "query":
		if len(args) < 2 {
			fmt.Println("Usage: query <collection> <field> [params_json]")
			return nil
		}
		paramsJSON := "{}"
		if len(args) >= 3 {
			paramsJSON = strings.Join(args[2:], " ")
		}
		return c.handleQuery(args[0], args[1], paramsJSON)

	case "quit", "exit":
		return c.handleQuit()

	default:
		fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		return nil
	}
}

func (c *Client) handlePing() error {
	req := protocol.Request{Command: protocol.PingCommand}
	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.Success {
		var pingResp protocol.PingResponse
		json.Unmarshal([]byte(resp.Data), &pingResp)
		fmt.Printf("✓ %s\n", pingResp.Message)
	} else {
		fmt.Printf("✗ Error: %s\n", resp.Error)
	}
	return nil
}

func (c *Client) handleCreateCollection(name, primaryKey string) error {
	payload := protocol.CreateCollectionRequest{
		Name:       name,
		PrimaryKey: primaryKey,
	}
	payloadJSON, _ := json.Marshal(payload)

	req := protocol.Request{
		Command: protocol.CreateCollectionCommand,
		Payload: string(payloadJSON),
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.Success {
		fmt.Printf("✓ %s\n", resp.Data)
	} else {
		fmt.Printf("✗ Error: %s\n", resp.Error)
	}
	return nil
}

func (c *Client) handleDeleteCollection(name string) error {
	payload := protocol.DeleteCollectionRequest{Name: name}
	payloadJSON, _ := json.Marshal(payload)

	req := protocol.Request{
		Command: protocol.DeleteCollectionCommand,
		Payload: string(payloadJSON),
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.Success {
		fmt.Printf("✓ %s\n", resp.Data)
	} else {
		fmt.Printf("✗ Error: %s\n", resp.Error)
	}
	return nil
}

func (c *Client) handleGetCollection(name string) error {
	payload := protocol.GetCollectionRequest{Name: name}
	payloadJSON, _ := json.Marshal(payload)

	req := protocol.Request{
		Command: protocol.GetCollectionCommand,
		Payload: string(payloadJSON),
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.Success {
		var info protocol.CollectionInfo
		json.Unmarshal([]byte(resp.Data), &info)
		fmt.Printf("✓ Collection: %s, Primary Key: %s\n", info.Name, info.PrimaryKey)
	} else {
		fmt.Printf("✗ Error: %s\n", resp.Error)
	}
	return nil
}

func (c *Client) handleListCollections() error {
	req := protocol.Request{Command: protocol.ListCollectionsCommand}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.Success {
		var listResp protocol.ListCollectionsResponse
		json.Unmarshal([]byte(resp.Data), &listResp)
		fmt.Printf("✓ Collections (%d):\n", len(listResp.Collections))
		for _, coll := range listResp.Collections {
			fmt.Printf("  - %s (primary key: %s)\n", coll.Name, coll.PrimaryKey)
		}
	} else {
		fmt.Printf("✗ Error: %s\n", resp.Error)
	}
	return nil
}

func (c *Client) handlePut(collection, jsonStr string) error {
	var doc documentstore.Document
	if err := json.Unmarshal([]byte(jsonStr), &doc); err != nil {
		fmt.Printf("✗ Invalid JSON: %s\n", err)
		return nil
	}

	payload := protocol.PutDocumentRequest{
		Collection: collection,
		Document:   doc,
	}
	payloadJSON, _ := json.Marshal(payload)

	req := protocol.Request{
		Command: protocol.PutDocumentCommand,
		Payload: string(payloadJSON),
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.Success {
		fmt.Printf("✓ %s\n", resp.Data)
	} else {
		fmt.Printf("✗ Error: %s\n", resp.Error)
	}
	return nil
}

func (c *Client) handleGet(collection, key string) error {
	payload := protocol.GetDocumentRequest{
		Collection: collection,
		Key:        key,
	}
	payloadJSON, _ := json.Marshal(payload)

	req := protocol.Request{
		Command: protocol.GetDocumentCommand,
		Payload: string(payloadJSON),
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.Success {
		var getResp protocol.GetDocumentResponse
		json.Unmarshal([]byte(resp.Data), &getResp)
		prettyJSON, _ := json.MarshalIndent(getResp.Document, "", "  ")
		fmt.Printf("✓ Document:\n%s\n", string(prettyJSON))
	} else {
		fmt.Printf("✗ Error: %s\n", resp.Error)
	}
	return nil
}

func (c *Client) handleDelete(collection, key string) error {
	payload := protocol.DeleteDocumentRequest{
		Collection: collection,
		Key:        key,
	}
	payloadJSON, _ := json.Marshal(payload)

	req := protocol.Request{
		Command: protocol.DeleteDocumentCommand,
		Payload: string(payloadJSON),
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.Success {
		fmt.Printf("✓ %s\n", resp.Data)
	} else {
		fmt.Printf("✗ Error: %s\n", resp.Error)
	}
	return nil
}

func (c *Client) handleList(collection string) error {
	payload := protocol.ListDocumentsRequest{Collection: collection}
	payloadJSON, _ := json.Marshal(payload)

	req := protocol.Request{
		Command: protocol.ListDocumentsCommand,
		Payload: string(payloadJSON),
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.Success {
		var listResp protocol.ListDocumentsResponse
		json.Unmarshal([]byte(resp.Data), &listResp)
		fmt.Printf("✓ Documents (%d):\n", len(listResp.Documents))
		for i, doc := range listResp.Documents {
			prettyJSON, _ := json.MarshalIndent(doc, "  ", "  ")
			fmt.Printf("  [%d]\n  %s\n", i+1, string(prettyJSON))
		}
	} else {
		fmt.Printf("✗ Error: %s\n", resp.Error)
	}
	return nil
}

func (c *Client) handleCreateIndex(collection, field string) error {
	payload := protocol.CreateIndexRequest{
		Collection: collection,
		FieldName:  field,
	}
	payloadJSON, _ := json.Marshal(payload)

	req := protocol.Request{
		Command: protocol.CreateIndexCommand,
		Payload: string(payloadJSON),
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.Success {
		fmt.Printf("✓ %s\n", resp.Data)
	} else {
		fmt.Printf("✗ Error: %s\n", resp.Error)
	}
	return nil
}

func (c *Client) handleDeleteIndex(collection, field string) error {
	payload := protocol.DeleteIndexRequest{
		Collection: collection,
		FieldName:  field,
	}
	payloadJSON, _ := json.Marshal(payload)

	req := protocol.Request{
		Command: protocol.DeleteIndexCommand,
		Payload: string(payloadJSON),
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.Success {
		fmt.Printf("✓ %s\n", resp.Data)
	} else {
		fmt.Printf("✗ Error: %s\n", resp.Error)
	}
	return nil
}

func (c *Client) handleQuery(collection, field, paramsJSON string) error {
	var params documentstore.QueryParams
	if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
		fmt.Printf("✗ Invalid params JSON: %s\n", err)
		return nil
	}

	payload := protocol.QueryRequest{
		Collection: collection,
		FieldName:  field,
		Params:     params,
	}
	payloadJSON, _ := json.Marshal(payload)

	req := protocol.Request{
		Command: protocol.QueryCommand,
		Payload: string(payloadJSON),
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.Success {
		var queryResp protocol.QueryResponse
		json.Unmarshal([]byte(resp.Data), &queryResp)
		fmt.Printf("✓ Query results (%d):\n", len(queryResp.Documents))
		for i, doc := range queryResp.Documents {
			prettyJSON, _ := json.MarshalIndent(doc, "  ", "  ")
			fmt.Printf("  [%d]\n  %s\n", i+1, string(prettyJSON))
		}
	} else {
		fmt.Printf("✗ Error: %s\n", resp.Error)
	}
	return nil
}

func (c *Client) handleQuit() error {
	req := protocol.Request{Command: protocol.QuitCommand}
	c.sendRequest(req)
	return fmt.Errorf("quit")
}

func main() {
	client, err := NewClient("localhost:8080")
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Println("Connected to document store server")
	fmt.Println("Type 'help' for available commands")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if err := client.handleCommand(line); err != nil {
			if err.Error() == "quit" {
				fmt.Println("Goodbye!")
				break
			}
			fmt.Printf("Error: %v\n", err)
		}
	}
}
