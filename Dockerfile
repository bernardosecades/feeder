# STEP 1: Build executable sever and purge binaries with UPX (it is an advanced executable file compressor)
FROM golang:1.16 AS builder
# Add Maintainer Info
LABEL maintainer="Bernardo Secades <bernardosecades@gmail.com>"

RUN apt-get update && \
	apt-get install --no-install-recommends -y \
	ca-certificates upx-ucl

RUN useradd appuser

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/bernardosecades/feeder
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build the Go app

# -ldflags="-w -s" reduce size of bnary
RUN cd cmd/feedersrv && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/feedersrv .

# Compress binary files
RUN upx /bin/feedersrv

# Expose port 3000 to the outside world
EXPOSE 3333

# Command to run the executable
ENTRYPOINT ["/bin/feedersrv"]

# STEP 2: Build a small image
FROM alpine:3.13 AS production

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# passwd file to appuser created in first stage
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /bin/main /go/bin/feedersrv

RUN  ls -la /go/bin/

# Use an unprivileged user.
USER appuser
# Expose port 3000 to the outside world
EXPOSE 3333

# Command to run the executable
ENTRYPOINT ["/go/bin/feedersrv"]
