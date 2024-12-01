package metadata

import (
	"context"
	"fmt"

	"github.com/webbsalad/storya-passport-backend/internal/model"
	"google.golang.org/grpc/metadata"
)

func GetDeviceID(ctx context.Context) (model.DeviceID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return model.DeviceID{}, fmt.Errorf("missing metadata")
	}

	strDeviceID := md.Get("device_id")
	if len(strDeviceID) == 0 {
		return model.DeviceID{}, fmt.Errorf("missing device id")
	}

	deviceID, err := model.DeviceIDFromString(strDeviceID[0])
	if err != nil {
		return model.DeviceID{}, fmt.Errorf("convert str to device id: %w", err)
	}

	return deviceID, nil
}
