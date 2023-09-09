package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"snowflake-id/snowflake"
)

func main() {
	workerID, err := getWorkerID()
	if err != nil {
		log.Fatalf("Failed to get worker ID: %v", err)
	}

	sf, err := snowflake.NewSnowflake(workerID)
	if err != nil {
		log.Fatalf("Failed to initialize snowflake: %v", err)
	}

	http.HandleFunc("/generate-id", logRequest(func(w http.ResponseWriter, r *http.Request) {
		id, err := sf.GenerateId()
		if err != nil {
			http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
			return
		}

		log.Printf("Successfully generated id: %d", id)

		// Create a map to hold the ID
		response := map[string]int64{
			"id": id,
		}

		// Set the content type to application/json
		w.Header().Set("Content-Type", "application/json")

		// Use the json package to encode the map to JSON and write to the response
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getWorkerID() (int64, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return 0, err
	}

	// Convert the hostname to a hashcode
	hashCode := int64(0)
	for _, char := range hostname {
		hashCode = 31*hashCode + int64(char)
	}

	return hashCode & snowflake.MaxWorkerId, nil // Ensure workerID is within bounds
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
	}
}
