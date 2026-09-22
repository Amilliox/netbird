package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVAIB7334_AddAllGroupLeavesPublicIDEmpty(t *testing.T) {
	account := &Account{
		Id:     "account-vaib-7334",
		Groups: map[string]*Group{},
	}

	require.NoError(t, account.AddAllGroup(true))
	require.Len(t, account.Groups, 1)

	for _, group := range account.Groups {
		require.Equal(t, GroupAllName, group.Name)
		require.Empty(t, group.PublicID,
			"v0.77.1 AddAllGroup creates a production group after startup migration without assigning public_id")
	}
}
