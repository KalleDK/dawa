package dawa

import (
	"bytes"
	"io"
	"reflect"
	"testing"

  "github.com/kalledk/dawa/url"
  "github.com/kalledk/dawa/uuid"
  "github.com/kalledk/dawa/time"
)

var csv_data = `id,status,oprettet,ændret,vejkode,vejnavn,husnr,etage,dør,supplerendebynavn,postnr,postnrnavn,kommunekode,kommunenavn,ejerlavkode,ejerlavnavn,matrikelnr,esrejendomsnr,etrs89koordinat_øst,etrs89koordinat_nord,wgs84koordinat_bredde,wgs84koordinat_længde,nøjagtighed,kilde,tekniskstandard,tekstretning,ddkn_m100,ddkn_km1,ddkn_km10,adressepunktændringsdato,adgangsadresseid,adgangsadresse_status,adgangsadresse_oprettet,adgangsadresse_ændret,kvhx,regionskode,regionsnavn,sognekode,sognenavn,politikredskode,politikredsnavn,retskredskode,retskredsnavn,opstillingskredskode,opstillingskredsnavn,zone,href
0a3f50b7-6545-32b8-e044-0003ba298018,1,2000-02-05T18:09:56.000,2000-02-16T21:58:33.000,0001,A Hansensvej,6,,,Vråby,6792,Rømø,0550,Tønder,1470852,"Kirkeby, Rømø",76,9097,470620,6105713,55.0972751504817,8.53959543878291,A,5,UF,200,100m_61057_4706,1km_6105_470,10km_610_47,2004-10-08T00:00:00.000,0a3f508c-3307-32b8-e044-0003ba298018,1,2000-02-05T18:09:56.000,2009-11-24T03:15:25.000,05500001___6_______,1083,Region Syddanmark,9062,Rømø,1464,Syd- og Sønderjyllands Politi,1147,Retten i Sønderborg,0051,Tønder,Landzone,http://dawa.aws.dk/adresser/0a3f50b9-68b1-32b8-e044-0003ba298018
0a3f50b7-6544-32b8-e044-0003ba298018,1,2000-02-05T18:09:53.000,2000-02-16T21:58:33.000,0001,A Hansensvej,5,,,Vråby,6792,Rømø,0550,Tønder,1470852,"Kirkeby, Rømø",720,8848,470531,6105718,55.0973148017997,8.53820028085436,A,5,UF,200,100m_61057_4705,1km_6105_470,10km_610_47,2004-10-08T00:00:00.000,0a3f508c-3306-32b8-e044-0003ba298018,1,2000-02-05T18:09:53.000,2009-11-24T03:15:25.000,05500001___5_______,1083,Region Syddanmark,9062,Rømø,1464,Syd- og Sønderjyllands Politi,1147,Retten i Sønderborg,0051,Tønder,Landzone,http://dawa.aws.dk/adresser/0a3f50b9-68b1-32b8-e044-0003ba298018
0a3f50b7-6547-32b8-e044-0003ba298018,1,2000-02-05T18:09:49.000,2004-02-12T16:05:28.000,0001,A Hansensvej,8,,,Vråby,6792,Rømø,0550,Tønder,1470852,"Kirkeby, Rømø",770,8559,470587,6105811,55.0981538216665,8.5390681928166,A,5,UF,200,100m_61058_4705,1km_6105_470,10km_610_47,2004-10-08T00:00:00.000,0a3f508c-3309-32b8-e044-0003ba298018,1,2000-02-05T18:09:49.000,2009-11-24T03:15:25.000,05500001___8_______,1083,Region Syddanmark,9062,Rømø,1464,Syd- og Sønderjyllands Politi,1147,Retten i Sønderborg,0051,Tønder,Landzone,http://dawa.aws.dk/adresser/0a3f50b9-68b1-32b8-e044-0003ba298018
`

