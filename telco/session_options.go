package telco

//#include <telco-core.h>
import "C"
import "unsafe"

// SessionOptions type is used to configure session
type SessionOptions struct {
	opts *C.TelcoSessionOptions
}

// NewSessionOptions create new SessionOptions with the realm and
// timeout to persist provided
func NewSessionOptions(realm Realm, persistTimeout uint) *SessionOptions {
	opts := C.telco_session_options_new()
	C.telco_session_options_set_realm(opts, C.TelcoRealm(realm))
	C.telco_session_options_set_persist_timeout(opts, C.guint(persistTimeout))

	return &SessionOptions{opts}
}

// Realm returns the realm of the options
func (s *SessionOptions) Realm() Realm {
	rlm := C.telco_session_options_get_realm(s.opts)
	return Realm(rlm)
}

// PersistTimeout returns the persist timeout of the script.s
func (s *SessionOptions) PersistTimeout() int {
	return int(C.telco_session_options_get_persist_timeout(s.opts))
}

// Clean will clean the resources held by the session options.
func (s *SessionOptions) Clean() {
	clean(unsafe.Pointer(s.opts), unrefTelco)
}
