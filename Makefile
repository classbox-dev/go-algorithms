all: build.Dockerfile run.Dockerfile
	docker build . -f build.Dockerfile -t stdlib-build
	docker build . -f run.Dockerfile -t stdlib-run
