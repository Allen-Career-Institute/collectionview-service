FROM 537984406465.dkr.ecr.ap-south-1.amazonaws.com/golang:1.22.7 AS builder

COPY . /src
WORKDIR /src

ARG GIT_TOKEN
ENV GOPRIVATE=github.com/Allen-Career-Institute/*
ENV GOSUMDB=off

COPY . /src
WORKDIR /src

RUN apt-get update && apt-get install git -y
RUN git config --global url."https://$GIT_TOKEN:@github.com/".insteadOf "https://github.com/"

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim
USER root

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        curl \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app
COPY --from=builder /src/configs /app
WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["./collectionview-service", "-conf", "/app/data/conf/config_stage.yaml"]
