package tests

import (
	"context"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/telnet"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/telnet/mocks"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
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
		Write(gomock.Any()).
		AnyTimes()

	mc.EXPECT().
		Read(gomock.Any()).
		AnyTimes()

	mcf := mocks.NewMockConnFactory(ctrl)
	mcf.EXPECT().
		NewConn(gomock.Any(), gomock.Any()).
		Return(mc, nil)

	mr := mocks.NewMockReader(ctrl)

	mr.EXPECT().
		ReadString(gomock.Any()).
		AnyTimes()

	mrf := mocks.NewMockReaderFactory(ctrl)

	mrf.EXPECT().
		NewReader(mc).
		Return(mr)

	c := telnet.NewConfig()
	l := logger.New(loggerEnv)

	tc := telnet.New(c, l, mcf, mrf)

	err := tc.Run(context.Background())

	assert.Nil(t, err)
}

func TestClient_Shutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mcf := mocks.NewMockConnFactory(ctrl)
	mrf := mocks.NewMockReaderFactory(ctrl)

	c := telnet.NewConfig()
	l := logger.New(loggerEnv)

	tc := telnet.New(c, l, mcf, mrf)

	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		cancelFunc()
	}()

	err := tc.Shutdown(ctx)

	assert.Nil(t, err)
}
