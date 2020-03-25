# gobble

Reads `systemd` journal for a given service, and looks for a Go panic
stacktrace. If one can be found, it is reported to Sentry.


## Why?

In Go, a `panic()` will usually crash a program. These crashes *can* be
intercepted, analyzed and submitted to external services, but only if
this recovering happens in the same Go routine; capturing panics in
different Go routines is not possible from a centralized handler.

The only viable option is to let the program crash and dump its stack
trace to stderr.

From there, we can use different tools, like [panicparse][] to analyze
the stack in a controlled environment (and not in a program where a fatal
error just happened, for whatever reason).

Gobble helps to automate this process, by using systemd's
[`ExecStopPost=`][ExecStopPost] hook.

[panicparse]: github.com/maruel/panicparse
[ExecStopPost]: https://www.freedesktop.org/software/systemd/man/systemd.service.html#ExecStopPost=


## Building

You can find binaries for tagged releases on the [Github release page][].
To build it yourself, you need a Go tool chain and systemd's development
package (on Debian, this is `libsystemd-dev`).

Then clone or download this project and run this command in the project
directory:

```console
$ go build -o gobble ./...
```

[Github release page]: https://github.com/digineo/gobble/releases


## Running

```console
$ ./gobble -dsn $SENTRY_DSN -service $SERVICE_NAME
$# or
$ export SENTRY_DSN="https://...@sentry.example.com/projects/42"
$ ./gobble -service $SERVICE_NAME
```

To integrate with `systemd`, you can add a `ExecStopPost` line to
the service's unit file:

```ini
[Unit]
Description=Template Service File
After=network.target

[Service]
PrivateTmp=true
EnvironmentFile=/etc/app/env
ExecStart=/usr/local/bin/app-binary
ExecStopPost=/usr/local/bin/gobble -dsn $SENTRY_DSN -service %n

[Install]
WantedBy=multi-user.target
```

In general, it is advisable to just add a stanza. For a service named
`yourapp`, create a file named `/etc/systemd/system/yourapp.service.d/gobble.conf`
(you might need to create the `yourapp.service.d` directory first) and
add the following contents:

```ini
[Service]
ExecStopPost=/usr/local/bin/gobble -dsn $SENTRY_DSN -service %n
```

The latter approach is advisable if you don't have the unit file under
your direct control and want to ensure the hooks stays in place if you
update that service.

In both cases, you'll need to reload systemd:

```console
# systemctl daemon-reload
```

## Known issues

- [ ] If the service runs with an unprivileged user, executing gobble
  in a post-exec-hook will fail with "Failed to connect to bus" when
  trying to read the journal. A workaround is to add the setuid bit to
  gobble (`chmod u+s /path/to/gobble`), but that has its own problems.


## License

MIT License, Copyright © 2020 Arthur Skowronek, Digineo GmbH
