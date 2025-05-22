SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)), /bin/ash, /bin/bash)



 
run:
	go run apis/services/sales/main.go | go run apis/tooling/logfmt/main.go

help:
	go run apis/services/sales/main.go --help

version:
	go run apis/services/sales/main.go --version

curl-test:
	curl -il -X GET http://localhost:3000/test

curl-live:
	curl -il -X GET http://localhost:3000/liveness

curl-ready:
	curl -il -X GET http://localhost:3000/readiness

curl-error:
	curl -il -X GET http://localhost:3000/testerror

curl-panic:
	curl -il -X GET http://localhost:3000/testpanic

admin:
	go run apis/tooling/admin/main.go
# ==============================================================================
# Define dependencies

GOLANG          := golang:1.24
ALPINE          := alpine:3.21
KIND            := kindest/node:v1.32.0
POSTGRES        := postgres:17.2
GRAFANA         := grafana/grafana:11.4.0
PROMETHEUS      := prom/prometheus:v3.0.0
TEMPO           := grafana/tempo:2.6.0
LOKI            := grafana/loki:3.3.0
PROMTAIL        := grafana/promtail:3.3.0

KIND_CLUSTER    := ardan-starter-cluster
NAMESPACE       := sales-system
SALES_APP       := sales
AUTH_APP        := auth
BASE_IMAGE_NAME := localhost/ardanlabs
# your version could be a shell command like 
# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"
VERSION         := 0.0.1
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)

# ==============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner


dev-status:
	watch -n 2 kubectl get pods -o wide --all-namespaces

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

# ------------------------------------------------------------------------------

dev-load:
	kind load docker-image $(SALES_IMAGE) --name $(KIND_CLUSTER)

dev-restart:
	kubectl rollout restart deployment $(SALES_APP) --namespace=$(NAMESPACE)


dev-apply:
	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(SALES_APP) --timeout=120s --for=condition=Ready

dev-describe-deployment:
	kubectl describe deployment $(SALES_APP) --namespace=$(NAMESPACE)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(SALES_APP)

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(SALES_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run apis/tooling/logfmt/main.go -service=$(SALES_APP)

dev-update: build dev-load dev-restart


dev-update-apply: build dev-load dev-apply
# ==============================================================================
# Building the images
build: sales

sales:
	docker build \
		-f zarf/docker/dockerfile.sales \
		-t $(SALES_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.
# ==============================================================================
# Module support
tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
# Metrics and Tracing

metrics-view-sc:
	~/go/bin/expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

metrics-view:
	~/go/bin/expvarmon -ports="localhost:4020" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

grafana:
	open http://localhost:3100/

statsviz:
	open http://localhost:3010/debug/statsviz

# ==============================================================================
# # RSA Keys
# 	To generate a private/public key PEM file.
# 	$ openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# 	$ openssl rsa -pubout -in private.pem -out public.pem
# 	$ ./admin genkey
	