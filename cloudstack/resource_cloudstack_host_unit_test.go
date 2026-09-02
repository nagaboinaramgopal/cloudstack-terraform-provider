// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package cloudstack

import "testing"

func TestResourceCloudStackHostHypervisorValidation(t *testing.T) {
	validate := resourceCloudStackHost().Schema["hypervisor"].ValidateFunc

	valid := []string{"xenserver", "kvm", "vmware", "baremetal", "simulator"}
	for _, v := range valid {
		if _, errs := validate(v, "hypervisor"); len(errs) != 0 {
			t.Errorf("supported hypervisor %q should be accepted, got errors: %v", v, errs)
		}
	}

	// Values that are not in the supported list must be rejected at plan time.
	invalid := []string{"foo", "docker", "esxi", "kvm2", ""}
	for _, v := range invalid {
		if _, errs := validate(v, "hypervisor"); len(errs) == 0 {
			t.Errorf("unsupported hypervisor %q should be rejected, but validation accepted it", v)
		}
	}
}
