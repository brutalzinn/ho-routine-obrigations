# This docker file is to be used without installing air
FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /main
CMD ["/main"]
EXPOSE 80