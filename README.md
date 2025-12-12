# uni-nut
This is a subset of [go-nut](https://github.com/dominikh/go-nut) meant specifically for connecting to and monitoring variables on Ubiquiti UniFi UPS NUT servers.

Ubiquiti UPS units seem to use a somewhat non-standard format for responses from the in-built NUT server. Because of this, many existing NUT clients (including the module this project was forked off of) fail to function with Ubiquiti UPS units.

## Usage
1. Enable the NUT server on your Ubiquiti UniFi-connected UPS unit.
    - The "Login Credential" setting is optional; if enabled, be sure to use the .Authenticate() method immediately after dialing!
2. This module can now be used to get variables from your UPS's NUT server as follows:
```go
package main

import (
	"fmt"

	nut "github.com/rwinkhart/uni-nut"
)

func main() {
	host := "IP:PORT"
	username := "username"
	password := "password"
	
	// connect to NUT server
	client, err := nut.Dial(host)
	if err != nil {
		panic(err)
	}
	
	// authenticate (if "Login Credential" is enabled on the NUT server)
	err = client.Authenticate(username, password)
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
