package amba

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"net"
)

// Message ids
// [Options]: https://github.com/deltaflyer4747/Xiaomi_Yi/blob/master/Standalone_scripts/options.txt
// [Codes #1]: https://www.rigacci.org/wiki/doku.php/doc/appunti/hardware/sjcam-8pro-ambarella-wifi-api#web_references
// [Codes #2]: https://github.com/jnordberg/yichan/issues/1
// [Codes #3]: https://cgg.mff.cuni.cz/gitlab/i3d/mi-camera_photo/-/blob/master/commands/AMBACommands.txt
const ambaStartSessionToken = 0
const ambaStartSession = 257
const ambaStopSession = 258

type Client struct {
	config *Config
	token  int
	logger *logger.Logger
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

func New(config *Config) *Client {
	return &Client{
		config: config,
		logger: logger.New(),
	}
}

func (c *Client) configureConn() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.config.host, c.config.port))

	if err != nil {
		return c.errWrap("configureConn", "net dial", err)
	}

	c.conn = &conn
	c.reader = bufio.NewReader(conn)

	return nil
}

func (c *Client) configureLogger() error {
	return c.logger.SetLevel(c.config.logLevel)
}

func (c *Client) errWrap(methodName, message string, err error) error {
	return fmt.Errorf("\n\tambaclient::%s: %s failed: %w\n", methodName, message, err)
}

func (c *Client) Run(ctx context.Context) error {
	go func() {
		err := c.Shutdown(ctx)
		if err != nil {
			c.logger.Errorf("Shutdown amba adapter failed error: %s", err.Error())
		}
	}()

	err := c.configureLogger()
	if err != nil {
		return c.errWrap("Run", "configure logger", err)
	}

	err = c.configureConn()
	if err != nil {
		return c.errWrap("Run", "configure connection", err)
	}

	c.logger.Infof("Run connection with %s:%s", c.config.host, c.config.port)

	err = c.startSession()

	if err != nil {
		return c.errWrap("Run", "session start", err)
	}

	return nil
}

func (c *Client) startSession() error {
	res, err := c.sendRequest(Request{
		MsgId: ambaStartSession,
		Token: ambaStartSessionToken,
	})

	if err != nil {
		return c.errWrap("startSession", "send start session request", err)
	}

	c.logger.Infof("Success amba session start with %s:%s", c.config.host, c.config.port)

	c.token = res.Param

	c.logger.Infof("Success fetch token: %d", c.token)

	return nil
}

func (c *Client) stopSession() error {
	_, err := c.sendRequest(Request{
		MsgId: ambaStopSession,
		Token: c.token,
	})

	if err != nil {
		return c.errWrap("stopSession", "send stop session request", err)
	}

	c.logger.Infof("Success session stop with %s:%s", c.config.host, c.config.port)

	return nil
}

func (c *Client) sendRequest(request Request) (res Response, err error) {
	rawRequest, err := json.Marshal(request)

	if err != nil {
		return res, c.errWrap("sendRequest", "json marshal", err)
	}

	_, err = fmt.Fprintf(*c.conn, string(rawRequest)+"\n")

	if err != nil {
		return res, c.errWrap("sendRequest", "fprintf to connection", err)
	}

	return c.fetchResponse()
}

func (c *Client) fetchResponse() (res Response, err error) {
	res = Response{}

	//TODO: добавить вычитывание по нужному msg_id (вычитывать, пока не получим нужную строку)
	rawRes, err := c.reader.ReadString('}')

	if err != nil {
		return res, c.errWrap("fetchResponse", "reader read string", err)
	}

	err = json.Unmarshal([]byte(rawRes), &res)

	if err != nil {
		err = c.errWrap("fetchResponse", "json unmarshal", err)
	}

	return res, err
}

func (c *Client) Shutdown(ctx context.Context) error {
	<-ctx.Done()

	if c.conn == nil {
		return nil
	}

	err := c.stopSession()

	if err != nil {
		return c.errWrap("Shutdown", "session stop", err)
	}

	err = (*c.conn).Close()

	if err != nil {
		return c.errWrap("Shutdown", "connection close", err)
	}

	c.logger.Infof("Success closing connection with %s:%s", c.config.host, c.config.port)

	return nil
}
