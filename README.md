# nats-transport
HTTP to NATS (https://nats.io/) golang transport (http.RoundTripper) implementation

inspired by and based on https://github.com/sohlich/nats-proxy

## dev

generate protobuf

```bash
protoc --go_opt=paths=source_relative --go_out=. protobuf.proto
```