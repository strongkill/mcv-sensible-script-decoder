#showMonitor
FROM golang:latest
ENV GO_FILE_PATH /mnt/app
MAINTAINER  kyle "keung@showpay.top"
WORKDIR ${GO_FILE_PATH}/goWork/src/github.com/ShowPay/script-decoder
ADD . ${GO_FILE_PATH}/goWork/src/github.com/ShowPay/script-decoder
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN go build
EXPOSE 9030
CMD ["go", "run", "main.go"]