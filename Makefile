all: build.Dockerfile run.Dockerfile
	docker build . -f build.Dockerfile -t stdlib-builder:latest
	docker build . -f run.Dockerfile -t stdlib-runner:latest

clean:
	find . -type f \( -name "*.test" -o -name "*.prof" \) -delete
