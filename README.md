# aws-lambda-copy-container-image-to-ecr

Copy container image to an ECR repository based on the [go-containerregistry](https://github.com/google/go-containerregistry.git) library

[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ekirmayer_aws-lambda-copy-container-image-to-ecr&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ekirmayer_aws-lambda-copy-container-image-to-ecr) 
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ekirmayer_aws-lambda-copy-container-image-to-ecr&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ekirmayer_aws-lambda-copy-container-image-to-ecr)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=ekirmayer_aws-lambda-copy-container-image-to-ecr&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=ekirmayer_aws-lambda-copy-container-image-to-ecr)

## Test

```shell
docker run -it --rm -v ~/.aws:/.aws -e AWS_DEFAULT_PROFILE=dev -p 9100:8080 \
        --entrypoint /usr/local/bin/aws-lambda-rie \
        docker-image:test ./main
```

```shell
curl --location --request GET 'http://localhost:9100/2015-03-31/functions/function/invocations' \
--header 'Content-Type: application/json' \
--data '{
    "body": "{\"src\":\"busybox:uclibc\",\"dest\":\"835389108797.dkr.ecr.eu-central-1.amazonaws.com/nginx:0.3\"}"
}'
```

## TO DO

- [ ] Add AWS IaC