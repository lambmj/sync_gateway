//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package auth

import (
	"github.com/couchbaselabs/sync_gateway/base"
	ch "github.com/couchbaselabs/sync_gateway/channels"
)

// A Principal is an abstract object that can have access to channels.
type Principal interface {
	// The Principal's identifier.
	Name() string

	// The set of channels the Principal belongs to, and what sequence access was granted.
	Channels() ch.TimedSet

	// The channels the Principal was explicitly granted access to thru the admin API.
	ExplicitChannels() ch.TimedSet

	// Sets the explicit channels the Principal has access to.
	SetExplicitChannels(ch.TimedSet)

	// Returns true if the Principal has access to the given channel.
	CanSeeChannel(channel string) bool

	// If the Principal has access to the given channel, returns the sequence number at which
	// access was granted; else returns zero.
	CanSeeChannelSince(channel string) uint64

	// Returns an error if the Principal does not have access to all the channels in the set.
	AuthorizeAllChannels(channels base.Set) error

	// Returns an appropriate HTTPError for unauthorized access -- a 401 if the receiver is
	// the guest user, else 403.
	UnauthError(message string) error

	docID() string
	accessViewKey() string
	validate() error
	setChannels(ch.TimedSet)
}

// Role is basically the same as Principal, just concrete. Users can inherit channels from Roles.
type Role interface {
	Principal
}

// A User is a Principal that can log in and have multiple Roles.
type User interface {
	Principal

	// The user's email address.
	Email() string

	// Sets the user's email address.
	SetEmail(string) error

	// If true, the user is unable to authenticate.
	Disabled() bool

	// Sets the disabled property
	SetDisabled(bool)

	// Authenticates the user's password.
	Authenticate(password string) bool

	// Changes the user's password.
	SetPassword(password string)

	// The set of Roles the user belongs to (including ones given to it by)
	RoleNames() []string

	// The roles the user was explicitly granted access to thru the admin API.
	ExplicitRoleNames() []string

	// Sets the explicit roles the user belongs to.
	SetExplicitRoleNames([]string)

	// Every channel the user has access to, including those inherited from Roles.
	InheritedChannels() ch.TimedSet

	// If the input set contains the wildcard "*" channel, returns the user's InheritedChannels;
	// else returns the input channel list unaltered.
	ExpandWildCardChannel(channels base.Set) base.Set

	// Returns a TimedSet containing only the channels from the input set that the user has access
	// to, annotated with the sequence number at which access was granted.
	FilterToAvailableChannels(channels base.Set) ch.TimedSet

	setRoleNames([]string)
}
