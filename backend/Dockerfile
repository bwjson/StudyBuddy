FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

FROM alpine:latest AS runner

WORKDIR /root

COPY --from=builder /app/main .
COPY .env .


CMD sh -c "sleep 5 && migrate -path ./schema -database 'postgres://db_iviy_user:89OmsjqnfM8ZQSlthfFZid08iX3MGF0z@dpg-cts369rqf0us73dofpkg-a:5432/db_iviy?sslmode=disable' up && ./main"