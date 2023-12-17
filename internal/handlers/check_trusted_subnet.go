package handlers

import (
	"fmt"
	"net/http"
	"net/netip"
)

// CheckTrustedSubnet - middleware для проверки подсети.
func CheckTrustedSubnet(subnet *netip.Prefix) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Real IP:", r.RemoteAddr)
			ip, err := netip.ParseAddr(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			if !subnet.Contains(ip) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
