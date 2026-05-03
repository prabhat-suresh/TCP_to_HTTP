package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func writeBadRequestResponse(w *response.Writer) {
	w.WriteStatusLine(response.StatusBadRequest)
	w.WriteHeaders(headers.Headers{"Content-Type": "text/html"})
	w.WriteBody([]byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`))
}

func proxy(w *response.Writer, req *request.Request) {
	suffix := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	resp, err := http.Get("https://httpbin.org/" + suffix)
	if err != nil {
		log.Print("Couldn't get http response from httpbin.org")
		writeBadRequestResponse(w)
		return
	}
	w.WriteStatusLine(response.StatusOk)
	req.Headers.Delete("Content-Length")
	req.Headers.Set("Transfer-Encoding", "chunked")
	req.Headers.Set("Trailer", "X-Content-SHA256")
	req.Headers.Set("Trailer", "X-Content-Length")
	w.WriteHeaders(req.Headers)
	body := resp.Body
	buf := make([]byte, 1024)
	fullResponse := make([]byte, 0, 1024)
	for {
		n, err := body.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Print("Error reading http response's body from httpbin.org")
			return
		}
		log.Printf("%v bytes read from httpbin.org:\n", n)
		w.WriteChunkedBody(buf[:n])
		fullResponse = append(fullResponse, buf...)
	}
	hash := sha256.Sum256(fullResponse)
	req.Headers.Set("X-Content-SHA256", fmt.Sprintf("%x", hash))
	req.Headers.Set("X-Content-Length", fmt.Sprintf("%d", len(fullResponse)))
	w.WriteChunkedBodyDone(req.Headers)
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxy(w, req)
		return
	}
	switch req.RequestLine.RequestTarget {
	case "/video":
		video, _ := os.ReadFile("/Users/prabhat.suresh/Golang/TCP_to_HTTP/assets/vim.mp4")
		w.WriteStatusLine(response.StatusOk)
		h := headers.Headers{"Content-Type": "video/mp4"}
		h.Set("Content-Length", fmt.Sprintf("%d", len(video)))
		w.WriteHeaders(h)
		w.WriteBody(video)
	case "/yourproblem":
		writeBadRequestResponse(w)
	case "/myproblem":
		w.WriteStatusLine(response.StatusInternalServerError)
		w.WriteHeaders(headers.Headers{"Content-Type": "text/html"})
		w.WriteBody([]byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`))
	default:
		w.WriteStatusLine(response.StatusOk)
		w.WriteHeaders(headers.Headers{"Content-Type": "text/html"})
		w.WriteBody([]byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`))
	}
}

func main() {
	server, err := server.Serve(handler, port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
