## Web App in Go

This is a sample web application developed in Go lang that can be run on various cloud /native platforms:

- Native
- Docker
- OpenShift
- Cloud Foundary

Features - Bootstrap4 / Go Template but bare bones and none necessary.

## 0. Clone Project

```
git clone https://github.com/ajarv/go-app.git ~/workspace/go-app
```

## 1. Run on Local Host

```
cd ~/workspace/go-app
#Install Go dependencies
go get github.com/thedevsaddam/gojsonq  github.com/go-redis/redis gopkg.in/yaml.v2 github.com/gorilla/mux
go run main.go
```

Access the app using curl

```
curl localhost:8080
```

## 2. Create docker Image

Two phaze docker build

```sh
docker build -t go-app .
```

## 3. Open shift

Assuming you have openshift installed somewhere
And you are logged into openshift with oc.

```
oc login ...
oc new-project gotum-city --display-name "Gotum City"

```

### 3.1 Deploy on OpenShift plain directly

```
oc new-app https://github.com/ajarv/go-app.git
```

### 3.2 Deploy on openshift using template

A compiled image of this repo has been put on docker.io/m7dock/go-app
Below openshift template uses that

```
oc process -f  kube-cfg/openshift-app-template.yaml -p APP_NAME=gcpd -p C1_NAME=batman -p C2_NAME=bruce  -pLBVIP=10.231.63.150  | oc create -f -
```

## 4 Cloud Foundary

### Cloudfoundry access

If you don't have access to cloudfoundry you can follow the tutorial [here](http://operator-workshop.cloudfoundry.org/agenda/) to provision it on a local linux box.
OR
you may get a trial account at [https://pivotal.io/platform](https://pivotal.io/platform)

```
cf login ..
# See running apps
cf apps
```

#### 4.1 White deployment

```
cd go-app
cf create-space gosham-city
cf target -s "gosham-city"

cf push -n joker

```

```bash
Showing health and status for app go-app in org ../ space gosham-city as .......
OK

requested state: started
instances: 1/1
usage: 1G x 1 instances
urls: joker.cfapps.io
last uploaded: Wed Dec 26 17:08:48 UTC 2018
stack: cflinuxfs2
buildpack: https://github.com/cloudfoundry/go-buildpack.git
```

Access the app at any of the ``urls` e.g. http://joker.cfapps.io/

##### Scale the app

```
cf scale joker -i 2
cf app joker
```

Modify the app and redeploy

#### 4.2 Blue deployment

```bash
sed 's/white/blue/'  templates/layout.html.t > templates/layout.html
cf push -n joker
```

Access the app at any of the ``urls` e.g. http://joker.cfapps.io/ in several seconds _when it becomes available_

Problem with blue deployment updates the existing app instances and while updating he existing instances are evicted.

Lets try a blue-green deployment

#### 4.3 Green deployment

Lets deploy a new version of the app which will become a new app instance group.i.e. it wont overwrite the existing app instances

```bash
sed 's/white/green/'  templates/layout.html.t > templates/layout.html
cf push green-joker -n green-joker
```

See if you can access the app at http://green-joker.cfapps.io/

Lets map the existing app URL to new app version

```
cf map-route green-joker cfapps.io -n joker

```
