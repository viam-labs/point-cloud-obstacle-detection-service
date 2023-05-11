# Obstacle Detection Service Using Pointcloud data #

## Using the Module ##
In order to use the obstacle detection service first you must clone the repo, build the module and provide the executable path for the module to your robot.

Build the module

```
cd pointcloud-obstacle-detection-service
go build
```

Configuring the Module
The command `go build` from earlier will create a binary. Provide the **full** path to the binary in `executable_path` and whatever you would like for the name.
```
{
    "executable_path": "/usr/local/bin/rplidar-module",
    "name": "rplidar_module"
}
```

Configuring the Service
Of the below params the only ones that need to be changed are `max_distance_mm` and `zero_position_mm`. for `max_distance_mm` provide the furthest distance at which you would like to consider something an obstacle. For `zero_position_mm` provide the value at which you would like to filter obstacles out because they are too close to the lidar. For example if your robot pokes out in front of your lidar, this value would be the distance from the edge of the robot to the center of the lidar.


```
{
      "name": "obstacle-vision-service",
      "namespace": "rdk",
      "type": "vision",
      "attributes": {
        "camera": "rplidar",
        "max_distance_mm": 750,
        "zero_position_mm": 150
      },
      "depends_on": [
        "rplidar"
      ],
      "model": "viamlabs:service:obstacle-detection"
}
```

See [sample-code/sample-config-linux.json](https://github.com/viam-labs/pointcloud-obstacle-detection-service/blob/main/sample-code/sample-config-linux.json) for a sample config of the pointcloud obstacle detection service detecting obstacles from an rplidar. See [here](https://github.com/viamrobotics/rplidar)  repo for more information on the rplidar