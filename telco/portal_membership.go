package telco

//#include <telco-core.h>
import "C"
import "unsafe"

// PortalMembership type is used to join portal with session.
type PortalMembership struct {
	mem *C.TelcoPortalMembership
}

// ID returns the ID of the membership
func (p *PortalMembership) ID() uint {
	return uint(C.telco_portal_membership_get_id(p.mem))
}

// Terminate terminates the session membership
func (p *PortalMembership) Terminate() error {
	var err *C.GError
	C.telco_portal_membership_terminate_sync(p.mem, nil, &err)
	if err != nil {
		return &FError{err}
	}
	return nil
}

// Clean will clean the resources held by the portal membership.
func (p *PortalMembership) Clean() {
	clean(unsafe.Pointer(p.mem), unrefTelco)
}
