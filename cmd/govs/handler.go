/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dpvs/govs"
)

var (
	Cmd       my_flag
	FirstCmd  *flag.FlagSet
	OthersCmd *flag.FlagSet
)

type my_flag struct {
	Name   int
	Action func(args interface{})
}

func init() {
	Cmd.Name = 0
	Cmd.Action = nil

	FirstCmd = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	OthersCmd = flag.NewFlagSet("", flag.ExitOnError)

	FirstCmd.BoolVar(&govs.FirstCmd.ADD, "A", false, "add virtual service with options")
	FirstCmd.BoolVar(&govs.FirstCmd.EDIT, "E", false, "edit virtual service with options")
	FirstCmd.BoolVar(&govs.FirstCmd.DEL, "D", false, "delete virtual service")
	FirstCmd.BoolVar(&govs.FirstCmd.ADDDEST, "a", false, "add real server with options")
	FirstCmd.BoolVar(&govs.FirstCmd.EDITDEST, "e", false, "edit real server with options")
	FirstCmd.BoolVar(&govs.FirstCmd.DELDEST, "d", false, "delete real server")
	FirstCmd.BoolVar(&govs.FirstCmd.FLUSH, "C", false, "clear the whole table")
	FirstCmd.BoolVar(&govs.FirstCmd.LIST, "L", false, "list the table")
	FirstCmd.BoolVar(&govs.FirstCmd.LIST, "l", false, "list the table")
	FirstCmd.BoolVar(&govs.FirstCmd.ZERO, "Z", false, "zero counters in a service or all services")
	FirstCmd.BoolVar(&govs.FirstCmd.TIMEOUT, "TAG_SET", false, "set connection timeout values")
	FirstCmd.BoolVar(&govs.FirstCmd.USAGE, "h", false, "display this help message")
	FirstCmd.BoolVar(&govs.FirstCmd.VERSION, "V", false, "get version")
	FirstCmd.BoolVar(&govs.FirstCmd.ADDLADDR, "P", false, "add local address")
	FirstCmd.BoolVar(&govs.FirstCmd.DELLADDR, "Q", false, "del local address")
	FirstCmd.BoolVar(&govs.FirstCmd.GETLADDR, "G", false, "get local address")
	FirstCmd.BoolVar(&govs.FirstCmd.STATUS, "s", false, "get dpvs status")

	OthersCmd.StringVar(&govs.CmdOpt.TCP, "t", "", "service-address is host[:port]")
	OthersCmd.StringVar(&govs.CmdOpt.UDP, "u", "", "service-address is host[:port]")
	OthersCmd.Var(&govs.CmdOpt.Netmask, "M", "netmask deafult 0.0.0.0")
	OthersCmd.StringVar(&govs.CmdOpt.Sched_name, "s", "", "scheduler name rr/wrr")
	OthersCmd.UintVar(&govs.CmdOpt.Flags, "flags", 0, "the service flags")
	OthersCmd.Var(&govs.CmdOpt.Daddr, "r", "server-address is host (and port)")
	OthersCmd.IntVar(&govs.CmdOpt.Weight, "w", -1, "capacity of real server")
	OthersCmd.UintVar(&govs.CmdOpt.U_threshold, "x", 0, "upper threshold of connections")
	OthersCmd.UintVar(&govs.CmdOpt.L_threshold, "y", 0, "lower threshold of connections")
	OthersCmd.Var(&govs.CmdOpt.Lip, "z", "local-address")
	OthersCmd.StringVar(&govs.CmdOpt.Typ, "type", "", "type of the stats name(io/w/we/dev/ctl/mem/falcon/vs)")
	OthersCmd.IntVar(&govs.CmdOpt.Id, "i", -1, "id of cpu worker")
	OthersCmd.StringVar(&govs.CmdOpt.Timeout_s, "set", "", "set <tcp,tcp_fin,udp>")
	OthersCmd.UintVar(&govs.CmdOpt.Conn_flags, "conn_flags", 0, "the conn flags")
	OthersCmd.BoolVar(&govs.CmdOpt.Print_detail, "detail", false, "print detail information")
	OthersCmd.BoolVar(&govs.CmdOpt.Print_all_worker, "all", false, "print all cpu worker")
	OthersCmd.Uint64Var(&govs.CmdOpt.Coefficient, "n", 1, "multiplication coefficient")
}

