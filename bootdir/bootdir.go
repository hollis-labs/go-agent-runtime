// Package bootdir defines app-extensible boot directory planting helpers.
package bootdir

import (
	"errors"
	"fmt"

	agentlaunch "github.com/hollis-labs/go-agent-launch/agentlaunch"
)

var ErrUnsafePath = errors.New("bootdir: unsafe bootdir-relative path")

type File struct {
	RelPath string
	Content string
	Mode    uint32
}

type DryRunPlan struct {
	NativeFiles []agentlaunch.NativeFile
	Overlays    []File
}

type Request struct {
	Provider    string
	Runtime     agentlaunch.RuntimeKind
	NativeFiles []agentlaunch.NativeFile
	Overlays    map[string]string
}

// BuildInjection validates caller-supplied native files and overlays and
// returns the go-agent-launch InjectionSpec. Apps own the content renderers;
// this package owns the safety boundary and dry-run shape.
func BuildInjection(req Request) (agentlaunch.InjectionSpec, DryRunPlan, error) {
	out := agentlaunch.InjectionSpec{
		NativeFiles:    append([]agentlaunch.NativeFile(nil), req.NativeFiles...),
		BootDirOverlay: map[string]string{},
	}
	plan := DryRunPlan{NativeFiles: append([]agentlaunch.NativeFile(nil), req.NativeFiles...)}
	for i := range out.NativeFiles {
		if err := out.NativeFiles[i].Validate(); err != nil {
			return agentlaunch.InjectionSpec{}, DryRunPlan{}, err
		}
	}
	for rel, content := range req.Overlays {
		if err := ValidateRelPath(rel); err != nil {
			return agentlaunch.InjectionSpec{}, DryRunPlan{}, fmt.Errorf("%w: %s", err, rel)
		}
		out.BootDirOverlay[rel] = content
		plan.Overlays = append(plan.Overlays, File{RelPath: rel, Content: content, Mode: 0o644})
	}
	if len(out.BootDirOverlay) == 0 {
		out.BootDirOverlay = nil
	}
	return out, plan, nil
}

func ValidateRelPath(rel string) error {
	if err := agentlaunch.ValidateBootDirRelPath(rel); err != nil {
		return ErrUnsafePath
	}
	return nil
}

// TaskBundle returns raw native files under root. The caller owns body
// rendering; the helper pins Torque/Tether's safe task bundle planting shape.
func TaskBundle(root string, files map[string]string) ([]agentlaunch.NativeFile, error) {
	if root == "" {
		root = "tasks"
	}
	if err := ValidateRelPath(root + "/README.md"); err != nil {
		return nil, err
	}
	out := make([]agentlaunch.NativeFile, 0, len(files))
	for rel, content := range files {
		path := root + "/" + rel
		nf := agentlaunch.NativeFile{Kind: agentlaunch.NativeFileRaw, RelPath: path, Content: content}
		if err := nf.Validate(); err != nil {
			return nil, err
		}
		out = append(out, nf)
	}
	return out, nil
}
