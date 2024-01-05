package telco

//#include <telco-core.h>
import "C"
import "unsafe"

type SnapshotOptions struct {
	opts *C.TelcoSnapshotOptions
}

// NewSnapshotOptions creates new snapshot options with warmup
// script and script runtime provided.
func NewSnapshotOptions(warmupScript string, rt ScriptRuntime) *SnapshotOptions {
	opts := C.telco_snapshot_options_new()
	warmupScriptC := C.CString(warmupScript)
	defer C.free(unsafe.Pointer(warmupScriptC))

	C.telco_snapshot_options_set_warmup_script(opts, warmupScriptC)
	C.telco_snapshot_options_set_runtime(opts, C.TelcoScriptRuntime(rt))

	return &SnapshotOptions{
		opts: opts,
	}
}

// WarmupScript returns the warmup script used to create the script options.
func (s *SnapshotOptions) WarmupScript() string {
	return C.GoString(C.telco_snapshot_options_get_warmup_script(s.opts))
}

// Runtime returns the runtime used to create the script options.
func (s *SnapshotOptions) Runtime() ScriptRuntime {
	return ScriptRuntime(int(C.telco_snapshot_options_get_runtime(s.opts)))
}

// Clean will clean the resources held by the snapshot options.
func (s *SnapshotOptions) Clean() {
	clean(unsafe.Pointer(s.opts), unrefTelco)
}
