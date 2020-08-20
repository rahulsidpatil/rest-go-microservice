# rest-go-microservice: Novice guide to build from the ground up & deploy sophisticated RestAPI-microservice application in golang
The application is a simple REST API server that will provide endpoints to allow creating, accessing and manipulating ‘messages’. The application also provides an endpoint to observe `application's runtime profiling data`

This application demonstrates how to:
1) Build golang project directory structure
2) Use go.mod for go package management
3) Swagger 2.0 integration with golang application for REST API documentation 
4) Build REST API handlers using `gorilla/mux` i.e. an http router/url matcher
5) Build Diagnostics capability by integration of `net/http/pprof` into the rest-go-microservice
6) Use database interface to ensure adaptability to any database 
7) Interact with mysql database using `go-sql-driver/mysql` i.e. MySQL driver for Go's (golang) database/sql package
8) Create Docker files for the rest-go-microservice and the related database service i.e. mysql server
9) Create Makefile with commands to swagger-update, build, deploy, stop and clean rest-go-microservice
10) Create a `docker-compose.yaml` to build and deploy rest-go-microservice into the local docker environment   

# Getting started
Following are the steps to run this application in local docker environment

## Prerequisites
- git (recommended version version 2.17.1)
- Go (recommended version version go1.13.10)
- docker (recommended version 19.03.12). 
Installation steps available at: https://docs.docker.com/engine/install/ubuntu/ https://docs.docker.com/engine/install/linux-postinstall/
Make sure to be able to run `docker as non root user` 
- docker-compose version (recommended version 1.26.2)
Installation steps available at: https://docs.docker.com/compose/install/ https://docs.docker.com/engine/security/rootless/
Make sure to be able to run `docker-compose as non root user`

## Building, Running and Accessing the application:
1) Clone the rest-go-microservice repository to your $(GOPATH) e.g. `~/go/src/github.com/rahulsidpatil`
```
git clone git@github.com:rahulsidpatil/rest-go-microservice.git
```
2) cd to rest-go-microservice dir
```
cd rest-go-microservice
```
3) To build and run rest-go-microservice; use follwoing command
```
make docker-deploy-up
```
It will a take a while to build and deploy the application on your local docker environment.
Once the deployment is successful; the console will display application access urls as:
```
echo "Server started at url: http://localhost:8080"
Server started at url: http://localhost:8080
echo "The API documentation is available at url: http://localhost:8080/swagger/"
The API documentation is available at url: http://localhost:8080/swagger/
echo "Server runtime profiling data available at url: http://localhost:8080/debug/pprof"
Server runtime profiling data available at url: http://localhost:8080/debug/pprof

```
4) To stop rest-go-microservice; use following command:
```
make docker-deploy-down
```

## rest-go-microservice in action:
1) Build and deploy rest-go-microservice:
![](./gifs/buildndeploy.gif)

2) Use rest-go-microservice:
![](./gifs/try-out-apis.gif)


# Overview of code structure:
Following is the rest-go-microservice directory structure. Each dir consists of a README.md which has detailed description of respective packages.
```
.
├── api
│   ├── docs
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   └── README.md
├── build
│   ├── db
│   │   └── mysql
│   │       ├── Dockerfile
│   │       └── scripts
│   │           ├── createDB.sql
│   │           └── initData.sql
│   ├── docker-deploy
│   │   └── docker-deploy.yaml
│   ├── k8s-deploy
│   │   └── k8sDeploy.yaml
│   └── README.md
├── cmd
│   ├── main.go
│   └── README.md
├── Dockerfile
├── gifs
│   ├── buildndeploy.gif
│   └── try-out-apis.gif
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
├── pkg
│   ├── dal
│   │   ├── dalintf.go
│   │   ├── entities.go
│   │   └── mysql.go
│   ├── handlers
│   │   └── app.go
│   ├── README.md
│   └── util
│       ├── httpErrorUtil.go
│       └── utils.go
└── README.md

```

# Licensing
rest-go-microservice is under the MIT License.
