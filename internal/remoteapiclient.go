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
	conn   *net.Conn
	token  int
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
	return &Client{ip, port, nil, 0, nil}
}

func (c *Client) Run() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.ip, c.port))

	if err != nil {
		log.Fatal(err)
	}

	c.conn = &conn
	c.reader = bufio.NewReader(conn)

	log.Printf("Run connection with %s:%s", c.ip, c.port)

	err = c.startSession()
	c.handleError(err, "Session start")
}

func (c *Client) startSession() error {
	res, err := c.sendRequest(Request{
		MsgId: ambaStartSession,
		Token: 0,
		Param: "",
	})

	if err != nil {
		return err
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
		Param: "",
	})

	if err != nil {
		return err
	}

	log.Printf("Success session stop with %s:%s", c.ip, c.port)

	return nil
}

func (c *Client) sendRequest(request Request) (res Response, err error) {
	rawRequest, err := json.Marshal(request)

	if err != nil {
		return
	}

	_, err = fmt.Fprintf(*c.conn, string(rawRequest)+"\n")

	if err != nil {
		return
	}

	return c.fetchResponse()
}

func (c *Client) fetchResponse() (res Response, err error) {
	res = Response{}

	rawRes, err := c.reader.ReadString('}')

	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(rawRes), &res)

	return res, err
}

func (c *Client) Stop() {
	err := c.stopSession()
	c.handleError(err, "Session stop")

	err = (*c.conn).Close()
	c.handleError(err, "Close connection")

	log.Printf("Success closing connection with %s:%s", c.ip, c.port)
}

func (c *Client) handleError(e error, prefix string) {
	if e != nil {
		log.Fatalf("%s error: %s", prefix, e)
	}
}
