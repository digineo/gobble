# gobble

Read `systemd` journal for a given process, look for a go panic stacktrace and report it to sentry, if the process crashed and dumped a stack trace.


## Building

```terminal
$ go build -o gobble ./...
```


## Running

```terminal
$ ./gobble -dsn $SENTRY_DSN -service $SERVICE_NAME
```

or to integrate with `systemd`:

```ini
[Unit]
Description=Template Service FIle
After=network.target

[Service]
PrivateTmp=true
EnvironmentFile=/etc/app/env
ExecStart=/usr/local/bin/app-binary
ExecStopPost=/usr/local/bin/gobble -dsn $SENTRY_DSN -service %n

[Install]
WantedBy=multi-user.target
```
