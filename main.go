package main

import (
	"fmt"
	"net"
	"strings"
)

type Result struct {
	FlagUp   string // [ up | down ]
	IFName   string // eth0
	Mac      string // A:B:C:D:E:F
	IpNet    string //  10.209.0.75/22
	Vlan     string // Ethernet39@sNJL364m
	Hostname string //  holmium
}

func (r *Result) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", r.FlagUp, r.IFName, r.Mac, r.IpNet, r.Vlan, r.Hostname)
}

func getStatus(flags net.Flags) (status string) {
	for _, flag := range strings.Split(flags.String(), "|") {
		if flag == "up" || flag == "down" {
			status = flag
		}
	}

	return
}

func resolveIp(ip string) (host string) {
	if resolvedHost, _ := net.LookupAddr(ip); len(resolvedHost) > 0 {
		host = resolvedHost[0]
	}

	return
}

func getCidrAndHost(netAddrs []net.Addr) string {
	for _, addr := range netAddrs {
		if _, ok := addr.(*net.IPNet); ok {
			if _, ipNet, err := net.ParseCIDR(addr.String()); err == nil {
				if ipAddr := ipNet.IP.To4(); ipAddr != nil {
					return ipAddr.String()
				}
			}
		}
	}

	return
}

func getAllInterfaces() (results []*Result) {
	interfaces, _ := net.Interfaces()

	for _, inter := range interfaces {
		myIF := &Result{FlagUp: getStatus(inter.Flags), IFName: inter.Name, Vlan: "?"}

		switch {
		case strings.Contains(inter.Flags.String(), "loopback"):
			myIF.Mac = "00:00:00:00:00:00"
			myIF.IpNet = "127.0.0.1/8"
		case inter.HardwareAddr != nil, myIF.FlagUp != "":
			myIF.Mac = fmt.Sprintf("%s", inter.HardwareAddr)

			if addrs, err := inter.Addrs(); err == nil {
				myIF.IpNet = getCidrAndHost(addrs)
			}
			// get dns names pointing to this host
			myIF.Hostname = resolveIp(strings.Split(myIF.IpNet, "/")[0])
		default:
			continue
		}

		results = append(results, myIF)
	}

	return
}

func main() {
	//w := tabwriter.NewWriter(os.Stdout, 0, 0,  4, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for _, result := range getAllInterfaces() {
		fmt.Printf("%s\n", result)
		//fmt.Fprintf(w, "a\tb\taligned\t")
	}

	//w.Flush()
}
