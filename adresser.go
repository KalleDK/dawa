package dawa

import (
	"bufio"
	"encoding/csv"
	"io"
	"strconv"

	"github.com/kalledk/dawa/time"
	"github.com/kalledk/dawa/url"
	"github.com/kalledk/dawa/uuid"
	"github.com/ugorji/go/codec"
)

type Adresse struct {
	Adgangsadresse    AdgangsAdresse `json:"adgangsadresse"`    // Adressens adgangsadresse
	Adressebetegnelse string         `json:"adressebetegnelse"` // (unknown)
	Dør               string         `json:"dør"`               // Dørbetegnelse. Tal fra 1 til 9999, små og store bogstaver samt tegnene / og -.
	Etage             string         `json:"etage"`             // Etagebetegnelse. Hvis værdi angivet kan den antage følgende værdier: tal fra 1 til 99, st, kl, kl2 op til kl9.
	Historik          Historik       `json:"historik"`          // Væsentlige tidspunkter for adressen
	Href              url.URL        `json:"href"`              // Adgangsadressens URL.
	ID                uuid.UUID      `json:"id"`                // Adressens unikke id, f.eks. 0a3f5095-45ec-32b8-e044-0003ba298018.
	Kvhx              string         `json:"kvhx"`              // KVHX-nøgle. 19 tegn bestående af 4 cifre der repræsenterer kommunekode, 4 cifre der repræsenterer vejkode, 4 tegn der repræsenter husnr, 3 tegn der repræsenterer etage og 4 tegn der repræsenter dør.
	Status            int            `json:"status"`            // Adressens status. 1 indikerer en gældende adresse, 3 indikerer en foreløbig adresse.
}

// AdresseIter is an Iterator that enable you to get individual entries.
type AdresseIter struct {
	closer
	a   chan Adresse
	err error
}

// Next will return addresses.
// It will return an error if that has been encountered.
// When there are not more entries nil, io.EOF will be returned.
func (a *AdresseIter) Next() (*Adresse, error) {
	v, ok := <-a.a
	if ok {
		return &v, nil
	}
	return nil, a.err
}

