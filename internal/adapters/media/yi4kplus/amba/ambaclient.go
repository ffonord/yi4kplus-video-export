package amba

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"log/slog"
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

type Conn interface {
	net.Conn
}

type ConnFactory interface {
	NewConn(host, port string) (Conn, error)
}

type Client struct {
	config      *Config
	token       int
	logger      *logger.Logger
	connFactory ConnFactory
	conn        *Conn
	reader      *bufio.Reader
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

func New(config *Config, logger *logger.Logger, connFactory ConnFactory) *Client {
	return &Client{
		config:      config,
		logger:      logger,
		connFactory: connFactory,
	}
}

func (c *Client) configureConn() error {
	const op = "AmbaClient.configureConn"

	conn, err := c.connFactory.NewConn(c.config.host, c.config.port)

	if err != nil {
		return c.errWrap(op, "net dial", err)
	}

	c.conn = &conn
	c.reader = bufio.NewReader(conn)

	return nil
}

func (c *Client) errWrap(methodName, message string, err error) error {
	return fmt.Errorf("%s: %s failed: %w", methodName, message, err)
}

func (c *Client) Run(ctx context.Context) error {
	const op = "AmbaClient.Run"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("config", c.config),
	)

	go func() {
		err := c.Shutdown(ctx)
		if err != nil {
			log.Info("Shutdown amba adapter failed: " + err.Error())
		}
	}()

	err := c.configureConn()
	if err != nil {
		return c.errWrap(op, "configure connection", err)
	}

	log.Info("Run connection")

	err = c.startSession()

	if err != nil {
		return c.errWrap(op, "session start", err)
	}

	return nil
}

func (c *Client) startSession() error {
	const op = "AmbaClient.startSession"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("config", c.config),
	)

	res, err := c.sendRequest(Request{
		MsgId: ambaStartSession,
		Token: ambaStartSessionToken,
	})

	if err != nil {
		return c.errWrap(op, "send start session request", err)
	}

	log.Info("Success amba session start")

	c.token = res.Param

	log.Info("Success fetch token")

	return nil
}

func (c *Client) stopSession() error {
	const op = "AmbaClient.stopSession"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("config", c.config),
	)

	_, err := c.sendRequest(Request{
		MsgId: ambaStopSession,
		Token: c.token,
	})

	if err != nil {
		return c.errWrap(op, "send stop session request", err)
	}

	log.Info("Success session stop")

	return nil
}

func (c *Client) sendRequest(request Request) (res Response, err error) {
	const op = "AmbaClient.sendRequest"

	rawRequest, err := json.Marshal(request)

	if err != nil {
		return res, c.errWrap(op, "json marshal", err)
	}

	_, err = fmt.Fprintf(*c.conn, string(rawRequest)+"\n")

	if err != nil {
		return res, c.errWrap(op, "fprintf to connection", err)
	}

	return c.fetchResponse()
}

func (c *Client) fetchResponse() (res Response, err error) {
	const op = "AmbaClient.fetchResponse"

	res = Response{}

	//TODO: добавить вычитывание по нужному msg_id (вычитывать, пока не получим нужную строку)
	rawRes, err := c.reader.ReadString('}')

	if err != nil {
		return res, c.errWrap(op, "reader read string", err)
	}

	err = json.Unmarshal([]byte(rawRes), &res)

	if err != nil {
		err = c.errWrap(op, "json unmarshal", err)
	}

	return res, err
}

func (c *Client) Shutdown(ctx context.Context) error {
	const op = "AmbaClient.Shutdown"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("config", c.config),
	)

	<-ctx.Done()

	if c.conn == nil {
		return nil
	}

	err := c.stopSession()

	if err != nil {
		return c.errWrap(op, "session stop", err)
	}

	err = (*c.conn).Close()

	if err != nil {
		return c.errWrap(op, "connection close", err)
	}

	log.Info("Success closing connection")

	return nil
}
