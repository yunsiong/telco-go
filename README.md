# telco-go
Go bindings for telco.

For the documentation, visit [https://pkg.go.dev/github.com/telco/telco-go/telco](https://pkg.go.dev/github.com/telco/telco-go/telco).

# Installation
* `GO111MODULE` needs to be set to `on` or `auto`.
* Download the _telco-core-devkit_ from the Telco releases [page](https://github.com/telco/telco/releases/) for you operating system and architecture.
* Extract the downloaded archive
* Copy _telco-core.h_ inside your systems include directory(inside /usr/local/include/) and _libtelco-core.a_ inside your lib directory (usually /usr/local/lib).

To use in your project, just execute: 
```bash
$ go get github.com/telco/telco-go/telco@latest
```

Supported OS:
- [x] MacOS
- [x] Linux
- [x] Android
- [ ] Windows

The reason why windows is not supported it the problem compiling telco-core with mingw because mingw(needed by cgo) can't link with MSVC .lib files.
If you manage to do it, feel free to submit your PR, also if you found any issues please submit new issue or create PR with the fix.

# Small example
```golang
package main

import (
  "bufio"
  "fmt"
  "github.com/telco/telco-go/telco"
  "os"
)

var script = `
Interceptor.attach(Module.getExportByName(null, 'open'), {
	onEnter(args) {
		const what = args[0].readUtf8String();
		console.log("[*] open(" + what + ")");
	}
});
Interceptor.attach(Module.getExportByName(null, 'close'), {
	onEnter(args) {
		console.log("close called");
	}
});
`

func main() {
  mgr := telco.NewDeviceManager()

  devices, err := mgr.EnumerateDevices()
  if err != nil {
    panic(err)
  }

  for _, d := range devices {
    fmt.Println("[*] Found device with id:", d.ID())
  }

  localDev, err := mgr.LocalDevice()
  if err != nil {
    fmt.Println("Could not get local device: ", err)
    // Let's exit here because there is no point to do anything with nonexistent device
    os.Exit(1)
  }

  fmt.Println("[*] Chosen device: ", localDev.Name())

  fmt.Println("[*] Attaching to Telegram")
  session, err := localDev.Attach("Telegram", nil)
  if err != nil {
	  fmt.Println("Error occurred attaching:", err)
	  os.Exit(1)
  }

  script, err := session.CreateScript(script)
  if err != nil {
    fmt.Println("Error occurred creating script:", err)
	os.Exit(1)
  }

  script.On("message", func(msg string) {
    fmt.Println("[*] Received", msg)
  })

  if err := script.Load(); err != nil {
    fmt.Println("Error loading script:", err)
    os.Exit(1)
  }

  r := bufio.NewReader(os.Stdin)
  r.ReadLine()
}

```

Build and run it, output will look something like this:
```bash
$ go build example.go && ./example
[*] Found device with id: local
[*] Found device with id: socket
[*] Chosen device:  Local System
[*] Attaching to Telegram
[*] Received {"type":"log","level":"info","payload":"[*] open(/Users/daemon1/Library/Application Support/Telegram Desktop/tdata/user_data/cache/0/25/0FDE3ED70BCA)"}
[*] Received {"type":"log","level":"info","payload":"[*] open(/Users/daemon1/Library/Application Support/Telegram Desktop/tdata/user_data/cache/0/8E/FD728183E115)"}
```
