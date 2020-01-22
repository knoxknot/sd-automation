!#/bin/bash

echo "Bootstrapping Kubernetes Cluster Master Node..."

# get the host IP of the master node
HOST_IP=`hostname -I | awk '{print $2}'` && echo $HOST_IP

# get the hostname of the master node
HOSTNAME=$(hostname -s)  && echo $HOSTNAME

# declare a pod network CIDR block
POD_CIDR='172.16.0.0/16' && echo $POD_CIDR

# initialize the cluster control plane on the master no
sudo kubeadm init --apiserver-advertise-address=$HOST_IP --apiserver-cert-extra-sans=$HOST_IP  --node-name=$HOSTNAME --pod-network-cidr=$POD_CIDR

# configure kube config
sudo \cp /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
export KUBECONFIG=/home/vagrant/.kube/config

# install weaveworks network addon
kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"

# make the master node workload schedulable
kubectl taint nodes $HOSTNAME node-role.kubernetes.io/master-

#  Wait 3 mins for node to be ready
sleep 180

# display the node details
kubectl get node