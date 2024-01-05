package telco

//#include <telco-core.h>
import "C"

import "unsafe"

// DeviceManager is the main structure which holds on devices available to Telco
// Single instance of the DeviceManager is created when you call telco.Attach() or telco.LocalDevice().
type DeviceManager struct {
	manager *C.TelcoDeviceManager
}

// NewDeviceManager returns new telco device manager.
func NewDeviceManager() *DeviceManager {
	manager := C.telco_device_manager_new()
	mgr := &DeviceManager{manager}
	return mgr
}

// Close method will close current manager.
func (d *DeviceManager) Close() error {
	var err *C.GError
	C.telco_device_manager_close_sync(d.manager, nil, &err)
	if err != nil {
		return &FError{err}
	}
	return nil
}

// EnumerateDevices will return all connected devices.
func (d *DeviceManager) EnumerateDevices() ([]*Device, error) {
	var err *C.GError
	deviceList := C.telco_device_manager_enumerate_devices_sync(d.manager, nil, &err)
	if err != nil {
		return nil, &FError{err}
	}

	numDevices := int(C.telco_device_list_size(deviceList))
	devices := make([]*Device, numDevices)

	for i := 0; i < numDevices; i++ {
		device := C.telco_device_list_get(deviceList, C.gint(i))
		devices[i] = &Device{device}
	}

	clean(unsafe.Pointer(deviceList), unrefTelco)
	return devices, nil
}

// LocalDevice returns the device with type DeviceTypeLocal.
func (d *DeviceManager) LocalDevice() (*Device, error) {
	return d.DeviceByType(DeviceTypeLocal)
}

// USBDevice returns the device with type DeviceTypeUsb.
func (d *DeviceManager) USBDevice() (*Device, error) {
	return d.DeviceByType(DeviceTypeUsb)
}

// RemoteDevice returns the device with type DeviceTypeRemote.
func (d *DeviceManager) RemoteDevice() (*Device, error) {
	return d.DeviceByType(DeviceTypeRemote)
}

// DeviceByID will return device with id passed or an error if it can't find any.
func (d *DeviceManager) DeviceByID(id string) (*Device, error) {
	idC := C.CString(id)
	defer C.free(unsafe.Pointer(idC))

	timeout := C.gint(defaultDeviceTimeout)

	var err *C.GError
	device := C.telco_device_manager_get_device_by_id_sync(d.manager, idC, timeout, nil, &err)
	if err != nil {
		return nil, &FError{err}
	}
	return &Device{device: device}, nil
}

// DeviceByType will return device or an error by device type specified.
func (d *DeviceManager) DeviceByType(devType DeviceType) (*Device, error) {
	var err *C.GError
	device := C.telco_device_manager_get_device_by_type_sync(d.manager,
		C.TelcoDeviceType(devType),
		1,
		nil,
		&err)
	if err != nil {
		return nil, &FError{err}
	}
	return &Device{device: device}, nil
}

// FindDeviceByID will try to find the device by id specified
func (d *DeviceManager) FindDeviceByID(id string) (*Device, error) {
	devID := C.CString(id)
	defer C.free(unsafe.Pointer(devID))

	timeout := C.gint(defaultDeviceTimeout)

	var err *C.GError
	device := C.telco_device_manager_find_device_by_id_sync(d.manager,
		devID,
		timeout,
		nil,
		&err)
	if err != nil {
		return nil, &FError{err}
	}

	return &Device{device: device}, nil
}

// FindDeviceByType will try to find the device by device type specified
func (d *DeviceManager) FindDeviceByType(devType DeviceType) (*Device, error) {
	timeout := C.gint(defaultDeviceTimeout)

	var err *C.GError
	device := C.telco_device_manager_find_device_by_type_sync(d.manager,
		C.TelcoDeviceType(devType),
		C.gint(timeout),
		nil,
		&err)
	if err != nil {
		return nil, &FError{err}
	}

	return &Device{device: device}, nil
}

// AddRemoteDevice add a remote device from the provided address with remoteOpts populated
func (d *DeviceManager) AddRemoteDevice(address string, remoteOpts *RemoteDeviceOptions) (*Device, error) {
	addressC := C.CString(address)
	defer C.free(unsafe.Pointer(addressC))

	var err *C.GError
	device := C.telco_device_manager_add_remote_device_sync(d.manager, addressC, remoteOpts.opts, nil, &err)
	if err != nil {
		return nil, &FError{err}
	}

	return &Device{device: device}, nil
}

// RemoveRemoteDevice removes remote device available at address
func (d *DeviceManager) RemoveRemoteDevice(address string) error {
	addressC := C.CString(address)
	defer C.free(unsafe.Pointer(addressC))

	var err *C.GError
	C.telco_device_manager_remove_remote_device_sync(d.manager,
		addressC,
		nil,
		&err)
	if err != nil {
		return &FError{err}
	}
	return nil
}

// Clean will clean the resources held by the manager.
func (d *DeviceManager) Clean() {
	clean(unsafe.Pointer(d.manager), unrefTelco)
}

// On connects manager to specific signals. Once sigName is triggered,
// fn callback will be called with parameters populated.
//
// Signals available are:
//   - "added" with callback as func(device *telco.Device) {}
//   - "removed" with callback as func(device *telco.Device) {}
//   - "changed" with callback as func() {}
func (d *DeviceManager) On(sigName string, fn any) {
	connectClosure(unsafe.Pointer(d.manager), sigName, fn)
}
