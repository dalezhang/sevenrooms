package optitable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"bindolabs/optitable_middleware/config"
	"bindolabs/optitable_middleware/log"
)

var client *Client

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
func (c *Client) refreshServe() {
	var (
		wait bool
	)

	duration := 5 * time.Second
	comTimer, nilTimer := time.NewTimer(duration), &time.Timer{C: nil}
	// nil chan 会一直 block
	timer := nilTimer
	// 马上停止计时器，否则第一次 timer.C 不会等待
	comTimer.Stop()

	for {
		select {
		case <-c.exitChan:
			comTimer.Stop()
			return
		case <-c.refreshChan:
			if wait {
				continue
			}

			timer, wait = comTimer, true
			timer.Reset(duration)
		case <-timer.C:
			c.lk.Lock()
			if c.doneChan != nil {
				close(c.doneChan)
			}
			c.lk.Unlock()

			// 保证所有等待的 doRequest 都跑起来后才设置  doneChan 为 nil
			// 否则就会有 goroutine 会被 block
			c.wg.Wait()
			c.doneChan = nil

			wait = false
		}
	}
}

func (c *Client) Get(api string, params url.Values, response interface{}) error {
	return c.doRequest(http.MethodGet, api, params, nil, response)
}
func Get(param *url.Values, resp interface{}) (err error) {
	err = client.Get(config.Conf.Setting.OpUrl, *param, &resp)
	if err != nil {
		return
	}
	return
}

func (c *Client) doRequest(method, api string, params url.Values, bodyParams interface{}, response interface{}) error {
	var (
		err error
		try int
	)

	for try < config.Conf.Setting.Retry {
		try++
		fmt.Println("\n try ====", try)
		err = doRequest(method, api, params, bodyParams, response)
		if err != nil {
			return err
		}

		return nil
	}

	return err
}

func doRequest(method, api string, params url.Values, bodyParams interface{}, response interface{}) error {
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
	if method != http.MethodGet {
		request.Header.Set("Content-Type", contentType)
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

	return json.Unmarshal(bodyBytes, response)
}
