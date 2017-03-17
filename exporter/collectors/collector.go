//   Copyright 2016 DigitalOcean
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package collectors

import (
	"github.com/prometheus/common/log"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "Zep"
)

// A ClusterUsageCollector is used to gather all the global stats about a given
// ceph cluster. It is sometimes essential to know how fast the cluster is growing
// or shrinking as a whole in order to zero in on the cause. The pool specific
// stats are provided separately.
type ZepClusterCollector struct {
	conn      Conn
	MetaCount prometheus.Gauge

	NodeCount prometheus.Gauge
}

// NewClusterUsageCollector creates and returns the reference to ClusterUsageCollector
// and internally defines each metric that display cluster stats.
func NewZepClusterCollector(conn Conn) *ZepClusterCollector {
	return &ZepClusterCollector{
		conn: conn,

		MetaCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "zep_meta_count",
			Help:      "zeppelin meta server count",
		}),
		NodeCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "zep_node_count",
			Help:      "zeppelin node server count",
		}),
	}
}

func (c *ZepClusterCollector) metricsList() []prometheus.Metric {
	return []prometheus.Metric{
		c.MetaCount,
		c.NodeCount,
	}
}

func (c *ZepClusterCollector) collect() error {

	/*
		conn := NewConn()
		//conn.mu.Lock()
		data, _ := conn.ListMeta()
		if data.Code.String() != "OK" {
			fmt.Println(*data.Msg)
			os.Exit(0)
	*/

	m, _ := c.conn.ListMeta()
	if m.Code.String() != "OK" {
		log.Error("get listmeta error: ", *m.Msg)
	}

	metas := len(m.ListMeta.Nodes.Followers) + 1

	n, _ := c.conn.ListNode()
	if n.Code.String() != "OK" {
		log.Error("get listnode error: ", *n.Msg)
	}
	nodes := len(n.ListNode.Nodes.Nodes)

	c.MetaCount.Set(float64(metas))
	c.NodeCount.Set(float64(nodes))
	return nil
}

// Describe sends the descriptors of each metric over to the provided channel.
// The corresponding metric values are sent separately.
func (c *ZepClusterCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metricsList() {
		ch <- metric.Desc()
	}
}

// Collect sends the metric values for each metric pertaining to the global
// cluster usage over to the provided prometheus Metric channel.
func (c *ZepClusterCollector) Collect(ch chan<- prometheus.Metric) {
	if err := c.collect(); err != nil {
		log.Error("[ERROR] failed collecting cluster usage metrics:", err)
		return
	}

	for _, metric := range c.metricsList() {
		ch <- metric
	}
}
