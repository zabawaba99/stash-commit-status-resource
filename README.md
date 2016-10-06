# Stash Commit Status Resource

[![Build Status](https://travis-ci.org/zabawaba99/stash-commit-status-resource.svg?branch=master)](https://travis-ci.org/zabawaba99/stash-commit-status-resource)

### Under Development. API may/will change.

A concourse resource that will set a build status on your commits.

## Source Configuration

* `host`: *Required.* The host (including the port) of your stash instance.

* `username`: *Required.* Username for HTTP(S) auth when setting and retrieving
  the build status on a commit.

* `password`: *Required.* Password for HTTP(S) auth when setting and retrieving
  the build status on a commit.

* `retry_attempts`: *Optional.* Number of times to attempt to set the build status.
  Using basic auth without SSL may cause an issue with the stash API where a retry
  needed.

### Example

```yaml
resource_types:
- name: commit-status
  type: docker-image
  source:
    repository: zabawaba99/stash-commit-status-resource

resources:
- name: status
  type: commit-status
  source:
    host: http://10.0.0.5:7990
    username: username
    password: password
- name: src
  type: git
  source:
    uri: git@github.com:zabawaba99/stash-commit-status-resource.git
    branch: master
    private_key: |
      -----BEGIN RSA PRIVATE KEY-----
      MIIEowIBAAKCAQEAtCS10/f7W7lkQaSgD/mVeaSOvSF9ql4hf/zfMwfVGgHWjj+W
      <Lots more text>
      DWiJL+OFeg9kawcUL6hQ8JeXPhlImG6RTUffma9+iGQyyBMCGd1l
      -----END RSA PRIVATE KEY-----

jobs:
- name: hello-world
  plan:
  - get: src
    trigger: true
  - put: status
    params:
      repository: src
      state: INPROGRESS
      description: "starting build"
  - task: foo
    config:
      platform: linux
      image_resource:
        type: docker-image
        source:
          repository: ruby
          tag: '2.1'
      run:
        path: ruby
        args:
          - -e
          - 'fail "nope" if Random.new.rand(1..100) % 2 == 0'
    on_success:
      put: status
      params:
        repository: src
        state: SUCCESSFUL
        description: "everything is ok"
    on_failure:
      put: status
      params:
        repository: src
        state: FAILED
        description: "something went wrong"
```

## Behavior

### `check`: Check for new commits.

Not implemented.

### `in`: Fetch the build status of the commit in question.

Not implemented.

### `out`: Set the build status of a commit.

Make a request to the [stash api]()
to set the build status on a commit.

#### Parameters

* `repository`: *Required* A git repository to get the last commit off.

* `state`: *Required.* The state of the build. Must be [INPROGRESS, SUCCESSFUL, FAILURE].

* `description`: *Optional.* A message given context to the build status.
