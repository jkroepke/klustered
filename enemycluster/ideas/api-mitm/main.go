package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type conf struct {
	Users []user `yaml:"users"`
}

type user struct {
	Name string          `yaml:"name"`
	User userCertificate `yaml:"user"`
}

type userCertificate struct {
	Cert string `yaml:"client-certificate-data"`
	Key  string `yaml:"client-key-data"`
}

func (c *conf) getConf() *conf {
	adminFile := os.Getenv("ADMIN_CONF")
	if len(adminFile) == 0 {
		adminFile = "/etc/kubernetes/admin.conf"
	}

	yamlFile, err := ioutil.ReadFile(adminFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {
	certDir := os.Getenv("CERT_DIR")
	if len(certDir) == 0 {
		certDir = "/etc/kubernetes/pki/"
	}

	var adminFile conf
	adminFile.getConf()

	// adminFile.Users[0].User.Cert
	// adminFile.Users[0].User.Key

	cert, err := base64.StdEncoding.DecodeString(adminFile.Users[0].User.Cert)
	if err != nil {
		log.Fatal("DecodeString:cert", err)
	}

	key, err := base64.StdEncoding.DecodeString(adminFile.Users[0].User.Key)
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

		if strings.HasPrefix(req.Header.Get("User-Agent"), "kubectl/") {
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

		fmt.Printf("[reverse proxy server] [%s] %s %s -> %s\n", req.Header.Get("User-Agent"), req.Method, req.RequestURI, req.URL.RequestURI())

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

	log.Fatal(http.ListenAndServeTLS(":2222", certDir+"apiserver.crt", certDir+"apiserver.key", reverseProxy))
}
