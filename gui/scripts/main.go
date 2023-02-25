package main

import (
  "fmt"
  "flag"
  "os"
  "path/filepath"
  "github.com/evanw/esbuild/pkg/api"
)

var outDirPath = filepath.Join("../", "public/build")

var mode = flag.String("mode", "dev", "build mode")

func main() {
  if err := os.RemoveAll(outDirPath); err != nil {
    fmt.Println("failed to clear", outDirPath)
    os.Exit(1)
  }

  flag.Parse()

  var result api.BuildResult
  switch *mode {
  case "dev":
    result = buildDev()
  case "prod":
    result = buildProd()
  default:
    fmt.Println("unknown mode", mode)
    os.Exit(1)
  }

  if len(result.Errors) != 0 {
    fmt.Println("build has errors")
    os.Exit(1)
  }
}

func buildDev() api.BuildResult {
  result := api.Build(api.BuildOptions{
    EntryPoints: []string{"./init.ts", "./style.css"},
    Bundle: true,
    Write: true,
    Outdir: outDirPath,
  })

  return result
}

func buildProd() api.BuildResult {
  result := api.Build(api.BuildOptions{
    EntryPoints: []string{"./init.ts", "./style.css"},
    Bundle: true,
    Write: true,
    MinifyWhitespace: true,
    MinifyIdentifiers: true,
    MinifySyntax: true,
    Outdir: outDirPath,
  })

  return result
}