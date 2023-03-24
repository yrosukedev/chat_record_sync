package use_case

import "context"

type UseCase interface {
	Run(ctx context.Context) error
}
