# NETCONF

[![GoDoc](https://godoc.org/github.com/adetalhouet/go-netconf/netconf?status.svg)](https://godoc.org/github.com/adetalhouet/go-netconf/netconf)
[![Report Card](https://goreportcard.com/badge/github.com/adetalhouet/go-netconf)](https://goreportcard.com/report/github.com/adetalhouet/go-netconf)
[![Build Status](https://travis-ci.org/adetalhouet/go-netconf.png)](https://travis-ci.org/adetalhouet/go-netconf)
[![codecov](https://codecov.io/gh/adetalhouet/go-netconf/branch/main/graph/badge.svg)](https://codecov.io/gh/adetalhouet/go-netconf)

This library is a simple NETCONF client :
- [RFC6241](http://tools.ietf.org/html/rfc6241): **Network Configuration Protocol (NETCONF)** 
    - Support for the following RPC: `lock`, `unlock`, `edit-config`, `comit`, `validate`,`get`, `get-config`
    - Support for custom RPC
- [RFC6242](http://tools.ietf.org/html/rfc6242): **Using the NETCONF Protocol over Secure Shell (SSH)**
    - Support for username/password
    - Support for pub key
- [RFC5277](https://datatracker.ietf.org/doc/html/rfc5277): **NETCONF Event Notifications**
    - Support for `create-subscription`
    - No support for notification filtering
- Partially [RFC8641](https://datatracker.ietf.org/doc/html/rfc8641) and [RFC8639](https://datatracker.ietf.org/doc/html/rfc8639): **Subscription to YANG Notifications for Datastore Updates**
    - Support for `establish-subscription`

#### Links
This client is an adaptation of the code taken from:
- https://github.com/andaru/netconf
- https://github.com/Juniper/go-netconf