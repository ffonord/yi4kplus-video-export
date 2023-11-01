package telnet

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"net"
	"time"
)

const (
	startFtpServerCmdFormat = "tcpsvd -u %s -vE 0.0.0.0 %s ftpd -w %s 1>/dev/null 2>&1"
)

type Client struct {
	config            *Config
	startFtpServerCmd string
	logger            *logger.Logger
	conn              *net.Conn
	reader            *bufio.Reader
}

func New(config *Config) *Client {
	startFtpServerCmd := fmt.Sprintf(
		startFtpServerCmdFormat,
		config.ftpServerUser,
		config.ftpServerPort,
		config.ftpMediaDir,
	)

	return &Client{
		config:            config,
		startFtpServerCmd: startFtpServerCmd,
		logger:            logger.New(),
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
	return fmt.Errorf("\n\ttelnetclient::%s: %s failed: %w\n", methodName, message, err)
}

func (c *Client) Run(ctx context.Context) error {
	go func() {
		err := c.Shutdown(ctx)
		if err != nil {
			c.logger.Errorf("Shutdown telnet adapter error: %s", err.Error())
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
		return c.errWrap("Run", "telnet session start", err)
	}

	err = c.startFtpServer()

	if err != nil {
		return c.errWrap("Run", "start ftp server by telnet", err)
	}

	c.logger.Infof("Success start ftp server by telnet on %s:%s", c.config.host, c.config.ftpServerPort)

	return nil
}

func (c *Client) login() error {
	_, _ = c.fetchResponse()

	time.Sleep(1 * time.Second)
	_, err := c.sendRequest(c.config.user)

	if err != nil {
		return c.errWrap("login", "send login request", err)
	}

	return nil
}

func (c *Client) startFtpServer() error {
	_, err := c.sendRequest(c.startFtpServerCmd)

	if err != nil {
		return c.errWrap("startFtpServer", "send start ftp request", err)
	}

	time.Sleep(1 * time.Second)
	return nil
}

func (c *Client) startSession() error {
	err := c.login()

	if err != nil {
		return c.errWrap("startSession", "login", err)
	}

	c.logger.Infof("Success telnet session start with %s:%s", c.config.host, c.config.port)

	return nil
}

func (c *Client) sendRequest(command string) (res string, err error) {
	_, err = fmt.Fprintf(*c.conn, command+"\n")

	if err != nil {
		return res, c.errWrap("sendRequest", "fprintf to connection", err)
	}

	return c.fetchResponse()
}

func (c *Client) fetchResponse() (res string, err error) {
	res, err = c.reader.ReadString('\n')

	if err != nil {
		err = c.errWrap("fetchResponse", "read string", err)
	}

	return res, err
}

func (c *Client) Shutdown(ctx context.Context) error {
	<-ctx.Done()

	if c.conn == nil {
		return nil
	}

	err := (*c.conn).Close()
	if err != nil {
		return c.errWrap("Shutdown", "connection close", err)
	}

	c.logger.Infof("Success closing telnet connection with %s:%s", c.config.host, c.config.port)

	return nil
}