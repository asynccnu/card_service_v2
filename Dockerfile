FROM golang:1.14.3
ENV GO111MODULE "on"
ENV GOPROXY "https://goproxy.cn"
WORKDIR /src/card_service_v2
COPY . /src/card_service_v2
RUN make
#FROM ubuntu
#COPY --from=0 /src/card_service_v2 .
EXPOSE 8080
CMD ["./main", "-c", "conf/config.yaml"]