func handler() {
	if len(os.Args) < 2 {
		Cmd.Action = list_handle
		Cmd.Name = govs.CMD_LIST
		return
	}
	FirstCmd.Parse(os.Args[1:2])
	switch {
	case govs.FirstCmd.ADD:
		Cmd.Action = add_handle
		Cmd.Name = govs.CMD_ADD
	case govs.FirstCmd.EDIT:
		Cmd.Action = edit_handle
		Cmd.Name = govs.CMD_EDIT
	case govs.FirstCmd.DEL:
		Cmd.Action = del_handle
		Cmd.Name = govs.CMD_DEL
	case govs.FirstCmd.ADDDEST:
		Cmd.Action = add_handle
		Cmd.Name = govs.CMD_ADDDEST
	case govs.FirstCmd.EDITDEST:
		Cmd.Action = edit_handle
		Cmd.Name = govs.CMD_EDITDEST
	case govs.FirstCmd.DELDEST:
		Cmd.Action = del_handle
		Cmd.Name = govs.CMD_DELDEST
	case govs.FirstCmd.ADDLADDR:
		Cmd.Action = add_handle
		Cmd.Name = govs.CMD_ADDLADDR
	case govs.FirstCmd.DELLADDR:
		Cmd.Action = del_handle
		Cmd.Name = govs.CMD_DELLADDR
	case govs.FirstCmd.GETLADDR:
		Cmd.Action = list_handle
		Cmd.Name = govs.CMD_GETLADDR
	case govs.FirstCmd.FLUSH:
		Cmd.Action = flush_handle
		Cmd.Name = govs.CMD_FLUSH
	case govs.FirstCmd.LIST:
		Cmd.Action = list_handle
		Cmd.Name = govs.CMD_LIST
	case govs.FirstCmd.STATUS:
		Cmd.Action = stats_handle
		Cmd.Name = govs.CMD_STATUS
	case govs.FirstCmd.TIMEOUT:
		Cmd.Action = timeout_handle
		Cmd.Name = govs.CMD_TIMEOUT
	case govs.FirstCmd.VERSION:
		Cmd.Action = version_handle
		Cmd.Name = govs.CMD_VERSION
	case govs.FirstCmd.ZERO:
		Cmd.Action = zero_handle
		Cmd.Name = govs.CMD_ZERO
	default:
		usage()
		return
	}
	OthersCmd.Parse(os.Args[2:])
	CmdCheck()
}

func CmdCheck() {
	var options uint
	OptCheck(&options)
	i := Cmd.Name - 1
	for j := 0; j < govs.NUMBER_OF_OPT; j++ {
		if options&(1<<uint(j+1)) == 0 {
			if govs.CMD_V_OPT[i][j] == '+' {
				log.Fatalf("\nYou need to supply the '%s' option for the '%s' command\n\n", govs.OPTNAMES[j], govs.CMDNAMES[i])
			}
		} else {
			if govs.CMD_V_OPT[i][j] == 'x' {
				log.Fatalf("\nIllegal '%s' option with the '%s' command\n\n", govs.OPTNAMES[j], govs.CMDNAMES[i])
			}
		}
	}

}

func OptCheck(options *uint) {
	if govs.CmdOpt.TCP != "" || govs.CmdOpt.UDP != "" {
		set_option(options, govs.OPT_SERVICE)
	}

	if govs.CmdOpt.Netmask != 0 {
		set_option(options, govs.OPT_NETMASK)
	}

	if govs.CmdOpt.Sched_name == "" {
		govs.CmdOpt.Sched_name = "rr"
	} else {
		set_option(options, govs.OPT_SCHEDULER)
	}

	if govs.CmdOpt.Flags != 0 {
		set_option(options, govs.OPT_FLAGS)
	}

	if govs.CmdOpt.Daddr.Ip != govs.Be32(0) {
		set_option(options, govs.OPT_REALSERVER)
	}

	if govs.CmdOpt.Weight == -1 {
		govs.CmdOpt.Weight = 0
	} else {
		set_option(options, govs.OPT_WEIGHT)
	}

	if govs.CmdOpt.U_threshold != 0 {
		set_option(options, govs.OPT_UTHRESHOLD)
	}

	if govs.CmdOpt.L_threshold != 0 {
		set_option(options, govs.OPT_LTHRESHOLD)
	}

	if govs.CmdOpt.Lip != 0 {
		set_option(options, govs.OPT_LADDR)
	}

	if govs.CmdOpt.Typ == "" {
		govs.CmdOpt.Typ = "io"
	} else {
		set_option(options, govs.OPT_TYPE)
	}

	if govs.CmdOpt.Id != -1 {
		set_option(options, govs.OPT_ID)
	}

	if govs.CmdOpt.Timeout_s != "" {
		set_option(options, govs.OPT_TIMEOUT)
	}

	if govs.CmdOpt.Conn_flags != 0 {
		set_option(options, govs.OPT_CONNFLAGS)
	}
	if govs.CmdOpt.Print_detail != false {
		set_option(options, govs.OPT_PRINTDETAIL)
	}
	if govs.CmdOpt.Print_all_worker != false {
		set_option(options, govs.OPT_PRINTALLWORKER)
	}
	if govs.CmdOpt.Coefficient != 1 {
		set_option(options, govs.OPT_COEFFICIENT)
	}
}

