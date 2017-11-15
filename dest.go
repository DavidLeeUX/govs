/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package govs

import (
	"fmt"
	"strings"
)

type Vs_dest_user struct {
	Nic         uint8
	Addr        Be32
	Port        Be16
	Conn_flags  uint
	Weight      int
	U_threshold uint32
	L_threshold uint32
}

type Vs_dest_user_r struct {
	Addr        Be32
	Port        Be16
	Conn_flags  uint
	Weight      int
	U_threshold uint32
	L_threshold uint32
	Activeconns uint32
	Inactconns  uint32
	Persistent  uint32
	Conns       uint64
	Inpkts      uint64
	Outpkts     uint64
	Inbytes     uint64
	Outbytes    uint64
}

const (
	fmt_dest_t        = "%5s %21s %8s %8s %15s %12s %12s %12s %7s %10s %10s %10s %10s"
	fmt_dest          = "%5s %21s %8s %8d %15s %12d %12d %12d %7d %10d %10d %10d %10d"
	fmt_dest_simple   = "%5s %-21s %8s %8d %10d %10d"
	fmt_dest_simple_t = "%5s %-21s %8s %8s %10s %10s"
)

func Dest_title(detail bool) string {
	var res string
	if detail == true {
		res += fmt.Sprintf(fmt_dest_t,
			"->", "Addr:Port", "Flags", "Weight", "threshold",
			"Activeconns", "Inactconns", "Persistent",
			"Conns", "Inpkts", "Outpkts", "Inbytes", "Outbytes")

	} else {
		res += fmt.Sprintf(fmt_dest_simple_t, "->", "RemoteAddress:Port", "Forward",
			"Weight", "ActiveConn", "InActConn")
	}
	return res
}

func (d *Vs_stats_dest) ListDestStats(detail bool) {
	var res string
	if detail == true {
		res += fmt.Sprintf(fmt_dest,
			"->", fmt.Sprintf("%s:%s", d.Addr.String(), d.Port.String()), ip_vs_fwd_name(d.Conn_flags),
			d.Weight, fmt.Sprintf("%d-%d", d.L_threshold, d.U_threshold),
			d.Activeconns, d.Inactconns, d.Persistent,
			d.Conns, d.Inpkts, d.Outpkts, d.Inbytes, d.Outbytes)

	} else {
		res += fmt.Sprintf(fmt_dest_simple, "->",
			fmt.Sprintf("%s:%s", d.Addr.String(), d.Port.String()), ip_vs_fwd_name(d.Conn_flags),
			d.Weight, d.Activeconns, d.Inactconns)
	}
	fmt.Println(res)
}

func (d Vs_dest_user_r) String() string {
	return fmt.Sprintf(fmt_dest,
		"->", fmt.Sprintf("%s:%s", d.Addr.String(), d.Port.String()), ip_vs_fwd_name(uint32(d.Conn_flags)),
		d.Weight, fmt.Sprintf("%d-%d", d.L_threshold, d.U_threshold),
		d.Activeconns, d.Inactconns, d.Persistent,
		d.Conns, d.Inpkts, d.Outpkts, d.Inbytes, d.Outbytes)
}

type Vs_list_dests_r struct {
	Code  int
	Msg   string
	Dests []Vs_dest_user_r
}

func (r Vs_list_dests_r) String() string {
	var s string
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	for _, dest := range r.Dests {
		s += fmt.Sprintf("%s\n", dest)
	}

	return strings.TrimRight(s, "\n")
}

func Get_dests(o *CmdOptions) (*Vs_list_dests_r, error) {
	var reply Vs_list_dests_r
	args := Vs_list_q{
		Cmd: VS_CMD_GET_DEST,
		Service: Vs_service_user{
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
			Protocol: uint8(o.Protocol),
			Number:   o.Number,
		},
	}

	err := client.Call("api", args, &reply)
	return &reply, err
}

type Vs_dest_q struct {
	Cmd     int
	Service Vs_service_user
	Dest    Vs_dest_user
}

func Set_adddest(o *CmdOptions) (*Vs_cmd_r, error) {
	var reply Vs_cmd_r
	args := Vs_dest_q{
		Cmd: VS_CMD_NEW_DEST,
		Service: Vs_service_user{
			Protocol: uint8(o.Protocol),
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
		},
		Dest: Vs_dest_user{
			Nic:         uint8(o.Dnic),
			Addr:        o.Daddr.Ip,
			Port:        o.Daddr.Port,
			Conn_flags:  o.Conn_flags | VS_CONN_F_FULLNAT,
			Weight:      o.Weight,
			U_threshold: uint32(o.U_threshold),
			L_threshold: uint32(o.L_threshold),
		},
	}

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Set_editdest(o *CmdOptions) (*Vs_cmd_r, error) {
	var reply Vs_cmd_r
	args := Vs_dest_q{
		Cmd: VS_CMD_SET_DEST,
		Service: Vs_service_user{
			Protocol: uint8(o.Protocol),
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
		},
		Dest: Vs_dest_user{
			Nic:         uint8(o.Dnic),
			Addr:        o.Daddr.Ip,
			Port:        o.Daddr.Port,
			Conn_flags:  o.Conn_flags | VS_CONN_F_FULLNAT,
			Weight:      o.Weight,
			U_threshold: uint32(o.U_threshold),
			L_threshold: uint32(o.L_threshold),
		},
	}

	err := client.Call("api", args, &reply)
	return &reply, err
}

func Set_deldest(o *CmdOptions) (*Vs_cmd_r, error) {
	var reply Vs_cmd_r
	args := Vs_dest_q{
		Cmd: VS_CMD_DEL_DEST,
		Service: Vs_service_user{
			Protocol: uint8(o.Protocol),
			Addr:     o.Addr.Ip,
			Port:     o.Addr.Port,
		},
		Dest: Vs_dest_user{
			Addr: o.Daddr.Ip,
			Port: o.Daddr.Port,
		},
	}

	err := client.Call("api", args, &reply)
	return &reply, err
}
