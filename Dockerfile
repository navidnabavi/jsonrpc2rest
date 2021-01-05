FROM golang:1.15
ADD ./go.mod /go.mod
RUN go mod download

ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o jsonrpc2rest cmd/main.go

FROM scratch
COPY --from=0 /app/jsonrpc2rest .
VOLUME [ "/var/lib/jsonrpc2rest" ]
CMD ["./jsonrpc2rest", "config.json"]


