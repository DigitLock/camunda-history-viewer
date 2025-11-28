package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	camundaBaseURL string
	camundaUser    string
	camundaPass    string
	serverPort     string
)

// Camunda data structures
type ProcessInstance struct {
	ID                       string  `json:"id"`
	State                    string  `json:"state"`
	StartTime                string  `json:"startTime"`
	DurationInMillis         *int64  `json:"durationInMillis"`
	ProcessDefinitionVersion int     `json:"processDefinitionVersion"`
	BusinessKey              *string `json:"businessKey"`
	DeleteReason             *string `json:"deleteReason"`
}

type ActivityInstance struct {
	ActivityName     string `json:"activityName"`
	ActivityType     string `json:"activityType"`
	StartTime        string `json:"startTime"`
	EndTime          string `json:"endTime"`
	DurationInMillis int64  `json:"durationInMillis"`
	Canceled         bool   `json:"canceled"`
}

var templates *template.Template

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Environment variables
	camundaBaseURL = getEnv("CAMUNDA_BASE_URL", "http://localhost:8080/engine-rest")
	camundaUser = getEnv("CAMUNDA_USER", "demo")
	camundaPass = getEnv("CAMUNDA_PASSWORD", "demo")
	serverPort = getEnv("SERVER_PORT", "3000")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Initialize templates with custom functions
	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"formatTime": func(timeStr string) string {
			t, err := time.Parse(time.RFC3339, timeStr)
			if err != nil {
				return timeStr
			}
			return t.Format("2006-01-02 15:04:05")
		},
		"formatDuration": func(ms *int64) string {
			if ms == nil {
				return "N/A"
			}
			return fmt.Sprintf("%dms", *ms)
		},
	}).ParseGlob("templates/*.html"))

	// Routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/process/", handleProcessDetail)

	// Start server
	fmt.Printf("Server starting on http://localhost:%s\n", serverPort)
	fmt.Printf("Connected to Camunda at %s\n", camundaBaseURL)
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}

// Home page with process list
func handleHome(w http.ResponseWriter, r *http.Request) {
	// Fetch last 10 processes from Camunda
	processes, err := fetchProcesses()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching processes: %v", err), http.StatusInternalServerError)
		return
	}

	// Render template
	err = templates.ExecuteTemplate(w, "home.html", processes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Process detail page
func handleProcessDetail(w http.ResponseWriter, r *http.Request) {
	// Extract process ID from URL
	processID := r.URL.Path[len("/process/"):]

	// Fetch activity history
	activities, err := fetchProcessHistory(processID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching process history: %v", err), http.StatusInternalServerError)
		return
	}

	// Render template
	data := map[string]interface{}{
		"ProcessID":  processID,
		"Activities": activities,
	}

	err = templates.ExecuteTemplate(w, "process.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Fetch process list from Camunda
func fetchProcesses() ([]ProcessInstance, error) {
	url := fmt.Sprintf("%s/history/process-instance?processDefinitionKey=ticket-refund&sortBy=startTime&sortOrder=desc&maxResults=10", camundaBaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(camundaUser, camundaPass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var processes []ProcessInstance
	err = json.Unmarshal(body, &processes)
	if err != nil {
		return nil, err
	}

	return processes, nil
}

// Fetch process activity history
func fetchProcessHistory(processID string) ([]ActivityInstance, error) {
	url := fmt.Sprintf("%s/history/activity-instance?processInstanceId=%s&sortBy=startTime&sortOrder=asc", camundaBaseURL, processID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(camundaUser, camundaPass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var activities []ActivityInstance
	err = json.Unmarshal(body, &activities)
	if err != nil {
		return nil, err
	}

	return activities, nil
}
