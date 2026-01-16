FROM golang:1.25.6

WORKDIR /reverie-api

RUN go install github.com/air-verse/air@v1.64.0
RUN go install github.com/bufbuild/buf/cmd/buf@v1.63.0
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0
RUN go install github.com/fullstorydev/grpcui/cmd/grpcui@v1.4.3

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1
