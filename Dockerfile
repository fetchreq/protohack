FROM golang:1.21-alpine as builder

ENV GOOS linux
ENV CGO_ENABLED 0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build 

FROM golang:1.21-alpine as runner

COPY --from=builder /app/protohack .

EXPOSE 10000

CMD ./protohack p4
