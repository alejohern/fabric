/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package nwo

// BasicConfig is a configuration with two organizations and one peer per org.
// This configuration does not specify a consensus type.
func BasicConfig() *Config {
	return &Config{
		Organizations: []*Organization{{
			Name:          "OrdererOrg",
			MSPID:         "OrdererMSP",
			Domain:        "example.com",
			EnableNodeOUs: false,
			Users:         0,
			CA:            &CA{Hostname: "ca"},
		}, {
			Name:          "Org1",
			MSPID:         "Org1MSP",
			Domain:        "org1.example.com",
			EnableNodeOUs: true,
			Users:         2,
			CA:            &CA{Hostname: "ca"},
		}, {
			Name:          "Org2",
			MSPID:         "Org2MSP",
			Domain:        "org2.example.com",
			EnableNodeOUs: true,
			Users:         2,
			CA:            &CA{Hostname: "ca"},
		}},
		Consortiums: []*Consortium{{
			Name: "SampleConsortium",
			Organizations: []string{
				"Org1",
				"Org2",
			},
		}},
		Consensus: &Consensus{
			BootstrapMethod: "file",
		},
		SystemChannel: &SystemChannel{
			Name:    "systemchannel",
			Profile: "TwoOrgsOrdererGenesis",
		},
		Orderers: []*Orderer{
			{Name: "orderer", Organization: "OrdererOrg"},
		},
		Channels: []*Channel{
			{Name: "testchannel", Profile: "TwoOrgsChannel"},
		},
		Peers: []*Peer{{
			Name:         "peer0",
			Organization: "Org1",
			Channels: []*PeerChannel{
				{Name: "testchannel", Anchor: true},
			},
		}, {
			Name:         "peer0",
			Organization: "Org2",
			Channels: []*PeerChannel{
				{Name: "testchannel", Anchor: true},
			},
		}},
		Profiles: []*Profile{
			{
				Name:     "TwoOrgsOrdererGenesis",
				Orderers: []string{"orderer"},
			},
			{
				Name:          "TwoOrgsChannel",
				Consortium:    "SampleConsortium",
				Organizations: []string{"Org1", "Org2"},
			},
		},
	}
}

// ThreeOrgEtcdRaft returns a simple configuration with three organizations instead
// of two.
func ThreeOrgEtcdRaft() *Config {
	config := BasicEtcdRaft()
	config.Organizations = append(
		config.Organizations,
		&Organization{
			Name:   "Org3",
			MSPID:  "Org3MSP",
			Domain: "org3.example.com",
			Users:  2,
			CA:     &CA{Hostname: "ca"},
		},
	)
	config.Consortiums[0].Organizations = append(
		config.Consortiums[0].Organizations,
		"Org3",
	)
	config.SystemChannel.Profile = "ThreeOrgsOrdererGenesis"
	config.Channels[0].Profile = "ThreeOrgsChannel"
	config.Peers = append(
		config.Peers,
		&Peer{
			Name:         "peer0",
			Organization: "Org3",
			Channels: []*PeerChannel{
				{Name: "testchannel", Anchor: true},
			},
		},
	)
	config.Profiles = []*Profile{{
		Name:     "ThreeOrgsOrdererGenesis",
		Orderers: []string{"orderer"},
	}, {
		Name:          "ThreeOrgsChannel",
		Consortium:    "SampleConsortium",
		Organizations: []string{"Org1", "Org2", "Org3"},
	}}

	return config
}

// FullEtcdRaft is a configuration with two organizations and two peers per org.
func FullEtcdRaft() *Config {
	config := BasicEtcdRaft()

	config.Peers = append(
		config.Peers,
		&Peer{
			Name:         "peer1",
			Organization: "Org1",
			Channels: []*PeerChannel{
				{Name: "testchannel", Anchor: false},
			},
		},
		&Peer{
			Name:         "peer1",
			Organization: "Org2",
			Channels: []*PeerChannel{
				{Name: "testchannel", Anchor: false},
			},
		},
	)

	return config
}

