package sevenroom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sync"

	"bindolabs/sevenrooms/config"
	"bindolabs/sevenrooms/log"
)

var (
	client          *Client
	ivalidToken     = "Token is not valid"
	InvalidTokenErr = fmt.Errorf(ivalidToken)
)

type Client struct {
	lk sync.Mutex
	wg sync.WaitGroup

	doneChan    chan struct{}
	refreshChan chan struct{}
	exitChan    chan struct{}
}

func Init() error {
	var err error
	client, err = NewClient()
	if err != nil {
		return err
	}

	return nil
}

func NewClient() (*Client, error) {
	c := &Client{
		refreshChan: make(chan struct{}),
		exitChan:    make(chan struct{}),
	}
	go c.refreshServe()
	return c, nil
}

func (c *Client) Get(api string, params url.Values, response interface{}, token string) error {
	return c.doRequest(http.MethodGet, api, params, nil, response, token)
}
func (c *Client) Post(api string, body interface{}, response interface{}, token string) error {
	return c.doRequest(http.MethodPost, api, nil, body, response, token)
}

func Get(param *url.Values, resp interface{}) (err error) {
	err = client.Get(config.Conf.Setting.OpUrl, *param, &resp, "")
	if err != nil {
		return
	}
	return
}
func PostWebhooks(store *config.Store, params *map[string]interface{}, resp interface{}) (err error) {
	var try int

	for try < config.Conf.Setting.Retry {
		try++
		fmt.Println("\n try ====", try)
		url := fmt.Sprintf("%svenues/%s/webhooks/%s/basket/updates", config.Conf.Setting.OpUrl, store.VenueID, config.Conf.Setting.PosID)
		err = client.Post(url, params, &resp, store.Token)
		r := bytes.NewReader([]byte(fmt.Sprintln(err)))
		mached, matcherr := regexp.MatchReader(".*Permission denied.*", r)
		if matcherr != nil {
			fmt.Printf("\n matcherr =============== %+v \n", matcherr)
		}
		if mached {
			fmt.Println("\n mached ===============")
			log.Logger.Warnf("invaild token try[%d]again", try)
			GetStoreToken(store)
			fmt.Println("\n store.Token ===============", store.Token)
			continue
		}
	}
	return
}

func (c *Client) doRequest(method, api string, params url.Values, bodyParams interface{}, response interface{}, token string) error {
	var (
		err error
	)
	err = doRequest(method, api, token, params, bodyParams, response)
	if err != nil {
		return err
	}
	return nil

	return err
}

func doRequest(method, api string, token string, params url.Values, bodyParams interface{}, response interface{}) error {
	var (
		body io.Reader
	)
	if bodyParams != nil {
		if data, err := json.Marshal(bodyParams); err != nil {
			return err
		} else {
			body = bytes.NewBuffer(data)
		}
	}

	contentType := "application/json"
	if method == http.MethodPost {
		if body == nil {
			contentType = "application/x-www-form-urlencoded"
		}
		if len(params) > 0 {
			body = bytes.NewBufferString(params.Encode())
		}
	}

	request, err := http.NewRequest(method, api, body)
	if err != nil {
		return err
	}

	if method == http.MethodGet || method == http.MethodDelete {
		if len(params) > 0 {
			request.URL.RawQuery = params.Encode()
		}
	}
	if method != http.MethodGet && token != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if token != "" {
		request.Header.Set("Authorization", token)
	}

	if config.Conf.Debug {
		log.Logger.Debugf("req: %+v", request)
	}

	resp, err := http.DefaultClient.Do(request)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if config.Conf.Debug {
		log.Logger.Debugf("resp: %+v", string(bodyBytes))
	}
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("req[%+v], resp[code:%d, %s]", request, resp.StatusCode, string(bodyBytes))
	}

	if response == nil {
		return nil
	}
	if string(bodyBytes) == "" {
		fmt.Printf("resp.Body=========%+v", string(bodyBytes))
		return nil
	}

	return json.Unmarshal(bodyBytes, response)
}
