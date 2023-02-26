package main

import (
  "fmt"
  "flag"
  "os"
  "path/filepath"

  dotenv "github.com/joho/godotenv"
  "github.com/evanw/esbuild/pkg/api"
)

var outDirPath = filepath.Join("../", "public/build")

var mode = flag.String("mode", "dev", "build mode")

func main() {
  if err := os.RemoveAll(outDirPath); err != nil {
    fmt.Println("failed to clear", outDirPath)
    os.Exit(1)
  }

  if err := dotenv.Load(); err != nil {
    fmt.Println("Error loading .env file")
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

func define() map[string]string {
  return map[string]string{
    "BACKEND_URL": os.Getenv("BACKEND_URL"),
    "SOCKET_PROTOCOL": os.Getenv("SOCKET_PROTOCOL"),
    "HTTP_PROTOCOL": os.Getenv("HTTP_PROTOCOL"),
    "SOCKET_PATH": os.Getenv("SOCKET_PATH"),
    "TIMEOUT_OPEN": os.Getenv("TIMEOUT_OPEN"),
  }
}

func buildDev() api.BuildResult {
  result := api.Build(api.BuildOptions{
    EntryPoints: []string{"./init.ts", "./style.css"},
    Define: define(),
    Bundle: true,
    Write: true,
    Outdir: outDirPath,
  })

  return result
}

func buildProd() api.BuildResult {
  result := api.Build(api.BuildOptions{
    EntryPoints: []string{"./init.ts", "./style.css"},
    Define: define(),
    Bundle: true,
    Write: true,
    MinifyWhitespace: true,
    MinifyIdentifiers: true,
    MinifySyntax: true,
    Outdir: outDirPath,
  })

  return result
}