package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	dbviewer "internal.com/db-viewer/internal/db-viewer"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type NeighborInfo struct {
	IpAddr string
	MacAddr string
	Oui string
}

func main() {
	var err error

	fs := http.FileServer(
		dbviewer.NeuteredFileSystem {
			http.Dir("./web/static"),
		},
	)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/nmap/", nmapHandler)

	logrus.Info("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		logrus.Fatal("http.ListenAndServe ", err)
	}
}

func nmapHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	lp := filepath.Join("web/templates", "layout.html")
	fp := filepath.Join("web/templates", "nmap.html")

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// log the detailed error
		logrus.Error("template.ParseFiles ", err)
		// return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	db, err := sql.Open("sqlite3", "./data/nmap.db")
	if err != nil {
		logrus.Error("sql.Open ", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rows, err := db.Query("SELECT IpAddr, MacAddr, Oui FROM neighbors")
	if err != nil {
		logrus.Error("db.Query ", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	switch r.URL.Path {
	case "/nmap/":
		err = tmpl.ExecuteTemplate(w, "layout", nil)
	case "/nmap/more":
		var ipAddr string
		var macAddr string
		var oui string
		var neighbors []NeighborInfo
		neighbors = append(
			neighbors,
			NeighborInfo {
				IpAddr: "IP",
				MacAddr: "MAC",
				Oui: "OUI",
			},
		)
		fmt.Printf("IP\t\tMAC\t\t\tOUI\n")

		for rows.Next() {
			err = rows.Scan(&ipAddr, &macAddr, &oui)
			if err != nil {
				logrus.Error("rows.Scan ", err)
				http.Error(w, http.StatusText(500), 500)
				return
			}
	
			fmt.Printf("%s\t%s\t%s\n", ipAddr, macAddr, oui)
	
			neighbors = append(
				neighbors,
				NeighborInfo {
					IpAddr: ipAddr,
					MacAddr: macAddr,
					Oui: oui,
				},
			)
		}

		db.Close()
		err = tmpl.ExecuteTemplate(w, "neighbors", neighbors)
	}
	if err != nil {
		logrus.Error("tmpl.ExecuteTemplate ", err)
		http.Error(w, http.StatusText(500), 500)
	}
}
