package utils

import (
	"context"
	"mesas-api/pkg/logging"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func ConfigAws(ctx context.Context, log *logging.Logger) aws.Config {
	configAws, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedCredentialsFiles([]string{"../pkg/driver/data/credentials.aws"}),
		config.WithSharedConfigFiles([]string{"../pkg/driver/data/config.aws"}),
	)
	if err != nil {
		log.HandleError("E", "Failed to load AWS configuration", err)
		panic(1)
	}

	return configAws
}
