package main

import (
	"context"
	"fmt"
	"image"
	"math"

	"github.com/edaniels/golog"
	"github.com/golang/geo/r3"
	"github.com/pkg/errors"
	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/robot"
	"go.viam.com/rdk/spatialmath"

	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/vision"
	viz "go.viam.com/rdk/vision"
	"go.viam.com/rdk/vision/classification"
	"go.viam.com/rdk/vision/objectdetection"
	goutils "go.viam.com/utils"
)

var (
	Model            = resource.NewModel("viamlabs", "service", "obstacle-detection")
	errUnimplemented = errors.New("unimplemented")
	API              = resource.APINamespaceRDK.WithServiceType("vision")
	obstacle         = r3.Vector{}
)

func newObstacleService(deps resource.Dependencies, conf resource.Config, logger golog.Logger) (vision.Service, error) {
	service := &ObstacleService{logger: logger}
	err := service.Reconfigure(context.Background(), deps, conf)
	return service, err
}

type ObstacleServiceConfig struct {
	MaxDistance  *float64 `json:"max_distance_mm"`
	ZeroPosition *float64 `json:"zero_position_mm"`
	Camera       string   `json:"camera"`
}

// Validates JSON configuration
func (cfg *ObstacleServiceConfig) Validate(path string) ([]string, error) {
	if cfg.MaxDistance == nil {
		return nil, fmt.Errorf(`expected "max_distance_mm" attribute for obstacle-service %q`, path)
	}
	if cfg.ZeroPosition == nil {
		return nil, fmt.Errorf(`expected "zero_position_mm" attribute for obstacle-service %q`, path)
	}
	if cfg.Camera == "" {
		return nil, fmt.Errorf(`expected "camera" attribute for obstacle-service %q`, path)
	}

	return []string{cfg.Camera}, nil
}

// Handles attribute reconfiguration
func (service *ObstacleService) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	obstacleConfig, err := resource.NativeConfig[*ObstacleServiceConfig](conf)
	if err != nil {
		return errors.New("Could not assert proper config for obstacle-detection")
	}

	service.maxDistance = obstacleConfig.MaxDistance
	service.zeroPosition = obstacleConfig.ZeroPosition
	service.camera, err = camera.FromDependencies(deps, obstacleConfig.Camera)
	if err != nil {
		return err
	}

	return nil
}

// Attributes of the service
type ObstacleService struct {
	maxDistance  *float64
	zeroPosition *float64
	camera       camera.Camera
	robot        robot.Robot
	logger       golog.Logger
}

// Implement the methods the Viam RDK defines for the vision service API (rdk:service:vision)
func (service *ObstacleService) DetectionsFromCamera(ctx context.Context, cameraName string, extra map[string]interface{}) ([]objectdetection.Detection, error) {
	return nil, errUnimplemented
}

func (service *ObstacleService) Detections(ctx context.Context, img image.Image, extra map[string]interface{}) ([]objectdetection.Detection, error) {
	return nil, errUnimplemented
}

func (service *ObstacleService) Name() resource.Name {
	return service.Name()
}

// SetVelocity: unimplemented
func (service *ObstacleService) ClassificationsFromCamera(
	ctx context.Context,
	cameraName string,
	n int,
	extra map[string]interface{},
) (classification.Classifications, error) {
	return nil, errUnimplemented
}

func (service *ObstacleService) Classifications(
	ctx context.Context,
	img image.Image,
	n int,
	extra map[string]interface{},
) (classification.Classifications, error) {
	return nil, errUnimplemented
}
func getDistance(p r3.Vector) float64 {
	return math.Pow(math.Pow(p.X, 2)+math.Pow(p.Y, 2)+math.Pow(p.X, 2), 0.5)
}

func (service *ObstacleService) checkForObstacles(p r3.Vector, d pointcloud.Data) bool {
	if p.X < *service.zeroPosition {
		return true
	}

	distance := getDistance(p)
	if distance < *service.maxDistance {
		obstacle = r3.Vector{X: p.X, Y: p.Y, Z: p.Z}
		return false
	}
	return true
}

func (service *ObstacleService) GetObjectPointClouds(ctx context.Context, cameraName string, extra map[string]interface{}) ([]*viz.Object, error) {
	obstacle = r3.Vector{}
	currentPointcloud, err := service.camera.NextPointCloud(ctx)
	if err != nil {
		return nil, err
	}
	currentPointcloud.Iterate(0, 0, service.checkForObstacles)
	pc := pointcloud.New()
	pc.Set(obstacle, nil)
	geo := spatialmath.NewPoint(obstacle, "obstacle")

	vizObject, err := viz.NewObject(pc)
	if err != nil {
		return nil, err
	}
	vizObject.Geometry = geo

	return []*viz.Object{vizObject}, nil
}

func (service *ObstacleService) Close(ctx context.Context) error {
	return errUnimplemented
}
func (service *ObstacleService) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, errUnimplemented
}

func main() {
	goutils.ContextualMain(mainWithArgs, golog.NewDevelopmentLogger("obstacle-service"))
}

func registerService() {
	resource.RegisterService(
		vision.API,
		Model,
		resource.Registration[vision.Service, *ObstacleServiceConfig]{
			Constructor: func(
				ctx context.Context,
				deps resource.Dependencies,
				conf resource.Config,
				logger golog.Logger,
			) (vision.Service, error) {
				return newObstacleService(deps, conf, logger)
			}})
}

func mainWithArgs(ctx context.Context, args []string, logger golog.Logger) (err error) {
	registerService()

	obstacleModule, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}

	err = obstacleModule.AddModelFromRegistry(ctx, vision.API, Model)
	if err != nil {
		panic(err)
	}

	err = obstacleModule.Start(ctx)
	defer obstacleModule.Close(ctx)

	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
