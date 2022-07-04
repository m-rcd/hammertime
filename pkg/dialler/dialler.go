package dialler

import (
	"context"
	"encoding/base64"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// New process the dial config and returns a grpc.ClientConn. The caller is
// responsible for closing the connection.
// TODO pass in auth stuff as config.
func New(address, basicAuthToken string) (*grpc.ClientConn, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if basicAuthToken != "" {
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(
			Basic(basicAuthToken),
		))
	}

	return grpc.Dial(
		address,
		dialOpts...,
	)
}

// TODO move to separate package or file.
type basicAuth struct {
	token string
}

// Basic creates a basicAuth with a token.
func Basic(t string) basicAuth { //nolint: revive // this will not be used
	return basicAuth{token: t}
}

// GetRequestMetadata fullfills the credentials.PerRPCCredentials interface,
// adding the basic auth token to the request authorization header.
func (b basicAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	enc := base64.StdEncoding.EncodeToString([]byte(b.token))

	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

// GetRequestMetadata fullfills the credentials.PerRPCCredentials interface.
func (b basicAuth) RequireTransportSecurity() bool {
	return false
}
