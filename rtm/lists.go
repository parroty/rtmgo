package rtm

type List struct {
	ID        string
	Name      string
	Deleted   string
	Locked    string
	Archived  string
	Position  string
	Smart     string
	SortOrder string `json:"sort_order"`
}

type Lists struct {
	Lists []List `json:"list"`
}

type ListResponse struct {
	Stat  string
	Err   ErrorResponse
	Lists Lists
}

type ListRootResponse struct {
	Rsp ListResponse
}

type ListOperationResultContent struct {
	Stat        string
	Err         ErrorResponse
	Transaction Transaction
	List        List
}

type ListOperationResult struct {
	Rsp ListOperationResultContent
}

type ListsService struct {
	HTTP *HTTP
}

func (s *ListsService) Add(name string, timeline string) (Transaction, error) {
	result := new(ListOperationResult)

	query := map[string]string{}
	query["name"] = name
	query["timeline"] = timeline

	err := s.HTTP.Request("rtm.lists.add", query, &result)
	err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
	if err != nil {
		return result.Rsp.Transaction, err
	}

	return result.Rsp.Transaction, nil
}

func (s *ListsService) Delete(list List, timeline string) (Transaction, error) {
	result := new(ListOperationResult)

	query := map[string]string{}
	query["list_id"] = list.ID
	query["timeline"] = timeline

	err := s.HTTP.Request("rtm.lists.delete", query, &result)
	err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
	return result.Rsp.Transaction, err
}

func (s *ListsService) GetList() ([]List, error) {
	lists := make([]List, 0)
	result := new(ListRootResponse)

	query := map[string]string{}

	err := s.HTTP.Request("rtm.lists.getList", query, &result)
	err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
	if err != nil {
		return lists, err
	}

	return result.Rsp.Lists.Lists, nil
}
