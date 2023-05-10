# Obstacle Detection Service Using Pointcloud data #

## Using the Module ##
In order to use the obstacle detection service first you must clone the repo, build the module and provide the executable path for the module to your robot.

Build the module

```
cd pointcloud-obstacle-detection-service
go build
```

Provide the config for both the module and the service. See [sample-code/sample-config-linux.json](https://github.com/viam-labs/pointcloud-obstacle-detection-service/blob/main/sample-code/sample-config-linux.json) for a sample config