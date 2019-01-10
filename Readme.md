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

### Preparing workspace

Preparing Space

```
cf login -u admin -p $CF_ADMIN_PASSWORD
cf create-org order-manager
cf target -o "order-manager"
cf create-space uat
cf target -o "order-manager" -s "uat"
```

Configure Domain names for exposing the applicaitons from virtual box based CF installation

```
cf create-domain order-manager cf.< Virtual box host IP >.nip.io
cf create-domain order-manager cf.< Your domain name suffix e.g    example.com  this needs  you to have a wildcard lb setup for this hostname.>
```

Configure Dev User

```
cf create-user marco polo
cf set-space-role marco "order-manager" uat  SpaceDeveloper
```

### App Management

```
git clone https://github.com/ajarv/go-app.git ~/workspace/go-app
cd ~/workspace/go-app

cf login -u marco -p polo
# See running apps
cf apps
```

#### 4.1 Base deployment

```
cf push orders-app -n orders

```

```bash
Showing health and status for app orders-app in org ajar_v / space uat as ajar_v@yahoo.com...
OK

requested state: started
instances: 1/1
usage: 32M x 1 instances
urls: orders.cfapps.io
last uploaded: Thu Jan 10 21:52:20 UTC 2019
stack: cflinuxfs2
```

Access the app at any of the `urls` e.g. http://orders.cfapps.io/

##### Scale the app

```
cf scale orders-app -i 2
cf app orders-app
```

Modify the app and redeploy

#### 4.2 Blue deployment

```bash
sed 's/white/blue/'  templates/layout.html.t > templates/layout.html
cf push -n batman
```

Access the app at any of the ``urls` e.g. http://batman.cfapps.io/ in several seconds _when it becomes available_

Problem with blue deployment updates the existing app instances and while updating he existing instances are evicted.

Lets try a blue-green deployment

#### 4.3 Green deployment

Lets deploy a new version of the app which will become a new app instance group.i.e. it wont overwrite the existing app instances

```bash
sed 's/white/green/'  templates/layout.html.t > templates/layout.html
cf push green-batman -n green-batman
```

See if you can access the app at http://green-batman.cfapps.io/

Lets map the existing app URL to new app version

```
cf map-route green-batman cfapps.io -n batman

```
