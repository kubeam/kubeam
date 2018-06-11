* Write a github trigger/pull integration with KubeAM to download application yaml.
   1. Github commit to master triggers webhook to Kubeam (api to be written)
   2. KubeAM hook receiver api clones tiggered github repo and symlinks /kubeam/<folder-list> to /applications. 
      for cloning kubeAM will need a GitHub read only token.
      find kubeam/ -type d -exec ln -s {} /applications/{} \;
   3. We will save in MySQL or Redis the hash of triggered commit.  Alternative. Use persistance volume to save git clones

* Merge code from Intuit KubeAM to OSS KubeAM and if possible sunset Intuit version.
* Demo of KubeAM using Minikube
   - Installation script or program to configure all the needed components
   - bootstrap Kubeam with it self (this will work in other kubernetes clusters)
   - KubeAM application definition that works with AWS and ECR or minikube
   - Create a BlueOcean build environent with KubeAM in minicube and a pipeline for a simple application.
* Documentation on deploying KubeAM
* Create a "Official Docker or KubeAM " in dockerhub.io (this should be our last step because it makes it official)

Pending KubeAM Work:
 * Create a non-bloquing equivalent of WaitFor KubeAM endpoint.
 * Make WaitFor use api.yaml so that it can lookup for multiple resources configured in api.yaml currently only looks for application-{{.env}}-{.clulster}}
 * Trigger event messages during activities. (ie create event when we create, deploy, and calculate time it to deploy during WaitFor) also notify if resource timedout during WaitFor

