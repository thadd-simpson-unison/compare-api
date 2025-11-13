package Handlers

import (
	"compare-api/Jsend"
	"compare-api/Models"
	"compare-api/Utilities"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const SUCCESS string = "success"
const RED string = "red"
const BLUE string = "blue"

func CompareHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		post(w, req)
	default:
		Jsend.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func post(w http.ResponseWriter, req *http.Request) {
	// TODO: get rid of this
	startTime := time.Now()

	file, _, err := req.FormFile("File")
	if err != nil {
		http.Error(w, "Invalid CSV file given.", http.StatusBadRequest)
		return
	}

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		http.Error(w, "Error reading the file.", http.StatusBadRequest)
		return
	}

	stringArray := []string{}
	for r := 0; r < len(records); r++ {
		fmt.Println(records[r][0])
		stringArray = append(stringArray, records[r][0])
	}

	// TODO: Call Red
	// TODO: Call Blue

	// TODO: get rid of this
	combinedResults := []Models.CombinedResult{
		{
			Source:      RED,
			Email:       "test1@unison.com",
			Status:      SUCCESS,
			RedResponse: nil,
			BlueData:    nil,
		},
		{
			Source:      BLUE,
			Email:       "test1@unison.com",
			Status:      SUCCESS,
			RedResponse: nil,
			BlueData:    nil,
		},
		{
			Source:      RED,
			Email:       "test2@unison.com",
			Status:      "error",
			RedResponse: nil,
			BlueData:    nil,
		},
		{
			Source:      BLUE,
			Email:       "test2@unison.com",
			Status:      "error",
			RedResponse: nil,
			BlueData:    nil,
		},
		{
			Source:      BLUE,
			Email:       "test3@unison.com",
			Status:      SUCCESS,
			RedResponse: nil,
			BlueData:    nil,
		},
	}
	duration := time.Since(startTime)

	//totalEmailCount := int32(len(stringArray))
	//combinedResults := slices.Concat(redResults, blueResults)
	res := buildResponse(3, combinedResults, duration, duration)

	Jsend.Success(w, res)

}

// Helper to make errors easier to handle
func createErrorResponse(email string, err error) Models.CombinedResult {
	return Models.CombinedResult{
		Source: RED,
		Email:  email,
		Status: err.Error(),
	}
}

// External API requiring multiple calls
func hitRedApi(emailArray []string) ([]Models.CombinedResult, time.Duration, error) {
	redBaseUrl := Utilities.GlobalConfig.RedBaseUrl
	results := []Models.CombinedResult{}

	// Hit the red api for all emails, and map responses as combinedResults
	start := time.Now()
	for e := 0; e < len(emailArray); e++ {
		email := emailArray[e]
		reader := strings.NewReader("")

		req, err := http.NewRequest("POST", redBaseUrl, reader)
		if err != nil {
			results = append(results, createErrorResponse(email, err))
			continue
		}

		if req.Response.StatusCode != 200 {
			results = append(results, createErrorResponse(email, err))
			continue
		}

		bytes, err := io.ReadAll(req.Body)
		if err != nil {
			results = append(results, createErrorResponse(email, err))
			continue
		}

		redResponse := Models.RedResponse{}
		err = json.Unmarshal(bytes, &redResponse)
		if err != nil {
			results = append(results, createErrorResponse(email, err))
			continue
		}

		results = append(results, Models.CombinedResult{
			Source:      RED,
			Email:       email,
			Status:      SUCCESS,
			RedResponse: &redResponse,
		})

	}
	redDuration := time.Since(start)

	return results, redDuration, nil
}

// External API that sends in single bulk response
func hitBlueApi() ([]Models.CombinedResult, time.Duration, error) {
	blueBaseUrl := Utilities.GlobalConfig.BlueBaseUrl
	results := []Models.CombinedResult{}

	// Hit the blue api for their CSV
	start := time.Now()
	reader := strings.NewReader("")
	req, err := http.NewRequest("POST", blueBaseUrl, reader)
	if err != nil {
		return []Models.CombinedResult{}, time.Duration(1), err
	}
	blueDuration := time.Since(start)

	// Read and convert blue's csv to combinedResults
	csvReader := csv.NewReader(req.Response.Body)
	records, err := csvReader.ReadAll()
	if err != nil {
		return []Models.CombinedResult{}, time.Duration(1), err
	}

	for r := 0; r < len(records); r++ {
		// Todo: make sure the indexes match the format of the document
		result := Models.CombinedResult{
			Source:   BLUE,
			Email:    records[r][0],
			Status:   SUCCESS,
			BlueData: &records[r],
		}

		results = append(results, result)
	}

	return results, blueDuration, nil
}

// Use the results from both to generate the report
func buildResponse(suppliedEmailCount int32, combinedResults []Models.CombinedResult, redDuration time.Duration, blueDuration time.Duration) Models.CompareResponse {
	// Count the successes and fails in each result list
	redCount := 0
	blueCount := 0
	redFails := []string{}
	blueFails := []string{}
	for i := 0; i < len(combinedResults); i++ {
		cr := combinedResults[i]
		if cr.Status == SUCCESS {
			// Add Counts
			if cr.Source == RED {
				redCount++
			}
			if cr.Source == BLUE {
				blueCount++
			}
		} else {
			// Add Fails
			if cr.Source == RED {
				redFails = append(redFails, cr.Email)
			}
			if cr.Source == BLUE {
				blueFails = append(blueFails, cr.Email)
			}
		}
	}

	return Models.CompareResponse{
		SuppliedEmailsCount: suppliedEmailCount,
		RedEmailCount:       int32(redCount),
		BlueEmailCount:      int32(blueCount),
		RedRequestsDuration: redDuration.Seconds(),
		BlueRequestDuration: blueDuration.Seconds(),
		RedEmailFails:       redFails,
		BlueEmailFails:      blueFails,
		CombinedResults:     combinedResults,
	}
}
