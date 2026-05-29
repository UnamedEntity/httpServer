package main

import (
	"crypto/sha256"
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
	"strconv"
	"strings"
	"syscall"
)

const port = 42069

func assetPath(name string) string {
	candidates := []string{name, "../../" + name, "../" + name}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return name
}

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
	case "/video":
		data, err := os.ReadFile(assetPath("assets/vim.mp4"))
		if err != nil {
			return
		}
		hdrs = response.GetDefaultHeaders(len(data))
		hdrs.Set("content-type", "video/mp4")
		err = w.WriteStatusLine(response.Code200)
		if err != nil {
			return
		}
		if err = w.WriteHeaders(hdrs); err != nil {
			return
		}
		_, _ = w.WriteBody(data)
		return
	case "/assets/exambankmultiplechoice.htm":
		data, err := os.ReadFile(assetPath("assets/exambankmultiplechoice.htm"))
		if err != nil {
			return
		}
		hdrs = response.GetDefaultHeaders(len(data))
		hdrs.Set("content-type", "text/html; charset=utf-8")
		err = w.WriteStatusLine(response.Code200)
		if err != nil {
			return
		}
		if err = w.WriteHeaders(hdrs); err != nil {
			return
		}
		_, _ = w.WriteBody(data)
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

			// Write headers without Content-Length, but with Transfer-Encoding: chunked and trailers
			hdrs = headers.NewHeaders()
			hdrs["transfer-encoding"] = "chunked"
			hdrs["connection"] = "close"
			hdrs["trailer"] = "X-Content-SHA256, X-Content-Length"
			if contentType, ok := resp.Header["Content-Type"]; ok && len(contentType) > 0 {
				hdrs["content-type"] = contentType[0]
			}
			_ = w.WriteHeaders(hdrs)

			// Read from upstream response and write chunked data to client
			bodyBuf := make([]byte, 0, 4096)
			buf := make([]byte, 1024)
			for {
				n, err := resp.Body.Read(buf)
				fmt.Printf("%d\n", n)
				if n > 0 {
					bodyBuf = append(bodyBuf, buf[:n]...)
					_, _ = w.WriteChunkedBody(buf[:n])
				}
				if err != nil {
					if err != io.EOF {
						log.Printf("Error reading upstream response: %v", err)
					}
					break
				}
			}

			// Calculate trailers from the full raw response body
			hash := sha256.Sum256(bodyBuf)
			trailers := headers.NewHeaders()
			trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", hash))
			trailers.Set("X-Content-Length", strconv.Itoa(len(bodyBuf)))

			// Write final chunk and trailer headers
			_ = w.WriteChunkedBodyDone(trailers)
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
