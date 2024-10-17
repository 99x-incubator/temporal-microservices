# build from the root level of the repo as we need to copy some other packages 
# docker build --tag admin_gateway -f microservices/admin_gateway/Dockerfile .

FROM golang:1.23

WORKDIR /app

#add local dependencies 
ADD temporal/workflows ./temporal/workflows/

#add admin_gateway sorce
ADD microservices/admin_gateway ./microservices/admin_gateway/

#set working directory to admin_gateway
WORKDIR /app/microservices/admin_gateway

# download module dependencies
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /admin_gateway

EXPOSE 8081

# Run
CMD ["/admin_gateway"]


#todo
# run a lean image

# build from the root level of the repo as we need to copy some other packages 
# docker build --tag admin_gateway -f microservices/admin_gateway/Dockerfile .
# docker build --progress=plain --no-cache --tag admin_gateway -f microservices/admin_gateway/Dockerfile .

