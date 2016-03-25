// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms covering this software.
// Package etcd_harness provides an integration test harness for running a local etcd server.

package etcd_harness

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"time"

	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

// Harness represents a running etcd server for an integraiton test environment.
type Harness struct {
	errWriter  io.Writer
	etcdServer *exec.Cmd
	etcdDir    string
	Client     etcd.Client
	Endpoint   string
}

// LocalEtcdAvailable returns true if an etcd binary is available on PATH.
func LocalEtcdAvailable() bool {
	_, err := exec.LookPath("etcd")
	return err == nil
}

// New initializes and returns a new Harness.
// Failures here will be indicated as errors.
func New(etcdErrWriter io.Writer) (*Harness, error) {
	s := &Harness{errWriter: etcdErrWriter}
	endpointAddress, err := allocateLocalAddress()
	if err != nil {
		return nil, fmt.Errorf("failed allocating endpoint addr: %v", err)
	}
	peerAddress, err := allocateLocalAddress()
	if err != nil {
		return nil, fmt.Errorf("failed allocating peer addr: %v", err)
	}
	etcdBinary, err := exec.LookPath("etcd")
	if err != nil {
		return nil, err
	}
	s.etcdDir, err = ioutil.TempDir("/tmp", "etcd_testserver")
	if err != nil {
		return nil, fmt.Errorf("failed allocating new dir: %v", err)
	}
	endpoint := "http://" + endpointAddress
	peer := "http://" + peerAddress
	s.etcdServer = exec.Command(
		etcdBinary,
		"--log-package-levels=etcdmain=WARNING,etcdserver=WARNING,raft=WARNING",
		"--force-new-cluster="+"true",
		"--data-dir="+s.etcdDir,
		"--listen-peer-urls="+peer,
		"--initial-cluster="+"default="+peer+"",
		"--initial-advertise-peer-urls="+peer,
		"--advertise-client-urls="+endpoint,
		"--listen-client-urls="+endpoint)
	s.etcdServer.Stderr = s.errWriter
	s.etcdServer.Stdout = ioutil.Discard
	s.Endpoint = endpoint
	if err := s.etcdServer.Start(); err != nil {
		s.Stop()
		return nil, fmt.Errorf("cannot start etcd: %v, will clean up", err)
	}
	s.Client, err = etcd.New(etcd.Config{Endpoints: []string{endpoint}})
	if err != nil {
		s.Stop()
		return s, fmt.Errorf("failed allocating client: %v, will clean up", err)
	}
	if err := s.pollEtcdForReadiness(); err != nil {
		s.Stop()
		return nil, fmt.Errorf("%v, will clean up", err)
	}
	return s, nil
}

func (s *Harness) pollEtcdForReadiness() error {
	api := etcd.NewKeysAPI(s.Client)
	// Actively poll for etcd coming up for 3 seconds every 50 milliseconds.
	for i := 0; i < 20; i++ {
		until := time.Now().Add(200 * time.Millisecond)
		ctx, _ := context.WithDeadline(context.TODO(), until)
		_, err := api.Get(ctx, "/", &etcd.GetOptions{})
		if err == nil {
			return nil
		}
		toSleep := until.Sub(time.Now())
		if toSleep > 0 {
			time.Sleep(toSleep)
		}
	}
	return fmt.Errorf("etcd didn't come up in 4000ms")
}

// Stop kills the harnessed etcd server and cleans up the data directory.
func (s *Harness) Stop() {
	var err error
	if s.etcdServer != nil {
		if err := s.etcdServer.Process.Kill(); err != nil {
			fmt.Printf("failed killing etcd process: %v", err)
		}
		// Just to make sure we actually finish it before continuing.
		s.etcdServer.Wait()
	}
	if s.etcdDir != "" {
		if err = os.RemoveAll(s.etcdDir); err != nil {
			fmt.Printf("failed clearing temporary dir: %v", err)
		}
	}
}

func allocateLocalAddress() (string, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer l.Close()
	return l.Addr().String(), nil
}
