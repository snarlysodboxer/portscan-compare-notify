# portscan-compare-notify
Use Golang to parse and notify about Nmap results compared to a list of expected ports
Meant to be run at regular intervals, via Cron or something similar

#### Usage

###### dependancies (installed in the docker image below)
* Golang ~1.4
* nmap

###### example build and usage
```
go build scan.go
# run it without flags to see the needed settings
./scan
Usage:
  -e, --expected="": space separated list of ports that are expected to be found unfiltered (unfiltered = open or closed)
  -f, --from="": the 'from' email address
  -h, --host="": IP or resolvable hostname
  -m, --parallelism="": the Nmap --min-parrallelism setting, E.G. '1024'
  -x, --password="": the SMTP user password
  -p, --range="": dash separated port range to scan, E.G. '1-65535'
  -s, --server="": the SMTP server address, E.G 'smtp.example.com:587'
  -t, --to="": space separated list of 'to' address(es)
  -u, --username="": the SMTP username or email address
```

###### example binary run
```
sudo ./scan \
  --host=me.example.com \
  --range=1-65535 \
  --expected="22 80 443" \
  --parallelism=1024 \
  --to="my-addy@example.com my-addy777@example.com" \
  --from=portscan@example.com \
  --server=smtp.example.net:587 \
  --username=my-addy@example.com \
  --password=my-pass
```

###### example using Docker
```
docker run --rm \
  snarlysodboxer/portscan-compare-notify:latest \
  --host=me.example.com \
  --range=1-65535 \
  --expected="22 80 443" \
  --parallelism=1024 \
  --to="my-addy@example.com my-addy777@example.com" \
  --from=portscan@example.com \
  --server=smtp.example.net:587 \
  --username=my-addy@example.com \
  --password=my-pass
```

###### TODO
* Document better
* Refactor
* Measure length of time for each scan
* Log to log file, not just email
* Convert to passing entire nmap command
