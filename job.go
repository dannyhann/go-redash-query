package go_redash_query

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	apiKey       = "" // REDASH_API_KEY
	redashApiUrl = "" // REDASH_URL
)

type Parameters struct {
	Id   int `json:"id" `
	Size int `json:"size"`
}

type QueryData struct {
	Parameters Parameters `json:"parameters"`
	MaxAge     int        `json:"max_age"`
}

type Job struct {
	Status        int    `json:"status"`
	Error         string `json:"error"`
	Id            string `json:"id"`
	QueryResultId int    `json:"query_result_id"`
	UpdatedAt     int    `json:"updated_at"`
}

type FetchedData struct {
	QueryResult struct {
		RetrievedAt time.Time `json:"retrieved_at"`
		QueryHash   string    `json:"query_hash"`
		Query       string    `json:"query"`
		Runtime     float64   `json:"runtime"`
		Data        struct {
			Rows    []interface{} `json:"rows"`
			Columns []struct {
				FriendlyName string  `json:"friendly_name"`
				Type         *string `json:"type"`
				Name         string  `json:"name"`
			} `json:"columns"`
		} `json:"data"`
		Id           int `json:"id"`
		DataSourceId int `json:"data_source_id"`
	} `json:"query_result"`
}

func (f *FetchedData) GetData() []interface{} {
	return f.QueryResult.Data.Rows
}

type JobInfo struct {
	Message string `json:"message,omitempty"`
	Job     `json:"job"`
}

func (j *JobInfo) isSuccess() bool {
	return j.Job.Status == 3
}

func (j *JobInfo) isWait() bool {
	return j.Job.Status == 1 || j.Job.Status == 2
}

func (j *JobInfo) isError() bool {
	return j.Job.Status == 4 || j.Job.Status == 5
}

func (j *JobInfo) check() error {
	url := fmt.Sprintf("%s/jobs/%s?api_key=%s", redashApiUrl, j.Job.Id, apiKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, j)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobInfo) getFetchedData() (*FetchedData, error) {
	fetchedData := FetchedData{}
	url := fmt.Sprintf("%s/query_results/%d.json?api_key=%s", redashApiUrl, j.QueryResultId, apiKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &fetchedData)
	if err != nil {
		return nil, err
	}

	return &fetchedData, err
}
