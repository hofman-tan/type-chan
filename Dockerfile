FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o typechan

FROM alpine:3.17

# display color output in terminal
# https://stackoverflow.com/questions/33493456/docker-bash-prompt-does-not-display-color-output
ENV TERM=xterm-256color

COPY --from=builder /app/typechan /usr/bin
ENTRYPOINT ["sh"]
