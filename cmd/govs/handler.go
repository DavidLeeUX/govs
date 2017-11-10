/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dpvs/govs"
	"github.com/yubo/gotool/flags"
)

func init() {
	flags.FirstCmd.BoolVar(&govs.FirstCmd.ADD, "A", false, "add virtual service with options")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.EDIT, "E", false, "edit virtual service with options")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.DEL, "D", false, "delete virtual service")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.ADDDEST, "a", false, "add real server with options")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.EDITDEST, "e", false, "edit real server with options")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.DELDEST, "d", false, "delete real server")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.FLUSH, "C", false, "clear the whole table")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.LIST, "L", false, "list the table")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.LIST, "l", false, "list the table")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.ZERO, "Z", false, "zero counters in a service or all services")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.TIMEOUT, "TAG_SET", false, "set connection timeout values")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.USAGE, "h", false, "display this help message")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.VERSION, "V", false, "get version")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.ADDLADDR, "P", false, "add local address")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.DELLADDR, "Q", false, "del local address")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.GETLADDR, "G", false, "get local address")
	flags.FirstCmd.BoolVar(&govs.FirstCmd.STATUS, "s", false, "get dpvs status")

	flags.OthersCmd.StringVar(&govs.CmdOpt.TCP, "t", "", "service-address is host[:port]")
	flags.OthersCmd.StringVar(&govs.CmdOpt.UDP, "u", "", "service-address is host[:port]")
	flags.OthersCmd.Var(&govs.CmdOpt.Netmask, "M", "netmask deafult 0.0.0.0")
	flags.OthersCmd.StringVar(&govs.CmdOpt.Sched_name, "s", "", "scheduler name rr/wrr")
	flags.OthersCmd.UintVar(&govs.CmdOpt.Flags, "flags", 0, "the service flags")
	flags.OthersCmd.Var(&govs.CmdOpt.Daddr, "r", "server-address is host (and port)")
	flags.OthersCmd.IntVar(&govs.CmdOpt.Weight, "w", -1, "capacity of real server")
	flags.OthersCmd.UintVar(&govs.CmdOpt.U_threshold, "x", 0, "upper threshold of connections")
	flags.OthersCmd.UintVar(&govs.CmdOpt.L_threshold, "y", 0, "lower threshold of connections")
	flags.OthersCmd.Var(&govs.CmdOpt.Lip, "z", "local-address")
	flags.OthersCmd.StringVar(&govs.CmdOpt.Typ, "type", "", "type of the stats name(io/w/we/dev/ctl/mem/falcon)")
	flags.OthersCmd.IntVar(&govs.CmdOpt.Id, "i", -1, "id of the stats object")
	flags.OthersCmd.StringVar(&govs.CmdOpt.Timeout_s, "set", "", "set <tcp,tcp_fin,udp>")
	flags.OthersCmd.UintVar(&govs.CmdOpt.Conn_flags, "conn_flags", 0, "the conn flags")
}

