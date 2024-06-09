SHELL=/bin/bash -o pipefail
Version=v1
BACKEND_APP_NAME=web-app-sample
BACKEND_SERVERJOB_NAME=web-app-sample-job

DOCKER_BUILD_IMAGE_TAG=web-app-sample-build

DOCKER_RUN_HOTFIX_MANAGER_IMAGE_TAG=web-app-sample-run


GitCount= $(shell git rev-list --count HEAD)
GitSHA1= $(shell git log -n1 --format=format:"%H")
branch=$(shell git symbolic-ref --short HEAD)
Branch=$(subst /,_,$(branch))

.PHONY: buildimage test clean vet deployimage dockertest

buildimage:
	docker build --build-arg appname=$(BACKEND_APP_NAME) -t $(DOCKER_BUILD_IMAGE_TAG) \
			-f ./Dockerfile .

	docker run --rm $(DOCKER_BUILD_IMAGE_TAG) \
		| docker build -t $(DOCKER_RUN_HOTFIX_MANAGER_IMAGE_TAG):${Version}_${Branch}_${GitCount} -f Dockerfile.run -
	docker tag $(DOCKER_RUN_HOTFIX_MANAGER_IMAGE_TAG):${Version}_${Branch}_${GitCount} $(DOCKER_RUN_HOTFIX_MANAGER_IMAGE_TAG):latest

test:
	go env -w GOPROXY=https://goproxy.cn,direct
	go env -w GOPRIVATE=gitlab.appshahe.com
	go test  -cover -race -count=1 -v ./... | sed '/PASS/s//$(shell printf "\033[32mPASS\033[0m")/' | sed '/FAIL/s//$(shell printf "\033[31mFAIL\033[0m")/' | sed '/coverage/s//$(shell printf "\033[32mcoverage\033[0m")/' | sed '/undefined/s//$(shell printf "\033[31mUNDEFINED\033[0m")/'

unit_test:
	docker build --build-arg appname=$(BACKEND_APP_NAME) -t $(DOCKER_BUILD_IMAGE_TAG) \
		-f ./Dockerfile.test .
	docker-compose up -d --build
	docker logs -f ${DOCKER_BUILD_IMAGE_TAG} > unit_test.log
	docker-compose down

vet :
	go vet ./...

deployimage:
	make buildimage
	docker-compose down
	docker-compose up -d

clean:

