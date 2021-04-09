package main

import (
	"regexp"
	"strconv"
	"strings"

	sysctl "github.com/lorenzosaino/go-sysctl"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type Exporter struct {
	grepPattern string
	skipPattern string
	prefix      string
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	rawSysctls, err := sysctl.GetAll()
	if err != nil {
		log.Error(err)
		return
	}
	for sysctlName, value := range rawSysctls {
		if sysctlNameIsFiltered(sysctlName, e.grepPattern, e.skipPattern) {
			continue
		}
		values := strings.Split(value, "\t")
		if len(values) == 1 {
			parsed, err := strconv.ParseFloat(values[0], 64)
			if err != nil {
				log.Debugf("%s value is not integer, skipped", sysctlName)
				continue
			}
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(prometheus.BuildFQName(e.prefix, "", "parameter"), "Values of sysctl parameters", []string{"param"}, nil),
				prometheus.GaugeValue, parsed, sysctlName,
			)
			continue
		}
		for i, value := range values {
			parsed, err := strconv.ParseFloat(value, 64)
			if err != nil {
				log.Debugf("%s value %d is not integer, skipped", sysctlName, i)
				continue
			}
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(prometheus.BuildFQName("sysctl", "", "parameter"), "Values of sysctl parameters", []string{"param", "column"}, nil),
				prometheus.GaugeValue, parsed, sysctlName, strconv.Itoa(i),
			)

		}
	}
}

func sysctlNameIsFiltered(sysctlName string, grepPattern string, skipPattern string) bool {
	matched, err := regexp.MatchString(grepPattern, sysctlName)
	if err != nil || matched == false {
		return true
	}
	if skipPattern == "" {
		return false
	}
	matched, err = regexp.MatchString(skipPattern, sysctlName)
	if err != nil || matched == true {
		return true
	}
	return false
}
