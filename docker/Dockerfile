FROM proxies-stack-source as builder
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o proxiesStack .

FROM scratch
COPY --from=builder /go/src/github.com/maxsid/proxiesStack/proxiesStack /app/
WORKDIR /app
EXPOSE 8080
CMD ["./proxiesStack"]