package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func addTestData() {
	testDeployments := []Deployment{
		{
			AppName:     "frontend-app",
			Environment: "dev",
			Status:      "in-progress",
			StartedAt:   time.Now(),
			CommitHash:  "abc123",
			Branch:      "main",
		},
		{
			AppName:     "backend-api",
			Environment: "staging",
			Status:      "successful",
			StartedAt:   time.Now().Add(-30 * time.Minute),
			CompletedAt: time.Now().Add(-25 * time.Minute),
			CommitHash:  "def456",
			Branch:      "develop",
		},
		{
			AppName:     "auth-service",
			Environment: "prod",
			Status:      "failed",
			StartedAt:   time.Now().Add(-1 * time.Hour),
			CompletedAt: time.Now().Add(-55 * time.Minute),
			CommitHash:  "ghi789",
			Branch:      "main",
		},
	}

	for _, deployment := range testDeployments {
		jsonData, err := json.Marshal(deployment)
		if err != nil {
			fmt.Printf("Error marshaling deployment: %v\n", err)
			continue
		}

		resp, err := http.Post("http://localhost:8080/api/deployments", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating deployment: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			fmt.Printf("Successfully added deployment for %s\n", deployment.AppName)
		} else {
			fmt.Printf("Failed to add deployment for %s. Status: %d\n", deployment.AppName, resp.StatusCode)
		}
	}
} 