package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const rssUrl = "https://nationalbank.kz/rss/get_rates.cfm"

func main() {
	cfg, err := LoadConfiguration("config.json")
	if err != nil {
		log.Println("Error loading configs: ", err.Error())
	}
	err = initDB(cfg)
	if err != nil {
		log.Println("Error connecting db: ", err.Error())
	} else {
		migrate(DB)
		r := mux.NewRouter()
		r.HandleFunc("/kmf/v1/currency/save/{date}", SaveCurrencyRates).Methods("GET")
		r.HandleFunc("/kmf/v1/currency", CurrencyList).Methods("GET")
		fmt.Println(http.ListenAndServe(":"+cfg.Port, r))
	}
}

func SaveCurrencyRates(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	//req, err := http.NewRequest("GET", rssUrl, nil)
	//if err != nil {
	//	log.Print(err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	fmt.Fprintf(w, "success=false")
	//}
	//q := req.URL.Query()
	//q.Add("fdate", date)
	//req.URL.RawQuery = q.Encode()
	fakeResp, err := rssMock(date) //http.Get(req.URL.String())//todo call real rss url
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "success=false")
		return
	}
	rates := new(Rates)
	err = xml.Unmarshal([]byte(fakeResp), rates)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "success=false")
		return
	}

	go InsertRates(date, rates)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "success=true")
}

func CurrencyList(w http.ResponseWriter, r *http.Request) {
	qmap := r.URL.Query()
	date := qmap.Get("date")
	code := qmap.Get("code")
	list, err := SelectRates(date, code)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "success=false")
		return
	}

	bytes, err := json.Marshal(list)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "success=false")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(bytes))

}
