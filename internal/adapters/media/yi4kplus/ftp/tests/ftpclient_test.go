package tests

import (
	"context"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/ftp"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/ftp/mocks"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	bftp "github.com/jlaffaye/ftp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

const loggerEnv = logger.EnvTest

func TestClient_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mc := mocks.NewMockConn(ctrl)

	mc.EXPECT().
		Login(gomock.Any(), gomock.Any())

	mcf := mocks.NewMockConnFactory(ctrl)
	mcf.EXPECT().
		NewConn(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mc, nil)

	c := ftp.NewConfig()
	l := logger.New(loggerEnv)

	fc := ftp.New(c, l, mcf)

	err := fc.Run(context.Background())

	assert.Nil(t, err)
}

func TestClient_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mc := mocks.NewMockConn(ctrl)

	filepathForDelete := "test.mp4"
	mc.EXPECT().
		Delete(filepathForDelete)

	mcf := mocks.NewMockConnFactory(ctrl)
	mcf.EXPECT().
		NewConn(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mc, nil)

	c := ftp.NewConfig()
	l := logger.New(loggerEnv)

	fc := ftp.New(c, l, mcf)

	err := fc.ConfigureConn()
	assert.Nil(t, err)

	err = fc.Delete(filepathForDelete)

	assert.Nil(t, err)
}

func TestClient_GetReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mc := mocks.NewMockConn(ctrl)

	filepathForRead := "/some/dir/video1.mp4"
	mc.EXPECT().
		Retr(filepathForRead).
		Return(&bftp.Response{}, nil)

	mcf := mocks.NewMockConnFactory(ctrl)
	mcf.EXPECT().
		NewConn(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mc, nil)

	c := ftp.NewConfig()
	l := logger.New(loggerEnv)

	fc := ftp.New(c, l, mcf)

	err := fc.ConfigureConn()
	assert.Nil(t, err)

	r, err := fc.GetReader(filepathForRead)

	assert.Nil(t, err)
	assert.NotNil(t, r)
}

func TestClient_Shutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mcf := mocks.NewMockConnFactory(ctrl)

	c := ftp.NewConfig()
	l := logger.New(loggerEnv)

	fc := ftp.New(c, l, mcf)

	ctx, cancelFunc := context.WithCancel(context.Background())

	go func() {
		time.Sleep(1 * time.Second)
		cancelFunc()
	}()

	err := fc.Shutdown(ctx)

	assert.Nil(t, err)
}

func TestClient_GetFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mc := mocks.NewMockConn(ctrl)

	mediaDirs := []string{
		"/some/dir",
	}

	mediaFiles := []*bftp.Entry{
		{Target: mediaDirs[0] + "/video1.mp4"},
		{Target: mediaDirs[0] + "/video2.mp4"},
	}

	mc.EXPECT().
		NameList("").
		Return(mediaDirs, nil)

	mc.EXPECT().
		List(mediaDirs[0]).
		Return(mediaFiles, nil)

	mcf := mocks.NewMockConnFactory(ctrl)
	mcf.EXPECT().
		NewConn(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mc, nil)

	c := ftp.NewConfig()
	l := logger.New(loggerEnv)

	fc := ftp.New(c, l, mcf)

	err := fc.ConfigureConn()
	assert.Nil(t, err)

	cf, err := fc.GetFiles(context.Background())

	fileNumber := 0
	for file := range cf {
		assert.Equal(t, mediaFiles[fileNumber].Target, file.Path)
		fileNumber++
	}

	assert.Nil(t, err)
}
