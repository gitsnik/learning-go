#!/usr/bin/env bash

# Note: Adapted from `Securing DevOps: Security in the Cloud`
#
# Which is highly recommended. ISBN-13: 978-1617294136

docker pull owasp/zap2docker-weekly
VMIP=$(ip address show dev docker0 | grep 'inet ' | awk '{print $2}' | sed 's/\/.*//')
cd `dirname $0`
docker run --mount type=bind,source="$(pwd)",target=/zap/wrk --rm -t owasp/zap2docker-weekly zap-baseline.py -c zapbaseline.conf -t http://$VMIP
exit $?
