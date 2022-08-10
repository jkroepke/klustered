//go:build ignore
// +build ignore

//go:generate strobfus -filename $GOFILE

package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var KUBERNETES_CLIENT_CERT = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURJVENDQWdtZ0F3SUJBZ0lJZE05U1BhbWg3YTR3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TWpBNE1UQXhOREkzTWpSYUZ3MHlNekE0TVRBeE5ESTNNalphTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXJZcmZFcGpaTEsyUldEbUMKS3AraWdKTVhSUG1uMWdhdmVFbmhRTmt4cjByRFhULzIvMFJUcFJLNzRkT2FickoxWmszZHFRamRCYjdGUkVqNQpnNU5IOTlzUjNOSStjZy80YjE3clcxeVdmVUJheC9zcWhWeXFRWmdVQUphTEkyNjFTNXJ2YWtiaEptS3JpRGVRCmtzTlk4Y0MyeExIU3NFa1RCbU1KZzNCRDR5eFZVc05sRXhENjgyR294MkdxVHY0R0lEV1FZUzU5dXdiTk5mR1YKWllIWCsvTG9sK3JMZ0RFcTF5RUgxOTVtUWY3Y3hEZ2JWVGhrZHRGMEIvZmdja0tJRkxLTkFLVHhiclBhWmp4dwpsbDZhbWU5MzJjVUJuYlFCQUZaOTRUUVJEbEtQZzFjdWNuV2c2Zk5VZ25CTUxPS0dhcEdHZGhnMUNXYjVRNGhSCkkza0l0UUlEQVFBQm8xWXdWREFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RBWURWUjBUQVFIL0JBSXdBREFmQmdOVkhTTUVHREFXZ0JUS3YvRUZRelNINmxHMEkxSTlEVnpVTW5EeApNVEFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBQnp3K24xSndBQlVpMEY4QVZsMUZZQmQwWlpQZTZGZ2JUUVZhCnhWOC81T0RMQ3hlZnZXdGlVY2Z6eTh5UmdhUGp2VWhiMnlNK0d2bkRqZGJ2MzV3VC9EUHZZVGtzK2JoTWtxNW4KS2dVMER4ZGM1V3NZRE1PUnVMRUhLck5HQ2trUG1pSXBPMG1Kd3p5c2kxMWY0SDd1VkJCMjUvYlhjUngwaWVHZgpONUdRTE5zL1hSdUV2czhTWFFnNlZZMkxwaDNZS012eHJ2dXllS0ljcnhwSnB2RityMjA4ME1IcHFQU2wxTUpjCk9iVm1IVUtzUTlDZFFycUlwVkIxMDJNWWhKRTNVZ3A1M0FMUlJZcnVHYW9qQ2s3Kzd4TzVuMVczdHpLRXBYTFEKK0tiV1kyekVUOVJnbmc0ZEFKdXZyYUYxbys4N2Q1SkRUUG5jK25lNDZBa2hFSFY5TXc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
var KUBERNETES_CLIENT_KEY = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb2dJQkFBS0NBUUVBcllyZkVwalpMSzJSV0RtQ0twK2lnSk1YUlBtbjFnYXZlRW5oUU5reHIwckRYVC8yCi8wUlRwUks3NGRPYWJySjFaazNkcVFqZEJiN0ZSRWo1ZzVOSDk5c1IzTkkrY2cvNGIxN3JXMXlXZlVCYXgvc3EKaFZ5cVFaZ1VBSmFMSTI2MVM1cnZha2JoSm1LcmlEZVFrc05ZOGNDMnhMSFNzRWtUQm1NSmczQkQ0eXhWVXNObApFeEQ2ODJHb3gyR3FUdjRHSURXUVlTNTl1d2JOTmZHVlpZSFgrL0xvbCtyTGdERXExeUVIMTk1bVFmN2N4RGdiClZUaGtkdEYwQi9mZ2NrS0lGTEtOQUtUeGJyUGFaanh3bGw2YW1lOTMyY1VCbmJRQkFGWjk0VFFSRGxLUGcxY3UKY25XZzZmTlVnbkJNTE9LR2FwR0dkaGcxQ1diNVE0aFJJM2tJdFFJREFRQUJBb0lCQUNQaXc1NGszVVBQNEc1Tgo5Z3k2VmZBZ2VuOVk0TXZ4TmZlNXowcUpueXlRV1RXL05HUTB6TmNsdUpSS0hYVW1rZ0JGdWNCcWhNbmJXUTkxCng2TGRvZFF2Q05LUTV6ak85S0NURURna1BUcEpSSHgyQTZUd05JUzczZWNCT21ScFVEUUNKZC9rS0VxM3ZLQysKWExiOGpqZnZrZHU2cWNhcVZiVE1aZnM0QzlHOTBjQ3h1ZHZpM0h0dnlZekFSa0xOSDlva3BtalovNmo3cDBBUQpkODhPcUh2R3M0N3djV1NTdWJZVjB1QXRTV3lpOWRSdis2OEltckY0MDMrVlB1SUNTaWdZL29BY3U1OVZJK1RLClpBdWgvSkdGUGtpK2p2RmxhVWdoYTZhbTVhano2UlZtQTJPRFFFaytJc2l1ZzRtVkpidXFqQUdoUldadzhiZHUKakpLU0NBRUNnWUVBeEJVeU1ic3I2dTdhRmZaUlUxWWx0VWpyQjBvK1h1UHVNZm9LVlI5QVBQbG9QNzJ0bVhrMwp6R3FpS3BhdktMQmUwRGRIS3ZFL0tTSnpyN1hBV3o2YUhUVzBGUmVTN2ZYc1lHTnlWNFV6MmNRbmgrVVNnMDJ2Cjlad2tRSFFUOXQwS0dKcmdUZkp0TVQyVWhYNUZXaGpLOXQzRU9yaTVheEFLVlZoczg1bjFrdFVDZ1lFQTRwSnUKaTZTNS9HZHUvMk1wK1RmektsR2RWdmJrZVhUMG9wR1lONFN5emVSVkxOT0RuSi9VMUc4NWV6Ym51cy8wb1dxdAppVmdhNWpWNGtwZ3NFR1ljUDhHZDdaeit3Q1ZZZmc0QlBkT3ltNWZKYzhXYllMR08zK293cUVTVG9QN3VEaWZ6CjF2TmRTM2hxRVJRcGdYMjJqYmlmTmIrSmt4OVRIZDNOblB5eHptRUNnWUJ5Ykh3U0VVdWJtUzZpeWs3Qzl1NmkKVDU3M2JoZmZmOXNzUnVGb3N3ZmxqUldNdkw5bFpCdHZxbnBmcC9jbkkyVHcxSkV2T2dERm5Ga3VIRDNZQVR3bAo5NFRUR2lLZnduYmgrS1pzOUVwQnRmbnJqMzJ5S2MrWTREazNjNFdDOVpKQ3NYNWJmakRDSDFGZ1pVTkxSRlNNCm92VXozMEEwZmZQSndnUXlVNUcrMFFLQmdCVE9JWHlOT2M3bHFKbW0vM20xRzQwdFJXZHc4SFgrdVdBY1FvQUcKbld5dXBPdWkySmtQVERuZHBNZWR1Ulc4ZHRoRHRYL0JLV2N1VGM0WVR5T0tYTm0xNjh5ZjkveW84VUZTQStjcgpnMkVxUlFOdWgrQVBMZkY5emM0RnpoQ2dtRGVRajZHVFkyUEV3T2lrazFNaXVocTFjMWs2SjJYdElITERwVmZmCkZHekJBb0dBT3AyenRNVE1GTi9iQWVuaVJJUWZUbWZDTk1odDFOcHY3d0FSdjJiYkc4OHF3V2FvdG1VZ20yaDAKSjR0aTlTT09NeXFCRENEcllCdy9FNTJDbnJQSForc1hFaEpUUHZkWXJBVGxvckVIM1IwTHVwU2EyVUNiRWIxegpJRFdkT1NPWldBUmdYYlBCN0xLVGU4R0dzS1BDL2VyN20xNFJXbXdWVEhGZmxNV3VuZ3M9Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg=="

