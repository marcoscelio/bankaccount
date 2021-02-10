FROM golang:1.14.3-alpine  AS build

RUN mkdir -p /go/src
COPY . /go/src
WORKDIR /go/src
RUN apk add --no-cache git
RUN git --version
ENV GOPATH /go/src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/bankaccount
COPY ./config.toml /bin/config.toml
ARG PORT  
ARG READ_TIMEOUT  
ARG WRITE_TIMEOUT  
ENV PORT "$PORT"
ENV READ_TIMEOUT "$READ_TIMEOUT"
ENV WRITE_TIMEOUT "$WRITE_TIMEOUT"
EXPOSE ${PORT}
RUN env
ENTRYPOINT ["bankaccount"]