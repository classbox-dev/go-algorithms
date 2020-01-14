all: Dockerfile
	docker build . -t hsecode-algorithms-build \
	    --build-arg USER_ID=$(shell id -u) \
	    --build-arg GROUP_ID=$(shell id -g)
