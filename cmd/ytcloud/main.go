package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"ytcloud/internal/codec"
	"ytcloud/internal/crypto"
	"ytcloud/internal/media"
	"ytcloud/internal/packer"
	"ytcloud/internal/server"
	"ytcloud/internal/vault"
)

func main() {
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		switch os.Args[1] {
		case "encode":
			runCLIEncode(os.Args[2:])
			return
		case "decode":
			runCLIDecode(os.Args[2:])
			return
		case "peek":
			runCLIPeek(os.Args[2:])
			return
		case "calibrate":
			runCLICalibrate(os.Args[2:])
			return
		case "version":
			fmt.Println("BrightVault v0.7 — Optical 4K/8K Encrypted Cloud Manager")
			return
		}
	}

	vaultPath := flag.String("vault", filepath.Join(media.DefaultAppDir(), "vault.ytdb"), "Caminho do arquivo de cofre criptografado")
	addr := flag.String("addr", "127.0.0.1:0", "Endereço local do servidor UI")
	noBrowser := flag.Bool("no-browser", false, "Não abrir a janela automaticamente")
	flag.Parse()

	store := vault.NewStore(*vaultPath)
	srv, err := server.New(store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro inicializando BrightVault: %v\n", err)
		os.Exit(1)
	}

	ln, appURL, err := srv.Start(*addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro iniciando servidor local: %v\n", err)
		os.Exit(1)
	}
	defer ln.Close()

	// Write startup info to local log file (useful when compiled with -H=windowsgui where no console exists)
	_ = os.MkdirAll(media.DefaultAppDir(), 0o755)
	logLine := fmt.Sprintf("[%s] BrightVault iniciado | Cofre: %s | URL: %s\n", time.Now().Format(time.RFC3339), *vaultPath, appURL)
	_ = os.WriteFile(filepath.Join(media.DefaultAppDir(), "brightvault.log"), []byte(logLine), 0o600)

	if !*noBrowser {
		server.LaunchDesktopWindow(appURL)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
	case <-srv.ShutdownChan():
	}
	store.Lock()
}

// --- Standalone CLI Subcommands (for scripting, automation & direct round-trip testing) ---

