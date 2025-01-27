package handlers

import (
  "elysium-backend/config"
  "log"
  "net/http"
  "os"
  "path/filepath"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request, uniqueID, filename string) {
  baseDir := config.GetEnv("OUTPUT_DIR", "./compiled_binaries")

  realPath := filepath.Join(baseDir, uniqueID, filename)

  log.Println("handlers.DownloadHandler -> Serving file:", realPath)

  info, err := os.Stat(realPath)
  if os.IsNotExist(err) || info.IsDir() {
    http.Error(w, "File not found", http.StatusNotFound)
    return
  }

  w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(realPath))
  w.Header().Set("Content-Type", "application/octet-stream")
  http.ServeFile(w, r, realPath)
}
