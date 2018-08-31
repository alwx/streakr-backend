IMAGE_NAME = alwxx/streakr-go:latest

build-production:
	docker build --rm -t $(IMAGE_NAME) -f ./deployments/Dockerfile .
	docker push alwxx/streakr-go:latest

run-development-deps:
	docker-compose up