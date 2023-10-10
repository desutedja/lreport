FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /api-service

COPY env/application.yaml.sample env/application.yaml

COPY lreport lreport

RUN ls -la

ENTRYPOINT [ "./lreport" ]