import java.text.SimpleDateFormat


podTemplate(label: 'golang-1.9', containers: [
    containerTemplate(name: 'golang', image: 'golang:1.9', ttyEnabled: true, command: 'cat', args: ''),
    containerTemplate(name: 'docker', image: 'docker:17.09', ttyEnabled: true, command: 'cat', args: '' )
    ],
    volumes: [hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock')]
  )

// properties [[$class: 'BuildDiscarderProperty', strategy: [$class: 'LogRotator', daysToKeepStr: '10', numToKeepStr: '15']], [$class: 'ScannerJobProperty', doNotScan: false]]

{

    def image = "kubeam"
    def repo = "sample-registry"
    node('golang-1.9') {
        def dateFormat = new SimpleDateFormat("yyyyMMddHHmm")
        def date = new Date()
        def tag = dateFormat.format(date)

        stage('Build kubeAM Docker') {
            git url: 'https://github.com/llimon/kubeam.git', credentialsId: "ecr:us-west-2:-put-your-own-"
            //, credentialsId: "llimon", branch: master
            container('golang') {
                stage('Build kubeAM') {
                    sh 'bash build.sh none'
                }
                stage('Download kubectl') {
                    sh "curl -L -o kubectl.linux https://storage.googleapis.com/kubernetes-release/release/v1.9.0/bin/linux/amd64/kubectl"
                }
            }

            container('docker') {
                stage( 'Build docker image') {
                    sh "docker build . -f Dockerfile-kubeam.dkr -t ${repo}/${image}:${tag}"
                }
                
            }
        }  
        stage( "Publish") {
            container("docker"){
                stage("->preprod"){
                    docker.withRegistry("https://${repo}", "ecr:us-west-2:your-registry-") {
                        sh "docker push ${repo}/${image}:${tag}"
                    }
                }
            }
        }
        stage("Deploy") {
            container("golang") {
                stage( "->preprod"){
                    withCredentials([string(credentialsId: 'secret_admin', variable: 'KUBEAM_SECRET')]) {                 
                        sh "curl -k -X POST https://admin:${KUBEAM_SECRET}@kubeam.-my-url-:443/v1/deploy/sample/qa/main/${tag}"
                    }
                }
            }
        }

        stage( "Cleanup") {
            container("docker") {
                sh "docker rmi ${repo}/${image}:${tag}"
            }

        }
    }
}
