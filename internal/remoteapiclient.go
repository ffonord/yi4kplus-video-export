package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// Message ids
// [Options]: https://github.com/deltaflyer4747/Xiaomi_Yi/blob/master/Standalone_scripts/options.txt
// [Codes #1]: https://www.rigacci.org/wiki/doku.php/doc/appunti/hardware/sjcam-8pro-ambarella-wifi-api#web_references
// [Codes #2]: https://github.com/jnordberg/yichan/issues/1
const ambaStartSession = 257
const ambaStopSession = 258

type Client struct {
	ip     string
	port   string
	token  int
	conn   *net.Conn
	reader *bufio.Reader
}

type Response struct {
	Rval  int `json:"rval"`
	MsgId int `json:"msg_id"`
	Param int `json:"param"`
}

type Request struct {
	MsgId int    `json:"msg_id"`
	Token int    `json:"token"`
	Param string `json:"param"`
}

func NewClient(ip, port string) *Client {
	return &Client{
		ip:   ip,
		port: port,
	}
}

func (c *Client) Run() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.ip, c.port))

	if err != nil {
		return fmt.Errorf("run: init connection failed: %w", err)
	}

	c.conn = &conn
	c.reader = bufio.NewReader(conn)

	log.Printf("Run connection with %s:%s", c.ip, c.port)

	err = c.startSession()

	if err != nil {
		return fmt.Errorf("run: session start failed: %w", err)
	}

	return nil
}

func (c *Client) startSession() error {
	res, err := c.sendRequest(Request{
		MsgId: ambaStartSession,
		Token: 0,
	})

	if err != nil {
		return fmt.Errorf("startSession: send request failed: %w", err)
	}

	log.Printf("Success session start with %s:%s", c.ip, c.port)

	c.token = res.Param

	log.Printf("Success fetch token: %d", c.token)

	return nil
}

func (c *Client) stopSession() error {
	_, err := c.sendRequest(Request{
		MsgId: ambaStopSession,
		Token: c.token,
	})

	if err != nil {
		return fmt.Errorf("stopSession: sendRequest failed: %w", err)
	}

	log.Printf("Success session stop with %s:%s", c.ip, c.port)

	return nil
}

func (c *Client) sendRequest(request Request) (res Response, err error) {
	rawRequest, err := json.Marshal(request)

	if err != nil {
		return res, fmt.Errorf("sendRequest: marshal failed: %w", err)
	}

	_, err = fmt.Fprintf(*c.conn, string(rawRequest)+"\n")

	if err != nil {
		return res, fmt.Errorf("sendRequest: fprintf to %s:%s failed: %w", c.ip, c.port, err)
	}

	return c.fetchResponse()
}

func (c *Client) fetchResponse() (res Response, err error) {
	res = Response{}

	rawRes, err := c.reader.ReadString('}')

	if err != nil {
		return res, fmt.Errorf("fetchResponse: read string failed: %w", err)
	}

	err = json.Unmarshal([]byte(rawRes), &res)

	if err != nil {
		err = fmt.Errorf("fetchResponse: unmarshal failed: %w", err)
	}

	return res, err
}

func (c *Client) Stop() error {
	err := c.stopSession()

	if err != nil {
		return fmt.Errorf("stop: stop session failed: %w", err)
	}

	err = (*c.conn).Close()
	if err != nil {
		return
	}

	log.Printf("Success closing connection with %s:%s", c.ip, c.port)

	return nil
}

func (c *Client) handleError(e error, prefix string) {
	if e != nil {
		log.Fatalf("%s error: %s", prefix, e)
	}
}
