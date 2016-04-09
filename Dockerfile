FROM fedora

RUN dnf -y install go
RUN dnf -y install git

ENV GOPATH /root/gopath

RUN go get github.com/subgraph/go-seccomp