func TestImportAdresserCSV(t *testing.T) {
	// We test one entry only to match
	csv_expect := []Adresse{
		Adresse{
			Adgangsadresse: AdgangsAdresse{
				DDKN: DDKN{
          Km1: "1km_6105_470",
          Km10: "10km_610_47",
          M100: "100m_61057_4706",
        },
				Adgangspunkt: Adgangspunkt{
          Kilde: 5,
          Koordinater: []float64{55.0972751504817, 8.53959543878291},
          Nøjagtighed: "A",
          Tekniskstandard: "UF",
          Tekstretning: 200,
          Ændret: time.MustParse("2004-10-08T00:00:00.000"),
				},
				Ejerlav: Ejerlav{
          Kode: 1470852,
          Navn: "Kirkeby, Rømø",
        },
				EsrEjendomsNr:     "9097",
				Historik:          Historik{Oprettet: time.MustParse("2000-02-05T18:09:56.000"), Ændret: time.MustParse("2009-11-24T03:15:25.000")},
				Href:              "",
				Husnr:             "6",
				ID:                "0a3f508c-3307-32b8-e044-0003ba298018",
				Kommune:           KommuneRef{Href: "", Kode: "0550", Navn: "Tønder"},
				Kvh:               "05500001___6",
				Matrikelnr:        "76",
				Opstillingskreds:  OpstillingskredsRef{Href: "", Kode: "0051", Navn: "Tønder"},
				Politikreds:       PolitikredsRef{Href: "", Kode: "1464", Navn: "Syd- og Sønderjyllands Politi"},
				Postnummer:        PostnummerRef{Href: "", Navn: "Rømø", Nr: "6792"},
				Region:            RegionRef{Href: "", Kode: "1083", Navn: "Region Syddanmark"},
				Retskreds:         RetskredsRef{Href: "", Kode: "1147", Navn: "Retten i Sønderborg"},
				Sogn:              SognRef{Href: "", Kode: "9062", Navn: "Rømø"},
				Status:            1,
				SupplerendeBynavn: "Vråby",
				Vejstykke:         VejstykkeRef{Href: "", Kode: "0001", Navn: "A Hansensvej"}, Zone: "Landzone",
			},
			Adressebetegnelse: "",
			Dør:               "",
			Etage:             "",
			Historik: Historik{
				Oprettet: time.MustParse("2000-02-05T18:09:56.000"),
				Ændret:   time.MustParse("2000-02-16T21:58:33.000"),
			},
			Href:   url.MustParse("http://dawa.aws.dk/adresser/0a3f50b9-68b1-32b8-e044-0003ba298018"),
			ID:     uuid.MustParse("0a3f50b7-6545-32b8-e044-0003ba298018"),
			Kvhx:   "05500001___6_______",
			Status: 1},
		Adresse{},
		Adresse{},
	}
	b := bytes.NewBuffer([]byte(csv_data))
	iter, err := ImportAdresserCSV(b)
	if err != nil {
		t.Fatalf("ImportAdresserCSV: %v", err)
	}
	for i, expect := range csv_expect {
		item, err := iter.Next()
		if err != nil {
			t.Fatalf("ImportAdresserCSV, iter.Next(): %v", err)
		}
		if item == nil {
			t.Fatalf("ImportAdresserCSV, iter.Next() returned nil value")
		}
		if i == 0 && !reflect.DeepEqual(*item, expect) {
			t.Fatalf("ImportAdresserCSV, value mismatch.\nGot:\n%#v\nExpected:\n%#v\n", *item, expect)
		}
		if i == 0 {
			if item.Historik.Oprettet.Time().Unix() != 949770596 {
				t.Fatalf("Timestamp mismatch, expected 949770596, was %d", item.Historik.Oprettet.Time().Unix())
			}
		}
	}
	// We should now have read all entries
	_, err = iter.Next()
	if err != io.EOF {
		t.Fatalf("ImportAdresserCSV: Expected io.EOF, got:%v", err)
	}
}

