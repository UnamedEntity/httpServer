package main

import (
	"bytes"
	"strings"
	"testing"

	"httpServer/internal/request"
	"httpServer/internal/response"
)

func TestHandleServesExamBankPage(t *testing.T) {
	var buf bytes.Buffer
	w := response.NewWriter(&buf)

	req := &request.Request{
		RequestLine: request.RequestLine{
			RequestTarget: "/assets/exambankmultiplechoice.htm",
		},
	}

	handle(w, req)

	got := buf.String()
	if !strings.Contains(got, "Exam Bank") {
		t.Fatalf("expected response to contain exam page content, got: %q", got)
	}
	if !strings.Contains(got, "HTTP/1.1 200 OK") {
		t.Fatalf("expected 200 OK response, got: %q", got)
	}
}
