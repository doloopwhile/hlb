package git

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

type RawRemoteConfig interface {
}

type RawRemote interface {
	Config() *RawRemoteConfig
}

type Remote struct {
	URL         string
	ServiceHost string
	Owner       string
	RepoName    string
}

func NewRemote(remoteUrl string) (*Remote, error) {
	var assigned *regexp.Regexp
	if strings.HasPrefix(remoteUrl, "http") {
		assigned = regexp.MustCompile(`https?://[.+]?(.+)/(.+)/(.+)$`)
	} else if strings.HasPrefix(remoteUrl, "git://") {
		assigned = regexp.MustCompile(`git://[.+]?(.+)/(.+)/(.+)$`)
	} else if strings.HasPrefix(remoteUrl, "git@") {
		assigned = regexp.MustCompile(`git@(.+):(.+)/(.+).git`)
	} else {
		return nil, errors.New("unknown remote: " + remoteUrl)
	}

	result := assigned.FindStringSubmatch(remoteUrl)

	if result == nil || len(result) < 4 {
		return nil, errors.New("unknown remoteUrl pattern: " + remoteUrl)
	}
	hostNames := strings.Split(result[1], "@")
	serviceHost := hostNames[len(hostNames)-1]
	return &Remote{
		URL:         remoteUrl,
		ServiceHost: serviceHost,
		Owner:       result[2],
		RepoName:    strings.Replace(result[3], ".git", "", -1),
	}, nil
}

func GetDefaultRemote(path string) (*Remote, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, errors.Wrap(err, "Error occurred in git.GetDefaultRemote")
	}

	remote, err := r.Remote(git.DefaultRemoteName)
	if err != nil {
		return nil, errors.Wrap(err, "Error occurred in git.GetDefaultRemote")
	}
	return NewRemote(remote.Config().URL)

}

func SetRemote(path, remoteName, remoteUrl string) (*git.Remote, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, errors.Wrap(err, "Error occurred when open local git repository in git.SetRemote")
	}

	remote, err := r.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URL:  remoteUrl,
	})

	if err != nil {
		return nil, errors.Wrap(err, "Error occurred when remote creating in git.SetRemote")
	}

	_, err = r.Remotes()
	if err != nil {
		return nil, errors.Wrap(err, "Error occurred when remotes getting in git.SetRemote")
	}
	return remote, nil
}
