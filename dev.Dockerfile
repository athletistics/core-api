FROM golang:1.19.0-alpine

WORKDIR /app

# Add packages required for environemnt variable injection
RUN apk add --no-cache bash curl && curl -1sLf \
'https://dl.cloudsmith.io/public/infisical/infisical-cli/setup.alpine.sh' | bash \
&& apk add infisical

# Install "hot-reloading" moudle
RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

EXPOSE 3000

CMD ["infisical", "run", "--", "air", "-c", ".air.toml"]