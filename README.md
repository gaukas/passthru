# Passthru
Passthru, a Protocol Omni-multiplexer

## Demo: Usage

### Build

```bash
$ go build ./cmd/passthru
```

### Run

```bash
$ ./passthru -c=<configfile> -w=<worker_num> -t=<timeout>
```

Note that in the current release, the `worker_num` and `timeout` are not used for simplicity in demo. You may manually enable it in the `main()` function.

#### Config

```json
{
    "version": "v0.2.0",
    "servers": {
        "127.0.0.1:443": {
            "TLS": {
                "SNI gaukas.wang": {
                    "action": "FORWARD",
                    "to_addr": "185.199.111.153:443"
                },
                "SNI google.com": {
                    "action": "FORWARD",
                    "to_addr": "142.250.72.46:443"
                },
                "SNI www.google.com": {
                    "action": "FORWARD",
                    "to_addr": "142.250.72.46:443"
                },
                "CATCHALL": {
                    "action": "REJECT"
                }
            },
            "CATCHALL": {
                "CATCHALL": {
                    "action": "FORWARD",
                    "to_addr": "neverssl.com:80"
                }
            }
        }
    }
}
```

## Packages

### Config

The config package is used to parse the config file. The config file follows the structure below: 

- Version
- Servers
    - ServerAddr1
        - Protocol1: defined in `protocol` package
            - Rule1
                - Action: either `FORWARD` or `REJECT`
                - ToAddr (`FORWARD` only): the address to forward to
            - Rule2
                - ...
            - CATCHALL (as a rule)
                - Action
                - ToAddr (optional)
        - Protocol2: defined in `protocol` package
            - Rule1
                - ...
            - ...
        - CATCHALL (as a protocol)
            - CATCHELL (as the only rule of this protocol)
                - ...
    - ServerAddr2
        - ...

### Handler

Handler defines the handler of all incoming connections to a certain address as a `Server`. 

When `Server` receives a connection, it will first read the first few bytes of the connection to determine the protocol. Then it will pass the connection to the corresponding protocol handler.

### Protocol

Protocol defines the identifier/detector of a certain `Protocol`. It also defines `ProtocolManager` to manage all known protocols. Upon receiving a buffer of bytes, `ProtocolManager` will try to match the buffer with all known protocols. If a match is found, the matching `Action` will be returned.