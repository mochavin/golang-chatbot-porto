# Instructions to run the Go Gemini streaming chatbot backend

1. Install Go 1.21 or newer.
2. Set your Gemini API key in the `.env` file (or as an environment variable `GEMINI_API_KEY`).
3. Place your CV PDF at `resource/CV.pdf`.
4. Run:

```
go mod tidy
go run main.go
```

5. Access the streaming endpoint:

```
curl -N "http://localhost:8080/chat?q=Tell me about your experience."
```

- The response will be streamed as Server-Sent Events (SSE).
- You can use this endpoint from a frontend for real-time chat.
