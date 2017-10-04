# send-annotate

A command line tool written in Go for sending annotations to InfluxDB
for use in Grafana.

## Usage

### Command line

```
Usage of send-annotate:
  -Tags string
        Comma separated list of InfluxDB tags (in progress). 
  -db string
        InfluxDB database name (default "annotations")
  -desc descr
        Annotation description. Saved to the descr field.
  -host string
        InfluxDB server URL. (default "http://localhost:8086")
  -m string
        InfluxDb measurement name. (default "events")
  -tags tags
        Comma separated list of annotation tags. Saved to the tags field.
  -title title
        Annotation title. Saved to the title field.

```

### Grafana

Use the following query to retrieve the annotations from InfluxDB.

```
SELECT title, descr, tags from events WHERE $timeFilter order by asc
```



## Dependencies

Requires InfluxDB Client V2 available at `github.com/influxdata/influxdb/client/v2`.

