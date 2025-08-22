ARG ARCHITECTURE

FROM registry.cloud.com/beacon/ci/golang:1.22.2-alpine3.19-${ARCHITECTURE} AS builder

ARG GIT_BRANCH
ARG GIT_HASH
ARG BUILD_TS
ARG PROJ_NAME

RUN mkdir -p /root/${PROJ_NAME}
WORKDIR /root/${PROJ_NAME}
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCHITECTURE} go build -ldflags "-X \"main.GitBranch=${GIT_BRANCH}\" -X \"main.GitHash=${GIT_HASH}\" -X \"main.BuildTS=${BUILD_TS}\"" -o ./bin/ ./...

FROM registry.cloud.com/beacon/ci/alpine:3.19-v2-aws-${ARCHITECTURE}

USER root

RUN mkdir /.aws && mkdir /root/cache
RUN chmod 777 /.aws && chmod 777 /root/cache

COPY ./aws/credentials /.aws

ARG PROJ_NAME

WORKDIR /root/

COPY --from=builder /root/${PROJ_NAME}/bin/ ./
COPY --from=builder /root/${PROJ_NAME}/conf/config.yaml conf/
COPY --from=builder /root/${PROJ_NAME}/locales/ locales/
