/**
 * Copyright 2019 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ------------------------------------------------------------------------------
 */

package main 

import (
	flags "github.com/jessevdk/go-flags"
	"fmt"
	"os"
	"strconv"
	"time"
)

// All subcommands implement this interface
type Command interface {
	Register(*flags.Command) error
	Name() string
	Run() error
	Compute(data []byte) error
}

type Opts struct {
	Version bool   `short:"V" long:"version" description:"Display version information"`
}

type CryptoAlgorithm struct {}

var DISTRIBUTION_VERSION string

func (c CryptoAlgorithm) Run error {
	pid := os.Getpid()
	fmt.Println("Start performance measuring tool against the process id: ", strconv.Itoa(pid))
	fmt.Println("Then press [ENTER] key to continue!")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	
	// Create data of uint8 to perform crypto algorithm
	var data [NUMBER_OF_INPUT_BYTES]byte
	start := time.Now()
	result, err := c.Compute(data)
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Total time taken for crypto operation: ", elapsed, "ms")
	fmt.Println("Result: ", result)
	
	return err
}

func init() {
	if len(DISTRIBUTION_VERSION) == 0 {
		DISTRIBUTION_VERSION = "Unknown"
	}
}

func main() {
	arguments := os.Args[1:]
	for _, arg := range arguments {
		if arg == "-V" || arg == "--version" {
			fmt.Println(DISTRIBUTION_NAME + " version " + DISTRIBUTION_VERSION)
			os.Exit(0)
		}
	}
	
	var opts Opts
	parser := flags.NewParser(&opts, flags.Default)
	parser.Command.Name = "go-crypto-bmark"
	
	commands := []Command{
		&Sha256{},
		// &Sha512{},
	}
	
	for _, cmd := range commands {
		err := cmd.Register(parser.Command)
		if err != nil {
			logger.Errorf("Couldn't register command %v: %v", cmd.Name(), err)
			os.Exit(1)
		}
	}
	
	remaining, err := parser.Parse()
	if e, ok := err.(*flags.Error); ok {
		if e.Type == flags.ErrHelp {
			return
		} else {
			os.Exit(1)
		}
	}
	
	if len(remaining) > 0 {
		fmt.Println("Error: Unrecognized arguments passed: ", remaining)
		os.Exit(2)
	}
	
	if parser.Command.Active == nil {
		os.Exit(2)
	}
	
	name := parser.Command.Active.Name
	for _, cmd := range commands {
		if cmd.Name() == name {
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
		}
	}
	
	fmt.Println("Error: Command not found: ", name)
}
