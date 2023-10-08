package ftp

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"os"
	"time"
)

type Client struct {
	config *Config
	logger *logger.Logger
	conn   *ftp.ServerConn
	reader *bufio.Reader
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

	//TODO: тут зациклить с прослушиванием контекста и слипом
	return c.fetchFiles(ctx)
}

func (c *Client) fetchFiles(ctx context.Context) error {
	mediaDirs, err := c.conn.NameList("")
	if err != nil {
		return c.errWrap("fetchFiles", "get media dir list", err)
	}

	storageDir, err := c.storageDir()
	if err != nil {
		return c.errWrap("fetchFiles", "get storage dir", err)
	}

	for _, dirName := range mediaDirs {
		fileNames, err := c.conn.NameList(dirName)
		if err != nil {
			return c.errWrap("fetchFiles", "get media files", err)
		}

		for _, fileName := range fileNames {
			//TODO: прикрутить прослушивание контекста, чтобы остановиться
			err := c.popFile(fileName, dirName, storageDir)
			if err != nil {
				return c.errWrap("fetchFiles", "pop file "+fileName, err)
			}
		}
	}

	return nil
}

func (c *Client) popFile(fileName, dirName, storageDir string) error {
	//TODO: не обрабатывать *.SEC файлы, их сразу удаляем
	storageFileName := storageDir + string(os.PathSeparator) + fileName
	ftpFileName := dirName + string(os.PathSeparator) + fileName

	if _, err := os.Stat(storageFileName); os.IsNotExist(err) {
		response, err := c.conn.Retr(ftpFileName)
		if err != nil {
			return c.errWrap("popFile", "file retr "+ftpFileName, err)
		}

		outFile, err := os.Create(storageFileName)
		if err != nil {
			return c.errWrap("popFile", "file create "+storageFileName, err)
		}

		c.logger.Infof("Start download %s", ftpFileName)

		_, err = io.Copy(outFile, response)
		if err != nil {
			return c.errWrap("popFile", "io copy "+ftpFileName, err)
		}

		if err := response.Close(); err != nil {
			return c.errWrap("popFile", "ftp response close "+ftpFileName, err)
		}

		if err := outFile.Close(); err != nil {
			return c.errWrap("popFile", "storage file close "+storageFileName, err)
		}
	}

	ftpFileSize, err := c.conn.FileSize(ftpFileName)
	if err != nil {
		return c.errWrap("popFile", "ftp file size "+ftpFileName, err)
	}

	fileInfo, err := os.Stat(storageFileName)
	if err != nil {
		return c.errWrap("popFile", "storage file info "+storageFileName, err)
	}

	if ftpFileSize == fileInfo.Size() {
		err := c.conn.Delete(ftpFileName)
		if err != nil {
			return c.errWrap("popFile", "ftp delete file "+ftpFileName, err)
		}

		c.logger.Infof("Success downloaded %s", ftpFileName)
	} else {
		err := os.Remove(storageFileName)
		if err != nil {
			return c.errWrap("popFile", "storage delete file "+storageFileName, err)
		}
	}

	return nil
}

func (c *Client) storageDir() (string, error) {
	currentDate := time.Now().Format("01-02-2006")
	storageDir := c.config.storageDir + string(os.PathSeparator) + currentDate

	if _, err := os.Stat(storageDir); os.IsNotExist(err) {

		err = os.Mkdir(storageDir, 0755)

		if err != nil {
			return "", c.errWrap("storageDir", "mkdir "+storageDir, err)
		}
	}

	return storageDir, nil
}

func (c *Client) login() error {
	err := c.conn.Login(c.config.user, c.config.password)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		return c.errWrap("login", "send login request", err)
	}

	return nil
}

func (c *Client) startSession() error {
	err := c.login()

	if err != nil {
		return c.errWrap("startSession", "login", err)
	}

	c.logger.Infof("Success ftp session start with %s:%s", c.config.host, c.config.port)

	return nil
}

func (c *Client) Shutdown(ctx context.Context) error {
	if c.conn == nil {
		return nil
	}

	if err := c.conn.Quit(); err != nil {
		return c.errWrap("Shutdown", "connection close ftp", err)
	}

	c.logger.Infof("Success closing ftp connection with %s:%s", c.config.host, c.config.port)

	return nil
}
