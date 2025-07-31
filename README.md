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

---

## Frontend Integration Guide

To consume the streaming chat endpoint from your frontend (e.g., using Alpine.js, plain JavaScript, or any framework):

**Endpoint:**  
`GET http://localhost:8080/chat?q=YOUR_QUESTION`

**Response:**  
Server-Sent Events (SSE), streaming text data.

### Example: Using JavaScript (EventSource)

```javascript
const source = new EventSource(
  "http://localhost:8080/chat?q=Tell me about your experience."
);
source.onmessage = function (event) {
  // Append event.data to your chat UI
  console.log(event.data);
};
source.onerror = function (err) {
  source.close();
};
```

### Example: Using Fetch API (with ReadableStream)

```javascript
fetch("http://localhost:8080/chat?q=Tell me about your experience.").then(
  (response) => {
    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    function read() {
      reader.read().then(({ done, value }) => {
        if (done) return;
        const chunk = decoder.decode(value);
        // Append chunk to your chat UI
        console.log(chunk);
        read();
      });
    }
    read();
  }
);
```

### Notes

- Ensure your frontend runs on a different port (e.g., 3000) and the backend allows CORS if needed.
- For Alpine.js, use the above JavaScript in an Alpine component or as a custom function.
- The endpoint streams responses in real time for a chat-like experience.

---

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
