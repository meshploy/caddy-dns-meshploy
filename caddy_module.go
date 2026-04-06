package caddydnsmeshploy

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/meshploy/caddy-dns-meshploy/meshploydns"
)

func init() {
	caddy.RegisterModule(ProviderWrapper{})
}

type ProviderWrapper struct {
	*meshploydns.Provider
}

func (ProviderWrapper) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.meshploy", // This is what you'll type in the Caddyfile
		New: func() caddy.Module { return &ProviderWrapper{Provider: &meshploydns.Provider{}} },
	}
}

func (w *ProviderWrapper) Provision(ctx caddy.Context) error {
	if w.ZoneFilePath == "" {
		w.ZoneFilePath = "/etc/coredns/zones/_acme-challenge.example.com"
	}
	return nil
}

func (w *ProviderWrapper) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "zone_file_path":
				if !d.NextArg() {
					return d.ArgErr()
				}
				w.ZoneFilePath = d.Val()
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	return nil
}

var (
	_ caddy.Provisioner     = (*ProviderWrapper)(nil)
	_ caddyfile.Unmarshaler = (*ProviderWrapper)(nil)
)