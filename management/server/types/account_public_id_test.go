package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddAllGroup_AssignsPublicIDs(t *testing.T) {
	t.Run("group and default policy", func(t *testing.T) {
		account := &Account{
			Id:     "account-public-id",
			Groups: map[string]*Group{},
		}

		require.NoError(t, account.AddAllGroup(false))
		require.Len(t, account.Groups, 1)
		require.Len(t, account.Policies, 1)

		for _, group := range account.Groups {
			require.Equal(t, GroupAllName, group.Name)
			require.NotEmpty(t, group.PublicID)
			require.NotEqual(t, group.ID, group.PublicID)
		}

		require.NotEmpty(t, account.Policies[0].PublicID)
		require.NotEqual(t, account.Policies[0].ID, account.Policies[0].PublicID)
	})

	t.Run("group when default policy disabled", func(t *testing.T) {
		account := &Account{
			Id:     "account-public-id-no-policy",
			Groups: map[string]*Group{},
		}

		require.NoError(t, account.AddAllGroup(true))
		require.Len(t, account.Groups, 1)
		require.Empty(t, account.Policies)

		for _, group := range account.Groups {
			require.NotEmpty(t, group.PublicID)
		}
	})
}
