package main

import (
	"httpServer/internal/headers"
	"httpServer/internal/request"
	"httpServer/internal/response"
	"httpServer/internal/server"
	"log"
	"os"
	"os/signal"
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
