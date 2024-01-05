package telco

//#include <telco-core.h>
import "C"
import "unsafe"

// Spawn represents spawn of the device.
type Spawn struct {
	spawn *C.TelcoSpawn
}

// PID returns process id of the spawn.
func (s *Spawn) PID() int {
	return int(C.telco_spawn_get_pid(s.spawn))
}

// Identifier returns identifier of the spawn.
func (s *Spawn) Identifier() string {
	return C.GoString(C.telco_spawn_get_identifier(s.spawn))
}

// Clean will clean the resources held by the spawn.
func (s *Spawn) Clean() {
	clean(unsafe.Pointer(s), unrefTelco)
}
