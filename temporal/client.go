package temporal

import (
	"os"

	"go.temporal.io/sdk/client"
)

func NewClient(options client.Options) (client.Client, error) {
	if options.HostPort == "" {
		options.HostPort = os.Getenv("TEMPORAL_GRPC_ENDPOINT")
	}

	// options.DataConverter = dataconverter.NewEncryptionDataConverter(
	// 	converter.GetDefaultDataConverter(),
	// 	dataconverter.DataConverterOptions{KeyID: os.Getenv("DATACONVERTER_ENCRYPTION_KEY_ID")},
	// )

	return client.NewClient(options)
}
