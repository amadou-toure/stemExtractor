package utils

import (
	"sync"
)


var (
	jobStatuses = make(map[string]string) // jobID -> status
	jobMutex    sync.Mutex
)

func SetJobStatus (status string, jobId string){
	jobMutex.Lock()
	defer jobMutex.Unlock()
	jobStatuses[jobId] = status
}

func GetJobStatus(jobID string) string {
	jobMutex.Lock()
	defer jobMutex.Unlock()
	if status, exists := jobStatuses[jobID]; exists {
		return status
	}
	return "not_found"
}
