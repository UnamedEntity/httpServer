package main

import (
	"fmt"
	"httpServer/internal/headers"
	"httpServer/internal/request"
	"httpServer/internal/response"
	"httpServer/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func main() {
	srv, err := server.Serve(port, handle)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer srv.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
func handle(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	// HTML responses
	const html400 = `<html>
	<head>
		<title>400 Bad Request</title>
	</head>
	<body>
		<h1>Bad Request</h1>
		<p>Your request honestly kinda sucked.</p>
	</body>
</html>`

	const html500 = `<html>
	<head>
		<title>500 Internal Server Error</title>
	</head>
	<body>
		<h1>Internal Server Error</h1>
		<p>Okay, you know what? This one is on me.</p>
	</body>
</html>`

	const html200 = `<html>
	<head>
		<title>200 OK</title>
	</head>
	<body>
		<h1>Success!</h1>
		<p>Your request was an absolute banger.</p>
	</body>
</html>`

	var (
		hdrs headers.Headers
		err  error
	)
	//pritns the htmlt response base on status code use switch statment
	switch target {
	case "/yourproblem":
		hdrs = response.GetDefaultHeaders(len(html400))
		hdrs.Set("content-type", "text/html")
		err = w.WriteStatusLine(response.Code400)
		if err != nil {
			return
		}
		if err = w.WriteHeaders(hdrs); err != nil {
			return
		}
		_, _ = w.WriteBody([]byte(html400 + "\n"))
		return
	case "/myproblem":
		hdrs = response.GetDefaultHeaders(len(html500))
		hdrs.Set("content-type", "text/html")
		err = w.WriteStatusLine(response.Code500)
		if err != nil {
			return
		}
		if err = w.WriteHeaders(hdrs); err != nil {
			return
		}
		_, _ = w.WriteBody([]byte(html500 + "\n"))
		return
	default:
		// Handle /httpbin/ proxy requests with chunked transfer encoding
		if strings.HasPrefix(target, "/httpbin/") {
			proxiedPath := strings.TrimPrefix(target, "/httpbin")
			upstreamURL := "https://httpbin.org" + proxiedPath

			resp, err := http.Get(upstreamURL)
			if err != nil {
				log.Printf("Error proxying request to %s: %v", upstreamURL, err)
				hdrs = response.GetDefaultHeaders(len(html500))
				hdrs.Set("content-type", "text/html")
				_ = w.WriteStatusLine(response.Code500)
				_ = w.WriteHeaders(hdrs)
				_, _ = w.WriteBody([]byte(html500 + "\n"))
				return
			}
			defer resp.Body.Close()

			// Write status line
			_ = w.WriteStatusLine(response.Code200)

			// Write headers without Content-Length, but with Transfer-Encoding: chunked
			hdrs = headers.NewHeaders()
			hdrs["transfer-encoding"] = "chunked"
			hdrs["connection"] = "close"
			if contentType, ok := resp.Header["Content-Type"]; ok && len(contentType) > 0 {
				hdrs["content-type"] = contentType[0]
			}
			_ = w.WriteHeaders(hdrs)

			// Read from upstream response and write chunked data to client
			buf := make([]byte, 1024)
			for {
				n, err := resp.Body.Read(buf)
				fmt.Printf("%d\n", n)
				if n > 0 {
					// Write chunk using chunked encoding format
					_, _ = w.WriteChunkedBody(buf[:n])
				}
				if err != nil {
					if err != io.EOF {
						log.Printf("Error reading upstream response: %v", err)
					}
					break
				}
			}

			// Write final chunk (size 0) to signal end of chunked response
			_, _ = w.WriteChunkedBodyDone()
			return
		}

		// Default handler
		hdrs = response.GetDefaultHeaders(len(html200))
		hdrs.Set("content-type", "text/html")
		err = w.WriteStatusLine(response.Code200)
		if err != nil {
			return
		}
		if err = w.WriteHeaders(hdrs); err != nil {
			return
		}
		_, _ = w.WriteBody([]byte(html200 + "\n"))
		return
	}
}