func BasicEtcdRaftWithIdemix() *Config {
	config := BasicEtcdRaft()

	// Add idemix organization
	config.Organizations = append(config.Organizations, &Organization{
		Name:          "Org3",
		MSPID:         "Org3MSP",
		MSPType:       "idemix",
		Domain:        "org3.example.com",
		EnableNodeOUs: false,
		Users:         0,
		CA:            &CA{Hostname: "ca"},
	})
	// Add idemix organization to consortium
	config.Consortiums[0].Organizations = append(config.Consortiums[0].Organizations, "Org3")
	config.Profiles[1].Organizations = append(config.Profiles[1].Organizations, "Org3")

	return config
}

func BasicEtcdRaft() *Config {
	config := BasicConfig()

	config.Consensus.Type = "etcdraft"
	config.Profiles = []*Profile{{
		Name:     "SampleDevModeEtcdRaft",
		Orderers: []string{"orderer"},
	}, {
		Name:          "TwoOrgsChannel",
		Consortium:    "SampleConsortium",
		Organizations: []string{"Org1", "Org2"},
	}}
	config.SystemChannel.Profile = "SampleDevModeEtcdRaft"

	return config
}

func MinimalRaft() *Config {
	config := BasicEtcdRaft()

	config.Peers[1].Channels = nil
	config.Channels = []*Channel{
		{Name: "testchannel", Profile: "OneOrgChannel"},
	}
	config.Profiles = []*Profile{{
		Name:     "SampleDevModeEtcdRaft",
		Orderers: []string{"orderer"},
	}, {
		Name:          "OneOrgChannel",
		Consortium:    "SampleConsortium",
		Organizations: []string{"Org1"},
	}}

	return config
}

func MultiChannelEtcdRaft() *Config {
	config := BasicConfig()

	config.Consensus.Type = "etcdraft"
	config.Channels = []*Channel{
		{Name: "testchannel", Profile: "TwoOrgsChannel"},
		{Name: "testchannel2", Profile: "TwoOrgsChannel"},
	}

	for _, peer := range config.Peers {
		peer.Channels = []*PeerChannel{
			{Name: "testchannel", Anchor: true},
			{Name: "testchannel2", Anchor: true},
		}
	}

	config.Profiles = []*Profile{{
		Name:     "SampleDevModeEtcdRaft",
		Orderers: []string{"orderer"},
	}, {
		Name:          "TwoOrgsChannel",
		Consortium:    "SampleConsortium",
		Organizations: []string{"Org1", "Org2"},
	}}
	config.SystemChannel.Profile = "SampleDevModeEtcdRaft"

	return config
}

func MultiNodeEtcdRaft() *Config {
	config := BasicEtcdRaft()
	config.Orderers = []*Orderer{
		{Name: "orderer1", Organization: "OrdererOrg"},
		{Name: "orderer2", Organization: "OrdererOrg"},
		{Name: "orderer3", Organization: "OrdererOrg"},
	}
	config.Profiles = []*Profile{{
		Name:     "SampleDevModeEtcdRaft",
		Orderers: []string{"orderer1", "orderer2", "orderer3"},
	}, {
		Name:          "TwoOrgsChannel",
		Consortium:    "SampleConsortium",
		Organizations: []string{"Org1", "Org2"},
	}}
	return config
}

func MultiNodeBFT() *Config {
	config := BasicConfig()

	config.Consensus.Type = "BFT"
	config.Orderers = []*Orderer{
		{Name: "orderer1", Organization: "OrdererOrg"},
		{Name: "orderer2", Organization: "OrdererOrg"},
		{Name: "orderer3", Organization: "OrdererOrg"},
	}
	config.Profiles = []*Profile{
		{
			Name:     "SampleDevModeBFT",
			Orderers: []string{"orderer1", "orderer2", "orderer3"},
		},
		{
			Name:          "TwoOrgsChannel",
			Consortium:    "SampleConsortium",
			Organizations: []string{"Org1", "Org2"},
		},
	}
	config.SystemChannel.Profile = "SampleDevModeBFT"

	return config
}
