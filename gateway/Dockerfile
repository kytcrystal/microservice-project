# Specify a parent image
FROM golang:1.21.4-alpine as golang

# Create an app directory
WORKDIR /app
 
# Copy go.mod and go.sum into /app
COPY go.mod go.sum ./

# Install Go dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . ./

# Build app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main main/main.go 

RUN chmod +x /main

FROM gcr.io/distroless/static-debian11

COPY --from=golang /main .
COPY --from=golang /app/configuration-deploy.yaml .

# Tell Docker which network port the container listens on
EXPOSE 3333
 
# Specify the executable command that runs when the container starts
CMD [ "/main" ]
