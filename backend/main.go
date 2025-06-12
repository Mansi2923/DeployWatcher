package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Deployment represents a deployment event
type Deployment struct {
	ID          string    `json:"id"`
	AppName     string    `json:"appName"`
	Environment string    `json:"environment"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt,omitempty"`
	CommitHash  string    `json:"commitHash"`
	Branch      string    `json:"branch"`
}

// In-memory store for deployments (replace with database in production)
var deployments []Deployment

func main() {
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/deployments", getDeployments).Methods("GET")
	r.HandleFunc("/api/deployments", createDeployment).Methods("POST")
	r.HandleFunc("/api/deployments/{id}", updateDeployment).Methods("PUT")
	r.HandleFunc("/api/deployments/{id}", getDeployment).Methods("GET")
	r.HandleFunc("/api/github-webhook", githubWebhookHandler).Methods("POST")

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Add test data
	go addTestData()

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}

func getDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployments)
}

func createDeployment(w http.ResponseWriter, r *http.Request) {
	var deployment Deployment
	if err := json.NewDecoder(r.Body).Decode(&deployment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deployment.ID = generateID() // Implement this function
	deployment.StartedAt = time.Now()
	deployments = append(deployments, deployment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deployment)
}

func updateDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedDeployment Deployment
	if err := json.NewDecoder(r.Body).Decode(&updatedDeployment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, d := range deployments {
		if d.ID == id {
			deployments[i] = updatedDeployment
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedDeployment)
			return
		}
	}

	http.Error(w, "Deployment not found", http.StatusNotFound)
}

func getDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for _, d := range deployments {
		if d.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(d)
			return
		}
	}

	http.Error(w, "Deployment not found", http.StatusNotFound)
}

func generateID() string {
	return time.Now().Format("20060102150405")
} 

func githubWebhookHandler(w http.ResponseWriter, r *http.Request) {
    // For now, just log the event and return 200 OK
    var payload map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }
    log.Printf("Received GitHub webhook: %+v\n", payload)
    w.WriteHeader(http.StatusOK)
}