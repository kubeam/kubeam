![ISaKongf](../../img/isakongf-logo.png)


Templating configuration utility used to setup docker images or Cloud application images. One a 
image or system is fully ready with all software installed. This application can take over and
configure a system in micro seconds! Since everything is "burned in". Variability is reduced to near 0


## Execution

Usage of ./ISaKonf :
  -v	Turn on verbose
  -vv
    	Turn on more verbose
    	
  Followed by configuration files to render:
  
## Sample:


` ./ISaKonf -v test-conf1.yaml test-conf2.yaml ...`

## Format of config yaml :


```description: Configuration files used to build a Application servers
templates:
  server.xml:
    template: server.xml.template
    destination: rendered/server.xml
    datasource: app.ini
    sourcetype: ini
    mode: 644
    action: create
  temp2.txt:
    source: /tmp/test2.txt
    destination: /tmp/test3.txt
    action: copy
    mode: 644
  temp.txt:
    destination: /tmp/test.txt
    action: delete
```

**action** can be [create], [copy] or [delete] 

[delete]

parameter:

 __destination__

[copy]

parameters:

 __source, destination__ optional __mode__

**sourcetype** Currently supports ini or environment variables, We will be adding support for env and possible consul in the near future. Also convine ( ini,env at same time) 

## Build

Build is easy just type:

` ./build.sh`

## Template files

Templates use a simple {{.key}} this is rendered with the value

Sample:
```
AppIdID={{.application_id}}
servlet.root=/main{{.cluster_number}}
servlet.zonename=main{{.cluster_number}}

# Akamai WAA
absolute.domain.enabled=true
www.dynamic.domain={{.application_name}}.{{.application_env}}.localhost
www.javascript.security_domain={{.application_name}}.{{.application_env}}.localhostj
```

### TODO

- [x] Create environment variable driver
  [ ] Real build (import dependencies)
  [ ] Create docker and publish to docker hub
  [ ] Better parameters (cobra/viper?)  --templates <template1> <template2> <template3>
