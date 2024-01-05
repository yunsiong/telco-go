package telco

//#include <telco-core.h>
import "C"
import (
	"unsafe"
)

// PeerOptions type represents struct used to setup p2p connection.
type PeerOptions struct {
	opts *C.TelcoPeerOptions
}

// NewPeerOptions creates new empty peer options.
func NewPeerOptions() *PeerOptions {
	opts := C.telco_peer_options_new()
	return &PeerOptions{opts}
}

// StunServer returns the stun server for peer options.
func (p *PeerOptions) StunServer() string {
	return C.GoString(C.telco_peer_options_get_stun_server(p.opts))
}

// ClearRelays removes previously added relays.
func (p *PeerOptions) ClearRelays() {
	C.telco_peer_options_clear_relays(p.opts)
}

// AddRelay adds new relay to use for peer options.
func (p *PeerOptions) AddRelay(relay *Relay) {
	C.telco_peer_options_add_relay(p.opts, relay.r)
}

// SetStunServer sets the stun server for peer options.
func (p *PeerOptions) SetStunServer(stunServer string) {
	stunC := C.CString(stunServer)
	defer C.free(unsafe.Pointer(stunC))
	C.telco_peer_options_set_stun_server(p.opts, stunC)
}

// Clean will clean the resources held by the peer options.
func (p *PeerOptions) Clean() {
	clean(unsafe.Pointer(p.opts), unrefTelco)
}
