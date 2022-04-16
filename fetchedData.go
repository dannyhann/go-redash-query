package go_redash_query

import "time"

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
