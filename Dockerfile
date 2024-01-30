FROM golang:1.20 as builder

WORKDIR /builds/ports
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-w -s" -o ports ./cmd/ports

FROM scratch
COPY --from=builder /builds/ports/ports .
CMD ["./ports"]