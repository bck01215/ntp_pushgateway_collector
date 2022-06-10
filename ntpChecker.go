package main

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/beevik/ntp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	address  = kingpin.Flag("pushgateway", "The url to push the metrics to").Required().String()
	insecure = kingpin.Flag("insecure", "Add flag if the url should skip verification").Bool()
	servers  = kingpin.Arg("servers", "Servers to check").Strings()
)

type Result struct {
	Time     int64
	Instance string
}

func getTime(servers []string) []Result {
	var wg sync.WaitGroup
	waitTil := make(chan struct{})
	Results := make([]Result, len(servers))
	serverLen := len(servers)

	wg.Add(serverLen)
	for i := 0; i < serverLen; i++ {
		go func(i int) {
			var result Result
			<-waitTil
			instance := servers[i]
			defer wg.Done()
			time, err := ntp.Time(instance)
			if err != nil {
				Results[i] = Result{
					Time:     -1,
					Instance: instance,
				}
				logrus.Error(err)
			} else {
				result.Time = time.UnixMilli()
				result.Instance = instance
				Results[i] = result
			}
		}(i)
	}
	logrus.Info("Scraped at ", time.Now().UnixMilli())
	close(waitTil)
	wg.Wait()
	return Results

}

func pushToGateway(instance string, values []Result, insecure bool) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}
	ntpTime := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ntp_unix_time",
		Help: "Unix time to milliseconds of NTP server.",
	}, []string{"instance"})

	for i := 0; i < len(values); i++ {
		ntpTime.WithLabelValues(values[i].Instance).Set(float64(values[i].Time))
	}
	if err := push.New(instance, "ntp").
		Collector(ntpTime).
		Push(); err != nil {
		logrus.Fatal("Could not push completion time to Pushgateway: ", err)
	}
}
func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	logrus.SetLevel(logrus.DebugLevel)
	results := getTime(*servers)
	pushToGateway(*address, results, *insecure)

}
