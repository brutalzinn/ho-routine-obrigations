FROM golang:latest
ENV TZ="America/Sao_Paulo"
RUN go install github.com/cosmtrek/air@latest
WORKDIR /app
ENTRYPOINT ["air"]