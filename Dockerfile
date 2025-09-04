FROM golang:1.25.0-alpine AS build

WORKDIR /app

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

# Clean up and vendor dependencies
RUN go mod tidy && go mod vendor

COPY . .
# Generate docs (simplified for Fiber)
RUN swag init --generalInfo ./cmd/api/main.go --output ./docs
# Entrypoint script
COPY migration_command/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Verify docs
RUN ls -la ./docs && [ -f ./docs/swagger.json ]

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /app/docs /app/docs
COPY --from=build /app/migration_command /app/migration_command
# Copy the goose binary from the build stage
COPY --from=build /app/internal/database/migrations /app/internal/database/migrations
COPY --from=build /go/bin/goose /app/bin/goose
RUN chmod +x /app/migration_command/entrypoint.sh
RUN chmod +x /app/bin/goose

# Update PATH to include goose
ENV PATH="/app/bin:${PATH}"
# Make sure your main application is configured to serve from ./docs
ENV SWAGGER_JSON=/app/docs/swagger.json
EXPOSE ${PORT}
ENTRYPOINT ["/app/migration_command/entrypoint.sh"]
CMD ["./main"]


