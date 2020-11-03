# Tendang!

Simple and secure deployment proxy.

### Run

PORT=8000 ./tendang

### Usage

```
curl --header "Content-Type: application/json" --request POST --data '{"name":"xyz","value":"123", "token": "RYWKkMSGK7tCb7jCSVZNmJzWneNDb2funq6kSLUPDVCgL8gAMPBfUWLyKtQdLp7A"}' http://localhost:8000
```

### Daemon

`/lib/systemd/system/tendang.service`

```
[Unit]
Description=tendang

[Service]
Environment="PORT=8000"
User=deployment
WorkingDirectory=/path/to/workdir
ExecStart=/usr/bin/tendang
ExecStop=/bin/kill -9 $MAINPID
StandardOutput=file:/var/log/tendang.log
StandardError=file:/var/log/tendang.log

[Install]
WantedBy=multi-user.target
```
