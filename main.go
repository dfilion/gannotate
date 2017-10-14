package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

var (
	// Version number.  Set in Makefile.
	Version string

	// Build number. Set in Makefile.
	Build string

	// Configuration settings.
	settings Settings
)

// Settings is a structure containing the values passed as commandline parameters.
type Settings struct {
	host            string // InfluxDB Host
	db              string // Database to write to
	username        string // InfluxDB user
	password        string // InfluxDB user's password
	measurement     string // Measurement name to write to
	tags            string // InfluxDB tags
	annotationTitle string // Annotation title
	annotationDescr string // Annotation description
	annotationTags  string // Annotation tags
}

// dbExists returns a boolean indicating if the name exists.
func dbExists(c client.Client, name string) (result bool, err error) {
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

// parseInfluxdbTags accepts a string of comma separated key=value pairs and
// parses them into a map that it returns.
func parseInfluxdbTags(tags string) (map[string]string, error) {

	var kv map[string]string

	// Initialize the map length bases on the number of key/value pairs
	kv = make(map[string]string, len(strings.Split(tags, ",")))

	for _, val := range strings.Split(tags, ",") {
		parts := strings.Split(val, "=")
		if len(parts) == 2 {
			kv[parts[0]] = parts[1]
		}
	}
	return kv, nil
}

// Print usage information.
func usage(exitCode int) {
	fmt.Println(`Usage of gannotate:
	-D string   InfluxDB database name. Default: annotations
	-H string   InfluxDB server URL.Default: http://localhost:8086
	-U username User name to authenticate with.
	-P password Username's password.
	-T string   Comma separated list of key=value InfluxDB tags.
	-M string   InfluxDb measurement name. Default: events
	-a tags	    Comma separated list of annotation tags. Saved to the tags field.
	-d descr    Annotation description. Saved to the descr field.
	-t title    Annotation title. Saved to the title field.
	-v          Print version information then exit.
	`)
	//fmt.Printf("Version: %s\tBuild: %s\n", Version, Build)
	printVersionInfo()
	os.Exit(exitCode)
}

// Print version and build information.
func printVersionInfo() {
	fmt.Printf("Version: %s\tBuild: %s\n", Version, Build)
}

// parseFlags parses the command line arguments populating the Settings structure.
func parseFlags() {
	var printVersion bool

	flag.StringVar(&settings.host, "H", "http://localhost:8086", "InfluxDB URL.")
	flag.StringVar(&settings.db, "D", "annotations", "Database name.")
	flag.StringVar(&settings.username, "U", "", "Username.")
	flag.StringVar(&settings.password, "P", "", "Password.")
	flag.StringVar(&settings.measurement, "M", "events", "Measurement name.")
	flag.StringVar(&settings.annotationTitle, "t", "", "Annotation title.")
	flag.StringVar(&settings.annotationDescr, "d", "", "Annotation description.")
	flag.StringVar(&settings.annotationTags, "a", "", "Snnotation tags.")
	flag.StringVar(&settings.tags, "T", "", "InfluxDB tags.")
	flag.BoolVar(&printVersion, "v", false, "Print usage information then exit.")
	flag.Usage = func() { usage(0) }
	flag.Parse()
	if !flag.Parsed() {
		flag.PrintDefaults()
	}

	if printVersion {
		printVersionInfo()
		os.Exit(0)
	}

	if settings.annotationTitle == "" || settings.annotationDescr == "" || settings.annotationTags == "" {
		fmt.Printf("error: -t -d  and -a are required\n\n")
		usage(1)
	}

	if (settings.username != "" && settings.password == "") || (settings.username == "" && settings.password != "") {
		fmt.Printf("error: You must specify a username and password\n\n")
		usage(1)
	}
}

func init() {
	if Version == "" {
		Version = "Unknown"
	}

	if Build == "" {
		Build = "Unknown"
	}
}

func main() {

	parseFlags()

	// Connect
	dbconnConfig := client.HTTPConfig{Addr: settings.host}
	if len(settings.username) > 0 {
		dbconnConfig.Username = settings.username
		dbconnConfig.Password = settings.password
	}
	dbconn, err := client.NewHTTPClient(dbconnConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer dbconn.Close()

	// Create the DB if needed
	exists, err := dbExists(dbconn, settings.db)
	if !exists {
		qry := client.NewQuery(fmt.Sprintf("CREATE DATABASE %s", settings.db), "", "")
		resp, err := dbconn.Query(qry)
		if err != nil {
			log.Fatal(err)
		}
		if resp.Error() != nil {
			log.Fatal(resp.Error())
		}
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
	tags, err := parseInfluxdbTags(settings.tags)

	fields := map[string]interface{}{
		"title": settings.annotationTitle,
		"descr": settings.annotationDescr,
		"tags":  settings.annotationTags,
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
