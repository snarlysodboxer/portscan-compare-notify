FROM ubuntu:14.04
MAINTAINER david amick <docker@davidamick.com>

RUN ["/bin/bash", "-c", "apt-get update -qq && apt-get install -qy nmap ca-certificates"]

COPY scan /scan

ENTRYPOINT ["/scan"]