// ImportAdresserCSV will import "adresser" from a CSV file, supplied to the reader.
// An iterator will be returned that return all addresses.
func ImportAdresserCSV(in io.Reader) (*AdresseIter, error) {
	r := csv.NewReader(in)
	r.Comma = ','

	// Read first line as headers
	name, err := r.Read()
	if err != nil {
		return nil, err
	}

	ret := &AdresseIter{a: make(chan Adresse, 100)}
	if err != nil {
		return nil, err
	}
	go func() {
		defer close(ret.a)
		v := make(map[string]string, len(name))
		for {
			records, err := r.Read()
			if err != nil {
				ret.err = err
				return
			}
			// Map to indexes, so we don't rely on index numbers, but on column names.
			for j := range records {
				v[name[j]] = records[j]
			}
			// PROCESS: id,status,oprettet,ændret,vejkode,vejnavn,husnr,etage,dør,supplerendebynavn
			a := Adresse{}
			a.ID, err = uuid.Parse(v["id"])
			if err != nil {
				ret.err = err
				return
			}

			a.Status, err = strconv.Atoi(v["status"])
			if err != nil {
				ret.err = err
				return
			}

			// Example 2000-02-16T21:58:33.000
			o, err := time.Parse(v["oprettet"])
			if err != nil {
				ret.err = err
				return
			}
			a.Historik.Oprettet = *o

			o, err = time.Parse(v["ændret"])
			if err != nil {
				ret.err = err
				return
			}

			h, err := url.Parse(v["href"])
			if err != nil {
				ret.err = err
				return
			}
			a.Href = h

			a.Historik.Ændret = *o

			a.Adgangsadresse.Vejstykke.Kode = v["vejkode"]
			a.Adgangsadresse.Vejstykke.Navn = v["vejnavn"]
			a.Adgangsadresse.Husnr = v["husnr"]
			a.Etage = v["etage"]
			a.Adgangsadresse.SupplerendeBynavn = v["supplerendebynavn"]

			// PROCESS: postnr,postnrnavn,kommunekode,kommunenavn,ejerlavkode,ejerlavnavn,matrikelnr,esrejendomsnr,etrs89koordinat_øst,etrs89koordinat_nord,wgs84koordinat_bredde,wgs84koordinat_længde,
			a.Adgangsadresse.Postnummer.Nr = v["postnr"]
			a.Adgangsadresse.Postnummer.Navn = v["postnrnavn"]
			a.Adgangsadresse.Kommune.Kode = v["kommunekode"]
			a.Adgangsadresse.Kommune.Navn = v["kommunenavn"]
			a.Adgangsadresse.Ejerlav.Kode, _ = strconv.Atoi(v["ejerlavkode"])
			a.Adgangsadresse.Ejerlav.Navn = v["ejerlavnavn"]
			a.Adgangsadresse.Matrikelnr = v["matrikelnr"]
			a.Adgangsadresse.EsrEjendomsNr = v["esrejendomsnr"]
			// ????
			// x,_ = strconv.ParseFloat("etrs89koordinat_øst")
			// x,_ = strconv.ParseFloat("etrs89koordinat_nord")
			a.Adgangsadresse.Adgangspunkt.Koordinater = make([]float64, 2)
			a.Adgangsadresse.Adgangspunkt.Koordinater[0], _ = strconv.ParseFloat(v["wgs84koordinat_bredde"], 64)
			a.Adgangsadresse.Adgangspunkt.Koordinater[1], _ = strconv.ParseFloat(v["wgs84koordinat_længde"], 64)

			// PROCESS: nøjagtighed,kilde,tekniskstandard,tekstretning,ddkn_m100,ddkn_km1,ddkn_km10,adressepunktændringsdato,adgangsadresseid,adgangsadresse_status
			a.Adgangsadresse.Adgangspunkt.Nøjagtighed = v["nøjagtighed"]
			a.Adgangsadresse.Adgangspunkt.Kilde, _ = strconv.Atoi(v["kilde"])
			a.Adgangsadresse.Adgangspunkt.Tekniskstandard = v["tekniskstandard"]
			a.Adgangsadresse.Adgangspunkt.Tekstretning, _ = strconv.ParseFloat(v["tekstretning"], 64)
			a.Adgangsadresse.DDKN.M100 = v["ddkn_m100"]
			a.Adgangsadresse.DDKN.Km1 = v["ddkn_km1"]
			a.Adgangsadresse.DDKN.Km10 = v["ddkn_km10"]
			o, err = time.Parse(v["adressepunktændringsdato"])
			if err != nil {
				ret.err = err
				return
			}
			a.Adgangsadresse.Adgangspunkt.Ændret = *o
			a.Adgangsadresse.ID = v["adgangsadresseid"]
			a.Adgangsadresse.Status, _ = strconv.Atoi(v["adgangsadresse_status"])

			// PROCESS: adgangsadresse_oprettet,adgangsadresse_ændret,kvhx,regionskode,regionsnavn,sognekode,sognenavn,politikredskode,politikredsnavn,retskredskode,retskredsnavn
			o, err = time.Parse(v["adgangsadresse_oprettet"])
			if err != nil {
				ret.err = err
				return
			}
			a.Adgangsadresse.Historik.Oprettet = *o
			o, err = time.Parse(v["adgangsadresse_ændret"])
			if err != nil {
				ret.err = err
				return
			}
			a.Adgangsadresse.Historik.Ændret = *o
			a.Kvhx = v["kvhx"]
			a.Adgangsadresse.Kvh = string([]byte(a.Kvhx)[:12])
			a.Adgangsadresse.Region.Kode = v["regionskode"]
			a.Adgangsadresse.Region.Navn = v["regionsnavn"]
			a.Adgangsadresse.Sogn.Kode = v["sognekode"]
			a.Adgangsadresse.Sogn.Navn = v["sognenavn"]
			a.Adgangsadresse.Politikreds.Kode = v["politikredskode"]
			a.Adgangsadresse.Politikreds.Navn = v["politikredsnavn"]
			a.Adgangsadresse.Retskreds.Kode = v["retskredskode"]
			a.Adgangsadresse.Retskreds.Navn = v["retskredsnavn"]

			// opstillingskredskode,opstillingskredsnavn,zone
			a.Adgangsadresse.Opstillingskreds.Kode = v["opstillingskredskode"]
			a.Adgangsadresse.Opstillingskreds.Navn = v["opstillingskredsnavn"]
			a.Adgangsadresse.Zone = v["zone"]
			ret.a <- a
		}
	}()
	return ret, nil
}

// ImportAdresserJSON will import "adresser" from a JSON input, supplied to the reader.
// An iterator will be returned that return all addresses.
func ImportAdresserJSON(in io.Reader) (*AdresseIter, error) {
	var h codec.JsonHandle
	h.DecodeOptions.ErrorIfNoField = JSONStrictFieldCheck
	// use a buffered reader for efficiency
	if _, ok := in.(io.ByteScanner); !ok {
		in = bufio.NewReader(in)
	}
	ret := &AdresseIter{a: make(chan Adresse, 100)}
	go func() {
		defer close(ret.a)
		var dec *codec.Decoder = codec.NewDecoder(in, &h)
		ret.err = dec.Decode(&ret.a)
		if ret.err == nil {
			ret.err = io.EOF
		}
	}()

	return ret, nil
}
