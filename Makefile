project_name = articpad
image_name = articpad:latest

run-local:
	go run main.go

run-local-frontend:
	cd ui && npm run dev

build-local:
	go build -o $(project_name) main.go

build-local-frontend:
	cd ui && npm run build

requirements:
	go mod tidy
	cd ui && npm install

clean-packages:
	go clean -modcache

up: 
	make up-silent
	make shell

build:
	docker build -t $(image_name) .

build-no-cache:
	docker build --no-cache -t $(image_name) .

up-silent:
	make delete-container-if-exist
	docker run -d -p 3000:3000 --name $(project_name) $(image_name)

delete-container-if-exist:
	docker stop $(project_name) || true && docker rm $(project_name) || true

shell:
	docker exec -it $(project_name) /bin/sh

stop:
	docker stop $(project_name)

start:
	docker start $(project_name)