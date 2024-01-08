package telnet

import (
	"context"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"log/slog"
	"time"
)

const (
	startFtpServerCmdFormat = "tcpsvd -u %s -vE 0.0.0.0 %s ftpd -w %s 1>/dev/null 2>&1"
)

type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

type ConnFactory interface {
	NewConn(host, port string) (Conn, error)
}

type Reader interface {
	ReadString(delim byte) (string, error)
}

type ReaderFactory interface {
	NewReader(Conn) Reader
}

type Client struct {
	config            *Config
	startFtpServerCmd string
	logger            *logger.Logger
	connFactory       ConnFactory
	readerFactory     ReaderFactory
	conn              Conn
	reader            Reader
}

func New(
	config *Config,
	logger *logger.Logger,
	connFactory ConnFactory,
	readerFactory ReaderFactory,
) *Client {
	startFtpServerCmd := fmt.Sprintf(
		startFtpServerCmdFormat,
		config.ftpServerUser,
		config.ftpServerPort,
		config.ftpMediaDir,
	)

	return &Client{
		config:            config,
		startFtpServerCmd: startFtpServerCmd,
		logger:            logger,
		connFactory:       connFactory,
		readerFactory:     readerFactory,
	}
}

func (c *Client) configureConn() error {
	if c.conn != nil {
		return nil
	}

	const op = "TelnetClient.configureConn"

	conn, err := c.connFactory.NewConn(c.config.host, c.config.port)

	if err != nil {
		return c.errWrap(op, "net dial", err)
	}

	c.conn = conn
	c.reader = c.readerFactory.NewReader(conn)

	_, err = c.fetchResponse()
	time.Sleep(1 * time.Second)

	return err
}

// Run is simple configuration and start telnet session and send command of start ftp server
func (c *Client) Run(ctx context.Context) error {
	const op = "TelnetClient.Run"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("config", c.config),
	)

	go func() {
		err := c.Shutdown(ctx)
		if err != nil {
			log.Error("Shutdown telnet adapter error: " + err.Error())
		}
	}()

	err := c.configureConn()
	if err != nil {
		return c.errWrap(op, "configure connection", err)
	}

	log.Info("Run connection")

	err = c.startSession()

	if err != nil {
		return c.errWrap(op, "telnet session start", err)
	}

	err = c.startFtpServer()

	if err != nil {
		return c.errWrap(op, "start ftp server by telnet", err)
	}

	log.Info("Success start ftp server by telnet")

	return nil
}

func (c *Client) login() error {
	const op = "TelnetClient.login"

	_, err := c.sendRequest(c.config.user)

	if err != nil {
		return c.errWrap(op, "send login request", err)
	}

	return nil
}

func (c *Client) startFtpServer() error {
	const op = "TelnetClient.startFtpServer"

	_, err := c.sendRequest(c.startFtpServerCmd)

	if err != nil {
		return c.errWrap(op, "send start ftp request", err)
	}

	time.Sleep(1 * time.Second)
	return nil
}

func (c *Client) startSession() error {
	const op = "TelnetClient.startSession"

	err := c.configureConn()
	if err != nil {
		return c.errWrap(op, "configure connection", err)
	}

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("config", c.config),
	)

	err = c.login()

	if err != nil {
		return c.errWrap(op, "login", err)
	}

	log.Info("Success telnet session start")

	return nil
}

func (c *Client) sendRequest(command string) (res string, err error) {
	const op = "TelnetClient.sendRequest"

	_, err = fmt.Fprintf(c.conn, command+"\n")

	if err != nil {
		return res, c.errWrap(op, "fprintf to connection", err)
	}

	return c.fetchResponse()
}

func (c *Client) fetchResponse() (res string, err error) {
	const op = "TelnetClient.fetchResponse"

	res, err = c.reader.ReadString('\n')

	if err != nil {
		err = c.errWrap(op, "read string", err)
	}

	return res, err
}

func (c *Client) Shutdown(ctx context.Context) error {
	const op = "TelnetClient.Shutdown"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("config", c.config),
	)

	<-ctx.Done()

	if c.conn == nil {
		return nil
	}

	err := (c.conn).Close()
	if err != nil {
		return c.errWrap(op, "connection close", err)
	}

	c.conn = nil
	c.reader = nil

	log.Info("Success closing telnet connection")

	return nil
}

func (c *Client) errWrap(methodName, message string, err error) error {
	return fmt.Errorf("%s: %s failed: %w", methodName, message, err)
}
