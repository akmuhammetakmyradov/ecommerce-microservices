DOCKER_COMPOSE = sudo docker-compose -f ./deploy/docker-compose.yml
DOCKER_IMAGE = akmuhammet/metrics-consumer:hw9

create_network: ## 🌐 Create shared Docker network if not exists
	@sudo docker network inspect services-network >/dev/null 2>&1 || { \
		echo "🔧 Creating shared network: services-network"; \
		sudo docker network create services-network; \
	}

docker_tag_push: ## 🚀 Build & push image to Docker Hub
	sudo docker build -t $(DOCKER_IMAGE) -f ./deploy/Dockerfile .
	sudo docker push $(DOCKER_IMAGE)

dev_build: ## 🛠 Build image
	sudo docker build -t $(DOCKER_IMAGE) -f ./deploy/Dockerfile .

dev_up: create_network ## 🚀 Start containers
	$(DOCKER_COMPOSE) up --build

dev_down: ## 🧹 Stop and remove containers
	$(DOCKER_COMPOSE) down -v

dev_restart: dev_down dev_up ## 🔁 Restart dev

dev_logs: ## 📜 View logs
	$(DOCKER_COMPOSE) logs -f
