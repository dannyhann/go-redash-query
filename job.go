package go_redash_query

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
