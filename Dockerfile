FROM ubuntu:14.04
MAINTAINER david amick <docker@davidamick.com>

RUN ["/bin/bash", "-c", "apt-get update -qq && apt-get install -qy nmap"]

COPY scan /scan

ENTRYPOINT ["/scan"]