func set_option(options *uint, option uint) {
	*options |= (1 << option)
}

func version_handle(arg interface{}) {
	if version, err := govs.Get_version(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(version)
	}
}

func info_handle(arg interface{}) {
	if info, err := govs.Get_version(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(info)
	}
}

func timeout_handle(arg interface{}) {
	opt := arg.(*govs.CallOptions)
	o := &opt.Opt

	if o.Timeout_s != "" {
		if timeout, err := govs.Set_timeout(o); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(timeout)
		}
	} else {
		if timeout, err := govs.Get_timeout(o); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(timeout)
		}
	}
}

func list_handle(arg interface{}) {
	opt := arg.(*govs.CallOptions)
	if err := govs.Parse_service(opt); err != nil {
		fmt.Println(err)
		return
	}
	o := &opt.Opt

	ret, err := govs.Get_stats_vs(o)
	if err != nil {
		fmt.Println(err)
		return
	}
	if ret.Code != 0 {
		fmt.Println(ret.Msg)
		return
	}

	// print title
	fmt.Println(govs.Svc_title(o.Print_detail))
	if !govs.FirstCmd.GETLADDR {
		fmt.Println(govs.Dest_title(o.Print_detail))
	} else {
		fmt.Println(govs.Laddr_title())
	}

	//print data
	for _, svc := range ret.Services {
		svc.ListVsStats(o.Print_detail, o.Coefficient)
		if !govs.FirstCmd.GETLADDR {
			for _, d := range svc.Dests {
				d.ListDestStats(o.Print_detail, o.Coefficient)
			}
			fmt.Println("")
		} else {
			o.Addr.Ip = svc.Addr
			o.Addr.Port = svc.Port
			o.Protocol = govs.Protocol(svc.Protocol)

			laddrs, err := govs.Get_laddrs(o)
			if err != nil || laddrs.Code != 0 || len(laddrs.Laddrs) == 0 {
				continue
			}
			fmt.Println(laddrs)
		}
	}
}

func flush_handle(arg interface{}) {
	if reply, err := govs.Set_flush(nil); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}
}

func zero_handle(arg interface{}) {
	opt := arg.(*govs.CallOptions)
	govs.Parse_service(opt)

	if reply, err := govs.Set_zero(&opt.Opt); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}

}

func add_handle(arg interface{}) {
	var err error
	var reply *govs.Vs_cmd_r

	opt := arg.(*govs.CallOptions)
	if err := govs.Parse_service(opt); err != nil {
		fmt.Println(err)
		return
	}
	o := &opt.Opt

	switch {
	case govs.FirstCmd.ADD:
		reply, err = govs.Set_add(o)
	case govs.FirstCmd.ADDDEST:
		reply, err = govs.Set_adddest(o)
	case govs.FirstCmd.ADDLADDR:
		reply, err = govs.Set_addladdr(o)
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}
}

func edit_handle(arg interface{}) {
	var err error
	var reply *govs.Vs_cmd_r

	opt := arg.(*govs.CallOptions)
	if err := govs.Parse_service(opt); err != nil {
		fmt.Println(err)
		return
	}
	o := &opt.Opt

	switch {
	case govs.FirstCmd.EDIT:
		reply, err = govs.Set_edit(o)
	case govs.FirstCmd.EDITDEST:
		reply, err = govs.Set_editdest(o)
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}
}

