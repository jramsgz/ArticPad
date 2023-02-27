project_name = articpad
image_name = articpad:latest

.PHONY: run-local
run-local:
	go run main.go

.PHONY: run-local-frontend
run-local-frontend:
	cd ui && npm run dev

.PHONY: build-local
build-local:
	go build -o $(project_name) main.go

.PHONY: build-local-frontend
build-local-frontend:
	cd ui && npm run build

.PHONY: requirements
requirements:
	go mod tidy
	cd ui && npm install

.PHONY: clean
clean-packages:
	go clean -modcache

.PHONY: up
up: 
	make up-silent
	make shell

.PHONY: build
build:
	make build-local
	make build-local-frontend
	docker build -t $(image_name) .

.PHONY: build-no-cache
build-no-cache:
	make build-local
	make build-local-frontend
	docker build --no-cache -t $(image_name) .

.PHONY: up-silent
up-silent:
	make delete-container-if-exist
	docker run -d -p 3000:3000 --name $(project_name) $(image_name)

.PHONY: delete-container-if-exist
delete-container-if-exist:
	docker stop $(project_name) || true && docker rm $(project_name) || true

.PHONY: shell
shell:
	docker exec -it $(project_name) /bin/sh

.PHONY: stop
stop:
	docker stop $(project_name)

.PHONY: start
start:
	docker start $(project_name)
