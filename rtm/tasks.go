package rtm

type Chunk struct {
	ID         string
	Due        string
	HasDueTime string `json:"has_due_time"`
	Added      string
	Completed  string
	Deleted    string
	Priority   string
	Postponed  string
	Estimate   string
}

type Task struct {
	ID           string
	Created      string
	Modified     string
	Name         string
	Source       string
	URL          string
	LocationID   string `json:"location_id"`
	Tags         interface{}
	Participants []string
	Notes        []string
	Chunks       []Chunk `json:"task"`
	ListID       string  // This field is not included in the API response, and should be manually set.
}

type TaskList struct {
	ID    string
	Tasks []Task `json:"taskseries"`
}

type TaskListResponse struct {
	Rev   string
	Lists []TaskList `json:"list"`
}

type TaskOperationResultContent struct {
	Stat        string
	Err         ErrorResponse
	Transaction Transaction
	List        TaskList
}

type TaskOperationResult struct {
	Rsp TaskOperationResultContent
}

type GetTasksResultContent struct {
	Stat  string
	Err   ErrorResponse
	Tasks TaskListResponse
}

type GetTasksResult struct {
	Rsp GetTasksResultContent
}

type TasksService struct {
	HTTP *HTTP
}

func (s *TasksService) Add(timeline string, name string) ([]Task, *Transaction, error) {
	result := new(TaskOperationResult)

	query := map[string]string{}
	query["name"] = name
	query["timeline"] = timeline

	err := s.HTTP.Request("rtm.tasks.add", query, &result)
	err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
	return result.Rsp.List.Tasks, &result.Rsp.Transaction, err
}

func (s *TasksService) AddTags(task Task, timeline string, tags string) ([]Transaction, error) {
	transactions := make([]Transaction, 0)

	for _, chunk := range task.Chunks {
		query := map[string]string{}
		query["list_id"] = task.ListID
		query["taskseries_id"] = task.ID
		query["task_id"] = chunk.ID
		query["timeline"] = timeline
		query["tags"] = tags

		result := new(TaskOperationResult)
		err := s.HTTP.Request("rtm.tasks.addTags", query, &result)
		err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, result.Rsp.Transaction)
	}

	return transactions, nil
}

func (s *TasksService) Complete(task Task, timeline string) ([]Transaction, error) {
	results := make([]Transaction, 0)

	for _, chunk := range task.Chunks {
		result := new(TaskOperationResult)

		query := map[string]string{}
		query["list_id"] = task.ListID
		query["task_id"] = chunk.ID
		query["taskseries_id"] = task.ID
		query["timeline"] = timeline

		err := s.HTTP.Request("rtm.tasks.complete", query, &result)
		err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
		if err != nil {
			return results, err
		}

		results = append(results, result.Rsp.Transaction)
	}
	return results, nil
}

func (s *TasksService) Delete(task Task, timeline string) ([]Transaction, error) {
	results := make([]Transaction, 0)

	for _, chunk := range task.Chunks {
		result := new(TaskOperationResult)

		query := map[string]string{}
		query["list_id"] = task.ListID
		query["task_id"] = chunk.ID
		query["taskseries_id"] = task.ID
		query["timeline"] = timeline

		err := s.HTTP.Request("rtm.tasks.delete", query, &result)
		err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
		if err != nil {
			return results, err
		}

		results = append(results, result.Rsp.Transaction)
	}
	return results, nil
}

func (s *TasksService) GetList(listID string, filter string, lastSync string) ([]Task, error) {
	tasks := make([]Task, 0)
	result := new(GetTasksResult)

	query := map[string]string{}
	if listID != "" {
		query["list_id"] = listID
	}
	if filter != "" {
		query["filter"] = filter
	}
	if lastSync != "" {
		query["last_sync"] = lastSync
	}

	err := s.HTTP.Request("rtm.tasks.getList", query, &result)
	err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
	if err != nil {
		return tasks, err
	}

	for _, list := range result.Rsp.Tasks.Lists {
		for _, task := range list.Tasks {
			task.ListID = list.ID
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (s *TasksService) SetDueDate(task Task, timeline string, due string) ([]Transaction, error) {
	results := make([]Transaction, 0)

	for _, chunk := range task.Chunks {
		result := new(TaskOperationResult)

		query := map[string]string{}
		query["list_id"] = task.ListID
		query["task_id"] = chunk.ID
		query["taskseries_id"] = task.ID
		query["timeline"] = timeline
		query["due"] = due

		err := s.HTTP.Request("rtm.tasks.setDueDate", query, &result)
		err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
		if err != nil {
			return results, err
		}
		results = append(results, result.Rsp.Transaction)
	}
	return results, nil
}

func (s *TasksService) Uncomplete(task Task, timeline string) ([]Transaction, error) {
	results := make([]Transaction, 0)

	for _, chunk := range task.Chunks {
		result := new(TaskOperationResult)

		query := map[string]string{}
		query["list_id"] = task.ListID
		query["task_id"] = chunk.ID
		query["taskseries_id"] = task.ID
		query["timeline"] = timeline

		err := s.HTTP.Request("rtm.tasks.uncomplete", query, &result)
		err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
		if err != nil {
			return results, err
		}

		results = append(results, result.Rsp.Transaction)
	}
	return results, nil
}
