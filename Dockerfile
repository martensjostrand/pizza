FROM golang@sha256:8dea7186cf96e6072c23bcbac842d140fe0186758bcc215acb1745f584984857 as builder

WORKDIR $GOPATH/src/github.com/martensjostrand/pizza-dough
COPY . .

RUN adduser -D -g '' appuser
RUN go build -o /go/bin/recipe

FROM scratch
COPY --from=builder /go/bin/recipe /go/bin/recipe
COPY --from=builder /etc/passwd /etc/passwd
USER appuser

ENTRYPOINT ["/go/bin/recipe"]