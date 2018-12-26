## Web App in Go

This is a sample web application developed in Go lang that can be run on various cloud /native platforms:

* Native
* Docker
* OpenShift
* Cloud Foundary

Features - Bootstrap4 / Go Template but bare bones and none necessary.

## 0. Clone Project

```
git clone https://github.com/ajarv/go-app-docker.git
```

## 1. Run on Local Host

```
cd go-app-docker
run app.go
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
oc new-project gosham-city --display-name "Gosham City"

```
### 3.1 Deploy on OpenShift plain directly


```
oc new-app <this git repo url>
```

### 3.2 Deploy on openshift using template

```
oc process -f  kube-cfg/openshift-app-template.yaml -p APP_NAME=joker  | oc create -f -
```


## 4 Cloud Foundary

### Cloudfoundry access
If you don't have access to cloudfoundry you can follow the tutorial [here](http://operator-workshop.cloudfoundry.org/agenda/) to provision it on a local linux box.


#### 4.1  White deployment
```
cd go-app-docker
cf login ..
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
Access the app at any of the ```urls`` e.g.  joker.cfapps.io 

Modify the app and redeploy
#### 4.2  Blue deployment

```bash
sed 's/white/blue/'  templates/layout.html.t > templates/layout.html
cf push -n joker  
```
Access the app at any of the ```urls`` e.g.  joker.cfapps.io  *when it becomes available*

Problem with blue dep is that while its updating the existing app becomes unavailable
Lets try a green deployment
#### 4.3  Green deployment

```bash
sed 's/white/green/'  templates/layout.html.t > templates/layout.html
cf push green-joker -n green-joker  
```




