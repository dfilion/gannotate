# gannotate

`gannotate` is a command line tool for sending Granafa annotations to InfluxDB for storage.


## Usage Summary

```
Usage of gannotate:
        -D dbname    InfluxDB database name.
        -H URL       InfluxDB server URL.
        -U username  InfluxDB user name.
        -P password  InfluxDB user password.
        -T tags      InfluxDB tags.
        -M name      InfluxDB measurement name.
        -a tags      Annotation tags.
        -d descr     Annotation description.
        -t title     Annotation title.
        -v           Print version information then exit.

        -a, -d and -t are required.
```

It is not currently possible to specify the annotation's timestamp.
The sending host's time is currently used.

### Arguments

#### `-D dbname`
InfluxDB database name.
Default: `annotations`

#### `-H URL`
InfluxDB server URL.
Default: `http://localhost:8086`

#### `-U username`
InfluxDB user name.
Optional.

#### `-P password`
InfluxDB password for `username`.
Optional.

##### Note: `gannotate` does not prompt for a password.

#### `-T tags`
Comma separated list of key=value pairs used as InfluxDB tags.
Optional. No default.

#### `-M name`
Measurement name used to record the annotations.
Default: events

#### `-a tags`
Comma separated list of key=value pairs used as Grafana annotation tags.
Saved to the `tags` field.

#### `-d descr`
String used as the Grafana annotation description.
Saved to the `descr` field.

#### `-t title`
String used as the Grafana annotation title. 
Saved to the `title` field.

#### `-v`
Print the applications version information and exit.


## Example

### Create an annotation using the default InfluxDB settings.
```
gannotate -t MyTitle -d MyDescription -a Tag1,Tag2
```

#### What gets created
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
time                descr              tags        title
----                -----              ----        -----
1512337131000000000 MyDescription      Tag1,Tag2   MyTitle

```

This shows the `annotations` database along with the `events` measurement were
created and an entry inserted with our values.  Remember that `tags` is a field
and not an InfluxDB tag.


### Create an entry that includes InfluxDB tags

A example that includes the host name as an InfluxDB tag for later searching or filtering.

```
gannotate -t GTitle -d GDesc -a GTag1,GTag2 -T host=host2.domain.com
```

Searching the events measurment we find our new entry along with the new InfluxDB tag.
```
> select * from events;
name: events

time                descr         host             tags        title
----                -----         ----             ----        -----
1512337131000000000 MyDescription                  Tag1,Tag2   MyTitle
1512338133000000000 GDesc         host2.domain.com GTag1,GTag2 GTitle
```

We can use the InfluxDB tag to filter our annotations.
```
> select * from events where host='host2.domain.com'
name: events
time                descr         host             tags        title
----                -----         ----             ----        -----
1512338133000000000 GDesc         host2.domain.com GTag1,GTag2 GTitle
```


## Common error messages

### `2017/10/09 16:02:59 unable to parse authentication credentials`

The InfluxDB server uses authentication and none was provided.
Call `gannotate` with the `-U` and `-P` options.

### `2017/10/09 16:03:22 authorization failed`

The InfluxDB server uses authentication and the provided username
and/or password were incorrect.  


## Configuring Grafana

Grafana cannot display InfluxDB tags in annotations.  Instead it looks for
tags in a dedicated field in a comma separated format.  This is why
`gannotate` has separate options for InfluxDB tags and annotation tags.

The following example is a query to retrieve all the annotations from InfluxDB.

```
SELECT title, descr, tags from events WHERE $timeFilter order by time asc
```

#### Note
Recent Grafana versions do not allow mapping of the `title` field, it is expected to exist.
You must still map the `tags` and `Text` fields.


### Combining templates and annotations

You can use Grafana templates to filter which annotations are selected.
The following example is based on annotations having been created with the
source host name as an InfluxDB tag.

#### Configuring the template

Create a Grafana template named `host` which uses the `host` tag key from the `events` measurement. 

```
show tag values from events with key = "host"
```

#### Configuring the annotations

Create an annotation that uses the InfluxDB `annotations` database with the following query.

```
SELECT title, descr, tags from events WHERE $timeFilter AND host =~ /$host$/ order by time asc
```

## Development

### Building

You can use the Makefile or the standard `go get; go install` combination.

Go versions 1.7.x, 1.8.x and 1.9.x has both been successfully used to build `gannotate`.

### Dependencies

Requires the InfluxDB Client V2 available at `github.com/influxdata/influxdb/client/v2`.

You may use `go get` to download the dependencies for you.

