// Simple script to build pather client and server. This is not required, but it will properly insert version date and commit
// metadata into the resulting binaries, which `go build` will not do by default.
package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/josephspurrier/goversioninfo"
	"github.com/mhenderson-so/go-blackbox/version"
)

var (
	shaFlag       = flag.String("sha", "", "SHA to embed. Omit to pull from current repository")
	buildOS       = flag.String("os", "windows", "OS to build for.")
	buildOfficial = flag.Bool("release", false, "Release build. Used in -version info only.")
	buildBranch   = flag.String("branch", "", "Branch name. Used in -version info only. Omit to pull from current repository.")
	buildVersion  = flag.String("version", "2.0.0", "Version number to embed into the build output")
	output        = flag.String("output", os.Getenv("GOBIN"), "Output directory")
	buildChoco    = flag.Bool("choco", false, "Package binary into chocolatey/nuget package")
)

func main() {
	flag.Parse()
	*buildOS = strings.ToLower(*buildOS)
	origBuildOS := os.Getenv("GOOS")
	os.Setenv("GOOS", *buildOS)
	defer os.Setenv("GOOS", origBuildOS)

	goPath := os.Getenv("GOPATH")

	// Get current commit SHA
	sha := *shaFlag
	if sha == "" {
		cmd := exec.Command("git", "rev-parse", "HEAD")
		cmd.Stderr = os.Stderr
		output, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		sha = strings.TrimSpace(string(output))
	}

	//Get current branch name
	branchName := *buildBranch
	if branchName == "" {
		cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		cmd.Stderr = os.Stderr
		output, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		branchName = strings.TrimSpace(string(output))
	}

	var buildRelease string
	if *buildOfficial || branchName == "master" {
		buildRelease = "true"
	}

	timeStr := time.Now().UTC().Format("20060102150405")
	ldFlags := fmt.Sprintf("-X github.com/mhenderson-so/go-blackbox/version.VersionNumber=%s -X github.com/mhenderson-so/go-blackbox/version.VersionSHA=%s -X github.com/mhenderson-so/go-blackbox/version.VersionDate=%s -X github.com/mhenderson-so/go-blackbox/version.OfficialBuild=%s -X github.com/mhenderson-so/go-blackbox/version.BuildBranch=%s", *buildVersion, sha, timeStr, buildRelease, branchName)

	var buildArgs, genArgs []string

	binaryName := "blackbox"
	if *buildOS == "windows" {
		binaryName += ".exe"
	}

	//Path to the config for goversioninfo
	versionInfoTemplateFile := filepath.Join(goPath, "src", "github.ds.stackexchange.com", "mhenderson", "go-blackbox", "build", "versioninfo-orig.json")
	versionInfoFile := filepath.Join(goPath, "src", "github.ds.stackexchange.com", "mhenderson", "go-blackbox", "cmd", "blackbox", "versioninfo.json")
	//Contents of the config
	versionInfoBytes, err := ioutil.ReadFile(versionInfoTemplateFile)

	if err == nil { //If we have a goversioninfo config file
		if version.BuildVersion.Version == "" {
			version.BuildVersion.Version = *buildVersion
		}
		versionNos := strings.Split(version.BuildVersion.Version, ".") //Get the version numbers
		var versionMajor int
		var versionMinor int
		var versionBuild int
		if len(versionNos) == 3 {
			versionMajor, _ = strconv.Atoi(versionNos[0]) //Major as string
			versionMinor, _ = strconv.Atoi(versionNos[1]) //Minor as string
			versionBuild, _ = strconv.Atoi(versionNos[2]) //Build as string
		}

		//Set these into the version pkg so that we have the same data here as we will in the compiled program
		version.BuildVersion.OfficialBuild = buildRelease
		version.BuildVersion.VersionDate = timeStr
		version.BuildVersion.VersionSHA = sha

		//This will hold the contents of versioninfo.json
		versionInfo := &goversioninfo.VersionInfo{}
		err = versionInfo.ParseJSON(versionInfoBytes) //Load the JSON into the struct

		//If we have a version number and a populated VersionInfo{}
		if err == nil && (versionMajor > 0 || versionMinor > 0) {
			//Set all the properties that we want to set for the binary
			versionInfo.Build()
			versionInfo.FixedFileInfo.FileVersion.Major = versionMajor
			versionInfo.FixedFileInfo.FileVersion.Minor = versionMinor
			versionInfo.FixedFileInfo.FileVersion.Build = versionBuild
			versionInfo.FixedFileInfo.ProductVersion.Major = versionMajor
			versionInfo.FixedFileInfo.ProductVersion.Minor = versionMinor
			versionInfo.FixedFileInfo.ProductVersion.Build = versionBuild
			versionInfo.StringFileInfo.ProductVersion = sha
			versionInfo.StringFileInfo.Comments = version.GetVersionInfo()
			versionInfo.StringFileInfo.OriginalFilename = binaryName
		}
		//We need to write the config JSON back to the original file so that it can be picked up by the linker
		versionInfoBytes, err := json.Marshal(versionInfo)
		if err == nil {
			err = ioutil.WriteFile(versionInfoFile, versionInfoBytes, 0644)
			if err != nil {
				fmt.Println(err)
			}
			//Add parameters to the "go generate" array

		}
	}

	binPath := filepath.Join(*output, binaryName)
	buildArgs = append(buildArgs, "build", "-o", binPath)
	buildArgs = append(buildArgs, "-ldflags", ldFlags, "github.com/mhenderson-so/go-blackbox/cmd/blackbox")
	genArgs = append(genArgs, "generate", "github.com/mhenderson-so/go-blackbox/cmd/blackbox")

	//We always want to do a go generate. No harm in doing this anyway.
	//Check that we have goversioninfo in the path
	goVI, _ := exec.Command("goversioninfo", "-?").CombinedOutput() //This should get the help output
	if len(goVI) == 0 {                                             //If we don't have any help
		fmt.Println("[missing goversioninfo, attempting to install]")
		outInstall, _ := exec.Command("go", "install", "github.com/mhenderson-so/go-blackbox/vendor/github.com/josephspurrier/goversioninfo/cmd/goversioninfo").CombinedOutput()
		fmt.Println(string(outInstall)) //Install it
	}

	fmt.Println(genArgs)                     //Should output like "generate gitlab.stackexhange... etc"
	genCmd := exec.Command("go", genArgs...) //Create a Go Generate command
	genCmd.Stdout = os.Stdout                //Pipe everything to the normal std. outputs
	genCmd.Stderr = os.Stderr
	genCmd.Run() //Generate

	fmt.Println("building", filepath.Join(*output, binaryName), "for", *buildOS)
	cmd := exec.Command("go", buildArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	//Package into nuget
	if *buildChoco {
		fmt.Printf("[package %s]\n", binaryName)

		//Search for choco, fail if not found
		chocoVersion, err := exec.Command("choco", "-v").CombinedOutput()
		if err != nil {
			fmt.Println("ERROR: chocolatey not found:", err)
			return
		}
		if len(chocoVersion) == 0 {
			fmt.Println("ERROR: No chocolatey version found")
			return
		}

		//Build the nuspec data from the template
		packagePath := filepath.Join("build", "nuget", "blackbox")
		nuPkgPath := filepath.Join(packagePath, fmt.Sprintf("blackbox.%s.nupkg", version.BuildVersion.Version))
		nuspecInPath := filepath.Join(packagePath, "blackbox.template")
		nuspecOutPath := filepath.Join(packagePath, "blackbox.nuspec")
		nuspecRaw, err := ioutil.ReadFile(nuspecInPath)
		if err != nil {
			fmt.Println("ERROR: Unable to read nuspec file:", err)
			return
		}

		//Modify the nuspec with the date and version numbers
		var nuspec nuspecPackage
		xml.Unmarshal(nuspecRaw, &nuspec)
		nuspec.Metadata.Version = version.BuildVersion.Version
		nuspec.Metadata.Copyright = time.Now().UTC().Format("2006-01-02")

		if buildRelease != "true" {
			//nuspecOutPath = filepath.Join(packagePath, fmt.Sprintf("patcher-%s_%s.nuspec", app.dirName, branchName))
			nuPkgPath = filepath.Join(packagePath, fmt.Sprintf("blackbox-%s.%s.nupkg", branchName, version.BuildVersion.Version))
			nuspec.Metadata.ID = fmt.Sprintf("%s-%s", nuspec.Metadata.ID, branchName)
			nuspec.Metadata.Title = fmt.Sprintf("%s (%s)", nuspec.Metadata.Title, branchName)
			nuspec.Metadata.Tags = fmt.Sprintf("%s %s", nuspec.Metadata.Tags, branchName)
			nuspec.Metadata.Summary = fmt.Sprintf("%s (%s)", nuspec.Metadata.Summary, branchName)
			nuspec.Metadata.Description = fmt.Sprintf("%s (%s)", nuspec.Metadata.Description, branchName)
		}
		//Write the finalised newspec into the correct filename
		nuspecOut, _ := xml.Marshal(&nuspec)
		ioutil.WriteFile(nuspecOutPath, nuspecOut, 0444)

		//Pop our built binary into the tools directory for inclusion with the package
		binCopyTo := filepath.Join(packagePath, "tools", binaryName)
		err = osCopyFile(binPath, binCopyTo)
		if err != nil {
			fmt.Println("ERROR: Copying binary:", err)
		}

		//Package
		cmd := exec.Command("choco", "pack")
		cmd.Dir = filepath.Join("build", "nuget", "blackbox")
		cmdOut, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("ERROR: Writing nuspec file", err)
			fmt.Println("ERROR:", string(cmdOut))
			return
		}

		//Check if package exists
		_, err = os.Stat(nuPkgPath)
		if os.IsNotExist(err) {
			fmt.Println("ERROR: Cannot find nupkg file", nuPkgPath)
			fmt.Println("ERROR:", string(cmdOut))
			return
		}
		fmt.Println("nupkg file at", nuPkgPath)
	}

}

