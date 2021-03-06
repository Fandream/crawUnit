FROM ubuntu:latest
WORKDIR .
COPY ./backend/ /app/backend/
COPY ./fe/dist /app/dist/
COPY ./docker_init.sh /app/
COPY ./backend/conf /conf
# set as non-interactive
ENV DEBIAN_FRONTEND noninteractive

# set CRAWLAB_IS_DOCKER
ENV CRAWLAB_IS_DOCKER Y

# install packages
RUN chmod 777 /tmp \
	&& apt-get update \
	&& apt-get install -y curl git net-tools iputils-ping ntp ntpdate python3 python3-pip nginx wget \
	&& ln -s /usr/bin/pip3 /usr/local/bin/pip \
	&& ln -s /usr/bin/python3 /usr/local/bin/python

# install dumb-init
RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.2/dumb-init_1.2.2_amd64
RUN chmod +x /usr/local/bin/dumb-init

# install backend
RUN pip install scrapy pymongo bs4 requests crawlab-sdk scrapy-splash


# copy nginx config files
COPY ./nginx/crawlab.conf /etc/nginx/conf.d


# timezone environment
ENV TZ Asia/Shanghai

# language environment
ENV LC_ALL C.UTF-8
ENV LANG C.UTF-8

# frontend port
EXPOSE 8080

# backend port
EXPOSE 8000

WORKDIR /app/backend/

# start backend
CMD ["/bin/bash", "/app/docker_init.sh"]
