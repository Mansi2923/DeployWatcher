package main

import (
	"encoding/json"
	"fmt"
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
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	log.Printf("Received GitHub webhook: %+v\n", payload)

	// Check if this is a deployment or deployment_status event
	eventType := r.Header.Get("X-GitHub-Event")
	if eventType == "deployment" || eventType == "deployment_status" {
		// Extract deployment details
		var deployment Deployment
		if eventType == "deployment" {
			// For deployment events, extract from payload
			deployment = Deployment{
				ID:          fmt.Sprintf("%v", payload["id"]),
				AppName:     fmt.Sprintf("%v", payload["repository"].(map[string]interface{})["name"]),
				Environment: fmt.Sprintf("%v", payload["environment"]),
				Status:      "queued", // Default status for new deployments
				StartedAt:   time.Now(),
				CommitHash:  fmt.Sprintf("%v", payload["sha"]),
				Branch:      fmt.Sprintf("%v", payload["ref"]),
			}
		} else if eventType == "deployment_status" {
			// For deployment_status events, update existing deployment
			deploymentID := fmt.Sprintf("%v", payload["deployment"].(map[string]interface{})["id"])
			status := fmt.Sprintf("%v", payload["state"])
			for i, d := range deployments {
				if d.ID == deploymentID {
					deployments[i].Status = status
					if status == "success" || status == "failure" {
						deployments[i].CompletedAt = time.Now()
					}
					log.Printf("Updated deployment %s status to %s\n", deploymentID, status)
					w.WriteHeader(http.StatusOK)
					return
				}
			}
			log.Printf("Deployment %s not found\n", deploymentID)
			w.WriteHeader(http.StatusOK)
			return
		}

		// Add new deployment to the list
		deployments = append(deployments, deployment)
		log.Printf("Added new deployment: %+v\n", deployment)
	}

	w.WriteHeader(http.StatusOK)
}