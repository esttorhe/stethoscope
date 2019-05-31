.PHONY: default build

default:	
	@echo "=============Building & pushing Docker image============="
	docker build -t stethoscope .
	docker tag stethoscope:latest $(DOCKER_STETHOSCOPE_USERNAME)/stethoscope:latest
	docker push $(DOCKER_STETHOSCOPE_USERNAME)/stethoscope:latest