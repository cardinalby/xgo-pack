package buildctx

type ArtifactsI interface {
	RegisterBuilder(kind Kind, builder func(ctx Context) (Artifact, error))
	AddAnonymous(disposable DisposableArtifact)
	Get(ctx Context, kind Kind) (Artifact, error)
	Dispose() error
}