func main() {
	debug := os.Getenv("DEBUG")

	userAgent := os.Getenv("INTERCEPT_USERAGENT")
	if len(userAgent) == 0 {
		userAgent = "kubectl/"
	}

	cert, err := base64.StdEncoding.DecodeString(KUBERNETES_CLIENT_CERT)
	if err != nil {
		log.Fatal("DecodeString:cert", err)
	}

	key, err := base64.StdEncoding.DecodeString(KUBERNETES_CLIENT_KEY)
	if err != nil {
		log.Fatal("DecodeString:key", err)
	}

	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatal("X509KeyPair", err)
	}

	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		req.Host = "127.0.0.53"
		req.URL.Host = "127.0.0.53:6443"
		req.URL.Scheme = "https"
		req.Header.Set("Accept-Encoding", "")

		if strings.HasPrefix(req.Header.Get("User-Agent"), userAgent) {
			if req.Method == "DELETE" {
				postBody := "{\"propagationPolicy\":\"Background\",\"dryRun\":[\"All\"]}\n"
				req.Body = ioutil.NopCloser(strings.NewReader(postBody))
				req.ContentLength = int64(len(postBody))
			} else {
				req.URL.Query().Set("dryRun", "All")
				req.URL.RawQuery = req.URL.RawQuery + "&dryRun=All"

				req.URL.Query().Del("watch")
				req.URL.RawQuery = strings.Replace(req.URL.RawQuery, "watch=true", "", 1)

				req.URL.RawQuery = strings.Replace(req.URL.RawQuery, "&&", "&", 1)
			}
		}

		if debug == "1" {
			fmt.Printf("[reverse proxy server] [%s] %s %s -> %s\n", req.Header.Get("User-Agent"), req.Method, req.RequestURI, req.URL.RequestURI())
		}

		if strings.Contains(req.RequestURI, "watch=true") {
			rw.WriteHeader(http.StatusOK)
			return
		}

		req.RequestURI = ""

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true, Certificates: []tls.Certificate{tlsCert}}

		// send a request to the origin server
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(rw, err)
			return
		}

		if strings.HasPrefix(req.Header.Get("Accept"), "application/json") {
			rw.Header().Set("Content-Type", "application/json")
		} else {
			rw.Header().Set("Content-Type", "application/octet-stream")
		}

		defer resp.Body.Close()
		reader := bufio.NewReader(resp.Body)

		if _, err := io.Copy(rw, reader); err != nil {
			log.Print("copy error", err)
		}
	})

	log.Fatal(http.ListenAndServeTLS(":853", "/etc/kubernetes/pki/apiserver.crt", "/etc/kubernetes/pki/apiserver.key", reverseProxy))
}
