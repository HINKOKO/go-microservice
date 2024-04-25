FROM golang:1.22.2-alpine as dependencies

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

FROM dependencies AS build
COPY . ./
RUN CGO_ENABLED=0 go build -o /main -ldflags="-w -s" .

FROM golang:1.22.2-alpine
COPY --from=build /main /main
CMD ["/main"]
