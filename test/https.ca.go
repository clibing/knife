package main

import (
    "crypto/tls"
    "fmt"
    "log"
)

func main() {
    conf := &tls.Config{
        InsecureSkipVerify: true,
    }

    conn, err := tls.Dial("tcp", "www.google.com:443", conf)
    if err != nil {
        log.Println("Error in Dial", err)
        return
    }
    defer conn.Close()
    certs := conn.ConnectionState().PeerCertificates
    for _, cert := range certs {
        fmt.Printf("Issuer Name: %s\n", cert.Issuer)
        fmt.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-January-02"))
        fmt.Printf("Common Name: %s \n", cert.Issuer.CommonName)

    }
}
