package ports

import "context"

type Exporter interface {
	ExportFiles(ctx context.Context) error
}
