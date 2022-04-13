scan_exclude = [
  # This is to silence the error:
  # unicode control characters: vendor/sigs.k8s.io/kind/pkg/internal/env/term.go#L75 ['\u200d']
  r'/vendor/sigs\.k8s\.io/kind/pkg/internal/env/term\.go$'
]
