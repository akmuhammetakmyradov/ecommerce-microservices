DOCKER_COMPOSE_DEV = sudo docker-compose -f ./deploy/dev/docker-compose.yml
DOCKER_COMPOSE_PROD = sudo docker-compose -f ./deploy/prod/docker-compose.yml
DOCKER_IMAGE = akmuhammet/stocks:hw11

.PHONY: dev_build dev_up dev_down dev_restart dev_logs \
        prod_pull prod_up prod_down prod_restart prod_logs \
        create_network

create_network: ## 🌐 Create shared Docker network if not exists
	@sudo docker network inspect services-network >/dev/null 2>&1 || { \
		echo "🔧 Creating shared network: services-network"; \
		sudo docker network create services-network; \
	}

docker_tag_push: ## 🚀 Build & push image to Docker Hub
	sudo docker build -t $(DOCKER_IMAGE) -f ./deploy/prod/Dockerfile .
	sudo docker push $(DOCKER_IMAGE)

### ==== DEV COMMANDS ====

dev_build: ## 🛠 Build dev image
	$(DOCKER_COMPOSE_DEV) build

dev_up: create_network ## 🚀 Start dev containers
	$(DOCKER_COMPOSE_DEV) up --build

dev_down: ## 🧹 Stop and remove dev containers
	$(DOCKER_COMPOSE_DEV) down -v

dev_restart: dev_down dev_up ## 🔁 Restart dev

dev_logs: ## 📜 View dev logs
	$(DOCKER_COMPOSE_DEV) logs -f

### ==== PROD COMMANDS ====

prod_pull: ## ⬇️ Pull production images
	$(DOCKER_COMPOSE_PROD) pull

prod_up: create_network ## 🚀 Start prod containers
	$(DOCKER_COMPOSE_PROD) up -d

prod_down: ## 🧹 Stop and remove prod containers
	$(DOCKER_COMPOSE_PROD) down -v

prod_restart: prod_down prod_up ## 🔁 Restart prod

prod_logs: ## 📜 View prod logs
	$(DOCKER_COMPOSE_PROD) logs -f
