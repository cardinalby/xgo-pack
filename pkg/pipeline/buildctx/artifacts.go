package buildctx

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cardinalby/xgo-pack/pkg/consts"
	"github.com/cardinalby/xgo-pack/pkg/util/logging"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Artifact interface {
	GetPath() string
}

type DisposableArtifact interface {
	Artifact
	Disposable
}

type Kind string

const (
	KindIcon                      Kind = "icon"
	KindDefaultPngIcon            Kind = "default_png_icon"
	KindPngIcon                   Kind = "png_icon"
	KindIcoIcon                   Kind = "ico_icon"
	KindWinManifest               Kind = "win_manifest"
	KindMacosPlist                Kind = "macos_plist"
	KindMacosIconSet              Kind = "macos_icon_set"
	KindLinuxDesktopEntry         Kind = "linux_desktop_entry"
	KindMacosCreateDmgDockerImage Kind = "macos_create_dmg_docker_image"
)

func WinSysoKind(arch consts.Arch) Kind {
	return Kind("win_syso_" + arch)
}

func BinKind(os consts.Os, arch consts.Arch) Kind {
	return Kind(fmt.Sprintf("bin_%s_%s", os, arch))
}

func MacosBundleKind(arch consts.Arch) Kind {
	return Kind("macos_bundle_" + arch)
}

func MacosDmgKind(arch consts.Arch) Kind {
	return Kind("macos_dmg_" + arch)
}

func LinuxDebKind(arch consts.Arch) Kind {
	return Kind("linux_deb_" + arch)
}

type artifactEntry struct {
	mu     sync.RWMutex
	result Artifact
	err    error
}

type Artifacts struct {
	mu               sync.RWMutex
	anonymousCounter int
	builders         map[Kind]func(ctx Context) (Artifact, error)
	built            *orderedmap.OrderedMap[Kind, *artifactEntry]
	logger           logging.Logger
}

func NewArtifacts(logger logging.Logger) *Artifacts {
	return &Artifacts{
		builders: make(map[Kind]func(ctx Context) (Artifact, error)),
		built:    orderedmap.New[Kind, *artifactEntry](),
		logger:   logger,
	}
}

func (ds *Artifacts) RegisterBuilder(kind Kind, builder func(ctx Context) (Artifact, error)) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.builders[kind] = builder
}

func (ds *Artifacts) AddAnonymous(disposable DisposableArtifact) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.built.Set(Kind(fmt.Sprintf("_anonimous_%d", ds.anonymousCounter)), &artifactEntry{
		result: disposable,
	})
	ds.anonymousCounter++
}

func (ds *Artifacts) Get(ctx Context, kind Kind) (Artifact, error) {
	ds.mu.RLock()
	builtArtifact, ok := ds.built.Get(kind)
	ds.mu.RUnlock()

	if ok {
		builtArtifact.mu.RLock()
		defer builtArtifact.mu.RUnlock()
		return builtArtifact.result, builtArtifact.err
	}

	ds.mu.Lock()

	// repeat the check after acquiring the lock (in case it was added while we were waiting for the lock)
	builtArtifact, ok = ds.built.Get(kind)
	if ok {
		ds.mu.Unlock()
		builtArtifact.mu.RLock()
		defer builtArtifact.mu.RUnlock()
		return builtArtifact.result, builtArtifact.err
	}

	builder, ok := ds.builders[kind]
	if !ok {
		ds.mu.Unlock()
		return nil, fmt.Errorf("no builder for kind %s", kind)
	}
	entry := &artifactEntry{
		result: nil,
		err:    nil,
	}
	entry.mu.Lock()

	ds.built.Set(kind, entry)
	ds.mu.Unlock()

	ds.logger.Printf("--------------- Building %s... -------------\n", kind)
	entry.result, entry.err = builder(ctx)
	if entry.err != nil {
		entry.err = fmt.Errorf("error building artifact for kind %s: %w", kind, entry.err)
	}
	entry.mu.Unlock()

	return entry.result, entry.err
}

func (ds *Artifacts) Dispose() error {
	var errs []error
	ds.mu.Lock()
	defer ds.mu.Unlock()
	notDisposed := orderedmap.New[Kind, *artifactEntry]()
	defer func() {
		ds.built = notDisposed
	}()
	// iterate in reverse order to dispose in the opposite order of adding
	for pair := ds.built.Newest(); pair != nil; pair = pair.Prev() {
		pair := pair
		pair.Value.mu.Lock()
		if pair.Value.result != nil {
			if disposable, ok := pair.Value.result.(Disposable); ok {
				if e := disposable.Dispose(); e != nil {
					errs = append(errs, e)
				} else {
					//goland:noinspection GoDeferInLoop
					defer func() {
						notDisposed.Set(pair.Key, pair.Value)
					}()
				}
			}
			pair.Value.mu.Unlock()
		}
	}

	return errors.Join(errs...)
}
