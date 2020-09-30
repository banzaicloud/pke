// Copyright © 2020 Banzai Cloud
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

package pipeline

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/banzaicloud/pke/.gen/pipeline"
	"github.com/banzaicloud/pke/cmd/pke/app/constants"
	"github.com/spf13/cobra"
)

type pipelineStatusReporter struct {
	client       *pipeline.APIClient
	disabled     bool
	clusterID    int32
	orgID        int32
	nodeName     string
	nodePoolName string
	ip           string
	processId    string

	cmd    *cobra.Command
	output io.Writer
}

func NewPipelineStatusReporter() *pipelineStatusReporter {
	return &pipelineStatusReporter{}
}

func (p *pipelineStatusReporter) Init(cmd *cobra.Command) error {
	if p.cmd != nil {
		return nil
	}
	p.cmd = cmd
	p.disabled = true
	if !Enabled(p.cmd) {
		return nil
	}
	endpoint, token, insecure, orgID, clusterID, err := CommandArgs(p.cmd)
	if err != nil {
		return err
	}
	p.output = p.cmd.OutOrStdout()
	p.client = Client(p.output, endpoint, token, insecure)
	if p.client == nil {
		return nil
	}
	p.orgID = orgID
	p.clusterID = clusterID
	p.disabled = false

	p.nodePoolName, _ = p.cmd.Flags().GetString(constants.FlagPipelineNodepool)
	p.nodeName, _ = os.Hostname()

	return nil
}

type StatusReporter interface {
	ReportStep(phase string, message string) error
	ReportError(phase string, err error) error
	ReportResult(phase string, err error, finished bool) error
}

type StatusReporterIniter interface {
	StatusReporter
	Init(cmd *cobra.Command) error
}

func (p *pipelineStatusReporter) ReportStep(phase string, message string) error {
	return p.report(pipeline.RUNNING, false, phase, message)
}

func (p *pipelineStatusReporter) ReportError(phase string, err error) error {
	return p.report(pipeline.FAILED, false, phase, err.Error())
}

func (p *pipelineStatusReporter) ReportResult(phase string, err error, final bool) error {
	message := "success"
	status := pipeline.FINISHED
	if err != nil {
		message = err.Error()
		status = pipeline.FAILED
	}
	return p.report(status, final, phase, message)
}

func (p *pipelineStatusReporter) report(status pipeline.ProcessStatus, final bool, phase string, message string) error {
	if p == nil || p.disabled {
		return nil
	}
	if p.client == nil {
		return errors.New("reporter uninitialized")
	}

	now := time.Now()
	req := pipeline.ReportPkeNodeStatusRequest{
		Name:      p.nodeName,
		NodePool:  p.nodePoolName,
		Ip:        p.ip,
		Status:    status,
		Final:     final,
		Phase:     phase,
		Message:   message,
		Timestamp: &now,
		ProcessId: p.processId,
	}
	resp, _, err := p.client.ClustersApi.ReportPKENodeStatus(context.TODO(), p.orgID, p.clusterID, req)

	if p.processId == "" {
		p.processId = resp.ProcessId
	}

	return err
}
