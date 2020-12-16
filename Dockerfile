# Create temporary frontend backend
FROM golang:1.15.5-alpine3.12 as backend
COPY . /build/src/
RUN apk add --no-cache --virtual .build-deps --update gcc musl-dev git make && \
    cd /build/src/ && \
    #export GOPROXY=http://192.168.105.9/repository/go-proxy && export GOSUMDB=off && \
    go mod download && \
    make buildprod

# Create prodaction application container
#FROM scratch as production
FROM alpine:3.12.1 as production
WORKDIR /usr/app
COPY --from=backend /build/src/build/salary /usr/app/salary
COPY --from=backend /build/src/conf /usr/app/conf
COPY --from=backend /build/src/data /usr/app/data
COPY --from=backend /build/src/public /usr/app/public
COPY --from=backend /build/src/templates /usr/app/templates
RUN chmod +x /usr/app/salary
#ENTRYPOINT [ "/usr/app/salary start" ]