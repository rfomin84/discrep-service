FROM golang:1.17-stretch

WORKDIR /src

RUN go install -tags 'clickhouse' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.1 && \
    go install github.com/x-motemen/gore/cmd/gore@latest

COPY . .

RUN go mod download

RUN go mod vendor

RUN go build -v -o server cmd/server/server.go
RUN go build -v -o scheduler cmd/scheduller/scheduller.go

#RUN ./server
#RUN ./scheduler