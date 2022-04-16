package go_redash_query

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Connector struct {
}

func (c *Connector) CreateJob(queryId int, offset int, size int) (*JobInfo, error) {
	url := fmt.Sprintf("%s/queries/%d/results?api_key=%s", redashApiUrl, queryId, apiKey)
	queryData := QueryData{
		MaxAge: 0,
		Parameters: Parameters{
			Id:   offset,
			Size: size,
		},
	}
	jobInfo := JobInfo{}
	jsonBytes, err := json.Marshal(queryData)

	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(jsonBytes)

	req, err := http.NewRequest("POST", url, buff)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &jobInfo)
	if err != nil {
		return nil, err
	}

	if jobInfo.Message != "" {
		return nil, errors.New(jobInfo.Message)
	}

	return &jobInfo, nil
}

func (c *Connector) Fetch(job *JobInfo) (*FetchedData, error) {

	for {
		err := job.check()
		if err != nil {
			return nil, err
		}

		switch true {
		case job.isSuccess():
			return job.getFetchedData()
		case job.isError():
			return nil, err
		case job.isWait():
			time.Sleep(300 * time.Millisecond)
		}
	}

}
