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

import (
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// filterSetForTest builds a *schema.Set of {name,value} filter blocks with a
// content-based hash so distinct filters are not collapsed together.
func filterSetForTest(pairs ...[2]string) *schema.Set {
	hash := func(i interface{}) int {
		m := i.(map[string]interface{})
		return schema.HashString(m["name"].(string) + "|" + m["value"].(string))
	}
	items := make([]interface{}, 0, len(pairs))
	for _, p := range pairs {
		items = append(items, map[string]interface{}{"name": p[0], "value": p[1]})
	}
	return schema.NewSet(hash, items)
}

func TestApplyPodFiltersAreAndedNotOred(t *testing.T) {
	pod := &cloudstack.Pod{Name: "pod-a", Allocationstate: "Enabled"}

	// name matches but allocation_state does not: with AND semantics this must NOT match.
	filters := filterSetForTest([2]string{"name", "pod-a"}, [2]string{"allocation_state", "Disabled"})
	match, err := applyPodFilters(pod, filters)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if match {
		t.Errorf("pod matched only one of two filters but was selected; filters must be ANDed, not ORed")
	}

	// all filters match: must be selected.
	filters = filterSetForTest([2]string{"name", "pod-a"}, [2]string{"allocation_state", "Enabled"})
	match, err = applyPodFilters(pod, filters)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !match {
		t.Errorf("pod matched all filters but was not selected")
	}
}

func TestApplyClusterFiltersAreAndedNotOred(t *testing.T) {
	cluster := &cloudstack.Cluster{Name: "cluster-a", Allocationstate: "Enabled"}

	filters := filterSetForTest([2]string{"name", "cluster-a"}, [2]string{"allocation_state", "Disabled"})
	match, err := applyClusterFilters(cluster, filters)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if match {
		t.Errorf("cluster matched only one of two filters but was selected; filters must be ANDed, not ORed")
	}

	filters = filterSetForTest([2]string{"name", "cluster-a"}, [2]string{"allocation_state", "Enabled"})
	match, err = applyClusterFilters(cluster, filters)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !match {
		t.Errorf("cluster matched all filters but was not selected")
	}
}
