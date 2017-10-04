package main

import (
	"flag"
	"fmt"
	"log"
	//"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

// Settings is a structure containing the values passed as commandline parameters.
type Settings struct {
	host             string // InfluxDB Host
	db               string // Database to write to
	measurement      string // Measurement name to write to
	tags             string // InfluxDB tags
	annotation_title string // Annotation title
	annotation_descr string // Annotation description
	annotation_tags  string // Annotation tags
}

// dbexists returns a boolean indicating if the name exists.
func dbexists(c client.Client, name string) (result bool, err error) {
	qry := client.NewQuery("SHOW DATABASES", "", "")
	response, err := c.Query(qry)
	if err != nil {
		return false, err
	}
	if response.Error() != nil {
		return false, response.Error()
	}
	for _, n := range response.Results[0].Series[0].Values {
		//fmt.Printf("%s == %s\n", n[0], name)
		if n[0] == name {
			return true, nil
		}
	}
	return
}

func main() {

	var settings Settings

	flag.StringVar(&settings.host, "host", "http://localhost:8086", "InfluxDB server URL.")
	flag.StringVar(&settings.db, "db", "annotations", "InfluxDB database name")
	flag.StringVar(&settings.measurement, "m", "events", "InfluxDb measurement name.")
	flag.StringVar(&settings.annotation_title, "title", "", "Annotation title. Saved to the `title` field.")
	flag.StringVar(&settings.annotation_descr, "desc", "", "Annotation description. Saved to the `descr` field.")
	flag.StringVar(&settings.annotation_tags, "tags", "", "Comma separated list of annotation tags. Saved to the `tags` field.")
	flag.StringVar(&settings.tags, "Tags", "", "Comma separated list of InfluxDB tags.")
	flag.Parse()
	if !flag.Parsed() {
		flag.PrintDefaults()
	}

	// Connect
	dbconn_config := client.HTTPConfig{Addr: settings.host}
	dbconn, err := client.NewHTTPClient(dbconn_config)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	// Create the DB if needed
	exists, err := dbexists(dbconn, settings.db)
	if !exists {
		qry := client.NewQuery(fmt.Sprintf("CREATE DATABASE %s", settings.db), "", "")
		resp, err := dbconn.Query(qry)
		if err != nil {
			log.Fatal(err)
		}
		if resp.Error() != nil {
			log.Fatal(resp.Error())
		}
		fmt.Println(resp)
	}

	// Create a Point Batch
	bpc := client.BatchPointsConfig{
		Precision: "s",
		Database:  settings.db,
	}
	bp, err := client.NewBatchPoints(bpc)
	if err != nil {
		log.Fatal(err)
	}

	// Create a data Point
	//tags := map[string]string{"aTag": "aVal", "bTag": "bVal"}
	tags := map[string]string{} // Annotations do not support influxdb style tags
	//tags := strings.Split(settings.tags, ",")

	fields := map[string]interface{}{
		"title": settings.annotation_title,
		"descr": settings.annotation_descr,
		"tags":  settings.annotation_tags,
	}

	pt, err := client.NewPoint(settings.measurement, tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	// Add the point to the batch
	bp.AddPoint(pt)

	// Write the batch to the server
	err = dbconn.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
}
