package tests

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFooBarLocalStack(t *testing.T) {
    t.Parallel()

    stackDir := "../live/staging/foo-bar"

    // Ensure destroy runs at the end, even if the test fails
    t.Cleanup(func() {
        cmdDestroy := exec.Command("terragrunt", "--non-interactive", "stack", "run", "destroy", "--no-stack-generate")
        cmdDestroy.Dir = stackDir
        out, err := cmdDestroy.CombinedOutput()
        require.NoError(t, err, "stack run destroy failed: %s", string(out))
    })

    // Generate
    cmdGenerate := exec.Command("terragrunt", "stack", "generate")
    cmdGenerate.Dir = stackDir
    out, err := cmdGenerate.CombinedOutput()
    require.NoError(t, err, "stack generate failed: %s", string(out))

    // Apply
    cmdApply := exec.Command("terragrunt", "--non-interactive", "stack", "run", "apply", "--no-stack-generate")
    cmdApply.Dir = stackDir
    out, err = cmdApply.CombinedOutput()
    require.NoError(t, err, "stack run apply failed: %s", string(out))

    require.Contains(t, string(out), `content = "Foo content: Hello from foo in live repo! (staging) (staging)"`)

}