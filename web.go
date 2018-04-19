package dep_radar

import "context"

//go:generate mockery -name=IWebClient -case=underscore

// IWebClient is an interface for getting file
type IWebClient interface {
	Get(context.Context, string) ([]byte, error)
}