func del_handle(arg interface{}) {
	var err error
	var reply *govs.Vs_cmd_r

	opt := arg.(*govs.CallOptions)
	if err := govs.Parse_service(opt); err != nil {
		fmt.Println(err)
		return
	}
	o := &opt.Opt

	switch {
	case govs.FirstCmd.DEL:
		reply, err = govs.Set_del(o)
	case govs.FirstCmd.DELDEST:
		reply, err = govs.Set_deldest(o)
	case govs.FirstCmd.DELLADDR:
		reply, err = govs.Set_delladdr(o)
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}
}

func stats_handle(arg interface{}) {
	id := govs.CmdOpt.Id

	switch govs.CmdOpt.Typ {
	case "io":
		relay, err := govs.Get_stats_io(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	case "w":
		relay, err := govs.Get_stats_worker(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	case "we":
		relay, err := govs.Get_estats_worker(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	case "dev":
		relay, err := govs.Get_stats_dev(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	case "ctl":
		relay, err := govs.Get_stats_ctl()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	case "mem":
		relay, err := govs.Get_stats_mem()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(relay)
	case "vs":
		opt := arg.(*govs.CallOptions)
		if err := govs.Parse_service(opt); err != nil {
			fmt.Println(err)
			return
		}
		o := &opt.Opt
		relay, err := govs.Get_stats_vs(o)
		if err != nil {
			fmt.Println(err)
			return
		}
		relay.PrintVsStats(o.Coefficient)

	case "falcon":
		falcon_handle(id)
	default:
		fmt.Println("govs -s -type io/w/we/dev/ctl/mem/falcon/vs")
	}
}

func usage() {
	program := os.Args[0]
	fmt.Println(
		"Usage:\n",
		program, "-A|E -t|u service-address [-s scheduler] [-M netmask] [-flags service-flags]\n",
		program, "-D -t|u service-address\n",
		program, "-C\n",
		program, "-a|e -t|u service-address -r server-address [-w weight] [-x upper-threshold] [-y lower-threshold] [-conn_flags conn-flags]\n",
		program, "-d -t|u service-address -r server-address\n",
		program, "-L|l [-t|u service-address] [-detail] [-i id] [-all] [-n coefficient]\n",
		program, "-Z [-t|u service-address]\n",
		program, "-P|Q -t|u service-address -z local-address\n",
		program, "-G [-t|u service-address] \n",
		program, "-TAG_SET [-set tcp/tcp_fin/udp]\n",
		program, "-V\n",
		program, "-s [-type stats-name] [-i id] [-all] [-n coefficient]\n",
		program, "-h\n",
	)
	fmt.Printf("Commands:\n")
	FirstCmd.PrintDefaults()
	fmt.Printf("\nOptions:\n")
	OthersCmd.PrintDefaults()
}

func falcon_handle(id int) {
	falcon_io(id)
	falcon_dev(id)
	falcon_mem(id)
	falcon_we(id)
}

func falcon_we(id int) {
	var ret string
	var items_we = []struct {
		name  string
		tag   string
		count int64
	}{
		{"conn_new_mbuf", "conn.new.mbuf", 0},
		{"conn_new_mbuf_fail", "conn.new.mbuf.fail", 0},
		{"conn_reuse_mbuf", "conn.reuse.mbuf", 0},
		{"conn_reuse_mbuf_fail", "conn.reuse.mbuf.fail", 0},
	}

	relay_we, err := govs.Get_estats_worker(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	if relay_we.Code != 0 {
		fmt.Printf("%s:%s", govs.Ecode(relay_we.Code), relay_we.Msg)
		return
	}

	for _, item := range items_we {
		for _, e := range relay_we.Worker {
			item.count += e[item.name]
		}
		ret += fmt.Sprintf("COUNTER %s %d\n", item.tag, item.count)
	}

	fmt.Println(ret)
}

func falcon_io(id int) {
	var ret string
	var rx_ring_pkts_drop int64

	//get io stats
	relay_io, err := govs.Get_stats_io(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	if relay_io.Code != 0 {
		fmt.Printf("%s:%s", govs.Ecode(relay_io.Code), relay_io.Msg)
		return
	}
	for _, e := range relay_io.Io {
		for i, _ := range e.Rx_rings_iters {
			rx_ring_pkts_drop += e.Rx_rings_drop_pkts[i]
		}
	}
	ret += fmt.Sprintf("COUNTER net.if.in.ring.drop.pkts %d\n", rx_ring_pkts_drop)

	fmt.Printf("%s", ret)
}

func falcon_dev(id int) {
	var ret string

	//get dev stats
	relay_dev, err := govs.Get_stats_dev(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	if relay_dev.Code != 0 {
		fmt.Printf("%s:%s", govs.Ecode(relay_dev.Code), relay_dev.Msg)
		return
	}
	for _, e := range relay_dev.Dev {
		ret += fmt.Sprintf("COUNTER net.if.in.packets %d iface=port%d\n", e.Ipackets, e.Port_id)
		ret += fmt.Sprintf("COUNTER net.if.in.bytes %d iface=port%d\n", e.Ibytes, e.Port_id)
		ret += fmt.Sprintf("COUNTER net.if.in.bits %d iface=port%d\n", e.Ibytes*8, e.Port_id)
		ret += fmt.Sprintf("COUNTER net.if.in.errors %d iface=port%d\n", e.Ierrors, e.Port_id)
		ret += fmt.Sprintf("COUNTER net.if.in.dropped %d iface=port%d\n", e.Imissed, e.Port_id)
		ret += fmt.Sprintf("COUNTER net.if.out.packets %d iface=port%d\n", e.Opackets, e.Port_id)
		ret += fmt.Sprintf("COUNTER net.if.out.bytes %d iface=port%d\n", e.Obytes, e.Port_id)
		ret += fmt.Sprintf("COUNTER net.if.out.bits %d iface=port%d\n", e.Obytes*8, e.Port_id)
		ret += fmt.Sprintf("COUNTER net.if.out.errors %d iface=port%d\n", e.Oerrors, e.Port_id)
	}

	fmt.Printf("%s", ret)
}

func falcon_mem(id int) {
	var ret string

	// get mem stats
	var res_mem struct {
		mbuf_used  int
		svc_used   int
		rs_used    int
		laddr_used int
		conn_used  int
	}

	relay_mem, err := govs.Get_stats_mem()
	if err != nil {
		fmt.Println(err)
		return
	}
	if relay_mem.Code != 0 {
		fmt.Printf("%s:%s", govs.Ecode(relay_mem.Code), relay_mem.Msg)
		return
	}

	for _, e := range relay_mem.Available {
		res_mem.mbuf_used += (relay_mem.Size.Mbuf - e.Mbuf)
		res_mem.svc_used += (relay_mem.Size.Svc - e.Svc)
		res_mem.rs_used += (relay_mem.Size.Rs - e.Rs)
		res_mem.laddr_used += (relay_mem.Size.Laddr - e.Laddr)
		res_mem.conn_used += (relay_mem.Size.Conn - e.Conn)
	}

	ret += fmt.Sprintf("GAUGE mbuf.num.used %d\n", res_mem.mbuf_used)
	ret += fmt.Sprintf("GAUGE svc.num.used %d\n", res_mem.svc_used)
	ret += fmt.Sprintf("GAUGE rs.num.used %d\n", res_mem.rs_used)
	ret += fmt.Sprintf("GAUGE laddr.num.used %d\n", res_mem.laddr_used)
	ret += fmt.Sprintf("GAUGE conn.num.used %d\n", res_mem.conn_used)

	ret += fmt.Sprintf("GAUGE mbuf.percent.used %f\n", float64(res_mem.mbuf_used)/float64(relay_mem.Size.Mbuf))
	ret += fmt.Sprintf("GAUGE svc.percent.used %f\n", float64(res_mem.svc_used)/float64(relay_mem.Size.Svc))
	ret += fmt.Sprintf("GAUGE rs.percent.used %f\n", float64(res_mem.rs_used)/float64(relay_mem.Size.Rs))
	ret += fmt.Sprintf("GAUGE laddr.percent.used %f\n", float64(res_mem.laddr_used)/float64(relay_mem.Size.Laddr))
	ret += fmt.Sprintf("GAUGE conn.percent.used %f\n", float64(res_mem.conn_used)/float64(relay_mem.Size.Conn))

	fmt.Printf("%s", ret)
}
