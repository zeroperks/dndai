package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var state = "base"

type Player struct {
	Name      string
	DiscordId string
	Settings  string
}

var player = Player{
	Name:      "Player1",
	DiscordId: "",
	Settings:  "",
}

var conf = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/callback",
	ClientID:     "",
	ClientSecret: "",
	Scopes:       []string{"identify", "email"},
	Endpoint:     discord.Endpoint,
}

func main() {
	// Log the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	fmt.Printf("Current working directory: %s\n", cwd)

	r := mux.NewRouter()

	// Serve static files from the "static" directory
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	r.HandleFunc("/auth/discord", func(w http.ResponseWriter, r *http.Request) {
		url := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		fmt.Print("Redirecting to: " + url)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	r.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if code == "" || state == "" {
			http.Error(w, "Missing 'code' or 'state' in query parameters", http.StatusBadRequest)
			return
		}

		log.Printf("Received code: %s, state: %s\n", code, state)

		// Exchange authorization code for access token
		token, err := exchangeCodeForToken(conf.ClientID, conf.ClientSecret, code)
		if err != nil {
			log.Printf("Error exchanging code for token: %v", err)
			http.Error(w, "Error processing authorization", http.StatusInternalServerError)
			return
		}

		log.Printf("Access token: %s\n", token)
		fmt.Fprintf(w, "Authorization successful: %s", token)

		// Get user info
		req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return
		}

		// Set the content type to application/x-www-form-urlencoded
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		log.Println("Sending request to Discord token endpoint")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error during request: %v", err)
			return
		}
		defer resp.Body.Close()
		if err != nil {
			log.Printf("Error getting user: %v", err)
			http.Error(w, "Error getting user", http.StatusInternalServerError)
			return
		}
		// print body
		var user map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			log.Printf("Error decoding response: %v", err)
			return
		}

		// Print the entire response body for debugging
		responseBody, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			log.Printf("Error marshalling response body: %v", err)
			return

		}
		log.Printf("Response body: %s", responseBody)
		fmt.Printf("User: %+v\n", user["id"])
		player.DiscordId = user["id"].(string)

		// start refresh token loop
		go func() {
			for {
				// sleep for 10 minutes
				time.Sleep(60 * time.Minute)
				newToken, err := refreshAccessToken(conf.ClientID, conf.ClientSecret, token)
				if err != nil {
					log.Printf("Error refreshing token: %v", err)
					return
				}
				log.Printf("Refreshed token: %s", newToken)
				token = newToken
			}
		}()
	})

	fmt.Println("Server running at http://localhost:8080/")
	http.ListenAndServe(":8080", r)
}

func exchangeCodeForToken(clientID, clientSecret, code string) (string, error) {
	tokenURL := "https://discord.com/api/oauth2/token"

	// Create a url.Values object to hold the POST data
	data := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"http://localhost:8080/auth/callback"},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", err
	}

	// Set the content type to application/x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	log.Println("Sending request to Discord token endpoint")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error during request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("Status code: %d", resp.StatusCode)

	// Parse the response to get the access token
	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		log.Printf("Error decoding response: %v", err)
		return "", err
	}

	// Print the entire response body for debugging
	responseBody, err := json.MarshalIndent(tokenResponse, "", "  ")
	if err != nil {
		log.Printf("Error marshalling response body: %v", err)
		return "", err
	}
	log.Printf("Response body: %s", responseBody)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	if accessToken, ok := tokenResponse["access_token"].(string); ok {
		return accessToken, nil
	}
	return "", fmt.Errorf("invalid response from Discord token endpoint")
}

func refreshAccessToken(clientID, clientSecret, refreshToken string) (string, error) {
	tokenURL := "https://discord.com/api/oauth2/token"

	data := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Error creating request for token refresh: %v", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	log.Println("Refreshing access token")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error during token refresh request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("Status code: %d", resp.StatusCode)

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		log.Printf("Error decoding token refresh response: %v", err)
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("token refresh failed: %d", resp.StatusCode)
	}

	if newAccessToken, ok := tokenResponse["access_token"].(string); ok {
		return newAccessToken, nil
	}
	return "", fmt.Errorf("invalid response from Discord token endpoint during refresh")
}
