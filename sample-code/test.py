import asyncio
import time
import random
import threading

from viam.robot.client import RobotClient
from viam.rpc.dial import Credentials, DialOptions
from viam.services.vision import VisionClient
from viam.components.base import Base
from viam.components.base import Vector3

ANGULAR_VELOCITY = 25
LINEAR_VELOCITY = 100
SECONDS_TO_RUN = 60 * 15

async def connect():
    creds = Credentials(
        type='robot-location-secret',
        payload='npd5y5ej86myjx44kkda5qq9fcptcnmznhfmniftfd7x0xem')
    opts = RobotClient.Options(
        refresh_interval=0,
        dial_options=DialOptions(credentials=creds)
    )
    return await RobotClient.at_address('vision-rover-main.l50o5rvufg.viam.cloud', opts)

def get_random_degrees():
    degrees = [30, 60, 90, 120, 150, 180]
    return degrees[random.randint(0, len(degrees)) - 1]

async def move_straight_and_avoid_obstacles(base: Base, obstacle_detection_service: VisionClient):
    obstacle = await obstacle_detection_service.get_object_point_clouds("rplidar")
    
    while len(obstacle) != 0:
            await base.stop()
            await base.spin(get_random_degrees(), ANGULAR_VELOCITY)
            obstacle = await obstacle_detection_service.get_object_point_clouds("rplidar")
    
    await base.set_velocity(linear=Vector3(x=0,y=LINEAR_VELOCITY,z=0), angular=Vector3(x=0,y=0,z=0))
     

async def main():    
    robot = await connect()
    base = Base.from_robot(robot, "viam_base")
    obstacle_detection_service = VisionClient.from_robot(robot, "obstacle-vision-service")


    t_end = time.time() + SECONDS_TO_RUN
    while time.time() < t_end:
        await move_straight_and_avoid_obstacles(base, obstacle_detection_service)

    await base.stop()
    await robot.close()

if __name__ == '__main__':
    asyncio.run(main())