func handler() {
	if len(os.Args) < 2 {
		flags.Cmd.Action = list_handle
		flags.Cmd.Name = govs.CMD_LIST
		return
	}
	flags.FirstCmd.Parse(os.Args[1:2])
	switch {
	case govs.FirstCmd.ADD:
		flags.Cmd.Action = add_handle
		flags.Cmd.Name = govs.CMD_ADD
	case govs.FirstCmd.EDIT:
		flags.Cmd.Action = edit_handle
		flags.Cmd.Name = govs.CMD_EDIT
	case govs.FirstCmd.DEL:
		flags.Cmd.Action = del_handle
		flags.Cmd.Name = govs.CMD_DEL
	case govs.FirstCmd.ADDDEST:
		flags.Cmd.Action = add_handle
		flags.Cmd.Name = govs.CMD_ADDDEST
	case govs.FirstCmd.EDITDEST:
		flags.Cmd.Action = edit_handle
		flags.Cmd.Name = govs.CMD_EDITDEST
	case govs.FirstCmd.DELDEST:
		flags.Cmd.Action = del_handle
		flags.Cmd.Name = govs.CMD_DELDEST
	case govs.FirstCmd.ADDLADDR:
		flags.Cmd.Action = add_handle
		flags.Cmd.Name = govs.CMD_ADDLADDR
	case govs.FirstCmd.DELLADDR:
		flags.Cmd.Action = del_handle
		flags.Cmd.Name = govs.CMD_DELLADDR
	case govs.FirstCmd.GETLADDR:
		flags.Cmd.Action = list_handle
		flags.Cmd.Name = govs.CMD_GETLADDR
	case govs.FirstCmd.FLUSH:
		flags.Cmd.Action = flush_handle
		flags.Cmd.Name = govs.CMD_FLUSH
	case govs.FirstCmd.LIST:
		flags.Cmd.Action = list_handle
		flags.Cmd.Name = govs.CMD_LIST
	case govs.FirstCmd.STATUS:
		flags.Cmd.Action = stats_handle
		flags.Cmd.Name = govs.CMD_STATUS
	case govs.FirstCmd.TIMEOUT:
		flags.Cmd.Action = timeout_handle
		flags.Cmd.Name = govs.CMD_TIMEOUT
	case govs.FirstCmd.VERSION:
		flags.Cmd.Action = version_handle
		flags.Cmd.Name = govs.CMD_VERSION
	case govs.FirstCmd.ZERO:
		flags.Cmd.Action = zero_handle
		flags.Cmd.Name = govs.CMD_ZERO
	default:
		Usage()
		flags.Usage()
		return
	}
	flags.OthersCmd.Parse(os.Args[2:])
	CmdCheck()
}

func CmdCheck() {
	var options uint
	OptCheck(&options)
	i := flags.Cmd.Name - 1
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

func list_svc_handle(o *govs.CmdOptions) {

	ret, err := govs.Get_service(o)
	if err != nil {
		fmt.Println(err)
		return
	}

	if ret.Code != 0 {
		fmt.Println(ret.Msg)
		return
	}

	fmt.Println(govs.Svc_title())
	if !govs.FirstCmd.GETLADDR {
		fmt.Println(govs.Dest_title())
	} else {
		fmt.Println(govs.Laddr_title())
	}

	fmt.Println(ret)
	if !govs.FirstCmd.GETLADDR {
		dests, err := govs.Get_dests(o)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(dests)
	} else {
		laddrs, err := govs.Get_laddrs(o)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(laddrs)
	}
	return
}

func list_svcs_handle(o *govs.CmdOptions) {

	ret, err := govs.Get_services(o)

	if err != nil {
		fmt.Println(err)
		return
	}

	if ret.Code != 0 {
		fmt.Println(ret.Msg)
		return
	}

	fmt.Println(govs.Svc_title())
	if !govs.FirstCmd.GETLADDR {
		fmt.Println(govs.Dest_title())
	} else {
		fmt.Println(govs.Laddr_title())
	}

	for _, svc := range ret.Services {
		fmt.Println(svc)
		o.Addr.Ip = svc.Addr
		o.Addr.Port = svc.Port
		o.Protocol = govs.Protocol(svc.Protocol)

		if !govs.FirstCmd.GETLADDR {
			dests, err := govs.Get_dests(o)
			if err != nil || dests.Code != 0 ||
				len(dests.Dests) == 0 {
				//fmt.Println(err)
				continue
			}
			fmt.Println(dests)
		} else {
			laddrs, err := govs.Get_laddrs(o)
			if err != nil || laddrs.Code != 0 ||
				len(laddrs.Laddrs) == 0 {
				//fmt.Println(err)
				continue
			}
			fmt.Println(laddrs)
		}
	}

}

func list_handle(arg interface{}) {
	opt := arg.(*govs.CallOptions)
	govs.Parse_service(opt)
	o := &opt.Opt

	if o.Addr.Ip != 0 {
		list_svc_handle(o)
		return
	}

	list_svcs_handle(o)
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
	case "falcon":
		falcon_handle(id)
	default:
		fmt.Println("govs stats -t io/worker/dev/ctl")
	}
}

