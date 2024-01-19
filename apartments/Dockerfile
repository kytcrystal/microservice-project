# Specify a parent image
FROM golang:1.21.4

# Create an app directory
WORKDIR /app
 
# Copy go.mod and go.sum into /app
COPY go.mod go.sum ./

# Install Go dependencies
RUN go mod download

# Copy source code
COPY . ./

# Build app
RUN CGO_ENABLED=0 GOOS=linux go build -o /apartments main/apartments.go 

RUN chmod +x /apartments

# Tell Docker which network port the container listens on
EXPOSE 3000
 
# Specify the executable command that runs when the container starts
CMD [ "/apartments" ]