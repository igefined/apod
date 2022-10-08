FROM golang:1.18-alpine as builder

# Create appuser
RUN adduser -D -g '' appuser

WORKDIR /app
COPY . .

ENV PORT=8080

EXPOSE 8080

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/betera github.com/igilgyrg/betera-test/cmd/http
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]