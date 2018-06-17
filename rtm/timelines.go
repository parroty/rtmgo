package rtm

type TimelinesService struct {
	HTTP *HTTP
}

type timelineResponse struct {
	Stat     string
	Err      ErrorResponse
	Timeline string
}

type timelineRootResponse struct {
	Rsp timelineResponse
}

func (s *TimelinesService) Create() (string, error) {
	result := new(timelineRootResponse)

	query := map[string]string{}

	err := s.HTTP.Request("rtm.timelines.create", query, &result)
	err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
	if err != nil {
		return "", err
	}

	return result.Rsp.Timeline, nil
}
