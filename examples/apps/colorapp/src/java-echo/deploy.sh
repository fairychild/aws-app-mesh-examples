#!/usr/bin/env bash
# vim:syn=sh:ts=4:sw=4:et:ai

set -ex



DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
ECR_URL="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com"
COLOR_JAVAECHO_IMAGE=${COLOR_JAVAECHO_IMAGE:-"${ECR_URL}/java-echo"}
# GO_PROXY=${GO_PROXY:-"https://proxy.golang.org"}
AWS_CLI_VERSION=$(aws --version 2>&1 | cut -d/ -f2 | cut -d. -f1)

ecr_login() {
    if [ $AWS_CLI_VERSION -gt 1 ]; then
        aws ecr get-login-password --region ${AWS_DEFAULT_REGION} | \
            docker login --username AWS --password-stdin ${ECR_URL}
    else
        $(aws ecr get-login --no-include-email)
    fi
}

describe_create_ecr_registry() {
    local repo_name=$1
    local region=$2
    aws ecr describe-repositories --repository-names ${repo_name} --region ${region} \
        || aws ecr create-repository --repository-name ${repo_name} --region ${region}
}

# build
docker build -t $COLOR_JAVAECHO_IMAGE ${DIR}

# push
ecr_login
describe_create_ecr_registry java-echo ${AWS_DEFAULT_REGION}
docker push $COLOR_JAVAECHO_IMAGE
