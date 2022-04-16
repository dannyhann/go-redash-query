package main

import (
	"fmt"
	"github.com/dannyhann/go-redash-query"
	"github.com/mitchellh/mapstructure"
	"log"
)

/* TEST REDASH QUERY

SELECT id, name, email
FROM test_user
WHERE id > {{id}}
ORDER BY id ASC
LIMIT {{size}}
*/

var (
	redashQueryId = 1 // Redash QueryId
	redashApiUrl  = ""
	redashAPiKey  = ""
)

type TestUserData struct {
	Id    int    `json:"id" mapstructure:"id"`
	Name  string `json:"name" mapstructure:"name"`
	Email string `json:"email" mapstructure:"email"`
}

func main() {

	client := go_redash_query.CreateClient(redashApiUrl, redashAPiKey)

	queryData := go_redash_query.QueryData{
		MaxAge: 10,
		Parameters: go_redash_query.Parameters{
			Id:   1,
			Size: 20,
		},
	}

	job, err := client.CreateJobWithQuery(redashQueryId, queryData)
	if err != nil {
		log.Fatal(err)
	}

	fetchedData, err := client.Fetch(job)
	if err != nil {
		log.Fatal(err)
	}

	users := make([]TestUserData, 0)
	rows := fetchedData.GetData()

	err = mapstructure.Decode(rows, &users)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		fmt.Printf("%d %s %s\n", user.Id, user.Name, user.Email)
	}
}