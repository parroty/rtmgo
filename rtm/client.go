package rtm

type Client struct {
	HTTP         *HTTP
	Auth         AuthService
	Timelines    TimelinesService
	Lists        ListsService
	Tasks        TasksService
	Transactions TransactionsService
}

func NewClient(token string, apiKey string) *Client {
	config := HTTP{
		Token:   token,
		ApiKey:  apiKey,
		BaseURL: "https://api.rememberthemilk.com"}
	client := &Client{HTTP: &config}
	client.Auth = AuthService{HTTP: &config}
	client.Timelines = TimelinesService{HTTP: &config}
	client.Lists = ListsService{HTTP: &config}
	client.Tasks = TasksService{HTTP: &config}
	client.Transactions = TransactionsService{HTTP: &config}
	return client
}

func (r *Client) GetAuthURL(frob string, perms []string) string {
	return r.HTTP.GetAuthURL(frob, perms)
}

func (r *Client) SetBaseURL(url string) {
	r.HTTP.BaseURL = url
}
