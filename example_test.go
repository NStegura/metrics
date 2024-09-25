package metrics

import (
	"context"
	"fmt"
	"github.com/NStegura/metrics/config"
	"github.com/NStegura/metrics/internal/app/metricsapi/httpserver"
	"github.com/NStegura/metrics/internal/business"
	"github.com/NStegura/metrics/internal/repo"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

const (
	addr = "http://127.0.0.1:8080"
)

func startServer(l *logrus.Logger) {
	r, err := repo.New(
		context.Background(),
		"",
		100,
		"",
		false,
		l)
	if err != nil {
		l.Fatal(err)
	}
	businessLayer := business.New(r, l)
	sConfig := config.NewSrvConfig()

	server, err := httpserver.New(sConfig, businessLayer, l)
	if err != nil {
		l.Fatal(err)
	}
	go func() {
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(time.Second * 2)
}

func Example_updateAndGetGaugeMetric() {
	l := logrus.New()
	startServer(l)

	gaugeValue := 123.1231
	url := fmt.Sprintf(
		"%s/update/gauge/test/%v",
		addr, gaugeValue)

	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = url

	_, err := req.Send()
	if err != nil {
		l.Fatal(err)
	}

	url = fmt.Sprintf(
		"%s/value/gauge/test",
		"http://127.0.0.1:8080")

	req = resty.New().R()
	req.Method = http.MethodGet
	req.URL = url

	res, err := req.Send()
	if err != nil {
		l.Fatal(err)
	}

	fmt.Println(res)

	// Output:
	// 123.1231
}
