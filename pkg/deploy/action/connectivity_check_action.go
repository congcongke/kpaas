// Copyright 2019 Shanghai JingDuo Information Technology co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package action

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kpaas-io/kpaas/pkg/deploy/consts"
	"github.com/kpaas-io/kpaas/pkg/deploy/machine/ssh"
	pb "github.com/kpaas-io/kpaas/pkg/deploy/protos"
)

const ActionTypeConnectivityCheck Type = "ConnectivityCheck"

// ConnectivityCheckItem an item representing one check item of checking wheter a node can connect to another by the protocol and port.
type ConnectivityCheckItem struct {
	Protocol    consts.Protocol
	Port        uint16
	CheckResult *pb.ItemCheckResult
}

// ConnectivityCheckActionConfig configuration of checking connectivity from soruce to destination.
type ConnectivityCheckActionConfig struct {
	SourceNode             *pb.Node
	DestinationNode        *pb.Node
	ConnectivityCheckItems []ConnectivityCheckItem
	LogFileBasePath        string
}

type ConnectivityCheckAction struct {
	Base

	SourceNode      *pb.Node
	DestinationNode *pb.Node
	CheckItems      []ConnectivityCheckItem
}

// NewConnectivityCheckAction creates an action to check connectivity from soruce to destination.
func NewConnectivityCheckAction(cfg *ConnectivityCheckActionConfig) (Action, error) {
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	if cfg == nil {
		err = fmt.Errorf("action config is nil")
		return nil, err
	}
	if cfg.SourceNode == nil {
		err = fmt.Errorf("source node in config is nil")
		return nil, err
	}
	if cfg.DestinationNode == nil {
		err = fmt.Errorf("destination node in config is nil")
		return nil, err
	}
	actionName := GenActionName(ActionTypeConnectivityCheck)
	return &ConnectivityCheckAction{
		Base: Base{
			Name:              actionName,
			ActionType:        ActionTypeConnectivityCheck,
			Status:            ActionPending,
			LogFilePath:       GenActionLogFilePath(cfg.LogFileBasePath, actionName, cfg.SourceNode.GetName()),
			CreationTimestamp: time.Now(),
			Node:              cfg.SourceNode,
		},
		SourceNode:      cfg.SourceNode,
		DestinationNode: cfg.DestinationNode,
		CheckItems:      cfg.ConnectivityCheckItems,
	}, nil
}

func init() {
	RegisterExecutor(ActionTypeConnectivityCheck, new(connectivityCheckExecutor))
}

type connectivityCheckExecutor struct{}

func (e *connectivityCheckExecutor) Execute(act Action) *pb.Error {
	connectivityCheckAction, ok := act.(*ConnectivityCheckAction)
	if !ok {
		return errOfTypeMismatched(new(ConnectivityCheckAction), act)
	}

	dstNode := connectivityCheckAction.DestinationNode
	srcNode := connectivityCheckAction.SourceNode
	// start SSH connection to destination node to dump packets
	sshClientDst, err := ssh.NewClient(dstNode.Ssh.Auth.Username, dstNode.Ip, dstNode.Ssh)
	if err != nil {
		return &pb.Error{
			Reason: "failed to start SSH client",
			Detail: fmt.Sprintf("Failed to create SSH connetion to %s by connecting to %s:%d, error %v",
				dstNode.Name, dstNode.Ip, dstNode.Ssh.Port, err),
			FixMethods: "configure no-password ssh login from deploy node",
		}
	}

	// start SSH connection to source node to send packets
	sshClientSrc, err := ssh.NewClient(srcNode.Ssh.Auth.Username, srcNode.Ip, srcNode.Ssh)
	if err != nil {
		return &pb.Error{
			Reason: "failed to start SSH client",
			Detail: fmt.Sprintf("Failed to create SSH connetion to %s by connecting to %s:%d, error %v",
				srcNode.Name, srcNode.Ip, srcNode.Ssh.Port, err),
			FixMethods: "configure no-password ssh login from deploy node",
		}
	}

	for _, checkItem := range connectivityCheckAction.CheckItems {
		randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
		srcPort := (randGen.Uint32() % 16384) + 45000
		sshSessionDst, _ := sshClientDst.NewSession()
		sshSessionSrc, _ := sshClientSrc.NewSession()
		captureCommand := []string{"timeout", "5",
			"tcpdump", "-nni", "any", "-c", "1",
			"src", srcNode.Ip, "and", "dst", dstNode.Ip,
		}
		sendCommand := []string{"nc", "-p", fmt.Sprintf("%d", srcPort),
			"-s", srcNode.Ip}
		switch checkItem.Protocol {
		case consts.ProtocolTCP:
			captureCommand = append(captureCommand, "and", "tcp",
				"dst", "port", "dst", "port", fmt.Sprintf("%d", checkItem.Port),
				"and", "src", "port", fmt.Sprintf("%d", srcPort))
			sendCommand = append(sendCommand, "-zv",
				dstNode.Ip, fmt.Sprintf("%d", checkItem.Port))
		case consts.ProtocolUDP:
			captureCommand = append(captureCommand, "and", "udp",
				"dst", "port", "dst", "port", fmt.Sprintf("%d", checkItem.Port),
				"and", "src", "port", fmt.Sprintf("%d", srcPort))
			sendCommand = append(sendCommand, "-zuv",
				dstNode.Ip, fmt.Sprintf("%d", checkItem.Port))
		default:
			return &pb.Error{
				Reason: "protocol not supported",
				Detail: fmt.Sprintf("protocol %s is not supported. supported protocols are: TCP, UDP",
					string(checkItem.Protocol)),
				FixMethods: "Use a supported protocol",
			}
		}
		if checkItem.CheckResult != nil {
			checkItem.CheckResult.Status = ItemActionDoing
		}
		// first, start the capturing on destination node.
		sshSessionDst.Start(strings.Join(captureCommand, " "))
		captureChan := make(chan error)
		go func(errCh chan error) {
			errCh <- sshSessionDst.Wait()
		}(captureChan)

		// sleep one second to make sure that the packet is sent after capturing started
		time.Sleep(time.Second)
		sshSessionSrc.Start(strings.Join(sendCommand, " "))
		err := <-captureChan

		if err != nil {
			checkErr := &pb.Error{
				Reason: "check connectivity failed",
				Detail: fmt.Sprintf("%s cannot connect to %s %s:%d",
					srcNode.Name, string(checkItem.Protocol), dstNode.Name, checkItem.Port),
				FixMethods: "configure network or firewall to allow these packets",
			}
			if checkItem.CheckResult != nil {
				checkItem.CheckResult.Status = ItemActionFailed
				checkItem.CheckResult.Err = checkErr
			}
			return checkErr
			// continue to check other items
		} else {
			if checkItem.CheckResult != nil {
				checkItem.CheckResult.Status = ItemActionDone
			}
		}
	}
	return nil
}
