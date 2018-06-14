brew update

MINIKUBE_VERSION="v0.25.2"
KUBECTL_VERSION="v1.9.0"
# Choises virtualbox, vmwarefusion, hyperkit, xhyve
#HYPERVISOR="xhyve"
#HYPERVISOR="vmwarefusion"
#HYPERVISOR="virtualbox"
HYPERVISOR="hyperkit"


function die()
{
    echo "ERROR: $1"
    exit 1
}


##
## Cisco Anyconnect Umbrella breaks minikube dns resolution
if [ -f /Library/LaunchDaemons/com.opendns.osx.RoamingClientConfigUpdater.plist ]; then
  if sudo launchctl list  | grep RoamingClientConfigUpdater; then
    echo "Dissabling Umbrella DNS for minikube startup"
    sudo launchctl unload /Library/LaunchDaemons/com.opendns.osx.RoamingClientConfigUpdater.plist
  fi
fi


##
## Delete cached configuration if chaning drivers
CURRENT_DRIVER=`cat ~/.minikube/machines/minikube/config.json |grep  DriverName |awk -F: '{print $2}' |tr -d "\",\ "`
if [ "$CURRENT_DRIVER" != "$HYPERVISOR" ]; then
    echo "INFO: Deleting .minikube configuration due to HYPERVISOR Change"
    rm -r -f ~/.minikube
fi
    

if [ "$HYPERVISOR" == "xhyve" ]; then
    echo "INFO: Installing xhyve"
    brew install --HEAD xhyve || die "Failed insalling xhyve"
    brew install docker-machine-driver-xhyve || die "Failed installing docker-machine-driver-xhyve"
    sudo chown root:wheel $(brew --prefix)/opt/docker-machine-driver-xhyve/bin/docker-machine-driver-xhyve
    sudo chmod u+s $(brew --prefix)/opt/docker-machine-driver-xhyve/bin/docker-machine-driver-xhyve
elif [ "$HYPERVISOR" == "hyperkit" ]; then
    curl -LO https://storage.googleapis.com/minikube/releases/latest/docker-machine-driver-hyperkit \
        && chmod +x docker-machine-driver-hyperkit \
        && sudo cp -fv docker-machine-driver-hyperkit /usr/local/bin/ \
        && sudo chown root:wheel /usr/local/bin/docker-machine-driver-hyperkit \
        && sudo chmod u+s /usr/local/bin/docker-machine-driver-hyperkit
fi



DOWNLOAD_MINIKUBE="false"
if [ -f /usr/local/bin/minikube ]; then
    INSTALLED_MINIKUBE_VERSION=`/usr/local/bin/minikube version |awk -F: '{print $2}'|tr -d "\ "`
    if [ "$INSTALLED_MINIKUBE_VERSION" != "$MINIKUBE_VERSION" ]; then
        echo "INFO: Backup pre-existing minikube"
        mv -vf /usr/local/bin/minikube /usr/local/bin/minikube.$(date +%Y%m%d%H%M%S)
        DOWNLOAD_MINIKUBE="true"
    fi
else
    DOWNLOAD_MINIKUBE="true"
fi
if [ "$DOWNLOAD_MINIKUBE" == "true" ]; then
    echo "INFO: Downloading minikube version $MINIKUBE_VERSION"
    curl -Lo minikube https://storage.googleapis.com/minikube/releases/$MINIKUBE_VERSION/minikube-darwin-amd64 && \
        chmod +x minikube && \
        sudo mv minikube /usr/local/bin/
    if [ $? -ne 0 ]; then
        die "Failed downloading minikube"
    fi
fi
if [ ! -f /usr/local/bin/kubectl ]; then
    echo "INFO: Downloading kubectl"
    curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/$KUBECTL_VERSION/bin/darwin/amd64/kubectl && \
        chmod +x kubectl && \
        sudo mv kubectl /usr/local/bin/
    if [ $? -ne 0 ]; then
        die "Failed downloading kubectl"
    fi
fi

echo "INFO: Running minikube $HYPERVISOR hypervisor"
if [ "$HYPERVISOR" == "xhyve" ]; then
    minikube start --logtostderr --v=10 --insecure-registry localhost:5000 --vm-driver=xhyve
elif [ "$HYPERVISOR" == "hyperkit" ]; then
   minikube start --logtostderr --v=10 --insecure-registry localhost:5000 --vm-driver=hyperkit
elif [ "$HYPERVISOR" == "virtualbox" ]; then
   minikube start --logtostderr --v=10 --insecure-registry localhost:5000 --vm-driver=virtualbox
elif [ "$HYPERVISOR" == "vmwarefusion" ]; then
   minikube start --logtostderr --insecure-registry localhost:5000 --vm-driver=vmwarefusion
else
   echo "ERROR: Unknown hypervisor selected. aborting."
fi

##
## Install registry
/usr/local/bin/kubectl config use-context minikube
echo "Waiting 10 seconds for minikube to be ready..."
sleep 10
/usr/local/bin/kubectl apply -f minikube-resitry.yaml 
