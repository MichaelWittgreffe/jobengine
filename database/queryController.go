package database

import (
	"fmt"
	"sort"
	"time"

	"github.com/MichaelWittgreffe/jobengine/crypto"
)

// QueryController defines an object used to make queries to the database
type QueryController interface {
	CreateQueue(name, accessKey string) error
	GetQueue(name, accessKey string) (*Queue, error)
	DeleteQueue(name, accessKey string) error
	AddJob(job *Job, queueName, accessKey string, sort bool) error
	GetJob(uid, queueName, accessKey string) (*Job, error)
	GetAllJobs(queueName, accessKey string) ([]*Job, error)
	UpdateJobStatus(uid, newStatus, queueName, accessKey string) error
}

// QueryControl object is used to make queries to the database
type QueryControl struct {
	db   *DBFile
	hash crypto.HashHandler
}

// NewQueryController is a constructor for the QueryController interface
func NewQueryController(db *DBFile, hasher crypto.HashHandler) QueryController {
	if db == nil {
		return nil
	}

	return &QueryControl{
		db:   db,
		hash: hasher,
	}
}

// CreateQueue creates a new queue entry
func (c *QueryControl) CreateQueue(name, accessKey string) error {
	if len(name) == 0 || len(accessKey) == 0 {
		return fmt.Errorf("Invalid Arg")
	}

	hashedKey, err := c.hash.Process(accessKey)
	if err != nil {
		return err
	}

	c.db.lock.Lock()
	defer c.db.lock.Unlock()

	if _, found := c.db.Queues[name]; found {
		return fmt.Errorf("Queue Exists")
	}

	c.db.Queues[name] = &Queue{
		Name:      name,
		AccessKey: hashedKey,
		Size:      0,
		Jobs:      make([]*Job, 0),
	}

	return nil
}

// GetQueue returns the given queue object entry or nil if it cannot be found
func (c *QueryControl) GetQueue(name, accessKey string) (*Queue, error) {
	if len(name) == 0 || len(accessKey) == 0 {
		return nil, fmt.Errorf("Invalid Arg")
	}

	hashedKey, err := c.hash.Process(accessKey)
	if err != nil {
		return nil, err
	}

	c.db.lock.Lock()
	defer c.db.lock.Unlock()

	if result, found := c.db.Queues[name]; found {
		if hashedKey == result.AccessKey {
			return result, nil
		}
		return nil, fmt.Errorf("Unauthorized")
	}

	return nil, nil
}

// DeleteQueue removes the given queue by name if the access token is correct
func (c *QueryControl) DeleteQueue(name, accessKey string) error {
	if len(name) == 0 || len(accessKey) == 0 {
		return fmt.Errorf("Invalid Arg")
	}

	hashedKey, err := c.hash.Process(accessKey)
	if err != nil {
		return err
	}

	c.db.lock.Lock()
	defer c.db.lock.Unlock()

	queue, found := c.db.Queues[name]
	if !found {
		return fmt.Errorf("Not Found")
	} else if queue.AccessKey != hashedKey {
		return fmt.Errorf("Unauthorized")
	}

	delete(c.db.Queues, name)
	return nil
}

// AddJob adds the given job to the given queue name in priority order (100 at head, 0 at tail)
func (c *QueryControl) AddJob(job *Job, queueName, accessKey string, sort bool) error {
	if job == nil || len(queueName) == 0 || len(accessKey) == 0 {
		return fmt.Errorf("Invalid Args")
	}

	hashedKey, err := c.hash.Process(accessKey)
	if err != nil {
		return err
	}

	c.db.lock.Lock()
	defer c.db.lock.Unlock()

	queue, found := c.db.Queues[queueName]
	if !found {
		return fmt.Errorf("Not Found")
	} else if queue.AccessKey != hashedKey {
		return fmt.Errorf("Unauthorized")
	}

	queue.Jobs = append(queue.Jobs, job)

	if sort {
		c.sortQueue(queue)
	}

	return nil
}

// GetJob returns the given job UID's entry, nil if job cannot be found
func (c *QueryControl) GetJob(uid, queueName, accessKey string) (*Job, error) {
	if len(uid) == 0 || len(queueName) == 0 || len(accessKey) == 0 {
		return nil, fmt.Errorf("Invalid Args")
	}

	hashedKey, err := c.hash.Process(accessKey)
	if err != nil {
		return nil, err
	}

	c.db.lock.Lock()
	defer c.db.lock.Unlock()

	queue, found := c.db.Queues[queueName]
	if !found {
		return nil, nil
	} else if queue.AccessKey != hashedKey {
		return nil, fmt.Errorf("Unauthorized")
	}

	for _, job := range queue.Jobs {
		if job.UID == uid {
			return job, nil
		}
	}

	return nil, nil
}

