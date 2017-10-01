# send-annotate

A command line tool written in Go for sending annotations to InfluxDB
for use in Grafana.

## Usage

### Command line

```
Usage of send-annotate:
  -db string
        Database name (default "annotations")
  -desc descr
        Annotation description. Saved to descr field.
  -host string
        InfluxDB server URL. (default "http://localhost:8086")
  -m string
        Measurement written to. (default "events")
  -tags string
        Comma separated list of tags.
  -title title
        Annotation title. Saved to title field.
```

### Grafana

Use the following query to retrieve the annotations from InfluxDB.

```
SELECT title, descr, tags from events WHERE $timeFilter order by asc
```



## Dependencies

Requires InfluxDB Client V2 available at `github.com/influxdata/influxdb/client/v2`.

