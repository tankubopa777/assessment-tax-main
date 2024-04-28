FROM golang:1.21 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o assessment-tax .

FROM alpine:latest as production
WORKDIR /root/
COPY --from=builder /app/assessment-tax .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./assessment-tax"]
