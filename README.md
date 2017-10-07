# send-annotate

`send-annotate` is a command line tool written in Go for sending annotations 
to InfluxDB for use with Grafana.

## Usage

### Summary

```
Usage of send-annotate:
        -D string       InfluxDB database name. Default: annotations
        -H string       InfluxDB server URL. Default: http://localhost:8086
        -T string       Comma separated list of key=value InfluxDB tags.
        -M string       InfluxDb measurement name. Default: events      
        -a tags         Comma separated list of annotation tags. Saved to the tags field.
        -d descr        Annotation description. Saved to the descr field.        
        -t title        Annotation title. Saved to the title field.
```

### Details

#### `-D string`
Specify the InfluxDB database name.  
Default: `annotations`

#### `-H string`
Specify the InfluxDB server URL. 
Default: `http://localhost:8086`

#### `-M string`
InfluxDb measurement name used to record the annotation. 
Default: events

#### `-T string`
Specify a comma separated list of key=value pairs to be used as InfluxDB tags.
Optional. No default.

#### `-t title`
Annotation title. Saved to the `title` field.

#### `-a tags`
Comma separated list of key=value pairs used in the Grafana annotation tags. 
Saved to the `tags` field.

#### `-d descr`
Description to be used the Grafana annotation description.
Saved to the `descr` field.


### Example



### Details

Grafana cannot display InfluxDB tags in annotations.  Instead it looks for 
tags in a dedicated field in a comma separated format.  However because InfluxDB
tags are useful for selecting data 

### Grafana

Use the following query to retrieve the annotations from InfluxDB.

```
SELECT title, descr, tags from events WHERE $timeFilter order by asc
```



## Dependencies

Requires InfluxDB Client V2 available at `github.com/influxdata/influxdb/client/v2`.

