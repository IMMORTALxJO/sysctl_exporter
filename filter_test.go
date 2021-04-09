package main

import "testing"

func TestSysctlNameFilter(t *testing.T) {
	var res bool

	// default params, all parameters
	parameters := []string{"abc", "bcd", "cde", "a.b.c", ".a_s", "...", ".*"}
	include := ".*"
	exclude := ""

	for _, param := range parameters {
		res = sysctlNameIsFiltered(param, include, exclude)
		if res {
			t.Errorf("String '%s' did not pass '%s' + '%s'", param, include, exclude)
		}
	}

	// simple match and skip
	parameters = []string{"a.b.c", "b.c.d", "bb", "c.d.e", "_", ".*"}
	answer := []bool{false, false, false, true, true, true}
	include = "b"
	for i, param := range parameters {
		res = sysctlNameIsFiltered(param, include, "")
		if res != answer[i] {
			t.Errorf("String '%s' with '%s' + '%s' returned %t", param, include, exclude, res)
		}
		res = sysctlNameIsFiltered(param, "", include)
		if res == answer[i] {
			t.Errorf("String '%s' with '%s' + '%s' returned %t", param, exclude, include, res)
		}
	}

	//  match and skip, complex
	parameters = []string{"a.b.c", "b.c.d", "bb", "c.d.e", "_e", ".*"}
	answer = []bool{true, true, false, true, false, true}
	include = "b|e"
	exclude = "c"
	for i, param := range parameters {
		res = sysctlNameIsFiltered(param, include, exclude)
		if res != answer[i] {
			t.Errorf("String '%s' with '%s' + '%s' returned %t", param, include, exclude, res)
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
	include = "ipv4.(udp|tcp).*mem"
	exclude = "udp_mem"
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
		res = sysctlNameIsFiltered(param, include, exclude)
		if res != answer[i] {
			t.Errorf("String '%s' with '%s' + '%s' returned %t", param, include, exclude, res)
		}
	}

}
