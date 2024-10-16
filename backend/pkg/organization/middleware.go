package organization

import (
	"context"

	"github.com/neoxelox/kit"
)

var (
	KeyRequestOrganization kit.Key = kit.KeyBase + "request:organization"
)

func RequestOrganization(ctx context.Context) *Organization {
	return ctx.Value(KeyRequestOrganization).(*Organization) // nolint:forcetypeassert,errcheck
}
