# aws-lambda-copy-container-image-to-ecr

Copy container image to an ECR repository based on the [go-containerregistry](https://github.com/google/go-containerregistry.git) library

[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ekirmayer_aws-lambda-copy-container-image-to-ecr&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ekirmayer_aws-lambda-copy-container-image-to-ecr) 
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ekirmayer_aws-lambda-copy-container-image-to-ecr&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ekirmayer_aws-lambda-copy-container-image-to-ecr)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=ekirmayer_aws-lambda-copy-container-image-to-ecr&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=ekirmayer_aws-lambda-copy-container-image-to-ecr)

## Build

```shell
docker build -t lambda/go/tester .
docker tag lambda/go/tester:latest 404105784577.dkr.ecr.us-east-1.amazonaws.com/lambda/go/tester:test.1
docker push 404105784577.dkr.ecr.us-east-1.amazonaws.com/lambda/go/tester:test.1
```

## Test

```shell
docker run -it --rm -v ~/.aws:/.aws -e AWS_DEFAULT_PROFILE=dev -p 9100:8080 \
        --entrypoint /usr/local/bin/aws-lambda-rie \
        lambda/go/tester ./main
```

```shell
curl --location --request GET 'http://localhost:9100/2015-03-31/functions/function/invocations' \
--header 'Content-Type: application/json' \
--data '{
    "body": "{\"src\":\"busybox:uclibc\",\"dest\":\"835389108797.dkr.ecr.eu-central-1.amazonaws.com/nginx:0.3\"}"
}'
```

### Scan IaC
```shell
trivy config --skip-dirs ./iac/.terraform  .
```

```shell
trivy image lambda/go/tester
