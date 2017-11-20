/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package govs

import "fmt"

const (
	VS_CTL_S_SYNC = iota
	VS_CTL_S_PENDING
	VS_CTL_S_LAST
)

type Vs_stats_ifa struct {
	Port_id    int
	Rx_packets int64
	Rx_dropped int64
	Tx_packets int64
	Tx_dropped int64
}

type Vs_stats_io_entry struct {
	Core_id int

	Rx_nic_queues_port  []int32
	Rx_nic_queues_queue []int32
	Rx_nic_queues_iters []int64
	Rx_nic_queues_pkts  []int64

	Rx_rings_iters      []int64
	Rx_rings_pkts       []int64
	Rx_rings_drop_iters []int64
	Rx_rings_drop_pkts  []int64
	Rx_rings_drop_count []int64

	Tx_nic_ports_port       []int32
	Tx_nic_ports_queue      []int32
	Tx_nic_ports_iters      []int64
	Tx_nic_ports_pkts       []int64
	Tx_nic_ports_drop_iters []int64
	Tx_nic_ports_drop_pkts  []int64

	Kni []Vs_stats_ifa
}

type Vs_stats_io_r struct {
	Code int
	Msg  string
	Io   []Vs_stats_io_entry
}

/*
id kni_deq kni_err  k_rx rx_d tx tx_d
*/
func (r Vs_stats_io_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	var ret string

	for _, e := range r.Io {
		ret += fmt.Sprintf("\n%-32s %-10s %10d\n", "#######", "core_id", e.Core_id)

		for i, _ := range e.Rx_nic_queues_iters {
			ret += fmt.Sprintf("%-32s %-10s %10d %-10s %10d\n",
				fmt.Sprintf("Rx_nic_port%d_queue%02d",
					e.Rx_nic_queues_port[i], e.Rx_nic_queues_queue[i]),
				"iters", e.Rx_nic_queues_iters[i],
				"packets", e.Rx_nic_queues_pkts[i],
			)
		}

		for i, _ := range e.Rx_rings_iters {
			ret += fmt.Sprintf("%-32s %-10s %10d %-10s %10d %-10s %10d %-10s %10d %-10s %10f\n",
				fmt.Sprintf("rx_ring_worker%02d", i),
				"iters", e.Rx_rings_iters[i],
				"packets", e.Rx_rings_pkts[i],
				"drop_iters", e.Rx_rings_drop_iters[i],
				"drop_pkts", e.Rx_rings_drop_pkts[i],
				"drop_cnt", float64(e.Rx_rings_drop_count[i])/float64(e.Rx_rings_drop_iters[i]),
			)
		}

		for i, _ := range e.Tx_nic_ports_iters {
			ret += fmt.Sprintf("%-32s %-10s %10d %-10s %10d %-10s %10d %-10s %10d\n",
				fmt.Sprintf("tx_nic_port%d", e.Tx_nic_ports_port[i]),
				"iters", e.Tx_nic_ports_iters[i],
				"packets", e.Tx_nic_ports_pkts[i],
				"drop_iters", e.Tx_nic_ports_drop_iters[i],
				"drop_pkts", e.Tx_nic_ports_drop_pkts[i],
			)
		}

		for _, kni := range e.Kni {
			ret += fmt.Sprintf("%-32s %-10s %10d %-10s %10d "+
				"%-10s %10d %-10s %10d\n",
				fmt.Sprintf("veth%d", kni.Port_id),
				"rx_packets", kni.Rx_packets,
				"rx_dropped", kni.Rx_dropped,
				"tx_packets", kni.Tx_packets,
				"tx_dropped", kni.Tx_dropped)
		}
	}
	return ret
}

type Vs_stats_worker_entry struct {
	Core_id              int
	Conns                int64
	Inpkts               int64
	Outpkts              int64
	Inbytes              int64
	Outbytes             int64
	Rings_in_iters       []int64
	Rings_in_pkts        []int64
	Rings_in_miss        []int64
	Rings_in_miss_count  []int64
	Rings_out_port       []int32
	Rings_out_iters      []int64
	Rings_out_pkts       []int64
	Rings_out_drop_iters []int64
	Rings_out_drop_pkts  []int64
	Vs_drop              []int64
}

type Vs_stats_worker_r struct {
	Code   int
	Msg    string
	Worker []Vs_stats_worker_entry
}

