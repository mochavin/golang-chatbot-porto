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
const source = new EventSource('http://localhost:8080/chat?q=Tell me about your experience.');
source.onmessage = function(event) {
  // Append event.data to your chat UI
  console.log(event.data);
};
source.onerror = function(err) {
  source.close();
};
```

### Example: Using Fetch API (with ReadableStream)

```javascript
fetch('http://localhost:8080/chat?q=Tell me about your experience.')
  .then(response => {
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
  });
```

### Notes

- Ensure your frontend runs on a different port (e.g., 3000) and the backend allows CORS if needed.
- For Alpine.js, use the above JavaScript in an Alpine component or as a custom function.
- The endpoint streams responses in real time for a chat-like experience.
