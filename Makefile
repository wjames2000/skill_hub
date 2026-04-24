.PHONY: build push deploy deploy-staging deploy-prod logs status migrate verify help

DOCKER_REGISTRY ?= ghcr.io
IMAGE_TAG ?= latest

help:
	@echo "SkillHub Pro - Makefile"
	@echo ""
	@echo "Deployment targets:"
	@echo "  build             Build all Docker images locally"
	@echo "  push              Push all Docker images to registry"
	@echo "  deploy            Deploy all services to K8s"
	@echo "  deploy-staging    Deploy to staging environment"
	@echo "  deploy-prod       Deploy to production environment"
	@echo "  status            Show deployment status in K8s"
	@echo "  logs <service>    Tail logs for a service"
	@echo "  migrate           Run database migrations"
	@echo "  verify            Verify deployment health"
	@echo "  infra-up          Start infrastructure services (Docker)"
	@echo "  infra-down        Stop infrastructure services"
	@echo "  mon-up            Start monitoring stack (Docker)"
	@echo "  mon-down          Stop monitoring stack"

# ---------- Docker Build ----------

build:
	docker build -t $(DOCKER_REGISTRY)/skill-hub/router-api:$(IMAGE_TAG) \
		-f deployments/Dockerfile.backend \
		--build-arg SERVICE=router-api backend
	docker build -t $(DOCKER_REGISTRY)/skill-hub/sync-worker:$(IMAGE_TAG) \
		-f deployments/Dockerfile.backend \
		--build-arg SERVICE=sync-worker backend
	docker build -t $(DOCKER_REGISTRY)/skill-hub/admin-api:$(IMAGE_TAG) \
		-f deployments/Dockerfile.backend \
		--build-arg SERVICE=admin-api backend
	docker build -t $(DOCKER_REGISTRY)/skill-hub/frontend:$(IMAGE_TAG) \
		-f deployments/Dockerfile.frontend frontend

push: build
	docker push $(DOCKER_REGISTRY)/skill-hub/router-api:$(IMAGE_TAG)
	docker push $(DOCKER_REGISTRY)/skill-hub/sync-worker:$(IMAGE_TAG)
	docker push $(DOCKER_REGISTRY)/skill-hub/admin-api:$(IMAGE_TAG)
	docker push $(DOCKER_REGISTRY)/skill-hub/frontend:$(IMAGE_TAG)

# ---------- K8s Deploy ----------

deploy:
	kubectl apply -f deployments/k8s/namespace.yaml
	kubectl apply -f deployments/k8s/configmap.yaml
	kubectl apply -f deployments/k8s/secret.yaml
	kubectl apply -f deployments/k8s/backend-service.yaml
	kubectl apply -f deployments/k8s/backend-deployment.yaml
	kubectl apply -f deployments/k8s/frontend-deployment.yaml
	kubectl apply -f deployments/k8s/ingress.yaml
	kubectl rollout status deployment/router-api -n skill-hub
	kubectl rollout status deployment/sync-worker -n skill-hub
	kubectl rollout status deployment/admin-api -n skill-hub
	kubectl rollout status deployment/frontend -n skill-hub

deploy-infra:
	kubectl apply -f deployments/k8s/infra-services.yaml
	kubectl apply -f deployments/k8s/pvc.yaml

# ---------- Docker Compose ----------

infra-up:
	docker compose -f deployments/docker-compose.infra.yml up -d

infra-down:
	docker compose -f deployments/docker-compose.infra.yml down

mon-up:
	docker compose -f deployments/monitoring/docker-compose.monitoring.yml up -d

mon-down:
	docker compose -f deployments/monitoring/docker-compose.monitoring.yml down

# ---------- Utils ----------

status:
	kubectl get all -n skill-hub

logs:
	kubectl logs -n skill-hub -l app=$(filter-out $@,$(MAKECMDGOALS)) --tail=100 -f

migrate:
	bash deployments/scripts/migrate.sh

verify:
	bash deployments/scripts/verify-deployment.sh

%:
	@:
