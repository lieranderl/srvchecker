# make docker file build  golang
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Run the build command

COPY . .

#add PORT
# ENV PORT=8080

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/aws/main.go


#FROM public.ecr.aws/lambda/provided:al2023

# run builded binary in golang AS runner

FROM amazonlinux:latest
ENV PORT=8000
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.8.4 /lambda-adapter /opt/extensions/lambda-adapter

# Install necessary packages
RUN yum -y update && \
    yum -y install ca-certificates mailcap shadow-utils && \
    yum clean all

# Create a group and user
# RUN groupadd -r app && useradd -r -g app app

# Tell docker that all future commands should run as the app user
# USER app
WORKDIR /var/task
COPY --from=builder /app/main /var/task/main
EXPOSE $PORT
CMD ["./main"]
