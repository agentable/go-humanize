package humanize_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// releaseFixture executes the actual scripts against disposable Git repositories.
type releaseFixture struct {
	t      *testing.T
	dir    string
	remote string
}

func newReleaseFixture(t *testing.T) releaseFixture {
	t.Helper()
	dir := t.TempDir()
	r := releaseFixture{t: t, dir: dir, remote: filepath.Join(t.TempDir(), "origin.git")}
	if err := os.Mkdir(filepath.Join(dir, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"check_changes.sh", "tag_release.sh"} {
		data, err := os.ReadFile(filepath.Join("scripts", name))
		if err != nil {
			t.Fatal(err)
		}
		r.write("scripts/"+name, string(data))
	}
	r.write("go.mod", "module example.invalid/release\n\ngo 1.27.0\n")
	r.write("main.go", "package release\n")
	r.write("README.md", "Release fixture\n")
	r.git("init", "-b", "main")
	r.git("config", "user.name", "Release Test")
	r.git("config", "user.email", "release@example.invalid")
	r.git("config", "commit.gpgsign", "false")
	r.git("config", "tag.gpgsign", "false")
	r.git("add", ".")
	r.git("commit", "-m", "baseline")
	return r
}

func (r releaseFixture) write(name, content string) {
	r.t.Helper()
	if err := os.WriteFile(filepath.Join(r.dir, name), []byte(content), 0o644); err != nil {
		r.t.Fatal(err)
	}
}

func (r releaseFixture) run(command string, args ...string) (string, error) {
	r.t.Helper()
	cmd := exec.Command(command, args...)
	cmd.Dir = r.dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+os.DevNull, "GIT_TERMINAL_PROMPT=0")
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func (r releaseFixture) git(args ...string) string {
	r.t.Helper()
	output, err := r.run("git", args...)
	if err != nil {
		r.t.Fatalf("git %v: %v\n%s", args, err, output)
	}
	return output
}

func TestCheckChanges(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		file    string
		content string
		commit  bool
		remove  bool
		noTag   bool
		want    int
	}{
		{name: "no tag", noTag: true, want: 0},
		{name: "unchanged", want: 1},
		{name: "module", file: "go.mod", content: "module example.invalid/changed\n\ngo 1.27.0\n", commit: true},
		{name: "go version", file: "go.mod", content: "module example.invalid/release\n\ngo 1.27.1\n", commit: true},
		{name: "dependency", file: "go.mod", content: "module example.invalid/release\n\ngo 1.27.0\n\nrequire example.invalid/dep v1.0.0\n", commit: true},
		{name: "sum added", file: "go.sum", content: "new sum\n", commit: true},
		{name: "sum changed", file: "go.sum", content: "changed sum\n", commit: true},
		{name: "sum removed", file: "go.sum", commit: true, remove: true},
		{name: "source", file: "main.go", content: "package changed\n", commit: true},
		{name: "readme", file: "README.md", content: "changed\n", commit: true},
		{name: "uncommitted module", file: "go.mod", content: "module dirty\n", want: 1},
		{name: "uncommitted sum", file: "go.sum", content: "dirty sum\n", want: 1},
		{name: "uncommitted source", file: "main.go", content: "package dirty\n", want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := newReleaseFixture(t)
			if tt.name == "sum changed" || tt.remove {
				r.write("go.sum", "original sum\n")
				r.git("add", "go.sum")
				r.git("commit", "-m", "initial sum")
			}
			if !tt.noTag {
				r.git("tag", "v0.1.0")
			}
			if tt.file != "" {
				if tt.remove {
					r.git("rm", tt.file)
				} else {
					r.write(tt.file, tt.content)
				}
				if tt.commit {
					r.git("add", "-A")
					r.git("commit", "-m", "change")
				}
			}
			output, err := r.run("bash", "scripts/check_changes.sh")
			code := 0
			if err != nil {
				var exit *exec.ExitError
				if errors.As(err, &exit) {
					code = exit.ExitCode()
				} else {
					t.Fatal(err)
				}
			}
			if code != tt.want {
				t.Errorf("check_changes exit = %d, want %d; output: %s", code, tt.want, output)
			}
		})
	}
}

