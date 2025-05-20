package orderscl

import "context"

type Client interface {
	SetOrderComposition(ctx context.Context, in SetOrderCompositionIn) (err error)
}
