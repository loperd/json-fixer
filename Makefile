PROJECT_NAME=json-fixer
IMAGE_TAG=$(PROJECT_NAME):latest

remove: # Remove old images with tag json-fixer
	@docker images | grep $(PROJECT_NAME) | awk '{ print $$3 }' | xargs docker rmi

build: remove # Build new image
	@docker build -t $(IMAGE_TAG) .

drop:
	@docker ps -a | grep $(PROJECT_NAME) | xargs docker rm -f

up: drop
	@docker run -it --rm -d --name $(PROJECT_NAME) $(IMAGE_TAG)