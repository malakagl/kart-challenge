
### Build an API server implementing our OpenAPI spec for food ordering API in Go.

Go version used: 1.25

API Spec Published at https://orderfoodonline.deno.dev/public/openapi.yaml

### How to run

```
make run
```

### Run unit tests

```
make test
```

### Useful commands

```
# install dev dependencies
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# to create new migration file
migrate create -ext sql -dir db/migrations -seq init_db

# to start the database
docker compose up -d postgres

# to stop the database
docker compose down postgres 

# apply migations
migrate -path db/migrations -database "postgres://user:password@localhost:5432/test?sslmode=disable" up

# remove migrations
migrate -path db/migrations -database "postgres://user:password@localhost:5432/test?sslmode=disable" down

# install lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 <- fixed with latest version

# check coverage
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# deploy
minikube start --memory 10240 --cpus 4
docker build -f ./docker/Dockerfile -t kart-challenge:latest .
minikube image load kart-challenge:latest

# mount config and data
minikube mount ./config:/mnt/config
minikube mount ./promocodes:/mnt/promocodes
minikube mount ./db:/mnt/db
kubectl create configmap postgres-init \
  --from-file=docker_postgres_init.sql=./scripts/docker_postgres_init.sql \
  -n kart-challenge
kubectl apply -k ./deployment/k8s/

# verify
kubectl get deployments -A
minikube service kart-challenge --url -n kart-challenge
minikube service kart-challenge --url -n postgres
kubectl get pods -n kart-challenge
kubectl logs kart-challenge-5d8ccdfcdd-qh8cb -n kart-challenge -f
kubectl exec -it -n kart-challenge kart-challenge-6d69f6b49-hn5fb -- /bin/sh

# clean up
kubectl delete service kart-challenge -n kart-challenge
kubectl delete deployment kart-challenge -n kart-challenge
minikube stop
minikube delete --all
```

### Notes
- Took around 25 minutes to load all coupon codes to the database.
- To read gz files via code and scan it took around 15 seconds.

### Tasks
- [x] Implement the API server
- [x] Implement the OpenAPI spec
- [x] Implement the unit tests
- [x] Implement promo code functionality
- [x] Implement the database
- [x] Implement the database migrations
- [x] Implement logging
- [ ] Implement the error handling
- [x] Implement the configuration management
- [ ] Implement the integration tests
- [x] Implement the Dockerfile
- [x] Implement the Makefile
- [x] Implement the README.md
- [ ] Implement the CI/CD pipeline
- [ ] Implement the GitHub Actions workflow
- [ ] Implement the GitHub Pull Requests
- [ ] Implement money package for handling money
- [ ] Implement the caching
- [ ] Implement the rate limiting
- [ ] Implement the security
- [ ] Implement the monitoring
- [x] Implement the tracing
- [x] Create private schema for postgres
- [x] Add lint