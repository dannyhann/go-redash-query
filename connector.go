package go_redash_query

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Connector struct {
	redashApiUrl string
	redashApiKey string
}

func (c *Connector) request(method, url string, buff io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %s", c.redashApiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