func TestTagReleaseCommit(t *testing.T) {
	t.Parallel()

	for _, mode := range []string{"main", "feature", "detached", "invalid", "existing"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			r := newReleaseFixture(t)
			r.git("init", "--bare", r.remote)
			r.git("remote", "add", "origin", r.remote)
			r.git("push", "origin", "main")
			baseline := r.git("rev-parse", "HEAD")
			version := "0.2.0"
			switch mode {
			case "feature":
				r.git("checkout", "-b", "feature")
			case "detached":
				r.git("checkout", "--detach")
			case "invalid":
				version = "bad-version"
			case "existing":
				r.git("tag", "v0.2.0")
			}
			if mode != "existing" {
				r.write("README.md", "release change\n")
				r.git("add", "README.md")
				r.git("commit", "-m", "release change")
			}
			head := r.git("rev-parse", "HEAD")
			localBefore := r.git("show-ref")
			remoteBefore := r.git("ls-remote", "origin")
			output, err := r.run("bash", "scripts/tag_release.sh", version)
			if mode == "main" {
				if err != nil {
					t.Fatalf("release failed: %v\n%s", err, output)
				}
				if got := r.git("rev-parse", "v0.2.0^{}"); got != head {
					t.Errorf("tag commit = %s, want %s", got, head)
				}
				if got := r.git("ls-remote", "origin", "refs/heads/main"); got != head+"\trefs/heads/main" {
					t.Errorf("remote main = %s, want commit %s (baseline %s)", got, head, baseline)
				}
				if got := r.git("ls-remote", "origin", "refs/tags/v0.2.0^{}"); got != head+"\trefs/tags/v0.2.0^{}" {
					t.Errorf("remote tag = %s, want commit %s", got, head)
				}
				return
			}
			if err == nil {
				t.Errorf("%s release succeeded, want rejection; output: %s", mode, output)
			}
			if got := r.git("show-ref"); got != localBefore {
				t.Errorf("rejected release changed local refs:\n%s\nwas:\n%s", got, localBefore)
			}
			if got := r.git("ls-remote", "origin"); got != remoteBefore {
				t.Errorf("rejected release changed remote refs:\n%s\nwas:\n%s", got, remoteBefore)
			}
		})
	}
}

func TestTagReleasePush(t *testing.T) {
	t.Parallel()

	for _, mode := range []string{"success", "reject tag", "reject main", "conflicting tag", "non fast forward", "no atomic support"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			r := newReleaseFixture(t)
			r.git("init", "--bare", r.remote)
			r.git("remote", "add", "origin", r.remote)
			r.git("push", "origin", "main")
			if mode == "conflicting tag" {
				r.git("tag", "v0.2.0")
				r.git("push", "origin", "refs/tags/v0.2.0")
				r.git("tag", "-d", "v0.2.0")
			}
			if mode == "non fast forward" {
				r.write("README.md", "remote change\n")
				r.git("add", ".")
				r.git("commit", "-m", "remote change")
				r.git("push", "origin", "main")
				r.git("reset", "--hard", "HEAD~1")
			}
			if mode == "reject tag" || mode == "reject main" {
				ref := "refs/tags/v0.2.0"
				if mode == "reject main" {
					ref = "refs/heads/main"
				}
				hook := "#!/usr/bin/env bash\nwhile read -r old new ref; do\n  if [[ $ref == " + ref + " ]]; then exit 1; fi\ndone\n"
				if err := os.WriteFile(filepath.Join(r.remote, "hooks", "pre-receive"), []byte(hook), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "no atomic support" {
				r.git("--git-dir="+r.remote, "config", "receive.advertiseAtomic", "false")
			}
			r.git("tag", "-a", "unrelated-local", "-m", "unrelated tag")
			r.git("config", "push.followTags", "true")
			r.write("README.md", "new release\n")
			r.git("add", ".")
			r.git("commit", "-m", "new release")
			head := r.git("rev-parse", "HEAD")
			before := r.git("ls-remote", "origin")
			output, err := r.run("bash", "scripts/tag_release.sh", "0.2.0")
			if mode == "success" {
				if err != nil {
					t.Fatalf("release failed: %v\n%s", err, output)
				}
				for _, ref := range []string{"refs/heads/main", "refs/tags/v0.2.0^{}"} {
					if got := r.git("ls-remote", "origin", ref); got != head+"\t"+ref {
						t.Errorf("remote %s = %s, want %s", ref, got, head)
					}
				}
			} else {
				if err == nil {
					t.Errorf("%s release succeeded, want rejection; output: %s", mode, output)
				}
				if got := r.git("ls-remote", "origin"); got != before {
					t.Errorf("failed release partially changed remote:\n%s\nwas:\n%s", got, before)
				}
				if got := r.git("rev-parse", "v0.2.0^{}"); got != head {
					t.Errorf("local failed-release tag = %s, want %s", got, head)
				}
				if !strings.Contains(output, "not published") || !strings.Contains(output, "local tag v0.2.0") {
					t.Errorf("missing failure/retry information: %s", output)
				}
			}
			if got := r.git("ls-remote", "origin", "refs/tags/unrelated-local"); got != "" {
				t.Errorf("unrelated tag published: %s", got)
			}
		})
	}
}
