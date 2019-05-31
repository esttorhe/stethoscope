# stethoscope
External websites monitoring tool

## Description

This is a very naive service that will load a set of [`Rule`][rule.go]s from a configuration file ([`rules.yml`][rules.yml]) that will later use to monitor the «heartbeat» of the websites.

### Composition

#### Counter

A counter is the definition of a [`Prometheus` counter][prom_counter] that will be used when a monitored website is down.
They are loaded from [`counters.yml`][counters.yml] and are defined as so:

```yaml
   - namespace: website
     subsystem: health
     name: page_down
     help: Tracks the number of times there's an error loading a website
     labels: 
       - 'website'
       - 'status_code'
```

| Field       | Description                                                                                                                                                    |
| :---------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `namespace` | `Prometheus` namespace is an application prefix relevant to the domain the metric belongs to.[^1]                                                              |
| `subsystem` | `Prometheus` subsystem is used to prepend the counter; it is prepended after the `namespace`                                                                   |
| `name`      | The identification name for the counter. This will also be the mapping key that should match the counter specified in the rules.                               |
| `help`      | Description of what this counter is and what is it tracking.                                                                                                   |
| `labels`    | Array of labels that will be added to the counter. Right now it always expects this 2 `website` & `status_code` but will make them configurable in the future. |

#### Rule

A rule is a basic definition o what should be monitored and how.
It can be defined like this:

```yaml
  - name: name
    website: website
    counter: counter_name
    interval: 3600000000000 # Check every hour
    timeout: 2000000000 # Timeout after 2 seconds
    use_head: true
```
| Field      | Description                                                                                                                                                                             |
| :--------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `name`     | The name of your rule. Used mainly for logging purposes                                                                                                                                 |
| `website`  | An environment variable that contains the value of the website you want to monitor. This was done like this in order to `OSS` the project without exposing the websites I'm monitoring. |
| `counter`  | The name of the counter that should be used when the website is down. Should be a match with the `name` field in the one of the counters from `counters.yml`                            |
| `interval` | Number of nanoseconds to wait before attempting to monitor the website.                                                                                                                 |
| `timeout`  | Number of nanoseconds to wait for the website before considering it as down (due to a timetout)                                                                                         |
| `use_head` | If the monitoring should be attempted using a `HEAD` method. If `false` the service will use `GET`.                                                                                     |


## Docker

To make my life easier I created a [`Dockerfile`][docker] that compiles the service into a binary using a `Go` `1.12` image.
Once the binary is created it uses a custom `Alpine` image with `CA` certificates to wrap the binary and uses it as the base for the executing image.

All in all pretty simple but I find it quite comfortable to build on a container (which gives me confidence that it will work since technically it will always use the same environment for building).

### Running it

Once the image is ready that «recommended» way of executing ti would be like so:

```bash
docker run -dt --name=stethoscope -p 7000:7000 stethoscope:latest
```

This obviously if you are running it locally; otherwise you need to specify your registry.

## Makefile

Because I'm also lazy I created a `Makefile` whose sole purpose is to compile the service, generate the `Docker` image and upload it to the registry.
This way I can `ssh` into my server and just pull the image and launch a new container with the latest image. 

**Easy peasy!**

## Author
__Esteban Torres__ 

- [![](https://img.shields.io/badge/twitter-esttorhe-brightgreen.svg)](https://twitter.com/esttorhe) 
- ✉ me@estebantorr.es

## License

`stethoscope` is available under the MIT license. See the [LICENSE](LICENSE) file for more info.


[rule.go]:./rule.go
[rules.yml]:./rules.yml
[counters.yml]:./counters.yml
[prom_counter]:https://prometheus.io/docs/concepts/metric_types/#counter
[^1]:https://prometheus.io/docs/practices/naming/#metric-names
[docker]:./Dockerfile