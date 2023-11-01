package ftp

import (
	"context"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/core/domain"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"github.com/jlaffaye/ftp"
	"io"
	"time"
)

type Client struct {
	config *Config
	logger *logger.Logger
	conn   *ftp.ServerConn
}

func New(config *Config) *Client {
	return &Client{
		config: config,
		logger: logger.New(),
	}
}

func (c *Client) configureConn() error {
	addr := fmt.Sprintf("%s:%s", c.config.host, c.config.port)
	conn, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))

	if err != nil {
		return c.errWrap("configureConn", "ftp dial", err)
	}

	c.conn = conn

	return nil
}

func (c *Client) configureLogger() error {
	return c.logger.SetLevel(c.config.logLevel)
}

func (c *Client) errWrap(methodName, message string, err error) error {
	return fmt.Errorf("\n\tftpclient::%s: %s failed: %w\n", methodName, message, err)
}

func (c *Client) Run(ctx context.Context) error {
	go func() {
		err := c.Shutdown(ctx)
		if err != nil {
			c.logger.Errorf("Shutdown ftp adapter failed error: %s", err.Error())
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
		return c.errWrap("Run", "ftp session start", err)
	}

	c.logger.Infof("Success connection with ftp server %s:%s", c.config.host, c.config.port)

	return nil
}

func (c *Client) GetFiles(ctx context.Context) (<-chan *domain.File, error) {
	mediaDirs, err := c.conn.NameList("")
	if err != nil {
		return nil, c.errWrap("GetFiles", "name list request", err)
	}

	fileChan := make(chan *domain.File)

	go func() {
		for _, dirName := range mediaDirs {

			select {
			case <-ctx.Done():
				close(fileChan)
				return
			default:

				entries, err := c.conn.List(dirName)

				if err != nil {
					c.logger.Errorf("list request failed: %s", err.Error())
					close(fileChan)
					return
				}

				for _, entry := range entries {
					select {
					case <-ctx.Done():
						close(fileChan)
						return
					default:

						fileChan <- domain.NewFile(
							entry.Name,
							entry.Target,
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
	response, err := c.conn.Retr(path)
	if err != nil {
		return nil, c.errWrap("GetReader", "send retr request: "+path, err)
	}

	return response, nil
}

func (c *Client) Delete(path string) error {
	err := c.conn.Delete(path)
	if err != nil {
		return c.errWrap("Delete", "send delete request file "+path, err)
	}

	return nil
}

func (c *Client) startSession() error {
	err := c.conn.Login(c.config.user, c.config.password)

	if err != nil {
		return c.errWrap("login", "send login request with user "+c.config.user, err)
	}

	c.logger.Infof("Success ftp session start with %s:%s", c.config.host, c.config.port)

	return nil
}

func (c *Client) Shutdown(ctx context.Context) error {
	<-ctx.Done()

	if c.conn == nil {
		return nil
	}

	if err := c.conn.Quit(); err != nil {
		return c.errWrap("Shutdown", "connection close ftp", err)
	}

	c.logger.Infof("Success closing ftp connection with %s:%s", c.config.host, c.config.port)

	return nil
}
