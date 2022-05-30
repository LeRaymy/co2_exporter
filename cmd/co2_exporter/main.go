package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type XMLData struct {
	XMLName                xml.Name `xml:"liste"`
	Text                   string   `xml:",chardata"`
	DateActuelle           string   `xml:"date_actuelle"`
	DateDebut              string   `xml:"date_debut"`
	DateFin                string   `xml:"date_fin"`
	DateConsolidee         string   `xml:"date_consolidee"`
	DateDefinitive         string   `xml:"date_definitive"`
	DateMinimaleCalendrier string   `xml:"date_minimale_calendrier"`
	Echantillon            string   `xml:"echantillon"`
	Mixtr                  struct {
		Text string `xml:",chardata"`
		Date string `xml:"date,attr"`
		Type struct {
			Text        string `xml:",chardata"`
			V           string `xml:"v,attr"`
			Perimetre   string `xml:"perimetre,attr"`
			Granularite string `xml:"granularite,attr"`
			Qual        string `xml:"qual,attr"`
			Valeur      []struct {
				Text    string `xml:",chardata"`
				Periode string `xml:"periode,attr"`
			} `xml:"valeur"`
		} `xml:"type"`
	} `xml:"mixtr"`
}

func get_co2_emission(XML_URL string) uint8 {
	resp, err := http.Get(XML_URL)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalln(resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var data XMLData

	err = xml.Unmarshal(body, &data)

	if err != nil {
		log.Fatalf("xml.Unmarshal failed with '%s'\n", err)
	}

	co2_values := data.Mixtr.Type.Valeur
	last_co2 := co2_values[len(co2_values)-1].Text

	co2_emission, err := strconv.ParseUint(last_co2, 10, 64)

	if err != nil {
		log.Fatalln(err)
	}

	return uint8(co2_emission)
}

func recordMetrics() {
	go func() {
		for {
			co2_emission := get_co2_emission(XML_URL)
			co2EmissionGauge.Set(float64(co2_emission))
			time.Sleep(15 * time.Minute)
		}
	}()
}

var (
	XML_URL          = "https://www.rte-france.com/themes/swi/xml/power-co2-emission-fr.xml"
	co2EmissionGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "co2_emission",
		Help: "CO2 emission of the France's electricity consumption",
	})
)

func main() {
	prometheus.MustRegister(co2EmissionGauge)
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
