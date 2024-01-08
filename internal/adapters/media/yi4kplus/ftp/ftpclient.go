package ftp

import (
	"context"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/core/domain/file"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"github.com/jlaffaye/ftp"
	"io"
	"log/slog"
	"time"
)

type Conn interface {
	NameList(path string) (entries []string, err error)
	List(path string) (entries []*ftp.Entry, err error)
	Retr(path string) (*ftp.Response, error)
	Delete(path string) error
	Login(user, password string) error
	Quit() error
}

type ConnFactory interface {
	NewConn(host, port string, timeout time.Duration) (Conn, error)
}

type Client struct {
	config      *Config
	logger      *logger.Logger
	connFactory ConnFactory
	conn        Conn
}

func New(config *Config, logger *logger.Logger, connFactory ConnFactory) *Client {
	return &Client{
		config:      config,
		logger:      logger,
		connFactory: connFactory,
	}
}

func (c *Client) ConfigureConn() error {
	if c.conn != nil {
		return nil
	}

	const op = "FtpClient.ConfigureConn"

	conn, err := c.connFactory.NewConn(c.config.host, c.config.port, 5*time.Second)

	if err != nil {
		return c.errWrap(op, "ftp dial", err)
	}

	c.conn = conn

	return nil
}

func (c *Client) Run(ctx context.Context) error {
	const op = "FtpClient.Run"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("config", c.config),
	)

	go func() {
		err := c.Shutdown(ctx)
		if err != nil {
			log.Info("Shutdown ftp adapter failed: " + err.Error())
		}
	}()

	err := c.ConfigureConn()
	if err != nil {
		return c.errWrap(op, "configure connection", err)
	}

	log.Info("Run connection", c.config.host, c.config.port)

	err = c.startSession()

	if err != nil {
		return c.errWrap(op, "ftp session start", err)
	}

	log.Info("Success connection")

	return nil
}

func (c *Client) GetFiles(ctx context.Context) (<-chan *file.File, error) {
	const op = "FtpClient.GetFiles"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("config", c.config),
	)

	mediaDirs, err := c.conn.NameList("")
	if err != nil {
		return nil, c.errWrap(op, "name list request", err)
	}

	fileChan := make(chan *file.File)

	go func() {
		for _, dirName := range mediaDirs {

			select {
			case <-ctx.Done():
				close(fileChan)
				return
			default:

				entries, err := c.conn.List(dirName)

				if err != nil {
					log.Error("list request failed: " + err.Error())
					close(fileChan)
					return
				}

				for _, entry := range entries {
					select {
					case <-ctx.Done():
						close(fileChan)
						return
					default:

						fileChan <- file.New(
							entry.Name,
							dirName,
							entry.Time,
							entry.Size,
						)
					}
				}
			}
		}

		close(fileChan)
	}()

	return fileChan, nil
}

func (c *Client) GetReader(path string) (io.ReadCloser, error) {
	const op = "FtpClient.GetReader"

	response, err := c.conn.Retr(path)
	if err != nil {
		return nil, c.errWrap(op, "send retr request: "+path, err)
	}

	return response, nil
}

func (c *Client) Delete(path string) error {
	const op = "FtpClient.Delete"

	err := c.conn.Delete(path)
	if err != nil {
		return c.errWrap(op, "send delete request file "+path, err)
	}

	return nil
}

func (c *Client) startSession() error {
	const op = "FtpClient.startSession"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("config", c.config),
	)

	err := c.conn.Login(c.config.user, c.config.password)

	if err != nil {
		return c.errWrap(op, "send login request with user "+c.config.user, err)
	}

	log.Info("Success ftp session start")

	return nil
}

func (c *Client) Shutdown(ctx context.Context) error {
	const op = "FtpClient.Shutdown"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("config", c.config),
	)

	<-ctx.Done()

	if c.conn == nil {
		return nil
	}

	if err := c.conn.Quit(); err != nil {
		return c.errWrap(op, "connection close ftp", err)
	}

	c.conn = nil

	log.Info("Success closing ftp connection")

	return nil
}

func (c *Client) errWrap(methodName, message string, err error) error {
	return fmt.Errorf("%s: %s failed: %w", methodName, message, err)
}
