FROM ubuntu:20.10

RUN apt-get update && \
apt-get install -y curl sudo zip
 
RUN curl https://raw.githubusercontent.com/terraform-linters/tflint/master/install_linux.sh | bash

WORKDIR /app

ENTRYPOINT ["tflint"]


