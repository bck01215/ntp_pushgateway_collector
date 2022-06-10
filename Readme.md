# NTP Pushgateway Check
### Use
Pretty simple to use. Use the  --insecure flag if your psuhgateway is not secure.

```
ntpChecker  --insecure --pushgateway https://pushgateway.example.com/  1.ntp.example.com 2.ntp.example.com ...
```
### Metrics
It returns metrics like:
```
ntp_unix_time{instance="1.ntp.example.com", job="ntp"} 1654875566017
ntp_unix_time{instance="2.ntp.example.com", job="ntp"} 1654875566006
ntp_unix_time{instance="2.ntp.example.com", job="ntp"} -1
```
Where unix time is sensitive to milliseconds. Nanoseconds are negligble based off of different speeds in go routines.
If the value is `-1` the scrape failed.

### Docker
The docker image is `bkauffman7/ntp_pushgateway_collector:v0.0.1`
