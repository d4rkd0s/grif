<img src="assets/grif.png" width="100" height="100">

# Grif

Grif keeps you on top of the machines you are responsible for. 
Without having to setup an entire monitoring system, or bother the ops team to monitor new hosts that may be temporary.

## Installation
- Goto the Releases page
- Download .zip, and extract.
- Edit hosts file (paste URLs to check, currently only http:// and https:// are available).
- Grif will remake this file if its lost, as a demo https://github.com/ is placed in hosts by default on the downloaded version
- Run grif.exe and look in your tray (Grif will bark when it's ready, and when it finds a host with problems)

## How it works

Grif reads protocol URIs from a `hosts` textfile located in the same directory as the grif binary/executable.
Grif digests this list at a specified interval so it can be updated while Grif is running.
If Grif detects an outage, it will alert the user.

## Easy to use format

The `hosts` file located in the grif home directory should contain a list of URLs you wish to monitor

```
https://google.com/
https://twitter.com/
http://something.com/
```

Grif will check the hosts for a valid http response from the hosts with HTTP or HTTPS. Grif can also validate SSL if desired (this feature is being tested).

## Building Grif

In the main directory of this project get all of the dependencies
```
go get -u github.com/faiface/beep
go get -u	github.com/faiface/beep/mp3
go get -u	github.com/faiface/beep/speaker
go get -u	github.com/gen2brain/beeep
go get -u	github.com/getlantern/systray
go get -u	github.com/sparrc/go-ping
```

To run a build that moves to system tray run:
```
go build -ldflags -H=windowsgui -o build/grif.exe
```

To run a build that stays in command prompt, for debugging purposes run:
```
go build -o build/grif.exe
```

The files in the repo are already generated and in place to ensure the exe is built fully, with icon, manifest, etc.
