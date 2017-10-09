# gannotate

`gannotate` is a command line tool written in Go for sending annotations
to InfluxDB for use with Grafana.

## Usage Summary

```
Usage of gannotate:
        -D string       InfluxDB database name. Default: annotations
        -H string       InfluxDB server URL. Default: http://localhost:8086
        -T string       Comma separated list of key=value InfluxDB tags.
        -M string       InfluxDb measurement name. Default: events
        -a tags         Comma separated list of annotation tags. Saved to the tags field.
        -d descr        Annotation description. Saved to the descr field.
        -t title        Annotation title. Saved to the title field.
        -v              Print version information then exit.
```

-a, -d and -t do not have defaults so they are required.

### Arguments

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

#### `-v`
Print the applications version information and exit.


## Example

### Create an annotation using the default InfluxDB settings.
```
gannotate -t aTitle1 -d aDesc1 -a aTag1,aTag2
```

What gets created.
```
> show databases;
name: databases
name
----
_internal
telegraf
annotations

> use annotations;
Using database annotations

> show measurements;
name: measurements
name
----
events

> select * from events;
name: events
time                descr  tags        title
----                -----  ----        -----
1507423169000000000 aDesc1 aTag1,aTag2 aTitle1

```

As you can see, the annotations database was created, the events measurement
created and an entry created with our values.  Remember that `tags` is a field
and not InfluxDB tags.


### Create an entry that includes InfluxDB tags

This example includes an InfluxDB tag along with the other information.

```
gannotate -t aTitle2 -d aDesc2 -a aTag3,aTag4 -T host=host2.domain.com
```

The results.
```
> select * from events;
name: events
time                descr  host             tags        title
----                -----  ----             ----        -----
1507423169000000000 aDesc1                  aTag1,aTag2 aTitle1
1507423577000000000 aDesc2 host2.domain.com aTag3,aTag4 aTitle2
```

We can use the InfluxDB tag we created to filter our annotations.

> select * from events where host='host2.domain.com'
name: events
time                descr  host             tags        title
----                -----  ----             ----        -----
1507423577000000000 aDesc2 host2.domain.com aTag3,aTag4 aTitle2


## Common error messages

### `2017/10/09 16:02:59 unable to parse authentication credentials`

The InfluxDB server uses authentication and none was provided.
Call `gannotate` with the `-U` and `-P` options.

### `2017/10/09 16:03:22 authorization failed`

The InfluxDB server uses authentication and the provided username
and/or password were incorrect.  

## Grafana

Grafana cannot display InfluxDB tags in annotations.  Instead it looks for
tags in a dedicated field in a comma separated format.  This is why
`gannotate` puts the annotation tags in a field.

The following example is a query to retrieve all the annotations from InfluxDB.

```
SELECT title, descr, tags from events WHERE $timeFilter order by asc
```


## Development

### Building

You can use the Makefile or the standard `go get; go install` combination.

Go versions 1.8.x and 1.9.x has both been successfully used to build `gannotate`.

### Dependencies

Requires the InfluxDB Client V2 available at `github.com/influxdata/influxdb/client/v2`.

