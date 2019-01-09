<img src="assets/grif.png" width="100" height="100">

# Grif

Grif keeps you on top of the machines you are responsible for. 
Without having to setup an entire monitoring system, or bother the ops team to monitor new hosts that may be temporary.

## How it works

Grif reads protocol URIs from a `hosts` textfile located in the same directory as the grif binary/executable.
Grif digests this list at a specified interval so it can be updated while Grif is running.
If Grif detects an outage, it will alert the user.

## Easy to use format

The `hosts` file located in the grif home directory could contain this

```
icmp://8.8.8.8
icmp://8.8.4.4
https://google.com/
https://twitter.com/
```

and Grif will Ping the hosts that begin with icmp:// and test for a valid http response from the hosts with HTTP or HTTPS. Grif can also validate SSL if desired (this feature is being tested).
