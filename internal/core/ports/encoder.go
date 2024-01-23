package ports

import "context"

type Encoder interface {
	Encode(ctx context.Context, dstDirName, srcDirName string) error
}
