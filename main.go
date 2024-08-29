package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"golang.org/x/oauth2"
)

func main() {
	// Use DefaultAzureCredential to authenticate with Managed Identity
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Failed to obtain a credential: %v", err)
	}

	// Get an access token for Microsoft Graph API
	token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{"https://graph.microsoft.com/.default"},
	})
	if err != nil {
		log.Fatalf("Failed to obtain token: %v", err)
	}

	// Create an OAuth2 token using the acquired token
	oauthToken := &oauth2.Token{
		AccessToken: token.Token,
		TokenType:   "Bearer",
	}

	// Create an HTTP client using the OAuth2 token
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(oauthToken),
		},
	}

	// Make a request to the Microsoft Graph API to get calendar events
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me/calendar/events", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Decode the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	// Print the calendar events
	events, ok := result["value"].([]interface{})
	if !ok {
		log.Fatalf("No events found in the response")
	}

	for _, event := range events {
		e := event.(map[string]interface{})
		fmt.Printf("Event: %s\n", e["subject"])
		fmt.Printf("Start: %s\n", e["start"].(map[string]interface{})["dateTime"])
		fmt.Printf("End: %s\n", e["end"].(map[string]interface{})["dateTime"])
		fmt.Println("---")
	}
}
