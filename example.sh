#!/bin/bash

# sudo ./scan \
docker run --rm \
  snarlysodboxer/portscan-compare-notify:latest \
  --nmapoptions="-PN -n -sS -p1-65535 example.com" \
  --expected="22 80 443" \
  --to="myuser@example.com" \
  --from=portscan@example.com \
  --smtpserver=smtp.example.net:587 \
  --username=mysmtpuser@example.com \
  --password=mypass

