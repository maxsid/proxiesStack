FROM alpine:3.10.3

RUN adduser -D user
USER user
COPY proxies-stack .

CMD ["./proxies-stack"]