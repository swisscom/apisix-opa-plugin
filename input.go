package main

import "net"

type OpaInputRequest struct {
	Host    string              `json:"host"`
	Path    string              `json:"path"`
	SrcIp   net.IP              `json:"src_ip"`
	Headers map[string][]string `json:"headers"`
}

type OpaInput struct {
	Input OpaInputRequest `json:"input"`
}
