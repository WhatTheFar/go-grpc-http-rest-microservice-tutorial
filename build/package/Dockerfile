# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.11
ARG GO_VERSION=1.11

# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
# Git is required for fetching the dependencies.
RUN apk add --no-cache ca-certificates git

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
# RUN go mod download

# Import the code from the context.
COPY ./ ./

# Build the executable to `/app`. Mark the build as statically linked.
# RUN CGO_ENABLED=0 go build \
#     -installsuffix 'static' \
#     -o /app/cmd/server cmd/server/main.go
RUN CGO_ENABLED=0 go build \
    -mod vendor \
    -installsuffix 'static' \
    -o /app/cmd/server cmd/server/main.go

# Install grpc_health_probe
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.0 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# Final stage: the running container.
FROM alpine AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the first stage.
COPY --from=builder /app/cmd/server /app/cmd/server/

COPY --from=builder /src/config/config.yaml /app/config/
COPY --from=builder /src/config/dev-config.yaml /app/config/
COPY --from=builder /src/config/staging-config.yaml /app/config/

COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

# Declare the port on which the webserver will be exposed.
# As we're going to run the executable as an unprivileged user, we can't bind
# to ports below 1024.
EXPOSE 8080

# Perform any further action as an unprivileged user.
USER nobody:nobody

ARG APP_ENVIRONMENT
ENV APP_ENVIRONMENT ${APP_ENVIRONMENT}

WORKDIR /app/cmd/server

# Run the compiled binary.
# CMD ./server -env=${APP_ENVIRONMENT}
ENTRYPOINT ["/app/cmd/server/server", "-env=develop"]