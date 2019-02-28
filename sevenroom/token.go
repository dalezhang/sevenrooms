package sevenroom

import (
	"bindolabs/sevenrooms/config"
	"bindolabs/sevenrooms/log"
	"net/http"
	"net/url"
	"time"
)

// {
//     "status": 200,
//     "msg": "Successfully authenticated",
//     "data": {
//         "token": "14227f15adb22bb06407b00bbd738fe16b19205a291ea0fecfaf0b70638a81f73a988495f56ced8570b491e03a81a53134f87e15f805fbb334de90cacf01b27c"
//     }
// }
type authResp struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   token  `json:"data"`
}
type token struct {
	Token string `json:"token"`
}

func getToken() (string, error) {
	var (
		resp   authResp
		params = url.Values{}
	)

	err := doRequest(http.MethodPost, config.Conf.Setting.OpUrl+"auth?client_id="+config.Conf.Setting.ClientID+"&client_secret="+config.Conf.Setting.ClientSecret, "", params, nil, &resp)

	if err != nil {
		return "", err
	}

	if err = config.Conf.SetToken(resp.Data.Token); err != nil {
		log.Logger.Warnf("set token err: %v", err)
	}
	return resp.Data.Token, nil
}

func (c *Client) wait() {
	c.lk.Lock()
	if c.doneChan == nil {
		c.doneChan = make(chan struct{})
	}
	c.lk.Unlock()

	c.wg.Add(1)
	<-c.doneChan
	c.wg.Done()
}

func (c *Client) done() {
	c.lk.Lock()
	if c.doneChan != nil {
		close(c.doneChan)
	}
	c.lk.Unlock()

	// 保证所有等待的 doRequest 都跑起来后才设置  doneChan 为 nil
	// 否则就会有 goroutine 会被 block
	c.wg.Wait()
	c.doneChan = nil
}

func (c *Client) refresh() {
	c.refreshChan <- struct{}{}
}

func (c *Client) refreshServe() {
	var (
		err   error
		token string
		wait  bool
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
			// 多个 http 请求出问题的时候，单位时间内只更新一次 token
			if wait {
				continue
			}

			timer, wait = comTimer, true
			timer.Reset(duration)
		case <-timer.C:
			timer = nilTimer
			token, err = getToken()
			if err != nil {
				log.Logger.Warnf("get token err: %v", err)
			} else {
				c.Token = token
				c.done()
			}

			wait = false
		}
	}
}
