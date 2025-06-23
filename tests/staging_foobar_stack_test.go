package tests

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFooBarLocalStack(t *testing.T) {
	t.Parallel()

	stackDir := "../live/staging/foo-bar"

	t.Logf("Running test in stackDir: %s", stackDir)

	// Ensure destroy runs at the end, even if the test fails
	t.Cleanup(func() {
		t.Log("Running terragrunt destroy as cleanup...")
		cmdDestroy := exec.Command("terragrunt", "--non-interactive", "stack", "run", "destroy", "--no-stack-generate")
		cmdDestroy.Dir = stackDir
		out, err := cmdDestroy.CombinedOutput()
		t.Logf("Destroy output:\n%s", out)
		require.NoError(t, err, "stack run destroy failed: %s", string(out))
	})

	// Generate
	t.Log("Running terragrunt stack generate...")
	cmdGenerate := exec.Command("terragrunt", "stack", "generate")
	cmdGenerate.Dir = stackDir
	out, err := cmdGenerate.CombinedOutput()
	t.Logf("Generate output:\n%s", out)
	require.NoError(t, err, "stack generate failed: %s", string(out))

	// Apply
	t.Log("Running terragrunt stack run apply...")
	cmdApply := exec.Command("terragrunt", "--non-interactive", "stack", "run", "apply", "--no-stack-generate")
	cmdApply.Dir = stackDir
	out, err = cmdApply.CombinedOutput()
	t.Logf("Apply output:\n%s", out)
	require.NoError(t, err, "stack run apply failed: %s", string(out))

	t.Log("Checking that output contains expected content...")
	require.Contains(t, string(out), `content = "Foo content: Hello from foo in live repo! (staging) (staging)"`)
}
