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
	"fmt"
	"stark/utils/log"
	"strconv"

	logging "github.com/op/go-logging"
	//"github.com/prometheus/common/log"
	"github.com/tinytub/zep-cli/zeppelin"

	"github.com/prometheus/client_golang/prometheus"
)

var logger = logging.MustGetLogger("collector")

const (
	namespace = "zeppelin"
)

// A ClusterUsageCollector is used to gather all the global stats about a given
// ceph cluster. It is sometimes essential to know how fast the cluster is growing
// or shrinking as a whole in order to zero in on the cause. The pool specific
// stats are provided separately.
type ZepClusterCollector struct {
	//conn      Conn
	addrs     []string
	MetaCount prometheus.Gauge

	NodeCount   prometheus.Gauge
	UpNodeCount prometheus.Gauge
	NodeUp      *prometheus.GaugeVec
	TableUsed   *prometheus.GaugeVec
	TableRemain *prometheus.GaugeVec
	TableQuery  *prometheus.GaugeVec
	TableQPS    *prometheus.GaugeVec
}

// NewClusterUsageCollector creates and returns the reference to ClusterUsageCollector
// and internally defines each metric that display cluster stats.
//func NewZepClusterCollector(conn Conn) *ZepClusterCollector {
func NewZepClusterCollector(addrs []string) *ZepClusterCollector {
	return &ZepClusterCollector{
		//conn: conn,
		addrs: addrs,
		MetaCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "meta_count",
				Help:      "zeppelin meta server count",
			}),
		NodeCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "node_count",
				Help:      "zeppelin node server count",
			}),
		UpNodeCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "up_node_count",
				Help:      "zeppelin up node server count",
			}),
		NodeUp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "NodeUp",
				Help:      "zeppelin node is up",
			},
			[]string{"node", "port"},
		),
		TableUsed: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "TableUsed",
				Help:      "zeppelin Table space used",
			},
			[]string{"table"},
		),
		TableRemain: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "TableRemain",
				Help:      "zeppelin Table space remain",
			},
			[]string{"table"},
		),
		TableQuery: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "TableQuery",
				Help:      "zeppelin Table Query",
			},
			[]string{"table"},
		),
		TableQPS: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "TableQPS",
				Help:      "zeppelin Table QPS",
			},
			[]string{"table"},
		),
	}
}

func (c *ZepClusterCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		c.MetaCount,
		c.NodeCount,
		c.UpNodeCount,
		c.NodeUp,
		c.TableUsed,
		c.TableRemain,
		c.TableQuery,
		c.TableQPS,
	}
}

func (c *ZepClusterCollector) collect() error {

	rawNodes, _ := zeppelin.ListNode(c.addrs)
	nodes := len(rawNodes)
	c.NodeCount.Set(float64(nodes))

	upnodes := 0
	for _, node := range rawNodes {
		if node.GetStatus() == 0 {
			c.NodeUp.WithLabelValues(node.Node.GetIp(), strconv.Itoa(int(node.Node.GetPort()))).Set(float64(node.GetStatus()))
			upnodes += 1
		} else {
			c.NodeUp.WithLabelValues(node.Node.GetIp(), strconv.Itoa(int(node.Node.GetPort()))).Set(float64(node.GetStatus()))
		}
	}
	c.UpNodeCount.Set(float64(upnodes))

	rawMetas, _ := zeppelin.ListMeta(c.addrs)
	metas := len(rawMetas.Followers) + 1
	c.MetaCount.Set(float64(metas))
	// listable --> space
	tablelist, _ := zeppelin.ListTable(c.addrs)
	fmt.Println(tablelist)

	for _, tablename := range tablelist.Name {
		used, remain, _ := zeppelin.Space(tablename, c.addrs)
		query, qps, _ := zeppelin.Stats(tablename, c.addrs)
		c.TableUsed.WithLabelValues(tablename).Set(float64(used))
		c.TableRemain.WithLabelValues(tablename).Set(float64(remain))
		c.TableQuery.WithLabelValues(tablename).Set(float64(query))
		c.TableQPS.WithLabelValues(tablename).Set(float64(qps))
	}

	return nil
}

// Describe sends the descriptors of each metric over to the provided channel.
// The corresponding metric values are sent separately.
func (c *ZepClusterCollector) Describe(ch chan<- *prometheus.Desc) {
	/*
		for _, metric := range c.metricsList() {
			ch <- metric.Desc()
		}
	*/
	for _, metric := range c.collectorList() {
		metric.Describe(ch)
	}
}

// Collect sends the metric values for each metric pertaining to the global
// cluster usage over to the provided prometheus Metric channel.
func (c *ZepClusterCollector) Collect(ch chan<- prometheus.Metric) {
	if err := c.collect(); err != nil {
		log.Error("[ERROR] failed collecting cluster usage metrics:", err)
		return
	}
	/*
		for _, metric := range c.metricsList() {
			ch <- metric
		}
	*/
	for _, metric := range c.collectorList() {
		metric.Collect(ch)
	}
}
