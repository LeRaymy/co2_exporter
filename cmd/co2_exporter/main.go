package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
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

func get_co2_emission(XML_URL string) string {
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
	last_co2 := co2_values[len(co2_values) - 1].Text

	return last_co2
}

func main() {

	XML_URL := "https://www.rte-france.com/themes/swi/xml/power-co2-emission-fr.xml"
	co2_emission := get_co2_emission(XML_URL)

	log.Println(co2_emission)
}
