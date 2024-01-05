package telco

import (
	"fmt"
)

func Example() {
	manager := NewDeviceManager()
	devices, err := manager.EnumerateDevices()
	if err != nil {
		panic(err)
	}

	fmt.Printf("[*] Telco version: %s\n", Version())
	fmt.Println("[*] Devices: ")
	for _, device := range devices {
		fmt.Printf("[*] %s => %s\n", device.Name(), device.ID())
	}
}