func (r Vs_stats_worker_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	var ret string

	for _, e := range r.Worker {
		ret += fmt.Sprintf("%-32s %-10s %10d %-10s %10d\n",
			"", "core_id", e.Core_id, "conns", e.Conns)

		ret += fmt.Sprintf("%-32s %-10s %10d %-10s %10d\n",
			"in", "packets", e.Inpkts, "bytes", e.Inbytes)

		ret += fmt.Sprintf("%-32s %-10s %10d %-10s %10d\n",
			"out", "packets", e.Outpkts, "bytes", e.Outbytes)

		for i, _ := range e.Rings_in_iters {
			ret += fmt.Sprintf("%-32s %-10s %10d %-10s %10d %-10s %10d %-10s %10d %-10s %10f\n",
				fmt.Sprintf("rings_in_io%02d", i),
				"iters", e.Rings_in_iters[i],
				"packets", e.Rings_in_pkts[i],
				"miss", e.Rings_in_miss[i],
				"miss_count", e.Rings_in_miss_count[i],
				"bsz", float64(e.Rings_in_miss_count[i])/float64(e.Rings_in_miss[i]),
			)
		}

		for i, _ := range e.Rings_out_iters {
			ret += fmt.Sprintf("%-32s %-10s %10d %-10s %10d %-10s %10d %-10s %10d\n",
				fmt.Sprintf("rings_out_port%02d", e.Rings_out_port[i]),
				"iters", e.Rings_out_iters[i],
				"packets", e.Rings_out_pkts[i],
				"drop_iters", e.Rings_out_drop_iters[i],
				"drop_pkts", e.Rings_out_drop_pkts[i],
			)
		}
		ret += "\n"
	}
	return ret
}

