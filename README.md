RFC2822 Parser For Go
---------------------
This a simple [rfc2822](http://www.ietf.org/rfc/rfc2822.txt) parser for 
[Go](http://golang.org).  It does not (yet) support the full standard.


### Installation
`go get github.com/trapped/rfc2822`

### Usage
```go
package main

import (
    "github.com/trapped/rfc2822"
    "os"
    "fmt"
)

func main() {
    var (
        msg   *rfc2822.Message
        err   os.Error
        value string
    )

    if msg, err = rfc2822.ReadFile("message.txt"); err != nil {
        fmt.Printf("error reading file: %s\n", err)
        os.Exit(2)
    }

    if value, err = msg.GetHeader("subject"); err != nil {
        fmt.Println(err)
        os.Exit(3)
    }

    /*...*/rfc2822.ReadString("an email")/*...*/

    fmt.Printf("Subject: %s\n", value)
}
```
