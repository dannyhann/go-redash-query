package go_redash_query

import (
	"github.com/mitchellh/mapstructure"
	"time"
)

type QueryResult struct {
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
}

type FetchedData struct {
	QueryResult `json:"query_result"`
}

func (f *FetchedData) GetData() []interface{} {
	return f.QueryResult.Data.Rows
}

func (f *FetchedData) GetDataWithStruct(output interface{}) error {
	err := mapstructure.Decode(f.GetData(), output)
	if err != nil {
		return err
	}
	return nil
}
