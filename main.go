package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dslipak/pdf"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

// Extracts plain text from a PDF file
func extractPDFText(path string) (string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not loaded")
	}
	// Load PDF context at startup
	cvText, err := extractPDFText("resource/CV.pdf")
	if err != nil {
		log.Fatalf("Failed to extract CV: %v", err)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set")
	}

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var wordBuffer bytes.Buffer
		ctx := context.Background()
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  apiKey,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			http.Error(w, "Failed to create Gemini client", 500)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Compose context: CV + user prompt/history
		contents := []*genai.Content{
			genai.NewContentFromText(`
				You are a smart, friendly assistant designed to help visitors learn more about the website owner's background.

				Please answer user questions using only the CV context provided below.
				Keep your responses:
				- Clear and relevant
				- Brief but informative
				- Friendly and professional

				Avoid making things up or including information not present in the CV.
				When information is missing, say so honestly.

				CV Context:
				`+cvText, "user"),
		}

		if r.Method == http.MethodPost {
			type Message struct {
				Role string `json:"role"`
				Text string `json:"text"`
			}
			type Payload struct {
				History []Message `json:"history"`
			}
			var payload Payload
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, "Invalid JSON payload", 400)
				return
			}
			for _, msg := range payload.History {
				contents = append(contents, genai.NewContentFromText(msg.Text, genai.Role(msg.Role)))
			}
		} else {
			// Legacy GET support
			prompt := r.URL.Query().Get("q")
			if prompt == "" {
				http.Error(w, "Missing 'q' query param", 400)
				return
			}
			contents = append(contents, genai.NewContentFromText(prompt, "user"))
		}

		stream := client.Models.GenerateContentStream(ctx, "gemini-2.5-flash", contents, nil)
		for resp := range stream {
			if resp == nil {
				log.Printf("stream returned nil response")
				continue
			}
			for _, cand := range resp.Candidates {
				if cand.Content != nil {
					for _, part := range cand.Content.Parts {
						if part.Text != "" {
							wordBuffer.WriteString(part.Text)
							text := wordBuffer.String()
							lastWordIdx := -1
							for i := len(text) - 1; i >= 0; i-- {
								if text[i] == ' ' || text[i] == '\n' || text[i] == '.' || text[i] == ',' || text[i] == '!' || text[i] == '?' {
									lastWordIdx = i + 1
									break
								}
							}
							if lastWordIdx > 0 {
								toSend := text[:lastWordIdx]
								fmt.Fprintf(w, "%s", toSend)
								w.(http.Flusher).Flush()
								wordBuffer.Reset()
								wordBuffer.WriteString(text[lastWordIdx:])
							}
						}
					}
				}
			}
		}
	})

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
