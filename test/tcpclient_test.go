// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.

package test

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/dafanshu/modbus"
)

const (
	tcpDevice = "localhost:5020"
)

func TestTCPClient(t *testing.T) {
	client := modbus.TCPClient(tcpDevice)
	ClientTestAll(t, client)
}

func TestTCPClientAdvancedUsage(t *testing.T) {
	handler := modbus.NewTCPClientHandler(tcpDevice)
	handler.Timeout = 5 * time.Second
	handler.SlaveId = 1
	handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
	handler.Connect()
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadDiscreteInputs(15, 2)
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
	results, err = client.WriteMultipleRegisters(1, 2, []byte{0, 3, 0, 4})
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
	results, err = client.WriteMultipleCoils(5, 10, []byte{4, 3})
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
}


func TestSend(t *testing.T) {
	var wg sync.WaitGroup
	handler := modbus.NewTCPClientHandler("192.168.1.168:502")
	handler.Timeout = 30 * time.Millisecond
	handler.SlaveId = 0x00
	//handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	err := handler.Connect()
	defer handler.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	client := modbus.NewClient(handler)
	index := 1
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			results, err := client.ReadCoils(96, 1)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(results)
			wg.Done()
		}()
		if index<2 {
			wg.Add(1)
			go func() {
				results, err := client.ReadHoldingRegisters(0, 1)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(results)
				wg.Done()
				index++
			}()
		}
		wg.Add(1)
		go func() {
			results, err := client.ReadHoldingRegisters(609, 1)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(results)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			results, err := client.ReadHoldingRegisters(555, 1)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(results)
			wg.Done()
		}()
		wg.Wait()
	}

}


func TestBlankRegister(t *testing.T) {
	var wg sync.WaitGroup
	handler := modbus.NewTCPClientHandler("127.0.0.1:503")
	handler.Timeout = 30 * time.Millisecond
	handler.SlaveId = 0x00
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	err := handler.Connect()
	defer handler.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	client := modbus.NewClient(handler)
	wg.Add(1)
	go func() {
		// 693
		results, err := client.ReadHoldingRegisters(693, 3)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(results)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		// 693
		results, err := client.ReadHoldingRegisters(693, 2)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(results)
		wg.Done()
	}()
	wg.Wait()
}
