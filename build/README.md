# Build Artifacts:
This directory contains the build artifacts of this project.

## Directory structure:
```
build
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

```
## db:
This directory contains mysql db Dockerfile and initial Database creation scripts. Following are mysql DDL files we use for the same. 
`createDB.sql`:
```
CREATE TABLE IF NOT EXISTS messageDB.messages (
id INT NOT NULL auto_increment PRIMARY KEY,
msg  VARCHAR(512)
)auto_increment = 1;

```
`initData.sql`
```
INSERT INTO messageDB.messages (msg) VALUES ('First message');
```

## docker-deploy:
This directory contains the docker-compose .yaml file that we use to deploy the application on local docker environment.
The `docker-deploy.yaml` specifies volume, service and environment variables for both `rest-go-microservice` and `sqldb`
`docker-deploy.yaml`
```
version: '3'
volumes:
  sqldb:
services:
  sqldb:
    build:
      context: ../db/mysql
    volumes:
      - "sqldb:/var/lib/mysql"
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: messageDB
      MYSQL_USER: docker
      MYSQL_PASSWORD: docker
    container_name: sqldb
    ports:
      - "3306:3306"
    tty: true
  rest-go-microservice:
    build:
      context: ../../
    environment:
      SVC_PORT: 8080
      SVC_VERSION: /v1
      SVC_PATH_PREFIX: messages
      STATS: 'on'
      DB_DRIVER: mysql
      DB_HOST: sqldb
      DB_PORT: 3306
      DB_USER: docker
      DB_PASSWD: docker
      DB_NAME: messageDB
    container_name: rest-go-microservice
    ports:
      - "8080:8080"
      - "6060:6060"
    tty: true
    depends_on:
      - sqldb
```

## k8s-deploy:
This directory contains artifacts to deploy `rest-go-microservice` on k8s clusters. All the required artifacts are yet to be added here. This is still a work in progress...
