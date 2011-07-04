package main
/* 
*  File: gonew_main.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 19:17:53 PDT 2011
*  Usage: gonew [options]
*/
import (
    "os"
    //"io"
    "log"
    "fmt"
    "flag"
    //"bufio"
    //"io/ioutil"
    //"path/filepath"
    //"github.com/hoisie/mustache.go"
    "github.com/kr/pretty.go"
)

const (
    DEBUG = true
	DEBUG_LEVEL = 0
)
var (
    usage = "gonew [options] TYPE NAME"
    printUsageHead = func () { fmt.Fprint(os.Stderr, "\n", usage, "\n\n") }
    name   string
    ptype  string
    repo   string
    host   string
    user   string
    target string
    help   bool
)

func setupFlags() *flag.FlagSet {
    var fs = flag.NewFlagSet("gonew", flag.ExitOnError)
    fs.StringVar(&repo,
        "repo", "git", "Repository type (e.g. 'git').")
    fs.StringVar(&host,
        "host", "", "Repository host if any (e.g. 'github').")
    fs.StringVar(&user,
        "user", "", "Repo host username.")
    fs.StringVar(&target,
        "target", "", "Makefile target. Default based on NAME.")
    fs.BoolVar(&help,
        "help", false, "Show this message.")
    var usageTemp = fs.Usage
    fs.Usage = func () {
        printUsageHead()
        usageTemp()
    }
    return fs
}

type Request int
const (
    NilRequest Request = iota
    ProjectRequest
    LibraryRequest
)


var RequestedFile   File
var RequestedProject Project

func parseArgs() Request {
    var fs = setupFlags()
    fs.Parse(os.Args[1:])
    if help {
        fs.Usage()
        os.Exit(0)
    }
    var narg  = fs.NArg()
    if narg < 1 {
        fmt.Fprint(os.Stderr, "missing TYPE argument")
        os.Exit(1)
    }
    if narg < 2 {
        fmt.Fprint(os.Stderr, "mising NAME argument")
        os.Exit(1)
    }
    ptype = fs.Arg(0)
    name = fs.Arg(1)

    if target == "" {
        target = DefaultTarget(name)
    }

    var (
        file = File{Name:name, Pkg: "main", Repo:AppConfig.Repo,
            User:AppConfig.HostUser, Host:AppConfig.Host}
        project = Project{
                Name:name, Target:target, Type: NilProjectType,
                Host:AppConfig.Host, User:AppConfig.HostUser,
                Repo:AppConfig.Repo}
        produceProject = true
        repoObj = NilRepoType
        hostObj = NilRepoHost
    )
    log.Printf("%v", file)
    switch ptype {
    case "cmd":
        project.Type = CmdType
    case "pkg":
        project.Type = PkgType
    case "lib":
        produceProject = false
    default:
        fmt.Fprintf(os.Stderr, "Unknown TYPE %s\n", ptype)
        os.Exit(1)
    }
    switch repo {
    case "":
        break
    case "git":
        repoObj = GitType
    default:
        fmt.Fprintf(os.Stderr, "Unknown REPO %s\n", repo)
        os.Exit(1)
    }
    switch host {
    case "":
        break
    case "github":
        hostObj = GitHubHost
        repoObj = GitType
    default:
        fmt.Fprintf(os.Stderr, "Unknown HOST %s\n", host)
        os.Exit(1)
    }
    if produceProject {
        // TODO check target for improper characters.
        if user != "" {
            project.User = user
        }
        if hostObj != NilRepoHost {
            project.Host = hostObj
        }
        if repoObj != NilRepoType {
            project.Repo = repoObj
        }
        RequestedProject = project
        return ProjectRequest
    } else {
        if narg < 3 {
            fmt.Fprint(os.Stderr, "mising PKG argument")
            os.Exit(1)
        }
        file.Pkg = fs.Arg(2)
        if user != "" {
            file.User = user
        }
        if hostObj != NilRepoHost {
            file.Host = hostObj
        }
        if repoObj != NilRepoType {
            file.Repo = repoObj
        }
        RequestedFile = file
        return LibraryRequest
    }
    return NilRequest
}

func main() {
    var errTouch = TouchConfig()
    if errTouch != nil {
        fmt.Print(errTouch.String(), "\n")
        os.Exit(1)
    }
	ReadConfig()
    var request = parseArgs()
    switch request {
    case ProjectRequest:
        if DEBUG {
            fmt.Printf("Project requested %v\n", pretty.Formatter(RequestedProject))
        }
        var errCreate = RequestedProject.Create()
        if errCreate != nil {
            fmt.Fprint(os.Stderr, errCreate.String(), "\n")
            os.Exit(1)
        }
    case LibraryRequest:
        if DEBUG {
            fmt.Printf("Library requested %v\n", pretty.Formatter(RequestedFile))
        }
        var errCreate = RequestedFile.Create()
        if errCreate != nil {
            fmt.Fprint(os.Stderr, errCreate.String(), "\n")
            os.Exit(1)
        }
    }
}
