FROM golang:alpine3.13 as builder
LABEL "com.datadoghq.ad.check_names"='["bootcamp-feature-driven"]'
LABEL "com.datadoghq.ad.init_configs"='[{}]'
LABEL "com.datadoghq.ad.instances"='[{"golang-bootcamp-api_status_url": "http://%%host%%:%%port%%/health"}]'
LABEL "com.datadoghq.ad.logs"='[{"source": "bootcamp-feature-driven", "service": "bootcamp-feature-driven"}]'
RUN apk update && apk upgrade && apk add build-base git make sed
RUN go get github.com/silenceper/gowatch

#swagger addition via swagger-ui inyection
WORKDIR /go/src/github.com/eduardohoraciosanto/bootcamp-feature-driven/swagger
COPY ./oas/oas.yml ./swagger.yml
RUN git clone https://github.com/swagger-api/swagger-ui && \
  cp -r swagger-ui/dist/. . && rm -r swagger-ui/ && sed -i 's+https://petstore.swagger.io/v2/swagger.json+/swagger/swagger.yml+g' index.html


WORKDIR /go/src/github.com/eduardohoraciosanto/bootcamp-feature-driven
COPY . .

RUN GIT_COMMIT=$(git rev-parse --short HEAD) && \
  go build -o service -ldflags "-X 'github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/config.serviceVersion=$GIT_COMMIT'" ./cmd/api

FROM alpine:3.13 

COPY --from=builder /go/src/github.com/eduardohoraciosanto/bootcamp-feature-driven/cmd/api/service /
COPY --from=builder /go/src/github.com/eduardohoraciosanto/bootcamp-feature-driven/swagger /swagger

ENTRYPOINT [ "./service" ]