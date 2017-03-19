# poddy
your custom podcasts feed

# features
- add new episoses by
  - read directory
  - upload form

# run container
`docker run -v /srv/services/poddy/storage:/poddy/storage -v /srv/dropbox/Apps/WatchDirectory:/poddy/watch -p 9002:8080 --name poddy rogierlommers/poddy`

# create container

First build binary `./release.sh`

Then push:

`docker build -t rogierlommers/poddy .`
`docker push rogierlommers/poddy:latest`
