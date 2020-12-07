# mock-ntpv5-client
This is an experimental NTPV5 mock client for survey

## Introduction
The IETF is beginning to standardize NTPv5.

- [Network Time Protocol Version 5](https://tools.ietf.org/html/draft-mlichvar-ntp-ntpv5-00)

I am interested in NTP ossification. So I survey how the currently deployed NTP server responds to NTPv5 packets.

## How to Use
Please enter the destination host name in the command line argument

```
go run ./mock-client.go SUB-COMMAND HOSTNAME
```

This CLI has subcommands
- `v4`: send nomal NTPv4 packet
- `v4-ue`: Not yet implemented
- `v4-5`: send NTPv4 format packet that version field is specified with 5
- `v5`: send NTPv5 format packet (draft-mlichvar-ntp-ntpv5-00)

Outputs the version field of the received NTP packet

Example:
```
$ go run ./mock-client.go v4-5 pool.ntp.org
pool.ntp.org response version: 5
```

## What these tool showed
In the first step survey, I sent packets to public NTP servers in the world.

these tools showed following results.

### v4-5: 
- 25% response: timeout
- 65% response: NTPv4 format packet that version field is specified with 5
- 10% response: NTPv4 or NTPv3

It shows that many servers are processing NTPv5 packets

### v5
- 10% response: NTPv4 format packet that version field is specified with 5
- 90% response: timeout

The result depends on the date of Transmit Timestamp and the absence of Extension Field.

In this test case we have added a dummy extension field.

## Is this a problem?
I'm not sure.

I know we suffered from anomalous behavior in tcp-fast-open and TLS 1.3 deployments
- https://mailarchive.ietf.org/arch/msg/tls/i9blmvG2BEPf1s1OJkenHknRw9c/
- https://archive.nanog.org/sites/default/files/Paasch_Network_Support.pdf

but, it may be not matter
- This response is not a problem for the client, Because it can be filtered
- It will be fixed as NTPv5 develops


