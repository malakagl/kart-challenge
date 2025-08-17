
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
- [ ] Implement the tracing
- [ ] Create private schema for postgres