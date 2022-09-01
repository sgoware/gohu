#!/bin/bash

docker_names=('oauth-api' 'oauth-rpc-token-enhancer' 'oauth-rpc-token-store' \
'user-api' 'user-rpc-crud' 'user-rpc-info' 'user-rpc-vip' 'notification-api' \
'notification-rpc-crud' 'notification-rpc-info' 'mq-asynq-scheduler' 'mq-asynq-processor' \
'mq-nsq-consumer')

function docker_build() {
  if [ "$1" -ef "" ]; then
    return 0
  fi

  array=$(echo "$1" | tr '-' '\n')
  path='./app/service'
  for var in $array
  do
    path="${path}""/""${var}"
  done
  docker build -t "$PROJECT_NAME""_""$1" "${path}"
  return 1
}

export PROJECT_NAME=$1

echo "${PROJECT_NAME}"

cd /www/site/"$PROJECT_NAME" || exit

echo "docker_images: ""${#docker_names[@]}"

for docker_name in ${docker_names[*]}
do
{
  docker_build "${docker_name}"
  echo "build ""${docker_name}"" complete"
} &
done

wait