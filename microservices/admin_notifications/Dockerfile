FROM golang:1.23

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /admin_notifications

EXPOSE 8082

# Run
CMD ["/admin_notifications"]


#todo
# run a lean image

# build as 
# docker build --tag admin_notifications .
# (runtime  - this is a we hub. clients will connect to this service, and this service will push messages to registred clinets.)

