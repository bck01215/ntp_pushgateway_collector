# NTP Pushgateway Check

Pretty simple to use. Use the  --insecure flag if your psuhgateway is not secure.

```
ntpChecker  --insecure --pushgateway https://pushgateway.example.com/  1.ntp.example.com 2.ntp.example.com ...
```
It returns a metric like 
```
ntp_unix_time{instance="1.ntp.example.com", job="ntp"} 1654875566017
ntp_unix_time{instance="2.ntp.example.com", job="ntp"} 1654875566006
```
Where unix time is sensitive to milliseconds. Nanoseconds are negligble based off of different speeds in go routines.

The docker image is `bkauffman7/ntp_pushgateway_collector:v0.0.1`