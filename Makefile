all: build.Dockerfile run.Dockerfile
	docker build . -f build.Dockerfile -t stdlib-builder:latest -t docker.pkg.github.com/mkuznets/stdlib/builder:latest
	docker build . -f run.Dockerfile -t stdlib-runner:latest -t docker.pkg.github.com/mkuznets/stdlib/runner:latest
