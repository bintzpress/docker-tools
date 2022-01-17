package buildConfig

type ImageBuildConfig struct {
	Tag     string
	Context string
	Args    map[string]string
}

type ImageConfig struct {
	Depends_on string
	Pull       string
	Build      ImageBuildConfig
	Push       map[string]map[string]string
}

type BuildConfig struct {
	Version    string
	Images     map[string]ImageConfig
	Registries map[string]map[string]string
}
