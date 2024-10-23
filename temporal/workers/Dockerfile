# build from the root level of the repo as we need to copy some other packages 
# docker build --tag admin_worker -f temporal/workers/Dockerfile .

FROM golang:1.23

WORKDIR /app

#add local dependencies 
ADD temporal/workflows ./temporal/workflows/

#add admin_worker sorce
ADD temporal/workers ./temporal/workers/

#set working directory to admin_worker
WORKDIR /app/temporal/workers

# download module dependencies
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /admin_worker

# Run
CMD ["/admin_worker"]


#todo
# run a lean image

# build from the root level of the repo as we need to copy some other packages 
# docker build --tag admin_worker -f temporal/workers/Dockerfile .
