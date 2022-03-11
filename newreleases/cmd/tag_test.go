// Copyright (c) 2022, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"context"

	"newreleases.io/newreleases"
)

type mockTagsService struct {
	tags []newreleases.Tag
	err  error
}

func newMockTagsService(tags []newreleases.Tag, err error) mockTagsService {
	return mockTagsService{tags: tags, err: err}
}

func (s mockTagsService) List(ctx context.Context) ([]newreleases.Tag, error) {
	return s.tags, s.err
}

func (s mockTagsService) Get(ctx context.Context, id string) (*newreleases.Tag, error) {
	if len(s.tags) == 0 {
		return nil, newreleases.ErrNotFound
	}
	return &s.tags[0], s.err
}

func (s mockTagsService) Add(ctx context.Context, name string) (*newreleases.Tag, error) {
	return &newreleases.Tag{
		ID:   "new",
		Name: name,
	}, s.err
}

func (s mockTagsService) Update(ctx context.Context, id, name string) (*newreleases.Tag, error) {
	if len(s.tags) == 0 {
		return nil, newreleases.ErrNotFound
	}
	tag := &s.tags[0]
	tag.Name = name
	return tag, s.err
}

func (s mockTagsService) Delete(ctx context.Context, id string) error {
	if len(s.tags) == 0 {
		return newreleases.ErrNotFound
	}
	return s.err
}

var fullTags = []newreleases.Tag{
	{
		ID:   "33f1db7254b9",
		Name: "Cool",
	},
	{
		ID:   "1d33b7254b9f",
		Name: "Awesome",
	},
}
