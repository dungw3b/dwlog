# dwlog
Log service developed in Go, it is very fast in sending message log since used channel with non-blocking sender. It also uses gRPC as communication protocol

Build:
---------------
```bash
    $ cd example
    $ ./build.sh
```

Run:
---------------
```bash
    $ cd example
    $ mkdir example/data
    $ ./bin/dwlog-linux-amd64 -c conf/dwlog.json
```