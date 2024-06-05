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
