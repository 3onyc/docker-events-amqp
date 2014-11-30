package main

import (
	docker "github.com/fsouza/go-dockerclient"
)

func connectDocker(url string) (*docker.Client, error) {
	dockerClient, err := docker.NewClient(url)
	if err != nil {
		return nil, err
	}

	if err := dockerClient.Ping(); err != nil {
		return nil, err
	}

	return dockerClient, nil
}