func osCopyFile(source, destination string) error {
	from, err := os.Open(source)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	return nil
}

//Nuspec structs
type nuspecPackage struct {
	XMLName   xml.Name       `xml:"package,omitempty"`
	Attrxmlns string         `xml:"xmlns,attr"`
	Metadata  nuspecMetadata `xml:"metadata,omitempty"`
	Files     []nuspecFile   `xml:"files>file,omitempty"`
}

type nuspecMetadata struct {
	ID                       string             `xml:"id,omitempty"`
	Version                  string             `xml:"version,omitempty"`
	PackageSourceURL         string             `xml:"packageSourceUrl,omitempty"`
	Owners                   string             `xml:"owners,omitempty"`
	Title                    string             `xml:"title,omitempty"`
	Authors                  string             `xml:"authors,omitempty"`
	ProjectURL               string             `xml:"projectUrl,omitempty"`
	IconURL                  string             `xml:"iconUrl,omitempty"`
	Copyright                string             `xml:"copyright,omitempty"`
	LicenseURL               string             `xml:"licenseUrl,omitempty"`
	RequireLicenseAcceptance string             `xml:"requireLicenseAcceptance,omitempty"`
	ProjectSourceURL         string             `xml:"projectSourceUrl,omitempty"`
	DocsURL                  string             `xml:"docsUrl,omitempty"`
	MailingListURL           string             `xml:"mailingListUrl,omitempty"`
	BugTrackerURL            string             `xml:"bugTrackerUrl,omitempty"`
	ReleaseNotes             string             `xml:"releaseNotes,omitempty"`
	Dependencies             []nuspecDependency `xml:"dependencies>dependency,omitempty"`
	Tags                     string             `xml:"tags,omitempty"`
	Summary                  string             `xml:"summary,omitempty"`
	Description              string             `xml:"description,omitempty"`
}

type nuspecDependency struct {
	XMLName xml.Name `xml:"dependency"`
	ID      string   `xml:"id,attr"`
	Version string   `xml:"version,attr"`
}

type nuspecFile struct {
	XMLName xml.Name `xml:"file"`
	Src     string   `xml:"src,attr"`
	Target  string   `xml:"target,attr"`
}
