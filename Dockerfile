FROM ubuntu
MAINTAINER Rogier Lommers <rogier@lommers.org>
LABEL description="Poddy, because you need podcasts"

# install dependencies
RUN apt-get update  
RUN apt-get install -y ca-certificates

# add binary
COPY bin/poddy /poddy/bin/poddy
RUN mkdir /poddy/storage
RUN mkdir /poddy/watch

# change to data dir and run bianry
WORKDIR "/poddy"
CMD ["/poddy/bin/poddy"]
