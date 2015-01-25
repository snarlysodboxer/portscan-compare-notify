FROM ubuntu:14.04
MAINTAINER david amick <docker@davidamick.com>

RUN ["/bin/bash", "-c", "apt-get update -qq && apt-get install -qy gcc libc6-dev libc6-dev-i386 git-core nmap"]
RUN ["/bin/bash", "-c", "git clone https://go.googlesource.com/go /opt/go"]
RUN ["/bin/bash", "-c", "cd /opt/go && git checkout go1.4.1"]
RUN ["/bin/bash", "-c", "cd /opt/go/src && ./all.bash"]
RUN ["/bin/bash", "-c", "mkdir -p /opt/gopath"]
ENV PATH=/opt/go/bin:$PATH
ENV GOPATH=/opt/gopath
RUN ["/bin/bash", "-c", "go get github.com/ogier/pflag"]

COPY scan.go /opt/scan.go
RUN ["/bin/bash", "-c", "go build -o /opt/scan /opt/scan.go"]

ENTRYPOINT ["/opt/scan"]
