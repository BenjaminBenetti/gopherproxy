FROM golang:1.23-alpine AS build

RUN mkdir -p /src/
WORKDIR /src/
COPY . . 
RUN go build -o ./bin/gopherproxyserver ./cmd/gopherproxyserver

FROM alpine:3.21.2 AS run 

COPY --from=build /src/bin/gopherproxyserver ./gopherproxyserver

ENV GIN_MODE=release
ENTRYPOINT [ "./gopherproxyserver" ]
