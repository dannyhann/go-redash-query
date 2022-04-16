package go_redash_query

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Client struct {
	Connector
}

func CreateClient(redashApiUrl string, redashApiKey string) *Client {
	return &Client{
		Connector{
			redashApiUrl: redashApiUrl,
			redashApiKey: redashApiKey,
		},
	}
}

func (c *Client) CreateJob(queryId int) (*JobInfo, error) {
	url := fmt.Sprintf("%s/queries/%d/results", c.redashApiUrl, queryId)

	jobInfo := JobInfo{}

	body, err := c.request("POST", url, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &jobInfo)
	if err != nil {
		return nil, err
	}

	if jobInfo.Message != "" {
		return nil, errors.New(jobInfo.Message)
	}

	return &jobInfo, nil
}

func (c *Client) CreateJobWithQuery(queryId int, queryData QueryData) (*JobInfo, error) {
	url := fmt.Sprintf("%s/queries/%d/results", c.redashApiUrl, queryId)

	jsonBytes, err := json.Marshal(queryData)

	jobInfo := JobInfo{}

	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(jsonBytes)

	body, err := c.request("POST", url, buff)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &jobInfo)
	if err != nil {
		return nil, err
	}

	if jobInfo.Message != "" {
		return nil, errors.New(jobInfo.Message)
	}

	return &jobInfo, nil
}

func (c *Client) Fetch(job *JobInfo) (*FetchedData, error) {

	for {
		err := c.check(job)
		if err != nil {
			return nil, err
		}

		switch true {
		case job.isSuccess():
			return c.getFetchedData(job)
		case job.isError():
			return nil, err
		case job.isWait():
			time.Sleep(300 * time.Millisecond)
		}
	}

}

func (c *Client) check(j *JobInfo) error {
	url := fmt.Sprintf("%s/jobs/%s", c.redashApiUrl, j.Job.Id)

	body, err := c.request("GET", url, nil)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, j)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getFetchedData(j *JobInfo) (*FetchedData, error) {
	fetchedData := FetchedData{}
	url := fmt.Sprintf("%s/query_results/%d.json", c.redashApiUrl, j.QueryResultId)

	body, err := c.request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &fetchedData)
	if err != nil {
		return nil, err
	}

	return &fetchedData, err
}
