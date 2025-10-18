# uni-nut
This is a subset of [go-nut](https://github.com/dominikh/go-nut) meant specifically for connecting to and monitoring variables on Ubiquiti UniFi UPS NUT servers.

Ubiquiti UPS units seem to use a somewhat non-standard format for responses from the in-built NUT server. Because of this, many existing NUT clients (including the module this project was forked off of) fail to function with Ubiquiti UPS units.

## Usage
1. Enable the NUT server on your Ubiquiti UniFi-connected UPS unit.
    - Do **NOT** enable the "Login Credential" option!
2. This module can now be used to get variables from your UPS's NUT server as follows:
```go
package main

import (
	"fmt"

	nut "github.com/rwinkhart/uni-nut"
)

func main() {
	host := "IP:PORT"
	upsID := "UPS-ID"
	client, err := nut.Dial(host)
	if err != nil {
		panic(err)
	}
	err = client.GetListVar(upsID)
	if err != nil {
		panic(err)
	}
	for key, value := range nut.NutKeyValMap {
		fmt.Println(key + ": " + value)
	}
}
```
