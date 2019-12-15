# ProxiesStack
[Docker Hub Repository](https://hub.docker.com/repository/docker/maxsid/proxies-stack)

This is a service grabbing a website page, looking for 
addresses of proxy servers and checking their status.

## Requirements

Redis require for data storage. The Redis host must be set in the REDIS_HOST 
environment variable.

## Environment

You can set next environment variables:
- **REDIS_HOST** - (require) Host and port of the Redis Server. For example: "192.168.0.1:6379"
- **TIMEOUT** - (10 sec) Timeout of HTTP requests in seconds. After expiration of the timeout 
a proxy server will be marked as not working.
- **NO_SCAN** - (False) Specify should avoid scanning and running only as API server
- **SCAN_INTERVAL** - (5 min) The Interval between the scanning
- **GRAB_ADDRESS** - An address of the page contains hosts of proxy servers
- **GRAB_PATTERN** - Regex pattern the host on the page. Must contain next groups: *host*, *port* and *https*.
- **SCAN_ADDRESS** - An address of the page which have to match The Regex Pattern in *SCAN_PATTERN*
- **SCAN_PATTERN** - Regex pattern. If a match is exists proxy host will mark as working

## API

At the moment API has next requests:
- **/info** - Return information about scanning status, amount proxy hosts and etc.
- **/working/pop** - Return one proxy host and mark it as not working until next scanning, if it will be working. 