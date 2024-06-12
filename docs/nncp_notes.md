# NNCP

NNCP Commands:

Configuration file commands

* nncp-cfgnew
    This needs to be used to setup a new node. Creates a hjson configuration file.

Maintenance, monitoring and debugging commands:

* nncp-stat
    Print current spool statistics

Nearly all commands share some common options:
-debug, -minsize, -nice, -spool, -progress/-noprogress

```markdown
Packets creation commands
• nncp-file:
• nncp-exec:
• nncp-freq:
• nncp-trns:
• nncp-ack:

Packets sharing commands
• nncp-xfer:
• nncp-bundle:

Checking and tossing commands
• nncp-toss:
• nncp-check:
• nncp-reass:


Online synchronization protocol commands
• nncp-daemon:
• nncp-call:
• nncp-caller:
• nncp-cronexpr:

Maintenance, monitoring and debugging commands:
• nncp-log:
• nncp-rm:
• nncp-pkt:
• nncp-hash:
```

## Further Reading

These are some resources I used while working on this.

https://en.wikipedia.org/wiki/Usenet
https://en.wikipedia.org/wiki/News_server

https://www.eyrie.org/~eagle/software/inn/
https://www.eyrie.org/~eagle/software/inn/docs-2.6/newsfeeds.html
https://datatracker.ietf.org/doc/html/rfc3977

https://github.com/dmah42/slurp (newsreader)
https://github.com/dustin/go-nntp/blob/f00d51cf8cc1/examples/couchserver/couchserver.go
https://github.com/InterNetNews/inn

https://pkg.go.dev/github.com/dustin/go-nntp@v0.0.0-20210723005859-f00d51cf8cc1
https://pkg.go.dev/net/textproto#pkg-overview

https://groups.google.com/g/comp.mail.uucp/c/E1iFGMkULiU

[nncp repo on Salsa](https://salsa.debian.org/go-team/packages/nncp)
[NNCP in Docker](https://salsa.debian.org/jgoerzen/docker-nncp)
[Usenet Wikipedia](https://en.wikipedia.org/wiki/Usenet)
[Eternal September](https://www.eternal-september.org/)
[Usenet over NNCP](https://www.complete.org/usenet-over-nncp/)
[Introduction to Usenet](https://www.binaries4all.com/beginners/downloading.php)
[Sabnzbd - Open source newsreader written in Python](https://github.com/sabnzbd/sabnzbd/tree/develop)
