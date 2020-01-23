!#/bin/bash

echo "Creating Self Signed Certificates"

# get the host IP of the machine
HOST_IP=`hostname -I | awk '{print $2}'` && echo $HOST_IP

# get the hostname of the machine
HOSTNAME=$(hostname -s)  && echo $HOSTNAME

openssl req -newkey rsa:2048 -nodes -keyout /etc/nginx/ssl/$HOST_IP.key -x509 -days 365 -out /etc/nginx/ssl/$HOST_IP.crt -subj "/C=NG/ST=Lagos/L=Oke-Odo/O=Samfil/OU=SD/CN=$HOST_IP,$HOSTNAME" 