# 源镜像
FROM golang:latest
# 作者
LABEL gajanlee="leejiazh6@gmail.com"

# 设置工作目录
WORKDIR $GOPATH/src/github.com/Posger

COPY . $GOPATH/src/github.com/Posger
RUN go get -v
RUN go build ./cmd/main.go

EXPOSE 8080

#最终运行docker的命令
ENTRYPOINT  ["./main"]
