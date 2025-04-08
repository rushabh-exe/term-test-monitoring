FROM golang:1.24
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
# RUN go run ./cmd/main.go
EXPOSE 5051
CMD ["go", "run", "./cmd/main.go"]
