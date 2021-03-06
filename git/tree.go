// Copyright 2016 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitindex

import (
	"path/filepath"

	git "github.com/libgit2/git2go"
)

type repoWalker struct {
	*git.Repository

	repos      map[git.Oid]*git.Repository
	tree       map[string]git.Oid
	err        error
	submodules bool
}

func newRepoWalker(r *git.Repository, submodules bool) *repoWalker {
	return &repoWalker{
		Repository: r,
		tree:       map[string]git.Oid{},
		repos:      map[git.Oid]*git.Repository{},
		submodules: submodules,
	}
}

// TreeToFiles fetches the SHA1s for a tree, recursing into
// submodules. In addition, it returns a mapping that indicates in
// which repo each SHA1 can be found.
func TreeToFiles(r *git.Repository, t *git.Tree, submodules bool) (map[string]git.Oid, map[git.Oid]*git.Repository, error) {
	ref := newRepoWalker(r, submodules)
	t.Walk(ref.cbInt)
	return ref.tree, ref.repos, ref.err
}

func (r *repoWalker) cb(n string, e *git.TreeEntry) error {
	p := filepath.Join(n, e.Name)
	if e.Type == git.ObjectCommit && r.submodules {
		sub, err := r.Repository.Submodules.Lookup(p)
		if err != nil {
			return err
		}
		subRepo, err := sub.Open()
		if err != nil {
			return err
		}

		obj, err := subRepo.Lookup(e.Id)
		if err != nil {
			return err
		}
		treeObj, err := obj.Peel(git.ObjectTree)
		if err != nil {
			return err
		}
		tree, err := treeObj.AsTree()
		if err != nil {
			return err
		}
		subFiles, subRepos, err := TreeToFiles(subRepo, tree, r.submodules)
		if err != nil {
			return err
		}
		for k, v := range subRepos {
			r.repos[k] = v
		}
		for k, v := range subFiles {
			r.tree[filepath.Join(p, k)] = v
		}
		return nil
	}

	switch e.Filemode {
	case git.FilemodeBlob, git.FilemodeBlobExecutable:
	default:
		return nil
	}

	if e.Type != git.ObjectBlob {
		return nil
	}
	r.tree[p] = *e.Id
	r.repos[*e.Id] = r.Repository
	return nil
}

func (r *repoWalker) cbInt(n string, e *git.TreeEntry) int {
	err := r.cb(n, e)
	if err != nil {
		r.err = err
		return 1
	}
	return 0
}
