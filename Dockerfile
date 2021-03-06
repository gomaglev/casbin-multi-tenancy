FROM golang:latest 
COPY . /go/src/gin-casbin
WORKDIR /go/src/gin-casbin
RUN go get ./cmd/gin-casbin
RUN go build -ldflags "-w -s" -o ./cmd/gin-casbin/gin-casbin ./cmd/gin-casbin
CMD ["gin-casbin", "start", "web", "-c ./configs/config.toml -m ./configs/model.conf --menu ./configs/menu.yaml"]
