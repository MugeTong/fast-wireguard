package templates

import (
	_ "embed"
)


var (
	// WireGuard server configuration file template
	//go:embed wg.conf.tpl
	WgConfTpl string
	// WireFuard peer configuration template
	//go:embed peer.conf.tpl
	PeerConfTpl string
	// WireGuard client configuration file template
	//go:embed client.conf.tpl
	ClientConfTpl string
)
