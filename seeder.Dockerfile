FROM golang:1.22-alpine

WORKDIR /app

# Add packages required for environemnt variable injection
RUN apk add --no-cache bash curl && curl -1sLf \
'https://dl.cloudsmith.io/public/infisical/infisical-cli/setup.alpine.sh' | bash \
&& apk add infisical=0.16.10

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

EXPOSE 3000

CMD ["infisical", "run", "--", "go", "run", ".", "--gseed"]