func runCLIEncode(args []string) {
	fs := flag.NewFlagSet("encode", flag.ExitOnError)
	inPath := fs.String("in", "", "Arquivo ou pasta de entrada")
	outPath := fs.String("out", "ytcloud_out.mp4", "Vídeo de saída (.mp4)")
	pass := fs.String("pass", "", "Senha de criptografia")
	res := fs.String("res", "4k", "Resolução: 4k, 8k, 1080p")
	prof := fs.String("profile", "y8-gray", "Perfil: y8-gray, y6-gray, y4-gray, y2-safe, y2-balanced, ycbcr-hex, ycbcr-dual, ycbcr-fast")
	fps := fs.Int("fps", 10, "Frames por segundo")
	crf := fs.Int("crf", 0, "Qualidade CRF FFmpeg (0 = automático ideal do perfil)")
	_ = fs.Parse(args)

	if *inPath == "" || *pass == "" {
		fmt.Fprintln(os.Stderr, "Uso: YTCloud.exe encode -in <arquivo/pasta> -out <video.mp4> -pass <senha> [-res 8k] [-profile ycbcr-dual]")
		os.Exit(1)
	}
	tools := media.DetectTools()
	if !tools.HasFFmpeg {
		fmt.Fprintln(os.Stderr, "FFmpeg não encontrado.")
		os.Exit(1)
	}
	inputs, err := packer.CollectInputsFromPath(*inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro lendo entrada: %v\n", err)
		os.Exit(1)
	}
	contUUID, _ := crypto.RandomHex(8)
	salt, _ := crypto.RandomBytes(16)
	contID := binary.LittleEndian.Uint64(salt[:8])

	manifest, rawPayload, err := packer.PackFiles(inputs, contUUID, *prof, *res, tools.FFmpegPath, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro empacotando: %v\n", err)
		os.Exit(1)
	}
	manifest.CreatedAt = time.Now().Format(time.RFC3339)
	mJSON, _ := packer.MarshalManifest(manifest)

	encM, err := crypto.EncryptWithPassword(*pass, mJSON, []byte("YTCloud-CLI-Manifest"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro criptografando manifesto: %v\n", err)
		os.Exit(1)
	}
	encP, err := crypto.EncryptWithPassword(*pass, rawPayload, []byte("YTCloud-CLI-Payload"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro criptografando payload: %v\n", err)
		os.Exit(1)
	}

	stats, err := codec.EncodeStreamToVideo(codec.VideoEncodeOptions{
		FFmpegPath:     tools.FFmpegPath,
		OutputPath:     *outPath,
		ResolutionCode: *res,
		ProfileCode:    *prof,
		FPS:            *fps,
		CRF:            *crf,
		ContainerID:    contID,
		OnProgress: func(f, t int) {
			fmt.Printf("\rRenderizando frame %d/%d...", f, t)
		},
	}, encM, encP)
	fmt.Println()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro codificando vídeo: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Sucesso! Vídeo: %s (%dx%d, %d frames, %d bytes, %d arquivo(s))\n",
		stats.OutputPath, stats.Width, stats.Height, stats.TotalFrames, stats.VideoSizeBytes, len(manifest.Files))
}

func runCLIDecode(args []string) {
	fs := flag.NewFlagSet("decode", flag.ExitOnError)
	inPath := fs.String("in", "", "Vídeo óptico de entrada (.mp4/.mkv/.webm)")
	outDir := fs.String("out", "./restaurado", "Pasta de destino")
	pass := fs.String("pass", "", "Senha de criptografia")
	_ = fs.Parse(args)

	if *inPath == "" || *pass == "" {
		fmt.Fprintln(os.Stderr, "Uso: YTCloud.exe decode -in <video.mp4> -out <pasta_destino> -pass <senha>")
		os.Exit(1)
	}
	tools := media.DetectTools()
	decRes, err := codec.DecodeVideoStream(tools.FFmpegPath, tools.FFprobePath, *inPath, false, func(f, t int) {
		fmt.Printf("\rDecodificando frame %d/%d...", f, t)
	})
	fmt.Println()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro decodificando: %v\n", err)
		os.Exit(1)
	}
	mJSON, err := crypto.DecryptWithPassword(*pass, decRes.EncryptedManifest, []byte("YTCloud-CLI-Manifest"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro descriptografando manifesto: %v\n", err)
		os.Exit(1)
	}
	manifest, err := packer.UnmarshalManifest(mJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Manifesto inválido: %v\n", err)
		os.Exit(1)
	}
	rawPayload, err := crypto.DecryptWithPassword(*pass, decRes.EncryptedPayload, []byte("YTCloud-CLI-Payload"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro descriptografando payload: %v\n", err)
		os.Exit(1)
	}
	extracted, err := packer.ExtractFiles(manifest, rawPayload, *outDir, nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro extraindo: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Restaurado com 100%% integridade SHA-256! %d arquivo(s) em %s (bytes corrigidos pelo FEC: %d)\n",
		len(extracted), *outDir, decRes.TotalCorrectedBytes)
}

func runCLIPeek(args []string) {
	fs := flag.NewFlagSet("peek", flag.ExitOnError)
	inPath := fs.String("in", "", "Vídeo óptico de entrada")
	pass := fs.String("pass", "", "Senha de criptografia")
	_ = fs.Parse(args)

	if *inPath == "" || *pass == "" {
		fmt.Fprintln(os.Stderr, "Uso: YTCloud.exe peek -in <video.mp4> -pass <senha>")
		os.Exit(1)
	}
	tools := media.DetectTools()
	decRes, err := codec.DecodeVideoStream(tools.FFmpegPath, tools.FFprobePath, *inPath, true, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro no peek: %v\n", err)
		os.Exit(1)
	}
	mJSON, err := crypto.DecryptWithPassword(*pass, decRes.EncryptedManifest, []byte("YTCloud-CLI-Manifest"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro descriptografando manifesto: %v\n", err)
		os.Exit(1)
	}
	manifest, _ := packer.UnmarshalManifest(mJSON)
	fmt.Printf("Contêiner %s (%s %s) — %d arquivo(s):\n", manifest.ContainerID, manifest.Resolution, manifest.ProfileCode, len(manifest.Files))
	for _, f := range manifest.Files {
		fmt.Printf("  - %s (%d bytes, %s)\n", f.RelPath, f.OriginalSize, f.ModTime)
	}
}

func runCLICalibrate(args []string) {
	fs := flag.NewFlagSet("calibrate", flag.ExitOnError)
	res := fs.String("res", "1080p", "Resolução: 8k, 4k, 1080p")
	crf := fs.Int("crf", 26, "CRF simulando recompressão YouTube")
	_ = fs.Parse(args)

	tools := media.DetectTools()
	if !tools.HasFFmpeg {
		fmt.Fprintln(os.Stderr, "FFmpeg não encontrado.")
		os.Exit(1)
	}
	profiles := []string{"y8-gray", "y6-gray", "y4-gray", "y2-safe", "y2-balanced", "ycbcr-hex", "ycbcr-dual", "ycbcr-fast"}
	for _, p := range profiles {
		r := codec.RunLocalCalibration(tools.FFmpegPath, tools.FFprobePath, os.TempDir(), *res, p, *crf)
		status := "PASSOU"
		if !r.Passed {
			status = "FALHOU (" + r.Error + ")"
		}
		fmt.Printf("%-12s | %7d B/frame | BER: %6.3f%% (%d/%d B) | %s\n",
			r.ProfileCode, r.PayloadPerFrame, r.ByteErrorRate, r.CorrectedBytes, r.RawBytes, status)
	}
}