// GetNextJob returns the next job in the queue from the head that is not marked as 'inprogress', nil if not found or none avalible
func (c *QueryControl) GetNextJob(queueName, accessKey string) (*Job, error) {
	if len(queueName) == 0 || len(accessKey) == 0 {
		return nil, fmt.Errorf("Invalid Args")
	}

	hashedKey, err := c.hash.Process(accessKey)
	if err != nil {
		return nil, err
	}

	c.db.lock.Lock()
	defer c.db.lock.Unlock()

	queue, found := c.db.Queues[queueName]
	if !found {
		return nil, fmt.Errorf("Not Found")
	} else if queue.AccessKey != hashedKey {
		return nil, fmt.Errorf("Unauthorized")
	}

	for _, job := range queue.Jobs {
		if job.State == Inprogress {
			return job, nil
		}
	}

	return nil, nil
}

// GetAllJobs returns all the jobs for a given queue
func (c *QueryControl) GetAllJobs(queueName, accessKey string) ([]*Job, error) {
	if len(queueName) == 0 || len(accessKey) == 0 {
		return nil, fmt.Errorf("Invalid Args")
	}

	hashedKey, err := c.hash.Process(accessKey)
	if err != nil {
		return nil, err
	}

	c.db.lock.Lock()
	defer c.db.lock.Unlock()

	queue, found := c.db.Queues[queueName]
	if !found {
		return nil, fmt.Errorf("Not Found")
	} else if queue.AccessKey != hashedKey {
		return nil, fmt.Errorf("Unauthorized")
	}

	return queue.Jobs, nil
}

// UpdateJobStatus updates the given jobs status
func (c *QueryControl) UpdateJobStatus(uid, newStatus, queueName, accessKey string) error {
	if !c.validStatus(newStatus) || len(uid) == 0 || len(queueName) == 0 || len(accessKey) == 0 {
		return fmt.Errorf("Invalid Args")
	}

	hashedKey, err := c.hash.Process(accessKey)
	if err != nil {
		return err
	}

	c.db.lock.Lock()
	defer c.db.lock.Unlock()

	queue, found := c.db.Queues[queueName]
	if !found {
		return fmt.Errorf("Not Found")
	} else if queue.AccessKey != hashedKey {
		return fmt.Errorf("Unauthorized")
	}

	for _, job := range queue.Jobs {
		if job.UID == uid {
			job.State = newStatus
			job.LastUpdated = time.Now().Unix()
			return nil
		}
	}

	return fmt.Errorf("Not Found")
}

// UpdateQueue sorts the given queue by name, removes any jobs that are timed out etc
func (c *QueryControl) UpdateQueue(queueName string) error {
	if len(queueName) == 0 {
		return fmt.Errorf("Invalid Arg")
	}

	c.db.lock.Lock()
	defer c.db.lock.Unlock()

	queue, found := c.db.Queues[queueName]
	if !found {
		return fmt.Errorf(("Not Found"))
	}

	currentTime := time.Now().Unix()
	indexToDelete := make([]int, 0)

	for i, job := range queue.Jobs {
		if (job.State == Complete || job.State == Failed) && job.LastUpdated < (currentTime-job.KeepMinutes) {
			//remove complete/failed jobs that are outside the keep window
			indexToDelete = append(indexToDelete, i)
		} else if job.State == Inprogress && (job.LastUpdated < (currentTime - job.TimeoutMinutes)) {
			//mark as failed if no update within the timeout cut-off
			job.State = Failed
			job.LastUpdated = currentTime
		} else if (job.State == Queued) && (currentTime > job.TimeoutTime) {
			//delete queued jobs that are timed out
			indexToDelete = append(indexToDelete, i)
		}
	}

	for _, indexToDelete := range indexToDelete {
		c.deleteJobAtIndex(queue, indexToDelete)
	}

	c.sortQueue(queue)
	return nil
}

// deleteJobAtIndex removes the given index from the job queue, cleans up memory during delete
func (c *QueryControl) deleteJobAtIndex(queue *Queue, i int) {
	queueLenMinus := len(queue.Jobs) - 1
	if i < (queueLenMinus) {
		copy(queue.Jobs[i:], queue.Jobs[i+1:])
	}
	queue.Jobs[queueLenMinus] = nil
	queue.Jobs = queue.Jobs[:queueLenMinus]
	queue.Size = len(queue.Jobs)
}

//sortQueue orders the queue by priority, any other ordering should be maintained - must handle Lock outside of this function
func (c *QueryControl) sortQueue(in *Queue) {
	sort.Slice(in.Jobs, func(i, j int) bool {
		return in.Jobs[i].Priority > in.Jobs[j].Priority
	})
}

// validStatus checks the given status against the ValidStatus list, returns bool whether its valid
func (c *QueryControl) validStatus(status string) bool {
	for _, s := range ValidStatus {
		if s == status {
			return true
		}
	}
	return false
}
