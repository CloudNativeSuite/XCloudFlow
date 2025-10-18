package modules

import "context"

// Resource represents an IAC resource for CMDB recording.
type Resource struct {
	ID   string
	Type string
	Name string
}

// CMDB defines a backend storing resource status in a graph database.
type CMDB interface {
	WriteResourceStatus(ctx context.Context, r Resource) error
	Export(ctx context.Context, dest string) error
}

var cmdb CMDB

// SetCMDB sets the global CMDB backend.
func SetCMDB(c CMDB) { cmdb = c }

// RecordResource writes resource status if backend is configured.
func RecordResource(ctx context.Context, r Resource) error {
	if cmdb != nil {
		return cmdb.WriteResourceStatus(ctx, r)
	}
	return nil
}

// ExportState triggers an export of current resource status.
func ExportState(ctx context.Context, dest string) error {
	if cmdb != nil {
		return cmdb.Export(ctx, dest)
	}
	return nil
}
