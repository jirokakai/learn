package bench

import (
	"flag"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/gyuho/learn/doc/go_network/jsonrpc_vs_grpc/demogrpc"
	"github.com/gyuho/learn/doc/go_network/jsonrpc_vs_grpc/demojsonrpc"
)

var (
	totalConns   = 1
	totalClients = 1
	// totalClients = 100

	size = 100000
	opt  = "grpc"

	keys = make([][]byte, size)
	vals = make([][]byte, size)
)

func init() {
	sizePt := flag.Int(
		"size",
		100000,
		"Size of keys to put",
	)
	optPt := flag.String(
		"opt",
		"grpc",
		"'grpc' or 'jsonrpc'",
	)
	flag.Parse()

	size = *sizePt
	opt = *optPt
	if opt != "grpc" && opt != "jsonrpc" {
		log.Fatalf("%s is unknown\n", opt)
	}
	log.Println("Size chosen:", size)
	log.Println("Option chosen:", opt)

	keys = make([][]byte, size)
	vals = make([][]byte, size)
	for i := range keys {
		keys[i] = randBytes(100)
		vals[i] = randBytes(100)
	}
	log.Println("Done with generating random data...")
}

var once sync.Once

func BenchmarkRun(b *testing.B) {
	b.StartTimer()
	b.ReportAllocs()

	oncebody := func() {
		switch opt {
		case "grpc":
			port := ":3500"
			endpoint := "localhost" + port
			demogrpc.Run(port, endpoint, keys, vals, totalConns, totalClients)
		case "jsonrpc":
			port := ":3501"
			endpoint := "localhost" + port
			demojsonrpc.Run(port, endpoint, keys, vals)
		}
	}
	once.Do(oncebody)
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randBytes(n int) []byte {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}