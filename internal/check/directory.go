package check

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

const (
	CdktfIndexDirectory = `cdktf`

	LegacyIndexDirectory       = `website/docs`
	LegacyDataSourcesDirectory = `d`
	LegacyGuidesDirectory      = `guides`
	LegacyResourcesDirectory   = `r`
	LegacyFunctionsDirectory   = `functions`

	RegistryIndexDirectory       = `docs`
	RegistryDataSourcesDirectory = `data-sources`
	RegistryGuidesDirectory      = `guides`
	RegistryResourcesDirectory   = `resources`
	RegistryFunctionsDirectory   = `functions`

	DocumentationGlobPattern = `{docs/index.md,docs/{,cdktf/}{data-sources,guides,resources,functions},website/docs}/**/*`

	// Terraform Registry Storage Limits
	// https://www.terraform.io/docs/registry/providers/docs.html#storage-limits
	RegistryMaximumNumberOfFiles = 2000
	RegistryMaximumSizeOfFile    = 500000 // 500KB

)

var ValidLegacyDirectories = []string{
	LegacyIndexDirectory,
	LegacyIndexDirectory + "/" + LegacyDataSourcesDirectory,
	LegacyIndexDirectory + "/" + LegacyGuidesDirectory,
	LegacyIndexDirectory + "/" + LegacyResourcesDirectory,
	LegacyIndexDirectory + "/" + LegacyFunctionsDirectory,
}

var ValidRegistryDirectories = []string{
	RegistryIndexDirectory,
	RegistryIndexDirectory + "/" + RegistryDataSourcesDirectory,
	RegistryIndexDirectory + "/" + RegistryGuidesDirectory,
	RegistryIndexDirectory + "/" + RegistryResourcesDirectory,
	RegistryIndexDirectory + "/" + RegistryFunctionsDirectory,
}

var ValidCdktfLanguages = []string{
	"csharp",
	"go",
	"java",
	"python",
	"typescript",
}

var ValidLegacySubdirectories = []string{
	LegacyIndexDirectory,
	LegacyDataSourcesDirectory,
	LegacyGuidesDirectory,
	LegacyResourcesDirectory,
}

var ValidRegistrySubdirectories = []string{
	RegistryIndexDirectory,
	RegistryDataSourcesDirectory,
	RegistryGuidesDirectory,
	RegistryResourcesDirectory,
}

func InvalidDirectoriesCheck(dirPath string) error {
	if IsValidRegistryDirectory(dirPath) {
		return nil
	}

	if IsValidLegacyDirectory(dirPath) {
		return nil
	}

	if IsValidCdktfDirectory(dirPath) {
		return nil
	}

	return fmt.Errorf("invalid Terraform Provider documentation directory found: %s", dirPath)

}

func MixedDirectoriesCheck(providerDir string) error {
	var legacyDirectoryFound bool
	var registryDirectoryFound bool
	err := fmt.Errorf("mixed Terraform Provider documentation directory layouts found, must use only legacy or registry layout")

	providerFs := os.DirFS(providerDir)

	files, globErr := doublestar.Glob(providerFs, DocumentationGlobPattern)
	if globErr != nil {
		return fmt.Errorf("error finding documentation files: %w", err)
	}

	for _, file := range files {
		directory := filepath.Dir(file)

		// Allow docs/ with other files
		if IsValidRegistryDirectory(directory) && directory != RegistryIndexDirectory {
			registryDirectoryFound = true

			if legacyDirectoryFound {
				return err
			}
		}

		if IsValidLegacyDirectory(directory) {
			legacyDirectoryFound = true

			if registryDirectoryFound {
				return err
			}
		}
	}

	return nil
}

// NumberOfFilesCheck verifies that documentation is below the Terraform Registry storage limit.
// This check presumes that all provided directories are valid, e.g. that directory checking
// for invalid or mixed directory structures was previously completed.
func NumberOfFilesCheck(directories map[string][]string) error {
	var numberOfFiles int

	for directory, files := range directories {
		// Ignore CDKTF files. The file limit is per-language and presumably there is one CDKTF file per source HCL file.
		if IsValidCdktfDirectory(directory) {
			continue
		}

		directoryNumberOfFiles := len(files)
		log.Printf("[TRACE] Found %d documentation files in directory: %s", directoryNumberOfFiles, directory)
		numberOfFiles = numberOfFiles + directoryNumberOfFiles
	}

	log.Printf("[DEBUG] Found %d documentation files with limit of %d", numberOfFiles, RegistryMaximumNumberOfFiles)
	if numberOfFiles >= RegistryMaximumNumberOfFiles {
		return fmt.Errorf("exceeded maximum (%d) number of documentation files for Terraform Registry: %d", RegistryMaximumNumberOfFiles, numberOfFiles)
	}

	return nil
}

func IsValidLegacyDirectory(directory string) bool {
	for _, validLegacyDirectory := range ValidLegacyDirectories {
		if directory == validLegacyDirectory {
			return true
		}
	}

	return false
}

func IsValidRegistryDirectory(directory string) bool {
	for _, validRegistryDirectory := range ValidRegistryDirectories {
		if directory == validRegistryDirectory {
			return true
		}
	}

	return false
}

func IsValidCdktfDirectory(directory string) bool {
	if directory == fmt.Sprintf("%s/%s", LegacyIndexDirectory, CdktfIndexDirectory) {
		return true
	}

	if directory == fmt.Sprintf("%s/%s", RegistryIndexDirectory, CdktfIndexDirectory) {
		return true
	}

	for _, validCdktfLanguage := range ValidCdktfLanguages {

		if directory == fmt.Sprintf("%s/%s/%s", LegacyIndexDirectory, CdktfIndexDirectory, validCdktfLanguage) {
			return true
		}

		if directory == fmt.Sprintf("%s/%s/%s", RegistryIndexDirectory, CdktfIndexDirectory, validCdktfLanguage) {
			return true
		}

		for _, validLegacySubdirectory := range ValidLegacySubdirectories {
			if directory == fmt.Sprintf("%s/%s/%s/%s", LegacyIndexDirectory, CdktfIndexDirectory, validCdktfLanguage, validLegacySubdirectory) {
				return true
			}
		}

		for _, validRegistrySubdirectory := range ValidRegistrySubdirectories {
			if directory == fmt.Sprintf("%s/%s/%s/%s", RegistryIndexDirectory, CdktfIndexDirectory, validCdktfLanguage, validRegistrySubdirectory) {
				return true
			}
		}
	}

	return false
}
