package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddAllGroup_AssignsPublicIDs(t *testing.T) {
	account := &Account{
		Id:     "account-public-id",
		Groups: map[string]*Group{},
	}

	require.NoError(t, account.AddAllGroup(false))
	require.Len(t, account.Groups, 1)
	require.Len(t, account.Policies, 1)

	for _, group := range account.Groups {
		require.NotEmpty(t, group.PublicID, "bootstrap All group must have a public ID")
	}
	require.NotEmpty(t, account.Policies[0].PublicID, "bootstrap default policy must have a public ID")
}
