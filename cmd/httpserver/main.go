package main

//http://localhost:42069/assets/exambankmultiplechoice.htm
//curl -I "https://httpbin.org"
//curl -X POST "https://httpbin.org/post"      -H "Content-Type: application/json"      -d '{"id": 101, "status": "active"}'
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
	// creates assest path
	candidates := []string{name, "../../" + name, "../" + name}
	// loops through assest to check if file found
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			// returns path if found
			return candidate
		}
	}
	// returns the name if not found
	return name
}

func main() {
	// runs the http server
	srv, err := server.Serve(port, handle)
	// error check
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	//close when functions ends
	defer srv.Close()
	//log
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
func handle(w *response.Writer, req *request.Request) {
	// find the target
	target := req.RequestLine.RequestTarget
	// HTML responses
	const html400 = `<html>
	<head>
		<title>400 Bad Request</title>
	</head>
	<body>
		<h1>Bad Request</h1>
		<p>malformed request.</p>
	</body>
</html>`

	const html500 = `<html>
	<head>
		<title>500 Internal Server Error</title>
	</head>
	<body>
		<h1>Internal Server Error</h1>
		<p>Server error.</p>
	</body>
</html>`

	const html200 = `<html>
	<head>
		<title>200 OK</title>
	</head>
	<body>
		<h1>Success!</h1>
		<p>Good Request.</p>
	</body>
</html>`

	var (
		hdrs headers.Headers
		err  error
	)
	//pritns the htmlt response base on status code use switch statment
	switch target {
	case "/yourproblem":
		// gets headers speacial headers
		hdrs = response.GetDefaultHeaders(len(html400))
		//sets content type
		hdrs.Set("content-type", "text/html")
		//writes the status line
		//check for errors
		err = w.WriteStatusLine(response.Code400)
		if err != nil {
			return
		}
		if err = w.WriteHeaders(hdrs); err != nil {
			return
		}
		// write the body
		_, _ = w.WriteBody([]byte(html400 + "\n"))
		return
	case "/myproblem":
		// get headers
		hdrs = response.GetDefaultHeaders(len(html500))
		// assign content type
		hdrs.Set("content-type", "text/html")
		// write the status line
		err = w.WriteStatusLine(response.Code500)
		//check for errors
		if err != nil {
			return
		}
		if err = w.WriteHeaders(hdrs); err != nil {
			return
		}
		// write the body
		_, _ = w.WriteBody([]byte(html500 + "\n"))
		return
	case "/video":
		// read video if avalible
		data, err := os.ReadFile(assetPath("assets/Replay.mp4"))
		// check for errors
		if err != nil {
			return
		}
		// get the speacial headers
		hdrs = response.GetDefaultHeaders(len(data))
		// set the content to mp4
		hdrs.Set("content-type", "video/mp4")
		// writes the status line wihtout error
		err = w.WriteStatusLine(response.Code200)
		// checks for errors
		if err != nil {
			return
		}
		if err = w.WriteHeaders(hdrs); err != nil {
			return
		}
		_, _ = w.WriteBody(data)
		return
	case "/assets/exambankmultiplechoice.htm":
		//Reads html in memory
		data, err := os.ReadFile(assetPath("assets/exambankmultiplechoice.htm"))
		//check
		if err != nil {
			return
		}
		// gets the defualt headers
		hdrs = response.GetDefaultHeaders(len(data))
		//set content type
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
			//porxyed path
			proxiedPath := strings.TrimPrefix(target, "/httpbin")
			//url
			// for test purposes
			upstreamURL := "https://httpbin.org" + proxiedPath
			//sends get request
			resp, err := http.Get(upstreamURL)
			//check for errors
			if err != nil {
				//error response
				log.Printf("Error proxying request to %s: %v", upstreamURL, err)
				hdrs = response.GetDefaultHeaders(len(html500))
				hdrs.Set("content-type", "text/html")
				_ = w.WriteStatusLine(response.Code500)
				_ = w.WriteHeaders(hdrs)
				_, _ = w.WriteBody([]byte(html500 + "\n"))
				return
			}
			//close the reader when function ends
			defer resp.Body.Close()

			// Write status line
			_ = w.WriteStatusLine(response.Code200)

			// write headers without Content-Length, but with Transfer-Encoding: chunked and trailers
			hdrs = headers.NewHeaders()
			//assignes defult headers
			hdrs["transfer-encoding"] = "chunked"
			hdrs["connection"] = "close"
			hdrs["trailer"] = "X-Content-SHA256, X-Content-Length"
			//check if content type is captial instead of lower case
			if contentType, ok := resp.Header["Content-Type"]; ok && len(contentType) > 0 {
				//assigns lowercase value
				hdrs["content-type"] = contentType[0]
			}
			// write the headers to connection
			_ = w.WriteHeaders(hdrs)

			// Read from upstream response and write chunked data to client
			// 4096 is standerd page size
			bodyBuf := make([]byte, 0, 4096)
			// buffer max length, 1024 is one binary kilobyte
			buf := make([]byte, 1024)
			for {
				// reads into buffer
				n, err := resp.Body.Read(buf)
				// prints bytes used
				fmt.Printf("%d\n", n)
				// add to bytes to bodybuff
				if n > 0 {
					bodyBuf = append(bodyBuf, buf[:n]...)
					_, _ = w.WriteChunkedBody(buf[:n])
				}
				//check for errors
				if err != nil {
					// if it is still parsing return err
					if err != io.EOF {
						log.Printf("Error reading upstream response: %v", err)
					}
					break
				}
			}

			// calculate trailers from the full raw response body
			// hash type
			hash := sha256.Sum256(bodyBuf)
			// create header type
			trailers := headers.NewHeaders()
			// calculate hashes from strings
			trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", hash))
			// assign content length
			trailers.Set("X-Content-Length", strconv.Itoa(len(bodyBuf)))
			// Write final chunk and trailer headers
			_ = w.WriteChunkedBodyDone(trailers)
			return
		}

		// default handler
		hdrs = response.GetDefaultHeaders(len(html200))
		//set content type
		hdrs.Set("content-type", "text/html")
		// writes status line
		err = w.WriteStatusLine(response.Code200)
		//checks for errors
		if err != nil {
			return
		}
		if err = w.WriteHeaders(hdrs); err != nil {
			return
		}
		//write the HTML body
		_, _ = w.WriteBody([]byte(html200 + "\n"))
		return
	}
}
