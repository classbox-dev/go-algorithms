all: build.Dockerfile
	docker build . -f build.Dockerfile -t stdlib-build \
	    --build-arg USER_ID=$(shell id -u) \
	    --build-arg GROUP_ID=$(shell id -g)
