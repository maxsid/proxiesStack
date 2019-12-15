FROM golang:1.13.4
ENV PACKAGE_DIR /go/src/github.com/maxsid/proxiesStack
RUN mkdir -p ${PACKAGE_DIR}
ADD . ${PACKAGE_DIR}
WORKDIR ${PACKAGE_DIR}
RUN go get