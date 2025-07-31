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

```sh
curl -N -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"history":[{"role":"user","text":"Tell me about your experience."},{"role":"model","text":"I have 5 years of experience in software engineering..."},{"role":"user","text":"What programming languages do you use?"}]}'
```

- The response will be streamed as Server-Sent Events (SSE).
- You can use this endpoint from a frontend for real-time chat.

## Using POST /chat with Context (History)

You can send a POST request to `/chat` with a JSON payload containing a `history` array for context-aware conversations.

**Endpoint:**  
`POST http://localhost:8080/chat`

**Payload Example:**

```json
{
  "history": [
    { "role": "user", "text": "Tell me about your experience." },
    {
      "role": "model",
      "text": "I have 5 years of experience in software engineering..."
    },
    { "role": "user", "text": "What programming languages do you use?" }
  ]
}
```

**Curl Example:**

```sh
curl -N -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"history":[{"role":"user","text":"Tell me about your experience."},{"role":"model","text":"I have 5 years of experience in software engineering..."},{"role":"user","text":"What programming languages do you use?"}]}'
```

- The response will be streamed as Server-Sent Events (SSE).
- The `history` array should alternate between `"user"` and `"model"` roles for multi-turn chat.
- The backend will always prepend the CV context automatically.

---

## Running with Docker

You can run the backend using Docker:

```sh
docker build -t chatbot-porto .
docker run --env-file .env -p 8080:8080 chatbot-porto
```

- Make sure your `.env` file and `resource/CV.pdf` are present in the build context.
- The container will expose port 8080 for API access.
- Adjust the `-p` flag if you want to use a different host port.
