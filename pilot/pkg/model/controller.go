// Copyright 2017 Istio Authors
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

package model

// Controller defines an event controller loop. These are the events that a
// platform registry is expected to generate in response to changes such as
// addition or removal of pods/vms/cfapps, services, etc. Unlike the dataplane,
// platform service registries can only emit two classes of
// events: ServiceChange, EndpointChange
//
// The "change" could mean add/delete/update.
//
// The controller guarantees the following consistency requirement: registry
// view in the controller is as AT LEAST as fresh as the moment notification
// arrives, but MAY BE more fresh (e.g. "delete" cancels an "add" event).  For
// example, an event for a service creation will see a service registry without
// the service if the event is immediately followed by the service deletion
// event.
//
// Handlers execute on the single worker queue in the order they are appended.
// Handlers receive the notification event and the associated object.  Note
// that all handlers must be appended before starting the controller.
type Controller interface {
	// ServiceChange notifies about changes to the list of services
	ServiceChange(f func(*Service, Event)) error

	// EndpointChange notifies about changes to the endpoints associated with a service.
	EndpointChange(f func(*ServiceInstance, Event)) error

	// Run until a signal is received
	Run(stop <-chan struct{})
}

// Event represents a registry update event
type Event int

const (
	// EventAdd is sent when an object is added
	EventAdd Event = iota

	// EventUpdate is sent when an object is modified
	// Captures the modified object
	EventUpdate

	// EventDelete is sent when an object is deleted
	// Captures the object at the last known state
	EventDelete
)

func (event Event) String() string {
	out := "unknown"
	switch event {
	case EventAdd:
		out = "add"
	case EventUpdate:
		out = "update"
	case EventDelete:
		out = "delete"
	}
	return out
}
