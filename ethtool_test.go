/*
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

package ethtool

import (
	"net"
	"reflect"
	"testing"
)

func TestGoString(t *testing.T) {
	testcases := map[string][]byte{
		"":      {},
		"0.5-4": {48, 46, 53, 45, 52, 0, 0, 0},
		"1.6":   {49, 46, 54, 0, 46, 48, 45, 50, 50, 45, 103, 101, 110, 101, 114, 105, 99},
	}

	for want, s := range testcases {
		got := goString(s[:])
		if got != want {
			t.Errorf("String conversion failed, got: %s, want: %s.", got, want)
		}
	}
}

func TestStats(t *testing.T) {
	intfs, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}

	// we expected to have at least one success
	empty := true
	for _, intf := range intfs {
		stats, err := Stats(intf.Name)
		if err != nil {
			continue
		}

		for range stats {
			empty = false
			break
		}
	}

	if empty {
		t.Fatal("Unable to retrieve stats from any interface of this system.")
	}
}

func TestDriverName(t *testing.T) {
	intfs, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}

	// we expected to have at least one success
	success := false
	for _, intf := range intfs {
		_, err := DriverName(intf.Name)
		if err == nil {
			success = true
		}
	}

	if !success {
		t.Fatal("Unable to retrieve driver from any interface of this system.")
	}
}

func TestBusInfo(t *testing.T) {
	intfs, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}

	// we expected to have at least one success
	success := false
	for _, intf := range intfs {
		_, err := BusInfo(intf.Name)
		if err == nil {
			success = true
		}
	}

	if !success {
		t.Fatal("Unable to retrieve bus info from any interface of this system.")
	}
}

func TestSupportedLinkModes(t *testing.T) {
	var cases = []struct {
		inputMask uint64
		expected  []string
	}{
		{0b01100010_11101111, []string{"10baseT_Half", "10baseT_Full", "100baseT_Half", "100baseT_Full", "1000baseT_Full"}},
	}

	for _, testcase := range cases {
		actual := SupportedLinkModes(testcase.inputMask)
		if !reflect.DeepEqual(actual, testcase.expected) {
			t.Error("Expected ", testcase.expected, " got ", actual)
		}
	}
}

func TestFeatures(t *testing.T) {
	et, err := NewEthtool()
	if err != nil {
		t.Fatal(err)
	}
	defer et.Close()

	feats, err := et.Features("lo")
	if err != nil {
		t.Fatal(err)
	}

	if len(feats) == 0 {
		// TOD0: do we have a sane subset of features we should check?
		t.Fatalf("expected features for loopback interface")
	}

	featsWithState, err := et.FeaturesWithState("lo")
	if err != nil {
		t.Fatal(err)
	}

	if len(feats) != len(featsWithState) {
		t.Fatalf("features mismatch: %d with state %d", len(feats), len(featsWithState))
	}

	fixed := 0
	for key, val := range feats {
		state, ok := featsWithState[key]
		if !ok || val != state.Active {
			t.Errorf("inconsistent feature: %q reported %v active %v", key, val, state.Active)
		}
		if !state.Available {
			fixed++
		}
	}

	if fixed == 0 {
		// the lo interface MUST have some non-available features, by design
		t.Fatalf("loopback interface reported all features available")
	}
}
