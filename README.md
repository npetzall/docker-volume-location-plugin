[![Travis](https://img.shields.io/travis/npetzall/docker-volume-location-plugin.svg?style=plastic)]() [![GitHub release](https://img.shields.io/github/release/npetzall/docker-volume-location-plugin.svg?style=plastic)]()  
# docker-volume-location-plugin
Simple Volume plugin for docker, so that you can put you volumes at a location of your choice

```
Usage: docker-volume-location-plugin [-location [alias=]/mnt/docker-volumes]
  -location value
    	[[alias=]path] can be declared multiple times
	omitting alias= sets default
  -profile
    	profile executions
  -version
    	Version of docker-volume-location-plugin
```

# docker
Usage is docker is pretty simple as well  

```
docker volume create -d vlp [-opt location=[alias]] [name]
docker run -v [name]:/mountpoint --volume-driver vlp [image] [command]
```
