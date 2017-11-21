
tips : must run govs as root

# add a service
```
govs -A -t|u IP:PORT [-s scheduler] [-M netmask] [-flags service-flags]
```
-t : tcp

-u : udp

-s : scheduler(rr,wrr)

-M : netmask

-flags : service-flags

> PORT is optional, 0 by default
> scheduler : rr - Round-Robin, wrr - Weighted Round-Robin


# edit a service
```
govs -E -t|u IP:PORT [-s scheduler] [-M netmask] [-flags service-flags]
```
-t : tcp

-u : udp

-s : scheduler(rr,wrr)

-M : netmask

-flags : service-flags


# delete a service
```
govs -D -t|u IP:PORT
```
-t : tcp

-u : udp


# clear the whole table
```
govs -C
```


# add a real-server for a service
```
govs -a -t|u IP:PORT -r IP:PORT [-w weight] [-x upper-threshold] [-y lower-threshold] [-conn_flags conn-flags]
```
-t : tcp

-u : udp

-r : real-server

-w : capacity of real server

-x : upper threshold of connections

-y : lower threshold of connections

-conn_flags : conn flags(Local/Tunnel/Route/FullNat)


# edit a real-server
```
govs -e -t|u IP:PORT -r IP:PORT [-w weight] [-x upper-threshold] [-y lower-threshold] [-conn_flags conn-flags]
```
-t : tcp

-u : udp

-r : real-server

-w : capacity of real server

-x : upper threshold of connections

-y : lower threshold of connections

-conn_flags : conn flags(Local/Tunnel/Route/FullNat)


# delete a real-server
```
govs -d -t|u IP:PORT -r IP:PORT
```
-t : tcp

-u : udp

-r : real-server


# list services and real-servers
```
govs -L|l [-t|u IP:PORT] [-detail] [-i id] [-all] [-n coefficient]
```
-t : tcp

-u : udp

-detail : print detail information

-i : id of cpu workers

-all : show all of the cpu workers

-n : multiply by a coefficient

> show No.0 cpu worker by default


# add a local-address for a service
```
govs -P -t|u IP:PORT -z IP:PORT
```
-t : tcp

-u : udp

-z : local-address


# delete a local-address
```
govs -Q -t|u IP:PORT -z IP:PORT
```
-t : tcp

-u : udp

-z : local-address


# get local-address
```
govs -G [-t|u IP:PORT]
```
-t : tcp

-u : udp


# get version
```
govs -V
```
> get version of dpvs


# show usage message
```
govs -h
```


# zero counters in a service or all services
```
govs -Z [-t|u IP:PORT]
```
-t : tcp

-u : udp


# get/set connection timeout values
```
-TAG_SET [-set tcp/tcp_fin/udp]
```
-set : conntction timeout

> show timeout values when tcp/tcp_fin/udp is null


# get dpvs status
```
-s [-type stats-name] [-i id] [-all] [-n coefficient]
```
-type : status type(io/w/we/dev/ctl/mem/vs/falcon)

-i : id of cpu workers

-all : show all of the cpu workers

-n : multiply by a coefficient

