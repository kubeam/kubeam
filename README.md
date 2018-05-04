![kubeam-logo](img/kubeam-logo.png)

Kubernetes Application and Workflow manager

Integrates directly with your CI (Jenkins BO/CircleCI/Bamboo) and allows you to have end to end CI/CD.

![jenkins-bo-and-kubeam](img/kubeam-jenkins-bo.png)


### TODO:


- [x] We need to Add available cluster detection to kubeAM. Redis (not fully featured yet)
- [ ] Add samples of Jenkinsfile and circle.yaml
- [ ] Add TLS mutual cert auth. If possible add suport for multiple certs validation (one per account)
- [ ] Add Logic to parse application.json and return in structure. This will drive the logic of each rest endpoint.
- [ ] convert from using kubectl to cliento-go For better error handling, retry with backoff and return information to caller rest endpoint. We have a PR for this need input on pros and cons on doing client-go way.
- [x] Better build system. Look at isakconf for a sample.
- [ ] Reporting of cluster aviability.
- [x] getStatus should loop for each resource and generate a report. Right now is issuing a single call. If one resource is missing. Causes report to fail.
- [x] Health-check of kubeam     
- [x] Add support for includes to templates. This is needed for suporting source CIDR restrictions. 
- [x] Switch to go template from fast-template
- [ ] Rename provision api to something with more meaning is supposed to be endpoint for resources that timeout.
- [ ] Make /v1/provision work with generic applications using api.yaml file
- [ ] make /v1/waitforready generic driven by api.yaml
- [ ] Remove self-deploy endpoint. Now it can self deploy using api.yaml just like any other app.
- [ ] /v1/waitforready should look for multiple resources and wait for all of them. This is for complex applications

# API

Create a brand new cluster<br>
`curl  -X POST -k https://admin:{password}@localhost:8443/v1/create/sampleapp/QA/67/rhel7-1718-02/300`

Recreates (Delete then create) a resource using cluster selector; A cluster will be selected for you (Resources have TTL) <br>
`curl  -X POST -k https://admin:{password}@localhost:8443/v1/provision/sampleapp/QA/rhel7-170918-02`

Deploys to expecified <br>
`curl  -X POST -k https://admin:{password}@localhost:8443/v1/deploy/sampleapp/qa/rhel7-170918-02`


Gets status of a cluster<br>
`curl  -X GET -k https://admin:{password}@localhost:8443/v1/provision/sampleapp/QA/67`

List active clusters from Redis<br>
`curl  -X GET -k https://admin:{password}@localhost:8443/v1/listclusters/sampleapp/QA`

Wait for resource to be active<br>
`wget --no-check-certificate --user=admin --password=$SECRET https:/localhost:443/v1/waitforready/sampleapp/QA/66 -q -O -`
PS. wget due to limitations in how curl does timeout

Get detail on running cluster<br>
`curl  -X GET -k https://admin:${SECRET}@localhost:443/v1/getclusterdetail/sampleapp/QA/66`

Get detail on all running clusters<br>
`curl  -X GET -k https://admin:${SECRET}@localhost:443/v1/getallclustersdetail/sampleapp/QA`


Event Tracker

```
PAYLOAD="
{
   name: "Deployed",
   msg: "Passed all tests",
   ts: "2018-01-01 01:01:01"
}
curl -X POST -k https://admin:${SECRET}@localhost:443/v1/event/sampleapp/QA/66/{docker-TAG}
```
TODO: TAGS
```
PAYLOAD="
{
   name: "BATS-OK",
   msg: "PASSED sampleapp-BATS",
   ts: "2018-01-01 01:01:01"
}
curl -X POST -k https://admin:${SECRET}@localhost:443/v1/settag/sampleapp/QA/66/{docker-TAG}

KEYS are:  <APPLICATION>/<ENVIRONMENT>/<SHARD>/<DOCKER-TAG>
```


Other:

Self deployment<br>
`/v1/provision/self/{environment}/{tag}`

Health-check
`TCP:8081/health-check`




