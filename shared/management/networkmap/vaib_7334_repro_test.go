package networkmap_test

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	goproto "google.golang.org/protobuf/proto"

	mgmtgrpc "github.com/netbirdio/netbird/management/internals/shared/grpc"
	nbnetworkmap "github.com/netbirdio/netbird/shared/management/networkmap"
	"github.com/netbirdio/netbird/shared/management/networkmap/nmdata"
	"github.com/netbirdio/netbird/shared/management/proto"
)

// TestVAIB7334_RealisticAllMembership checks the missing-public-id hypothesis
// while preserving the built-in All-group invariant: the target peer belongs
// to All. This is the production-shaped counterpart to PR #7615's regression
// fixture, whose All group excludes the target peer.
func TestVAIB7334_RealisticAllMembership(t *testing.T) {
	c, localPeerKey := buildSmokeComponents(t)

	c.Groups = map[string]*nmdata.Group{
		"group-custom": {
			PublicID: "",
			Name:     "dns-clients",
			Peers:    []string{"peer-A"},
		},
		"group-all": {
			PublicID: "",
			Name:     "All",
			Peers:    []string{"peer-A", "peer-B"},
		},
	}
	c.NameServerGroups = []*nmdata.NameServerGroup{{
		ID:       "nsg-internal",
		PublicID: "",
		NameServers: []nmdata.NameServer{{
			IP:     c.Peers["peer-B"].IP,
			NSType: 1,
			Port:   5353,
		}},
		Groups:  []string{"group-custom"},
		Primary: true,
		Enabled: true,
	}}

	envelope := mgmtgrpc.EncodeNetworkMapEnvelope(mgmtgrpc.ComponentsEnvelopeInput{
		Components: c,
		DNSDomain:  "netbird.cloud",
	})

	// Force the same overwrite order used by PR #7615's regression:
	// the All group is decoded last for the duplicate empty wire ID.
	full := envelope.GetFull()
	require.Len(t, full.Groups, 2)
	sort.SliceStable(full.Groups, func(i, j int) bool {
		return !full.Groups[i].IsAll && full.Groups[j].IsAll
	})

	wire, err := goproto.Marshal(envelope)
	require.NoError(t, err)

	var decoded proto.NetworkMapEnvelope
	require.NoError(t, goproto.Unmarshal(wire, &decoded))

	result, err := nbnetworkmap.EnvelopeToNetworkMap(context.Background(), &decoded, localPeerKey, "netbird.cloud")
	require.NoError(t, err)
	require.True(t, result.NetworkMap.DNSConfig.ServiceEnable)
	require.Len(t, result.NetworkMap.DNSConfig.NameServerGroups, 1,
		"with a valid All group, duplicate empty group public IDs alone must not reproduce #7334")
}