var json_input = `
[
{
  "id": "0a3f50b9-68b1-32b8-e044-0003ba298018",
  "kvhx": "05630110___1__1____",
  "status": 1,
  "href": "http://dawa.aws.dk/adresser/0a3f50b9-68b1-32b8-e044-0003ba298018",
  "historik": {
    "oprettet": "2000-02-05T18:30:56.000",
    "ændret": "2000-02-16T22:02:44.000"
  },
  "etage": "1",
  "dør": null,
  "adressebetegnelse": "A B C Sti 1, 1., Nordby, 6720 Fanø",
  "adgangsadresse": {
    "href": "http://dawa.aws.dk/adgangsadresser/0a3f508d-d915-32b8-e044-0003ba298018",
    "id": "0a3f508d-d915-32b8-e044-0003ba298018",
    "kvh": "05630110___1",
    "status": 1,
    "vejstykke": {
      "href": "http://dawa.aws.dk/vejstykker/563/110",
      "navn": "A B C Sti",
      "kode": "0110"
    },
    "husnr": "1",
    "supplerendebynavn": "Nordby",
    "postnummer": {
      "href": "http://dawa.aws.dk/postnumre/6720",
      "nr": "6720",
      "navn": "Fanø"
    },
    "kommune": {
      "href": "http://dawa.aws.dk/kommuner/563",
      "kode": "0563",
      "navn": "Fanø"
    },
    "ejerlav": {
      "kode": 1351151,
      "navn": "Odden By, Nordby"
    },
    "esrejendomsnr": "10045",
    "matrikelnr": "320",
    "historik": {
      "oprettet": "2000-02-05T18:30:56.000",
      "ændret": "2009-11-24T03:15:25.000"
    },
    "adgangspunkt": {
      "koordinater": [
        8.40179905638495,
        55.4454386963562
      ],
      "nøjagtighed": "A",
      "kilde": 1,
      "tekniskstandard": "TK",
      "tekstretning": 125.9,
      "ændret": "2000-09-18T00:00:00.000"
    },
    "DDKN": {
      "m100": "100m_61445_4621",
      "km1": "1km_6144_462",
      "km10": "10km_614_46"
    },
    "sogn": {
      "kode": "8923",
      "navn": "Nordby",
      "href": "http://dawa.aws.dk/sogne/8923"
    },
    "region": {
      "kode": "1083",
      "navn": "Region Syddanmark",
      "href": "http://dawa.aws.dk/regioner/1083"
    },
    "retskreds": {
      "kode": "1151",
      "navn": "Retten i Esbjerg",
      "href": "http://dawa.aws.dk/retskredse/1151"
    },
    "politikreds": {
      "kode": "1464",
      "navn": "Syd- og Sønderjyllands Politi",
      "href": "http://dawa.aws.dk/politikredse/1464"
    },
    "opstillingskreds": {
      "kode": "0052",
      "navn": "Esbjerg By",
      "href": "http://dawa.aws.dk/opstillingskredse/52"
    },
    "zone": "Byzone"
  }
}
,
{
  "id": "0a3f50b9-7be9-32b8-e044-0003ba298018",
  "kvhx": "05639895__27_______",
  "status": 1,
  "href": "http://dawa.aws.dk/adresser/0a3f50b9-7be9-32b8-e044-0003ba298018",
  "historik": {
    "oprettet": "2000-02-05T18:31:07.000",
    "ændret": "2002-09-27T12:25:29.000"
  },
  "etage": null,
  "dør": null,
  "adressebetegnelse": "Østre Klitvej 27, Rindby Strand, 6720 Fanø",
  "adgangsadresse": {
    "href": "http://dawa.aws.dk/adgangsadresser/0a3f508d-eb8e-32b8-e044-0003ba298018",
    "id": "0a3f508d-eb8e-32b8-e044-0003ba298018",
    "kvh": "05639895__27",
    "status": 1,
    "vejstykke": {
      "href": "http://dawa.aws.dk/vejstykker/563/9895",
      "navn": "Østre Klitvej",
      "kode": "9895"
    },
    "husnr": "27",
    "supplerendebynavn": "Rindby Strand",
    "postnummer": {
      "href": "http://dawa.aws.dk/postnumre/6720",
      "nr": "6720",
      "navn": "Fanø"
    },
    "kommune": {
      "href": "http://dawa.aws.dk/kommuner/563",
      "kode": "0563",
      "navn": "Fanø"
    },
    "ejerlav": {
      "kode": 1351152,
      "navn": "Rindby By, Nordby"
    },
    "esrejendomsnr": "21101",
    "matrikelnr": "25bc",
    "historik": {
      "oprettet": "2000-02-05T18:31:07.000",
      "ændret": "2009-11-24T03:15:25.000"
    },
    "adgangspunkt": {
      "koordinater": [
        8.39446887121117,
        55.4150234465972
      ],
      "nøjagtighed": "A",
      "kilde": 1,
      "tekniskstandard": "TK",
      "tekstretning": 187.46,
      "ændret": "2000-09-18T00:00:00.000"
    },
    "DDKN": {
      "m100": "100m_61411_4616",
      "km1": "1km_6141_461",
      "km10": "10km_614_46"
    },
    "sogn": {
      "kode": "8923",
      "navn": "Nordby",
      "href": "http://dawa.aws.dk/sogne/8923"
    },
    "region": {
      "kode": "1083",
      "navn": "Region Syddanmark",
      "href": "http://dawa.aws.dk/regioner/1083"
    },
    "retskreds": {
      "kode": "1151",
      "navn": "Retten i Esbjerg",
      "href": "http://dawa.aws.dk/retskredse/1151"
    },
    "politikreds": {
      "kode": "1464",
      "navn": "Syd- og Sønderjyllands Politi",
      "href": "http://dawa.aws.dk/politikredse/1464"
    },
    "opstillingskreds": {
      "kode": "0052",
      "navn": "Esbjerg By",
      "href": "http://dawa.aws.dk/opstillingskredse/52"
    },
    "zone": "Landzone"
  }
}, {
  "id": "0a3f50b9-7be7-32b8-e044-0003ba298018",
  "kvhx": "05639895__23_______",
  "status": 1,
  "href": "http://dawa.aws.dk/adresser/0a3f50b9-7be7-32b8-e044-0003ba298018",
  "historik": {
    "oprettet": "2000-02-05T18:31:07.000",
    "ændret": "2000-02-16T22:02:53.000"
  },
  "etage": null,
  "dør": null,
  "adressebetegnelse": "Østre Klitvej 23, Rindby Strand, 6720 Fanø",
  "adgangsadresse": {
    "href": "http://dawa.aws.dk/adgangsadresser/0a3f508d-eb8c-32b8-e044-0003ba298018",
    "id": "0a3f508d-eb8c-32b8-e044-0003ba298018",
    "kvh": "05639895__23",
    "status": 1,
    "vejstykke": {
      "href": "http://dawa.aws.dk/vejstykker/563/9895",
      "navn": "Østre Klitvej",
      "kode": "9895"
    },
    "husnr": "23",
    "supplerendebynavn": "Rindby Strand",
    "postnummer": {
      "href": "http://dawa.aws.dk/postnumre/6720",
      "nr": "6720",
      "navn": "Fanø"
    },
    "kommune": {
      "href": "http://dawa.aws.dk/kommuner/563",
      "kode": "0563",
      "navn": "Fanø"
    },
    "ejerlav": {
      "kode": 1351152,
      "navn": "Rindby By, Nordby"
    },
    "esrejendomsnr": "21071",
    "matrikelnr": "25ba",
    "historik": {
      "oprettet": "2000-02-05T18:31:07.000",
      "ændret": "2009-11-24T03:15:25.000"
    },
    "adgangspunkt": {
      "koordinater": [
        8.39321938836807,
        55.4150944451575
      ],
      "nøjagtighed": "A",
      "kilde": 1,
      "tekniskstandard": "TK",
      "tekstretning": 276.6,
      "ændret": "2000-09-18T00:00:00.000"
    },
    "DDKN": {
      "m100": "100m_61411_4615",
      "km1": "1km_6141_461",
      "km10": "10km_614_46"
    },
    "sogn": {
      "kode": "8923",
      "navn": "Nordby",
      "href": "http://dawa.aws.dk/sogne/8923"
    },
    "region": {
      "kode": "1083",
      "navn": "Region Syddanmark",
      "href": "http://dawa.aws.dk/regioner/1083"
    },
    "retskreds": {
      "kode": "1151",
      "navn": "Retten i Esbjerg",
      "href": "http://dawa.aws.dk/retskredse/1151"
    },
    "politikreds": {
      "kode": "1464",
      "navn": "Syd- og Sønderjyllands Politi",
      "href": "http://dawa.aws.dk/politikredse/1464"
    },
    "opstillingskreds": {
      "kode": "0052",
      "navn": "Esbjerg By",
      "href": "http://dawa.aws.dk/opstillingskredse/52"
    },
    "zone": "Landzone"
  }
}
]
`

