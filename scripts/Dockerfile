# Create temporary frontend backend
FROM golang:rc-alpine as backend
COPY . /build/src/
RUN apk add --no-cache --virtual .build-deps --update gcc musl-dev git make\
    && cd /build/src/ && export GOPROXY=http://192.168.105.9/repository/go-proxy \
    && export GOSUMDB=off &&  go mod download && make buildprod

# Create prodaction application container
FROM alpine:latest as production
WORKDIR /usr/app
COPY --from=backend /build/src/build/salary ./salary
COPY --from=backend /build/src/conf ./conf
COPY --from=backend /build/src/data ./data
COPY --from=backend /build/src/public ./public
COPY --from=backend /build/src/templates ./templates
RUN chmod +x ./salary
ENTRYPOINT [ "./salary start" ]