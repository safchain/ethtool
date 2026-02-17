//go:build !linux

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

// Package ethtool  aims to provide a library giving a simple access to the
// Linux SIOCETHTOOL ioctl operations. It can be used to retrieve informations
// from a network device like statistics, driver related informations or
// even the peer of a VETH interface.
package ethtool

// CmdGet returns the interface settings in the receiver struct
// and returns speed
func (ecmd *EthtoolCmd) CmdGet(intf string) (uint32, error) {
	e, err := NewEthtool()
	if err != nil {
		return 0, err
	}
	defer e.Close()
	return e.CmdGet(ecmd, intf)
}

// CmdSet sets and returns the settings in the receiver struct
// and returns speed
func (ecmd *EthtoolCmd) CmdSet(intf string) (uint32, error) {
	e, err := NewEthtool()
	if err != nil {
		return 0, err
	}
	defer e.Close()
	return e.CmdSet(ecmd, intf)
}

// CmdGet returns the interface settings in the receiver struct
// and returns speed
func (e *Ethtool) CmdGet(ecmd *EthtoolCmd, intf string) (uint32, error) {
	return 0, errOSUnsupported
}

// CmdSet sets and returns the settings in the receiver struct
// and returns speed
func (e *Ethtool) CmdSet(ecmd *EthtoolCmd, intf string) (uint32, error) {
	return 0, errOSUnsupported
}

// CmdGetMapped returns the interface settings in a map
func (e *Ethtool) CmdGetMapped(intf string) (map[string]uint64, error) {
	return nil, errOSUnsupported
}

// CmdGetMapped returns the interface settings in a map
func CmdGetMapped(intf string) (map[string]uint64, error) {
	return nil, errOSUnsupported
}
