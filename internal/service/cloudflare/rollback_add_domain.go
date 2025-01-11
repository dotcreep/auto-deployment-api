package cloudflare

import (
	"context"
	"log"
)

func (c *Cloudflare) RollbackAddDomain(ctx context.Context, data *Subdomains) {
	if data.Domain == "" {
		log.Println("domain is required")
	}
	var tunnelInfo bool = true
	var dnsRecordInfo bool = true
	_, err := c.DeleteDomainFromTunnelConfiguration(ctx, data)
	if err != nil {
		log.Println(err)
		tunnelInfo = false
	}
	_, err = c.DeleteDomainDNSRecords(ctx, data)
	if err != nil {
		log.Println(err)
		dnsRecordInfo = false
	}
	if tunnelInfo && dnsRecordInfo {
		log.Println("rollback domain success")
	} else {
		log.Printf("rollback domain failed, tunnel status success: %v, dns record status success: %v\n", tunnelInfo, dnsRecordInfo)
	}
}
