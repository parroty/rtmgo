package rtm

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

type HTTP struct {
	Token   string
	ApiKey  string
	BaseURL string
}

type ErrorResponse struct {
	Msg  string
	Code string
}

func (err *ErrorResponse) Error() string {
	return fmt.Sprintf("%s (code=%s)", err.Msg, err.Code)
}

const baseURL = "https://api.rememberthemilk.com/services/rest"
const authURL = "https://api.rememberthemilk.com/services/auth"

func (c *HTTP) GetAuthURL(frob string, perms []string) string {
	query := map[string]string{}
	query["api_key"] = c.ApiKey
	query["perms"] = strings.Join(perms, ",")
	query["frob"] = frob

	sig := signParams(query)

	m := url.Values{}
	for k, v := range query {
		m.Add(k, v)
	}
	m.Add("api_sig", sig)

	return authURL + "?" + m.Encode()
}

func (c *HTTP) VerifyResponse(err error, stat string, resp ErrorResponse) error {
	if err != nil {
		return err
	}
	if stat == "fail" {
		return &ErrorResponse{Msg: resp.Msg, Code: resp.Code}
	}
	return nil
}

func (c *HTTP) Request(method string, params map[string]string, result interface{}) error {
	resp, err := c.doRequest(method, params, c.Token, c.ApiKey)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if result != nil {
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatal(err)
			return err
		}
	}

	//TODO: support debug mode
	//fmt.Println(string(body))
	return nil
}

func (c *HTTP) doRequest(method string, params map[string]string, token string, apiKey string) (*http.Response, error) {
	baseParams := map[string]string{
		"api_key":    apiKey,
		"auth_token": token,
		"format":     "json",
		"method":     method,
	}

	apiParams := mergeParams(baseParams, params)
	signature := signParams(apiParams)
	apiParams["api_sig"] = signature
	url := c.createURL(apiParams)

	return http.Get(url)
}

func (c *HTTP) createURL(params map[string]string) string {
	m := url.Values{}
	for k, v := range params {
		m.Add(k, v)
	}

	return c.BaseURL + "/services/rest?" + m.Encode()
}

func signParams(params map[string]string) string {
	sharedSecret := os.Getenv("RTM_SHARED_SECRET")

	// Extract keys and sort them
	keys := make([]string, 0)
	for k := range params {
		keys = append(keys, k)
	}
	sort.Sort(sort.StringSlice(keys))

	// Concatenate all keys and values
	items := make([]string, 0)
	for _, k := range keys {
		item := k + params[k]
		items = append(items, item)
	}

	text := sharedSecret + strings.Join(items, "")
	return getMD5(text)
}

func getMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func mergeParams(map1 map[string]string, map2 map[string]string) map[string]string {
	m := make(map[string]string, len(map1)+len(map2))
	for k, v := range map1 {
		m[k] = v
	}
	for k, v := range map2 {
		m[k] = v
	}
	return m
}