func TestImportAdresserJSON(t *testing.T) {
	// We test one entry only to match
	var json_expect = []Adresse{
		Adresse{Adgangsadresse: AdgangsAdresse{
			DDKN:             DDKN{Km1: "1km_6144_462", Km10: "10km_614_46", M100: "100m_61445_4621"},
			Adgangspunkt:     Adgangspunkt{Kilde: 1, Koordinater: []float64{8.40179905638495, 55.4454386963562}, Nøjagtighed: "A", Tekniskstandard: "TK", Tekstretning: 125.9, Ændret: time.MustParse("2000-09-18T00:00:00.000")}, // "ændret": "2000-09-18T00:00:00.000"
			Ejerlav:          Ejerlav{Kode: 1351151, Navn: "Odden By, Nordby"},
			EsrEjendomsNr:    "10045",
			Historik:         Historik{Oprettet: time.MustParse("2000-02-05T18:30:56.000"), Ændret: time.MustParse("2009-11-24T03:15:25.000")}, //       "oprettet": "2000-02-05T18:30:56.000", "ændret": "2009-11-24T03:15:25.000"
			Href:             "http://dawa.aws.dk/adgangsadresser/0a3f508d-d915-32b8-e044-0003ba298018",
			Husnr:            "1",
			ID:               "0a3f508d-d915-32b8-e044-0003ba298018",
			Kommune:          KommuneRef{Href: "http://dawa.aws.dk/kommuner/563", Kode: "0563", Navn: "Fanø"},
			Kvh:              "05630110___1",
			Matrikelnr:       "320",
			Opstillingskreds: OpstillingskredsRef{Href: "http://dawa.aws.dk/opstillingskredse/52", Kode: "0052", Navn: "Esbjerg By"},
			Politikreds:      PolitikredsRef{Href: "http://dawa.aws.dk/politikredse/1464", Kode: "1464", Navn: "Syd- og Sønderjyllands Politi"},
			Postnummer:       PostnummerRef{Href: "http://dawa.aws.dk/postnumre/6720", Navn: "Fanø", Nr: "6720"},
			Region:           RegionRef{Href: "http://dawa.aws.dk/regioner/1083", Kode: "1083", Navn: "Region Syddanmark"},
			Retskreds:        RetskredsRef{Href: "http://dawa.aws.dk/retskredse/1151", Kode: "1151", Navn: "Retten i Esbjerg"},
			Sogn:             SognRef{Href: "http://dawa.aws.dk/sogne/8923", Kode: "8923", Navn: "Nordby"}, Status: 1, SupplerendeBynavn: "Nordby",
			Vejstykke: VejstykkeRef{Href: "http://dawa.aws.dk/vejstykker/563/110", Kode: "0110", Navn: "A B C Sti"}, Zone: "Byzone"},
			Adressebetegnelse: "A B C Sti 1, 1., Nordby, 6720 Fanø",
			Dør:               "",
			Etage:             "1",
			Historik:          Historik{Oprettet: time.MustParse("2000-02-05T18:30:56.000"), Ændret: time.MustParse("2000-02-16T22:02:44.000")}, //     "oprettet": "2000-02-05T18:30:56.000", "ændret": "2000-02-16T22:02:44.000"
			Href:              url.MustParse("http://dawa.aws.dk/adresser/0a3f50b9-68b1-32b8-e044-0003ba298018"),
			ID:                uuid.MustParse("0a3f50b9-68b1-32b8-e044-0003ba298018"),
			Kvhx:              "05630110___1__1____",
			Status:            1,
		},
		Adresse{},
		Adresse{},
	}

	b := bytes.NewBuffer([]byte(json_input))
	iter, err := ImportAdresserJSON(b)
	if err != nil {
		t.Fatalf("ImportAdresserJSON: %v", err)
	}
	for i, expect := range json_expect {
		item, err := iter.Next()
		if err != nil {
			t.Fatalf("ImportAdresserJSON, iter.Next(): %v", err)
		}
		if item == nil {
			t.Fatalf("ImportAdresserJSON, iter.Next() returned nil value")
		}
		if i == 0 && !reflect.DeepEqual(*item, expect) {
			t.Fatalf("ImportAdresserJSON, value mismatch.\nGot:\n%#v\nExpected:\n%#v\n", *item, expect)
		}
		// Since we leak time parsing abstraction, we need to test a value.
		if i == 0 {
			if item.Historik.Oprettet.Time().Unix() != 949771856 {
				t.Fatalf("Timestamp mismatch, expected 949771856, was %d", item.Historik.Oprettet.Time().Unix())
			}
		}
	}
	// We should now have read all entries
	_, err = iter.Next()
	if err != io.EOF {
		t.Fatalf("ImportAdresserJSON: Expected io.EOF, got:%v", err)
	}
}
