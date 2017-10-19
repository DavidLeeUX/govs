/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package flags

import (
	"flag"
	"fmt"
	"os"
)

type my_flag struct {
	Name   int
	Action func(args interface{})
}

var Cmd my_flag

var FirstCmd *flag.FlagSet
var OthersCmd *flag.FlagSet

func Usage() {
	fmt.Printf("Commands:\n")
	FirstCmd.PrintDefaults()
	fmt.Printf("\nOptions:\n")
	OthersCmd.PrintDefaults()
}

func init() {
	Cmd.Name = 0
	Cmd.Action = nil

	FirstCmd = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	OthersCmd = flag.NewFlagSet("", flag.ExitOnError)
}
