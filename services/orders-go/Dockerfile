FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /orders .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /orders /orders
EXPOSE 8080
ENTRYPOINT ["/orders"]
