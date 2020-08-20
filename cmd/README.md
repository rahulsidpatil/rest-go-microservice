# rest-go-microservice Entry point

## Directory structure:
```
cmd
│   ├── main.go
│   └── README.md

```

## Overview

The main.go file is an entry point for rest-go-microservice application. It invokes application handler.
```
package main

import (
	"github.com/rahulsidpatil/rest-go-microservice/pkg/handlers"
)

func main() {
	a := handlers.App{}
	a.Initialize()
	a.Run()
}

```
