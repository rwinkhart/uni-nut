# uni-nut
This is a subset of [go-nut](https://github.com/dominikh/go-nut) meant specifically for connecting to and monitoring variables on Ubiquiti UniFi UPS NUT servers. It may work for NUT servers on other devices, but it was made to make it as easy as possible to monitor UniFi UPS units.

## Usage
1. Enable the NUT server on your Ubiquiti UniFi-connected UPS unit.
    - The "Login Credential" setting has no impact on this module; we only need RO access.
2. This module can now be used to get variables from your UPS's NUT server as follows:
```go
package main

import (
	"fmt"

	nut "github.com/rwinkhart/uni-nut"
)

func main() {
	host := "IP:PORT"
	
	// connect to NUT server
	client, err := nut.Dial(host)
	if err != nil {
		panic(err)
	}
	
	// auto-detect UPS ID
	err = client.AutomaticallySetID()
	if err != nil {
		panic(err)
	}
	
	// list all variables from your UPS
	err = client.ListVar()
	if err != nil {
		panic(err)
	}
	for key, value := range nut.NutKeyValMap {
		fmt.Println(key + ": " + value)
	}
}
```

This module also features a .GetVar() method for retrieving the value of only one variable.