type Vs_stats_dest struct {
	Addr        Be32
	Port        Be16
	Conn_flags  uint32
	Weight      int32
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

type Vs_stats_svc struct {
	Protocol    uint8
	Addr        Be32
	Port        Be16
	Sched_name  string
	Flags       uint32
	Timeout     uint32
	Netmask     Be32
	Est_timeout uint32
	Conns       uint64
	Inpkts      uint64
	Outpkts     uint64
	Inbytes     uint64
	Outbytes    uint64
	Num_dests   uint32
	Num_laddrs  uint32
	Dests       []Vs_stats_dest
}

type Vs_stats_vs_r struct {
	Code     int
	Msg      string
	Services []Vs_stats_svc
}

func ip_vs_fwd_name(flags uint32) string {
	switch flags & VS_CONN_F_FWD_MASK {
	case VS_CONN_F_LOCALNODE:
		return "Local"
	case VS_CONN_F_TUNNEL:
		return "Tunnel"
	case VS_CONN_F_DROUTE:
		return "Route"
	case VS_CONN_F_FULLNAT:
		return "FullNat"
	default:
		return "Masq"
	}

}

func (r *Vs_stats_vs_r) PrintVsStats(coefficient uint64) {
	var ret string
	ret += fmt.Sprintf("Prot LocalAddress:Port Scheduler Flags\n")

	ret += fmt.Sprintf("  -> RemoteAddress:Port Forward Weight ActiveConn InActConn\n")
	for _, svc := range r.Services {
		ret += fmt.Sprintf("%s  %08X:%04X %s \n", get_protocol_name(svc.Protocol),
			Ntohl(svc.Addr), Ntohs(svc.Port), svc.Sched_name)
		for _, dest := range svc.Dests {
			ret += fmt.Sprintf("  -> %08X:%04X      %-7s %-6d %-10d %-10d\n",
				Ntohl(dest.Addr), Ntohs(dest.Port), ip_vs_fwd_name(dest.Conn_flags),
				dest.Weight, uint64(dest.Activeconns)*coefficient, uint64(dest.Inactconns)*coefficient)
		}

	}

	fmt.Println(ret)
}

func (r Vs_stats_vs_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	var ret string

	ret += fmt.Sprintf("Prot LocalAddress:Port Scheduler Flags\n")
	//                  01234567890123456++++01234567890123456++++01234567890123456
	ret += fmt.Sprintf("  -> RemoteAddress:Port Forward Weight ActiveConn InActConn\n")
	for _, svc := range r.Services {
		ret += fmt.Sprintf("%s  %08X:%04X %s \n", get_protocol_name(svc.Protocol),
			Ntohl(svc.Addr), Ntohs(svc.Port), svc.Sched_name)
		for _, dest := range svc.Dests {
			ret += fmt.Sprintf("  -> %08X:%04X      %-7s %-6d %-10d %-10d\n",
				Ntohl(dest.Addr), Ntohs(dest.Port), ip_vs_fwd_name(dest.Conn_flags),
				dest.Weight, dest.Activeconns, dest.Inactconns)
		}

	}

	return ret
}

var estats_names = []string{
	"core_id",
	"fullnat_add_toa_ok",
	"fullnat_add_toa_fail_len",
	"fullnat_add_toa_head_full",
	"fullnat_add_toa_fail_mem",
	"fullnat_add_toa_fail_proto",
	"fullnat_conn_reused",
	"fullnat_conn_reused_close",
	"fullnat_conn_reused_timewait",
	"fullnat_conn_reused_finwait",
	"fullnat_conn_reused_closewait",
	"fullnat_conn_reused_lastack",
	"fullnat_conn_reused_estab",
	"synproxy_rs_error",
	"synproxy_null_ack",
	"synproxy_bad_ack",
	"synproxy_ok_ack",
	"synproxy_syn_cnt",
	"synproxy_ackstorm",
	"synproxy_synsend_qlen",
	"synproxy_conn_reused",
	"synproxy_conn_reused_close",
	"synproxy_conn_reused_timewait",
	"synproxy_conn_reused_finwait",
	"synproxy_conn_reused_closewait",
	"synproxy_conn_reused_lastack",
	"defence_ip_frag_drop",
	"defence_ip_frag_gather",
	"defence_tcp_drop",
	"defence_udp_drop",
	"fast_xmit_reject",
	"fast_xmit_pass",
	"fast_xmit_skb_copy",
	"fast_xmit_no_mac",
	"fast_xmit_synproxy_save",
	"fast_xmit_dev_lost",
	"rst_in_syn_sent",
	"rst_out_syn_sent",
	"rst_in_established",
	"rst_out_established",
	"gro_pass",
	"lro_reject",
	"xmit_unexpected_mtu",
	"conn_sched_unreach",
	"conn_new_mbuf",
	"conn_new_mbuf_fail",
	"conn_reuse_mbuf",
	"conn_reuse_mbuf_fail",
}

type Vs_estats_worker_r struct {
	Code   int
	Msg    string
	Worker []map[string]int64
}

func (r Vs_estats_worker_r) String() (ret string) {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}

	if len(r.Worker) == 0 {
		return fmt.Sprintf("No such object")
	}

	for _, name := range estats_names {
		ret += fmt.Sprintf("%-32s", name)
		for _, e := range r.Worker {
			ret += fmt.Sprintf(" %10d", e[name])
		}
		ret += "\n"
	}
	return ret
}

type Vs_stats_dev_entry struct {
	Port_id   int
	Ipackets  int64
	Opackets  int64
	Ibytes    int64
	Obytes    int64
	Imissed   int64
	Ierrors   int64
	Oerrors   int64
	Rx_nombuf int64
}

type Vs_stats_dev_r struct {
	Code int
	Msg  string
	Dev  []Vs_stats_dev_entry
}

func (r Vs_stats_dev_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	ret := fmt.Sprintf("%-10s %10s %10s %10s %10s %10s %10s %10s %10s\n",
		"port_id", "ipackets", "opackets", "ibytes", "obytes",
		"imissed", "ierrors", "oerrors", "rx_nombuf")
	for _, e := range r.Dev {
		ret += fmt.Sprintf("%-10d %10d %10d %10d %10d %10d %10d %10d %10d\n",
			e.Port_id, e.Ipackets, e.Opackets, e.Ibytes, e.Obytes,
			e.Imissed, e.Ierrors, e.Oerrors, e.Rx_nombuf)
	}
	return ret
}

type Vs_stats_ctl_r struct {
	Code         int
	Msg          string
	Num_services int
	Seq          int
	Wait_group   int
	Workers      []struct {
		Worker_id    int
		Num_services int
		Seq          int
	}
}

