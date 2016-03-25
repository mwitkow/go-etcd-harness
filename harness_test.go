// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms covering this software.

package etcd_harness_test

import (
	"os"
	"testing"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/mwitkow/go-etcd-harness"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type HarnessTestSuite struct {

	keys etcd.KeysAPI
}

func (s *HarnessTestSuite) SetupSuite() {
	_, err := s.keys.Set(newContext(), "/testdir2", "", &etcd.SetOptions{Dir: true})
	require.NoError(s.T(), err, "creating the test directory must never fail.")
}

func (s *HarnessTestSuite) TestReadWrite() {
	_, err := s.keys.Set(newContext(), "/testdir/somevalue", "SomeContent", &etcd.SetOptions{})
	require.NoError(s.T(), err, "set must succeed")
	resp, err := s.keys.Get(newContext(), "/testdir/somevalue", &etcd.GetOptions{})
	require.NoError(s.T(), err, "get must succeed")
	assert.Equal(s.T(), "SomeContent", resp.Node.Value)
}

func TestHarnessTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skipf("HarnessTestSuite is a long integration test suite. Skipping due to test short.")
	}
	if !etcd_harness.LocalEtcdAvailable() {
		t.Skipf("etcd is not available in $PATH, skipping suite")
	}

	harness, err := etcd_harness.New(os.Stderr)
	if err != nil {
		t.Fatalf("failed starting etcd harness: %v", err)
	}
	t.Logf("will use etcd harness endpoint: %v", harness.Endpoint)
	defer func() {
		harness.Stop()
		t.Logf("cleaned up etcd harness")
	}()
	suite.Run(t, &HarnessTestSuite{keys: etcd.NewKeysAPI(harness.Client)})
}

func newContext() context.Context {
	c, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	return c
}
