package jobm_support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var JobManager *JobManagerService

type JobManagerService struct {
	host string
	app  string
	key  string
}

type JobManagerPayload struct {
	Event string                 `json:"event"`
	Body  map[string]interface{} `json:"body"`
}

type JobManagerPayloadResponse struct {
	Return      map[string]interface{} `json:"return"`
	Status      string                 `json:"status"`
	Status_code string                 `json:"status_code"`
}

func (c *JobManagerService) Payload(props JobManagerPayload) (*string, error) {
	jsonPayload, err := json.Marshal(props)
	if err != nil {
		// Handle JSON marshaling error
		fmt.Println("json.Marshal - err :: ", err)
		return nil, err
	}
	var payload = bytes.NewBuffer(jsonPayload)
	request, err := http.NewRequest("POST", c.host+"/api/front/job_record/create", payload)
	if err != nil {
		// Handle request creation error
		fmt.Println("http.NewRequest - err :: ", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	// Add Bearer token to Authorization header
	token := c.app + "|" + c.key
	request.Header.Set("Authorization", "Bearer "+token)

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("client.Do - err :: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Process response
	// For example, read and decode JSON response
	var responseData JobManagerPayloadResponse
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("json.NewDecoder - err :: ", err)
		return nil, err
	}
	Return := responseData.Return
	UUID := Return["uuid"].(string)
	return &UUID, nil
}

type JobManagerPayloadGetLogs struct {
	Take int64  `json:"take"`
	Uuid string `json:"uuid"`
}

type JobManagerPayloadGetLogsResponse struct {
	Return      []string `json:"return"`
	Status      string   `json:"status"`
	Status_code string   `json:"status_code"`
}

func (c *JobManagerService) GetLastLogs(props JobManagerPayloadGetLogs) ([]string, error) {
	jsonPayload, err := json.Marshal(props)
	if err != nil {
		// Handle JSON marshaling error
		fmt.Println("json.Marshal - err :: ", err)
		return nil, err
	}
	var payload = bytes.NewBuffer(jsonPayload)
	request, err := http.NewRequest("POST", c.host+"/api/front/job_log/job_logs", payload)
	if err != nil {
		// Handle request creation error
		fmt.Println("http.NewRequest - err :: ", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	// Add Bearer token to Authorization header
	token := c.app + "|" + c.key
	request.Header.Set("Authorization", "Bearer "+token)

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("client.Do - err :: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Process response
	// For example, read and decode JSON response
	var responseData JobManagerPayloadGetLogsResponse
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("json.NewDecoder - err :: ", err)
		return nil, err
	}
	Return := responseData.Return
	return Return, nil
}

type JobManagerJobStopResponse struct {
	Return      string `json:"return"`
	Status      string `json:"status"`
	Status_code string `json:"status_code"`
}

func (c *JobManagerService) Stop(job_id string) (*string, error) {
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", c.host+"/api/front/job_record/"+job_id+"/stop", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Add Bearer token to Authorization header
	token := c.app + "|" + c.key
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	result := JobManagerJobStopResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("json.Unmarshal - err :: ", err)
		return nil, err
	}

	// Print the response body
	return &result.Return, nil
}

func (c *JobManagerService) Terminate(job_id string) (*string, error) {
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", c.host+"/api/front/job_record/"+job_id+"/terminate", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Add Bearer token to Authorization header
	token := c.app + "|" + c.key
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	result := JobManagerJobStopResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("json.Unmarshal - err :: ", err)
		return nil, err
	}

	// Print the response body
	return &result.Return, nil
}

// New own type

func InitJobManager() JobManagerService {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fmt.Println("Connecting job manager!")

	// Access the environment variables
	jobMApp := os.Getenv("JOB_M_APP")
	JobMKey := os.Getenv("JOB_M_KEY")
	jobMHost := os.Getenv("JOB_M_HOST")

	jobManagerService := JobManagerService{
		host: jobMHost,
		app:  jobMApp,
		key:  JobMKey,
	}

	JobManager = &jobManagerService

	// Your application logic here...
	fmt.Println("Successfully connected to the database!")

	return jobManagerService
}
