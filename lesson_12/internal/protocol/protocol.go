package protocol

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

func SendRequest(w *bufio.Writer, req Request) error {
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	if _, err := w.WriteString(encoded + "\n"); err != nil {
		return fmt.Errorf("failed to write request: %w", err)
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func ReceiveRequest(r *bufio.Reader) (*Request, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("failed to read request: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(line[:len(line)-1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	var req Request
	if err := json.Unmarshal(decoded, &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	return &req, nil
}

func SendResponse(w *bufio.Writer, resp Response) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	if _, err := w.WriteString(encoded + "\n"); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func ReceiveResponse(r *bufio.Reader) (*Response, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(line[:len(line)-1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var resp Response
	if err := json.Unmarshal(decoded, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}
