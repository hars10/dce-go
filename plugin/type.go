/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

//go:generate go-extpoints . ComposePlugin PodStatusHook
package plugin

import (
	"context"

	"github.com/mesos/mesos-go/executor"
	mesos "github.com/mesos/mesos-go/mesosproto"
)

type ComposePlugin interface {
	// Get the name of the plugin
	Name() string

	// execute some tasks before the Image is pulled
	LaunchTaskPreImagePull(ctx *context.Context, composeFiles *[]string, executorId string, taskInfo *mesos.TaskInfo) error

	// execute some tasks after the Image is pulled
	LaunchTaskPostImagePull(ctx *context.Context, composeFiles *[]string, executorId string, taskInfo *mesos.TaskInfo) error

	// execute the tasks after the pod is launched
	PostLaunchTask(ctx *context.Context, composeFiles []string, taskInfo *mesos.TaskInfo) (string, error)

	// execute the task before we send a Kill to Mesos
	PreKillTask(taskInfo *mesos.TaskInfo) error

	// execute the task after we send a Kill to Mesos
	PostKillTask(taskInfo *mesos.TaskInfo) error

	// execute the task to shutdown the pod
	Shutdown(executor.ExecutorDriver) error
}

// PodStatusHook allows custom implementations to be plugged when a Pod (mesos task) status changes. Currently this is
// designed to be executed on task status changes during LaunchTask.
type PodStatusHook interface {
	// Execute is invoked when the pod.taskStatusCh channel has a new status. It returns an error on failure,
	// and also a flag "failExec" indicating if the error needs to fail the execution when a series of hooks are executed
	// This is to support cases where a few hooks can be executed in a best effort manner and need not fail the executor
	Execute(podStatus string, data interface{}) (failExec bool, err error)
}
