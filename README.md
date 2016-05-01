# portscan-compare-notify
Use Golang to parse and notify about Nmap results compared to a list of expected ports
Meant to be run at regular intervals, via Cron or something similar

#### Usage

###### dependancies (installed in the docker image below)
* Golang ~1.4
* nmap

###### example build and usage, also see `example.sh`
```
go build scan.go
# run it without flags to see the needed settings
./scan
Usage:
  -o, --nmapoptions="": options to pass to NMAP
  -e, --expected="": space separated list of ports that are expected to be found unfiltered (unfiltered = open or closed)
  -s, --smtpserver="": the SMTP server address, E.G 'smtp.example.com:587'
  -u, --username="": the SMTP username or email address
  -x, --password="": the SMTP user password
  -t, --to="": space separated list of 'to' address(es)
  -f, --from="": the 'from' email address
```

###### example binary run
```
sudo ./scan \
  --nmapoptions="-PN -n -sS -p1-65535 example.com" \
  --expected="22 80 443" \
  --smtpserver=smtp.example.net:587 \
  --username=my-addy@example.com \
  --password=my-pass
  --to="my-addy@example.com my-addy777@example.com" \
  --from=portscan@example.com \
```

###### example using Docker
```
docker run --rm \
  snarlysodboxer/portscan-compare-notify:latest \
  --nmapoptions="-PN -n -sS -p1-65535 example.com" \
  --expected="22 80 443" \
  --smtpserver=smtp.example.net:587 \
  --username=my-addy@example.com \
  --password=my-pass
  --to="my-addy@example.com" \
  --from=portscan@example.com \
```

###### TODO
* Document better
* Refactor
