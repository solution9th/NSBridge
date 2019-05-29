HOST=172.16.120.103
APP=ns_bridge
PROJECTPATH=/production/api/src/gitlab.zlibs.com/${APP}

rsync -a . ${HOST}:${PROJECTPATH}
# rsync -a * ${HOST}:${PROJECTPATH} --exclude ./.git

ssh ${HOST} 'cat >  /production/jenkins_call.sh << EOF

APP=ns_bridge
PROJECT=gitlab.zlibs.com/\$APP
PROJECTPATH=/production/api/src/\$PROJECT
echo \$APP

export GOROOT=/usr/local/go
export PATH=\$PATH:\$GOROOT/bin
export GOPATH=/production/api

mkdir -p /etc/ns_bridge/keys/grpc
mkdir -p /var/log/ns_bridge

cd \$PROJECTPATH
ls -al
make build

cp -f \$APP /usr/local/bin/
cp -f \$PROJECTPATH/blueware-agent.production.ini /etc/ns_bridge/blueware-agent.ini
cp -f \$PROJECTPATH/config.dev.yml /etc/ns_bridge/config.yml
cp -f \$PROJECTPATH/keys/grpc/* /etc/ns_bridge/
cp -f \$APP.service /lib/systemd/system/
systemctl restart \$APP.service
EOF'

ssh ${HOST} 'cat /production/jenkins_call.sh | sh > /production/log.log'