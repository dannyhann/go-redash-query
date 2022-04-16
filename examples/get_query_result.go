package examples

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	go_redash_query "go-redash-query"
	"log"
)

type IdOffset struct {
	Offset int
	Size   int
}

const queryId = 1 // Redash QueryId

type TestUserData struct {
	Id    int    `json:"id" mapstructure:"id"`
	Name  string `json:"name" mapstructure:"name"`
	Email string `json:"email" mapstructure:"email"`
}

func main() {
	offset := IdOffset{
		Offset: 1,
		Size:   20,
	}
	connector := go_redash_query.Connector{}
	job, err := connector.CreateJob(queryId, offset.Offset, offset.Size)
	if err != nil {
		log.Fatal(err)
	}

	fetchedData, err := connector.Fetch(job)
	if err != nil {
		log.Fatal(err)
	}

	users := make([]TestUserData, 0)
	for _, row := range fetchedData.GetData() {
		user := TestUserData{}
		err := mapstructure.Decode(row, &user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	fmt.Println(users)
}
