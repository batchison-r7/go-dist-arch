package godist

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/paketo-buildpacks/packit/v2/postal"
)

type Sources struct {
	PURL         string `toml:"purl"`
	SHA256       string `toml:"sha256"`
	Source       string `toml:"source"`
	SourceSHA256 string `toml:"source_sha256"`
	URI          string `toml:"uri"`
}

// postal.Dependency with extra fields defining architecture
type ArchDependency struct {
	postal.Dependency
	AMD64 Sources `toml:"amd64"`
	ARM64 Sources `toml:"arm64"`
}

// parse buildpack.toml and return ArchDependency for the given version
func parseBuildpackWithArch(path string, genericDependecy postal.Dependency) (ArchDependency, error) {
	file, err := os.Open(path)
	if err != nil {
		return ArchDependency{}, fmt.Errorf("failed to parse buildpack.toml: %w", err)
	}

	var buildpack struct {
		Metadata struct {
			DefaultVersions map[string]string `toml:"default-versions"`
			Dependencies    []ArchDependency  `toml:"dependencies"`
		} `toml:"metadata"`
	}
	_, err = toml.NewDecoder(file).Decode(&buildpack)
	if err != nil {
		return ArchDependency{}, fmt.Errorf("failed to parse buildpack.toml: %w", err)
	}

	genericDependecyStacks := strings.Join(genericDependecy.Stacks, " ")

	for _, dependency := range buildpack.Metadata.Dependencies {
		dependecyStacks := strings.Join(dependency.Stacks, " ")
		if dependency.Version == genericDependecy.Version && dependecyStacks == genericDependecyStacks {
			return dependency, nil
		}
	}
	return ArchDependency{}, fmt.Errorf(
		"failed to find dependency for version %s and stack %s",
		genericDependecy.Version,
		genericDependecy.Stacks,
	)
}
