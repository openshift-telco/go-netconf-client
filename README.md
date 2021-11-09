# NETCONF

[![GoDoc](https://godoc.org/github.com/adetalhouet/go-netconf/netconf?status.svg)](https://godoc.org/github.com/adetalhouet/go-netconf/netconf)
[![Report Card](https://goreportcard.com/badge/github.com/adetalhouet/go-netconf/netconf)](https://goreportcard.com/report/github.com/adetalhouet/go-netconf/netconf)
[![Build Status](https://travis-ci.org/adetalhouet/go-netconf.png)](https://travis-ci.org/adetalhouet/go-netconf)

This library is a simple NETCONF client based on :
- [RFC6241](http://tools.ietf.org/html/rfc6241) Network Configuration Protocol (NETCONF) 
- [RFC6242](http://tools.ietf.org/html/rfc6242) Using the NETCONF Protocol over Secure Shell (SSH)
- [RFC5277](https://datatracker.ietf.org/doc/html/rfc5277) NETCONF Event Notifications

## Features
* Support for SSH transport using go.crypto/ssh. (Other transports are planned).
* Built in RPC support (in progress).
* Support for custom RPCs.
* Independent of XML library.  Free to choose encoding/xml or another third party library to parse the results.

#### Links
This client is an adaptation of the code taken from:
- https://github.com/andaru/netconf
- https://github.com/Juniper/go-netconf