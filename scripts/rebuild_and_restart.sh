#!/bin/bash
IMAGE_NAME="tenseconds-go"
CONTAINER_NAME="tenseconds-go"

container_id=$(docker ps -aqf "name=${CONTAINER_NAME}")
if [ ! -z "${container_id}" ]; then
    docker stop ${CONTAINER_NAME}
    docker rm "${container_id}"
fi

# 이미지 빌드
docker build -t "${IMAGE_NAME}" .

# 새 컨테이너 실행
docker run -d --name "${CONTAINER_NAME}" -p 8001:8001 -v /home/ec2-user/.aws:/root/.aws "${IMAGE_NAME}"
