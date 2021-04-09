package main

import "testing"

func TestSysctlNameFilter(t *testing.T) {
	var res bool

	// default params, all parameters
	parameters := []string{"abc", "bcd", "cde", "a.b.c", ".a_s", "...", ".*"}
	pattern := ".*"
	skip_pattern := ""

	for _, param := range parameters {
		res = sysctlNameIsFiltered(param, pattern, skip_pattern)
		if res {
			t.Errorf("String '%s' did not pass '%s' + '%s'", param, pattern, skip_pattern)
		}
	}

	// simple match and skip
	parameters = []string{"a.b.c", "b.c.d", "bb", "c.d.e", "_", ".*"}
	answer := []bool{false, false, false, true, true, true}
	pattern = "b"
	for i, param := range parameters {
		res = sysctlNameIsFiltered(param, pattern, "")
		if res != answer[i] {
			t.Errorf("String '%s' with '%s' + '%s' returned %t", param, pattern, skip_pattern, res)
		}
		res = sysctlNameIsFiltered(param, "", pattern)
		if res == answer[i] {
			t.Errorf("String '%s' with '%s' + '%s' returned %t", param, skip_pattern, pattern, res)
		}
	}

	//  match and skip, complex
	parameters = []string{"a.b.c", "b.c.d", "bb", "c.d.e", "_e", ".*"}
	answer = []bool{true, true, false, true, false, true}
	pattern = "b|e"
	skip_pattern = "c"
	for i, param := range parameters {
		res = sysctlNameIsFiltered(param, pattern, skip_pattern)
		if res != answer[i] {
			t.Errorf("String '%s' with '%s' + '%s' returned %t", param, pattern, skip_pattern, res)
		}
	}

	//  match and skip, complex
	parameters = []string{
		"net.ipv4.tcp_tso_win_divisor",
		"net.ipv4.tcp_tw_reuse",
		"net.ipv4.tcp_window_scaling",
		"net.ipv4.tcp_wmem",
		"net.ipv6.tcp_wmem",
		"net.ipv4.tcp_workaround_signed_windows",
		"net.ipv4.udp_early_demux",
		"net.ipv4.udp_l3mdev_accept",
		"net.ipv4.udp_mem",
		"net.ipv4.udp_rmem_min",
		"net.ipv4.udp_wmem_min",
	}
	pattern = "ipv4.(udp|tcp).*mem"
	skip_pattern = "udp_mem"
	answer = []bool{
		true,
		true,
		true,
		false,
		true,
		true,
		true,
		true,
		true,
		false,
		false,
	}
	for i, param := range parameters {
		res = sysctlNameIsFiltered(param, pattern, skip_pattern)
		if res != answer[i] {
			t.Errorf("String '%s' with '%s' + '%s' returned %t", param, pattern, skip_pattern, res)
		}
	}

}
