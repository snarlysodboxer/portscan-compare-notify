# portscan-compare-notify
Use Golang to parse and notify about Nmap results compared to a list of expected ports

#### Usage

E.G. build and run
```
go build scan.go
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

E.G. using Docker
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

TODO
* Account for unshown ports that are 'closed'
* Document better
