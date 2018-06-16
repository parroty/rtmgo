package rtm

type GetFrobResultContent struct {
	Stat string
	Err  ErrorResponse
	Frob string
}

type GetFrobResult struct {
	Rsp GetFrobResultContent
}

type User struct {
	ID       string
	Username string
	Fullname string
}

type Auth struct {
	Token string
	Perms string
	User  User
}

type GetTokenResultContent struct {
	Stat string
	Err  ErrorResponse
	Auth Auth
}

type GetTokenResult struct {
	Rsp GetTokenResultContent
}

type AuthService struct {
	HTTP *HTTP
}

func (s *AuthService) GetFrob() (string, error) {
	result := new(GetFrobResult)

	query := map[string]string{}

	err := s.HTTP.Request("rtm.auth.getFrob", query, &result)
	err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
	if err != nil {
		return "", err
	}

	return result.Rsp.Frob, nil
}

func (s *AuthService) GetToken(frob string) (Auth, error) {
	result := new(GetTokenResult)

	query := map[string]string{}
	query["frob"] = frob

	err := s.HTTP.Request("rtm.auth.getToken", query, &result)
	err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
	if err != nil {
		return result.Rsp.Auth, err
	}
	return result.Rsp.Auth, nil
}
