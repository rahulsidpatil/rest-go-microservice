# MIT License

# Copyright (c) 2020 rahulsidpatil

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

all: clean swagger-update build

.PHONY: build

swagger-update:
	swag init --parseDependency -d ./cmd/ -o ./api/docs

# ensure the changes are buildable
build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build  -o ./rest-go-microservice cmd/main.go

# build images (environment images are not included)
images:
	docker build -t rahulsidpatil/rest-go-microservice:latest -f ./Dockerfile .
	docker build -t rahulsidpatil/sqldb:latest ./build/db/mysql/.

docker-deploy-up:
	docker-compose -f ./build/docker-deploy/docker-deploy.yaml up --build -d
	echo "Server started ....!!"
	echo "The API documentation is available at url: http://localhost:8080/swagger/"
	echo "Server runtime profiling data available at url: http://localhost:8080/debug/pprof"

docker-deploy-down:
	docker-compose -f ./build/docker-deploy/docker-deploy.yaml down

clean:
	@rm -f ./rest-go-microservice