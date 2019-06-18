FROM golang

RUN apt -y update  && apt -y install unzip

RUN mkdir /srv/ -p
WORKDIR /srv
# Install protoc
RUN mkdir -p local && \
	PROTOC_URL=https://github.com/google/protobuf/releases/download/v3.4.0/protoc-3.4.0-linux-x86_64.zip && \
	curl -fSsL $PROTOC_URL -o protoc.zip && \
	unzip protoc.zip -d local

RUN go get -u github.com/googleapis/gnostic-grpc

# Install oapi-codegen
RUN go get -u github.com/deepmap/oapi-codegen/cmd/oapi-codegen

ADD . /code
WORKDIR /code
RUN make go-build
ENTRYPOINT ["make", "run"]



