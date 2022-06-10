FROM golang:1.18-buster AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /telemetry-service -buildvcs=false

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /telemetry-service /telemetry-service
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/telemetry-service"]