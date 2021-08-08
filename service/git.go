package service

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CheckoutRepo(dir, repoURL, branch string) error {
	repoName, err := getRepoName(repoURL)
	if err != nil {
		return err
	}
	fullDir := path.Join(dir, repoName)
	repo, err := openRepo(fullDir)
	if err != nil {
		return err
	}

	if repo == nil {
		_, err = cloneRepo(fullDir, repoURL, branch)
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed getting worktree for repo: %v, from path: %v, error: %v", repoURL, fullDir, err)
	}

	localRef := plumbing.NewBranchReferenceName(branch)
	remoteRef := plumbing.NewRemoteReferenceName("origin", branch)
	err = repo.Storer.SetReference(plumbing.NewSymbolicReference(localRef, remoteRef))
	if err != nil {
		return fmt.Errorf("failed setting reference for branch: %v for repo: %v in path: %v, error: %v", branch, repoURL, fullDir, err)
	}
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: localRef,
	})
	if err != nil {
		return fmt.Errorf("failed checking out branch: %v for repo: %v in path: %v, error: %v", branch, repoURL, fullDir, err)
	}

	err = pullBranch(localRef, worktree)
	if err != nil {
		return fmt.Errorf("failed getting submodules for branch: %v for repo: %v in path: %v, error: %v", branch, repoURL, fullDir, err)
	}

	err = updateSubmodules(worktree)
	if err != nil {
		return fmt.Errorf("failed updating submodules for branch: %v for repo: %v in path: %v, error: %v", branch, repoURL, fullDir, err)
	}
	return nil
}

func pullBranch(localRef plumbing.ReferenceName, worktree *git.Worktree) error {
	err := worktree.Pull(&git.PullOptions{
		ReferenceName:     localRef,
		RemoteName:        "origin",
		Progress:          os.Stdout,
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		if !strings.Contains(err.Error(), git.NoErrAlreadyUpToDate.Error()) {
			return fmt.Errorf("failed pulling changes for local ref: %v, error: %v", localRef, err)
		}
	}
	return nil
}

func updateSubmodules(worktree *git.Worktree) error {
	submodules, err := worktree.Submodules()
	if err != nil {
		return fmt.Errorf("failed getting submodules, error: %v", err)
	}
	err = submodules.Update(&git.SubmoduleUpdateOptions{
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return fmt.Errorf("failed updating submodules , error: %v", err)
	}
	return nil
}

func getBranchRef(repo *git.Repository, branch string) (*plumbing.Reference, error) {
	refs, err := repo.References()
	if err != nil {
		return nil, fmt.Errorf("failed getting references, err: %v", err)
	}
	defer refs.Close()
	for {
		ref, err := refs.Next()
		if err != nil {
			return nil, fmt.Errorf("failed getting reference for branch: %v, error: %v", branch, err)
		}
		if ref.Type() == plumbing.HashReference {
			segments := strings.Split(ref.Name().Short(), "/")
			if segments[len(segments)-1] == branch {
				return ref, nil
			}
		}
	}
}

func openRepo(path string) (*git.Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		if strings.Contains(err.Error(), git.ErrRepositoryNotExists.Error()) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return repo, nil
}

func cloneRepo(dir, repoURL, branch string) (*git.Repository, error) {
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               repoURL,
		Progress:          os.Stdout,
		ReferenceName:     plumbing.NewBranchReferenceName(branch),
		SingleBranch:      false,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to clone repo: %v and branch: %v, error: %v", repoURL, branch, err)
	}
	return repo, nil
}

func getRepoName(repoURL string) (string, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse repo url: %v, error: %v", repoURL, err)
	}
	segments := strings.Split(u.Path, "/")
	return segments[len(segments)-1], nil
}
