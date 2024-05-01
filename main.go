package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/netip"

	"github.com/behnamgolds/clf-ips-loop-test/models"
	// could use the following package instead of net/netip,
	// it has checks for network/broadcast address and much more
	// github.com/seancfoley/ipaddress-go/ipaddr
)

const clfIpApiUrl = "https://api.cloudflare.com/client/v4/ips"

func main() {
	resp, err := http.Get(clfIpApiUrl)
	logFatal(err)

	bs, _ := io.ReadAll(resp.Body)
	apiRsp := models.ApiResponse{}
	logFatal(json.Unmarshal(bs, &apiRsp))

	if apiRsp.Success {
		// loop through the first received prefix Ipv4Cidrs[0]
		pfx, err := netip.ParsePrefix(apiRsp.Result.Ipv4Cidrs[0])
		logFatal(err)

		addr := pfx.Addr() // this gets the first address in range aka the netowork id
		for {
			addr = addr.Next() // since we skip the first address, the isNetworkID()
			//is irrelevant here and should be removed
			if pfx.Contains(addr) && !isNetworkID(pfx, addr) && !isBroadcast(pfx, addr) {
				fmt.Printf("%v/%v\n", addr, pfx.Bits())
			} else {
				break
			}
		}

	} else {
		fmt.Println(apiRsp.Errors)
		logFatal(errors.New("something went wrong"))
	}
}

// func logWarn(err error) {
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// }

func logFatal(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// isNetworkID returns true if it is the first ip in the range aka the network ID,
// so it is not routable on cloudflare's network
func isNetworkID(p netip.Prefix, a netip.Addr) bool {
	return p.Addr() == a
}

// isBroadcast returns true if it is the last ip in range aka the broadcast address,
// so it is not routable on cloudflare's network
func isBroadcast(p netip.Prefix, a netip.Addr) bool {
	return !p.Contains(a.Next())
}
