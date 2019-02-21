{% if include.header %}
{% assign header = include.header %}
{% else %}
{% assign header = "###" %}
{% endif %}
Cleanup project images from images repo

{{ header }} Syntax

```bash
werf images cleanup [options]
```

{{ header }} Options

```bash
      --dir='':
            Change to the specified directory to find werf.yaml config
      --disable-pretty-log=false:
            Disable emojis, auto line wrapping and replace log process border characters with 
            spaces (default $WERF_DISABLE_PRETTY_LOG).
      --docker-config='':
            Specify docker config directory path. Default $WERF_DOCKER_CONFIG or $DOCKER_CONFIG or 
            ~/.docker (in the order of priority).
            Command needs granted permissions to delete images from the specified images repo.
      --dry-run=false:
            Indicate what the command would do without actually doing that
      --git-commit-strategy-expiry-days=-1:
            Keep images published with the git-commit tagging strategy in the images repo for the 
            specified maximum days since image published. Republished image will be kept specified 
            maximum days since new publication date. No days limit by default, -1 disables the 
            limit. Value can be specified by the $WERF_GIT_COMMIT_STRATEGY_EXPIRY_DAYS.
      --git-commit-strategy-limit=-1:
            Keep max number of images published with the git-commit tagging strategy in the images 
            repo. No limit by default, -1 disables the limit. Value can be specified by the 
            $WERF_GIT_COMMIT_STRATEGY_LIMIT.
      --git-tag-strategy-expiry-days=-1:
            Keep images published with the git-tag tagging strategy in the images repo for the 
            specified maximum days since image published. Republished image will be kept specified 
            maximum days since new publication date. No days limit by default, -1 disables the 
            limit. Value can be specified by the $WERF_GIT_TAG_STRATEGY_EXPIRY_DAYS.
      --git-tag-strategy-limit=-1:
            Keep max number of images published with the git-tag tagging strategy in the images 
            repo. No limit by default, -1 disables the limit. Value can be specified by the 
            $WERF_GIT_TAG_STRATEGY_LIMIT.
  -h, --help=false:
            help for cleanup
      --home-dir='':
            Use specified dir to store werf cache files and dirs (default $WERF_HOME or ~/.werf)
  -i, --images-repo='':
            Docker Repo to store images (default $WERF_IMAGES_REPO)
      --insecure-repo=false:
            Allow usage of insecure docker repos
      --kube-config='':
            Kubernetes config file path
      --kube-context='':
            Kubernetes config context (default $WERF_KUBE_CONTEXT)
      --log-color-mode='auto':
            Set log color mode.
            Supported on, off and auto (based on the stdout's file descriptor referring to a 
            terminal) modes.
            Default $WERF_LOG_COLOR_MODE or auto mode.
      --tmp-dir='':
            Use specified dir to store tmp files and dirs (default $WERF_TMP or system tmp dir)
      --without-kube=false:
            Do not skip deployed kubernetes images
```
