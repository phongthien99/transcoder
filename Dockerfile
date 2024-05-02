
FROM golang:1.21-alpine AS builder
# It is important that these ARG's are defined after the FROM statement

ARG APP
ARG VERSION

WORKDIR /src

COPY apps/$APP/go.mod apps/$APP/go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    --ldflags="-X 'github.com/sigmaott/gest/package/technique/version.Version=$VERSION' \
     -X 'github.com/sigmaott/gest/package/technique/version.Date=$(date)'"\
    -o app \
    ./apps/$APP/cli/api/main.go


# RUN mv app /src



# Final stage: the running container.
FROM scratch AS final
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the compiled executable from the first stage.
COPY --from=builder /src/app /app
# Perform any further action as an unprivileged user.
CMD ["/app","-c","/src/config/default.yaml"]