/*
id                seq      n_svc      state
-                   0          0          -

0                   0          0          s
*/
func (r Vs_stats_ctl_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	ret := fmt.Sprintf("%-10s %10s %10s %10s\n",
		"id", "seq", "n_svc", "wait_group")
	ret += fmt.Sprintf("%-10s %10d %10d %10d\n\n",
		"-", r.Seq, r.Num_services, r.Wait_group)
	for _, e := range r.Workers {
		ret += fmt.Sprintf("%-10d %10d %10d\n",
			e.Worker_id, e.Seq, e.Num_services)

	}
	return ret
}

type Vs_stats_mem_r struct {
	Code int
	Msg  string
	Size struct {
		Mbuf  int
		Svc   int
		Rs    int
		Laddr int
		Conn  int
	}
	Available []struct {
		Socket_id int
		Mbuf      int
		Svc       int
		Rs        int
		Laddr     int
		Conn      int
	}
}

func (r Vs_stats_mem_r) String() string {
	if r.Code != 0 {
		return fmt.Sprintf("%s:%s", Ecode(r.Code), r.Msg)
	}
	ret := fmt.Sprintf("%-3s %20s %20s %20s %20s %20s (used/total)\n",
		"id", "mbuf", "svc", "rs", "laddr", "conn")
	for _, e := range r.Available {
		ret += fmt.Sprintf("%-3d %20s %20s %20s %20s %20s\n",
			e.Socket_id,
			fmt.Sprintf("%d/%d", r.Size.Mbuf-e.Mbuf, r.Size.Mbuf),
			fmt.Sprintf("%d/%d", r.Size.Svc-e.Svc, r.Size.Svc),
			fmt.Sprintf("%d/%d", r.Size.Rs-e.Rs, r.Size.Rs),
			fmt.Sprintf("%d/%d", r.Size.Laddr-e.Laddr, r.Size.Laddr),
			fmt.Sprintf("%d/%d", r.Size.Conn-e.Conn, r.Size.Conn),
		)
	}
	return ret
}

type Vs_stats_q struct {
	Type     int
	Id       int /* worker id */
	Protocol uint8
	Addr     Be32
	Port     Be16
	Detail   uint8
}

func Get_stats_io(id int) (*Vs_stats_io_r, error) {
	args := Vs_stats_q{Type: VS_STATS_IO, Id: id}
	reply := &Vs_stats_io_r{}

	err := client.Call("stats", args, reply)
	return reply, err
}

func Get_stats_worker(id int) (*Vs_stats_worker_r, error) {
	args := Vs_stats_q{Type: VS_STATS_WORKER, Id: id}
	reply := &Vs_stats_worker_r{}

	err := client.Call("stats", args, reply)
	return reply, err
}

func Get_estats_worker(id int) (*Vs_estats_worker_r, error) {
	args := Vs_stats_q{Type: VS_ESTATS_WORKER, Id: id}
	reply := &Vs_estats_worker_r{}

	err := client.Call("stats", args, reply)
	return reply, err
}

func Get_stats_dev(id int) (*Vs_stats_dev_r, error) {
	args := Vs_stats_q{Type: VS_STATS_DEV, Id: id}
	reply := &Vs_stats_dev_r{}

	err := client.Call("stats", args, reply)
	return reply, err
}

func Get_stats_ctl() (*Vs_stats_ctl_r, error) {
	args := Vs_stats_q{Type: VS_STATS_CTL}
	reply := &Vs_stats_ctl_r{}

	err := client.Call("stats", args, reply)
	return reply, err
}

func Get_stats_mem() (*Vs_stats_mem_r, error) {
	args := Vs_stats_q{Type: VS_STATS_MEM}
	reply := &Vs_stats_mem_r{}

	err := client.Call("stats", args, reply)
	return reply, err
}

func Get_stats_vs(o *CmdOptions) (*Vs_stats_vs_r, error) {
	var reply Vs_stats_vs_r
	args := Vs_stats_q{
		Type:     VS_STATS_IP_VS_INFO,
		Id:       o.Id, /* worker id */
		Addr:     o.Addr.Ip,
		Port:     o.Addr.Port,
		Protocol: uint8(o.Protocol),
		Detail:   bool2u8(o.Print_detail),
	}

	if o.Print_all_worker == true {
		args.Id = -1
	}

	err := client.Call("stats", args, &reply)
	return &reply, err
}