func Usage() {
	program := os.Args[0]
	fmt.Println(
		"Usage:\n",
		program, "-A|E -t|u service-address [-s scheduler] [-M netmask] [-flags service-flags]\n",
		program, "-D -t|u service-address\n",
		program, "-C\n",
		program, "-a|e -t|u service-address -r server-address [-w weight] [-x upper-threshold] [-y lower-threshold] [-conn_flags conn-flags]\n",
		program, "-d -t|u service-address -r server-address\n",
		program, "-L|l [-t|u service-address]\n",
		program, "-Z [-t|u service-address]\n",
		program, "-P|Q -t|u service-address -z local-address\n",
		program, "-G [-t|u service-address] \n",
		program, "-TAG_SET [-set tcp/tcp_fin/udp]\n",
		program, "-V\n",
		program, "-s [-type stats-name] [-i id]\n",
		program, "-h\n",
	)
}

func falcon_handle(id int) {
	falcon_io(id)
	falcon_dev(id)
	falcon_mem(id)
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
	ret += fmt.Sprintf("net.if.in.ring.drop.pkts %d\n", rx_ring_pkts_drop)

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
		ret += fmt.Sprintf("net.if.in.packets %d iface=port%d\n", e.Ipackets, e.Port_id)
		ret += fmt.Sprintf("net.if.in.bytes %d iface=port%d\n", e.Ibytes, e.Port_id)
		ret += fmt.Sprintf("net.if.in.bits %d iface=port%d\n", e.Ibytes*8, e.Port_id)
		ret += fmt.Sprintf("net.if.in.errors %d iface=port%d\n", e.Ierrors, e.Port_id)
		ret += fmt.Sprintf("net.if.in.dropped %d iface=port%d\n", e.Imissed, e.Port_id)
		ret += fmt.Sprintf("net.if.out.packets %d iface=port%d\n", e.Opackets, e.Port_id)
		ret += fmt.Sprintf("net.if.out.bytes %d iface=port%d\n", e.Obytes, e.Port_id)
		ret += fmt.Sprintf("net.if.out.bits %d iface=port%d\n", e.Obytes*8, e.Port_id)
		ret += fmt.Sprintf("net.if.out.errors %d iface=port%d\n", e.Oerrors, e.Port_id)
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

	ret += fmt.Sprintf("mbuf.num.used %d\n", res_mem.mbuf_used)
	ret += fmt.Sprintf("svc.num.used %d\n", res_mem.svc_used)
	ret += fmt.Sprintf("rs.num.used %d\n", res_mem.rs_used)
	ret += fmt.Sprintf("laddr.num.used %d\n", res_mem.laddr_used)
	ret += fmt.Sprintf("conn.num.used %d\n", res_mem.conn_used)

	ret += fmt.Sprintf("mbuf.percent.used %f\n", float64(res_mem.mbuf_used)/float64(relay_mem.Size.Mbuf))
	ret += fmt.Sprintf("svc.percent.used %f\n", float64(res_mem.svc_used)/float64(relay_mem.Size.Svc))
	ret += fmt.Sprintf("rs.percent.used %f\n", float64(res_mem.rs_used)/float64(relay_mem.Size.Rs))
	ret += fmt.Sprintf("laddr.percent.used %f\n", float64(res_mem.laddr_used)/float64(relay_mem.Size.Laddr))
	ret += fmt.Sprintf("conn.percent.used %f\n", float64(res_mem.conn_used)/float64(relay_mem.Size.Conn))

	fmt.Printf("%s", ret)
}
