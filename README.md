# sessionid_server

Token server implemented using [github.com/powerpuffpenguin/sessionid](https://github.com/powerpuffpenguin/sessionid)

Support http and grpc interfaces

* grpc -> Please see the proto definition
* http -> Browser access [/document/](https://app.swaggerhub.com/apis-docs/zuiwuchang/sessionid_server/1.0.0) to get swagger

# client

```
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/powerpuffpenguin/sessionid"
	"github.com/powerpuffpenguin/sessionid_server/client"
	"google.golang.org/grpc"
)

func main() {
	runGRPC()
}

type User struct {
	ID   int
	Name string
}

func runGRPC() {
	cc, e := grpc.Dial(
		`127.0.0.1:9000`,
		grpc.WithInsecure(),
	)
	if e != nil {
		log.Fatalln(e)
	}
	var manager sessionid.Manager = client.NewManager(cc, sessionid.JSONCoder{})
	ctx := context.Background()
	session, refresh, e := manager.Create(ctx, `1`, sessionid.Pair{
		Key: `user`,
		Value: User{
			ID:   123,
			Name: `kate`,
		},
	})
	if e != nil {
		log.Fatalln(e)
	}
	fmt.Println(`access`, session.Token())
	fmt.Println(`refresh`, refresh)
	var user User
	session.Get(ctx, `user`, &user)
	fmt.Println(user)
}
```