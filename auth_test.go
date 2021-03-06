// Copyright 2015 tsuru-admin authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"net/http"

	"github.com/tsuru/tsuru/cmd"
	"github.com/tsuru/tsuru/cmd/cmdtest"
	"gopkg.in/check.v1"
)

func (s *S) TestListUsersInfo(c *check.C) {
	expected := &cmd.Info{
		Name:    "user-list",
		MinArgs: 0,
		Usage:   "user-list",
		Desc:    "List all users in tsuru.",
	}
	c.Assert((&listUsers{}).Info(), check.DeepEquals, expected)
}

func (s *S) TestListUsersRun(c *check.C) {
	var stdout, stderr bytes.Buffer
	context := cmd.Context{
		Stdout: &stdout,
		Stderr: &stderr,
	}
	manager := cmd.NewManager("glb", "0.2", "ad-ver", &stdout, &stderr, nil, nil)
	result := `[{"email": "test@test.com","teams":["team1", "team2", "team3"]}]`
	trans := cmdtest.ConditionalTransport{
		Transport: cmdtest.Transport{Message: result, Status: http.StatusOK},
		CondFunc: func(req *http.Request) bool {
			return req.Method == "GET" && req.URL.Path == "/users"
		},
	}
	expected := `+---------------+---------------------+
| User          | Teams               |
+---------------+---------------------+
| test@test.com | team1, team2, team3 |
+---------------+---------------------+
`
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	command := listUsers{}
	err := command.Run(&context, client)
	c.Assert(err, check.IsNil)
	c.Assert(stdout.String(), check.Equals, expected)
}
