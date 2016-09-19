# Stash Commit Status Resource

[![Build Status](https://travis-ci.org/zabawaba99/stash-commit-status-resource.svg?branch=master)](https://travis-ci.org/zabawaba99/stash-commit-status-resource)

A concourse resource that will set a build status on your commits.

## Source Configuration

* `host`: *Required.* The host (including the port) of your stash instance.

* `username`: *Required.* Username for HTTP(S) auth when setting and retrieving
  the build status on a commit.

* `password`: *Required.* Password for HTTP(S) auth when setting and retrieving
  the build status on a commit.

* `project`: *Required.* The project the repository that is being tracked lives in.

* `repository`: *Required.* The repository to track commits on.

* `branch`: *Optional.* Scopes the tracked commits to this branch. Defaults to master.

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
    project: foo
    repository: bar
    branch: master

jobs:
- name: hello-world
  plan:
  - get: status
    trigger: true
  - put: status
    params:
      name: status
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
        name: status
        state: SUCCESSFUL
        description: "everything is ok"
    on_failure:
      put: status
      params:
        name: status
        state: FAILED
        description: "something went wrong"
```

## Behavior

### `check`: Check for new commits.

A request is made to the [stash API](https://developer.atlassian.com/static/rest/stash/3.11.6/stash-rest.html#idp2461680),
and any commits from the given version on are returned. If no version is given, the ref
for `HEAD` is returned.


### `in`: Fetch the build status of the commit in question.

Requests the [build status](https://developer.atlassian.com/static/rest/stash/3.11.6/stash-build-integration-rest.html#idp57632)
of the commit that is being ran against.

Writes out the commit sha to `<resource-name>/commit` which is used by the put step.

### `out`: Set the build status of a commit.

Make a request to the [stash api]()
to set the build status on a commit.

#### Parameters

* `status`: *Required* The name you gave to the stash-commit-status-resource

* `state`: *Required.* The state of the build. Must be [INPROGRESS, SUCCESSFUL, FAILURE]

* `description`: *Optional.* A message given context to the build status.
