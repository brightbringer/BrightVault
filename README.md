<p align="center">
  <img src="assets/brightvault.png" alt="BrightVault Logo" width="128" />
</p>

<h1 align="center">BrightVault v0.701</h1>

<p align="center">
  <strong>Cofre pessoal criptografado em vídeo óptico de alta densidade com Blindagem Cross-Frame RS</strong><br/>
  Executável único para Windows · Criptografia Pós-Quântica em Cascata · Encode/Decode 4× Paralelo · Atualização Automática
</p>

---

## 🔬 Intuito Científico e Pesquisa White-Hat

O **BrightVault** e o **BrightWeb** (ecossistema **BrightNet**) são projetos abertos de pesquisa científica e engenharia experimental dedicados ao estudo da **resistência, integridade e transmissão de dados binários através de canais ruidosos com compressão severa com perdas** — especificamente codecs modernos de vídeo (H.264, VP9 e AV1) sob transcodificação agressiva, redução de bitrate variável (VBR) e subamostragem cromática (`YUV420p`). A investigação abrange teoria da informação aplicada, códigos corretores de erros Reed-Solomon intra-frame e inter-frame (*Cross-Frame RS*), decodificação *soft-decision* (GMD), esteganografia óptica resiliente e compressão híbrida por camadas com dicionários pré-treinados.

Todo o ecossistema é desenvolvido sob rigorosos princípios **White-Hat**, com código aberto, transparência total e finalidade estritamente educacional e defensiva: toda a criptografia (*Argon2id + HKDF-SHA512 + AES-256-GCM + XChaCha20-Poly1305*) e o processamento de sinais ocorrem 100% localmente no dispositivo do usuário, sem servidores centrais, sem telemetria oculta e sem coleta de dados. O objetivo científico é demonstrar, documentar e avançar técnicas abertas de preservação de integridade e segurança defensiva frente à compressão extrema — sempre incentivando o uso ético e o respeito aos Termos de Serviço das plataformas.

---

## ✨ Novidades da v0.701

- **🛡️ Blindagem Anti-Apagão YouTube (Cross-Frame RS — `YTCF`):** correção de erros inter-frame em `GF(2^8)` com modos *Automático*, *Máximo* (`K=6, M=2` — sobrevive até 25% do vídeo destruído), *Alto* (`K=8, M=2`) e *Desligado*, com parada antecipada assim que os stripes ficam recuperáveis.
- **⚡ Encode 4× e Decode 4× Particionados + `yt-dlp` Paralelo (`-N 8`):** codificação VP9/AV1 e decodificação de vídeos longos divididas automaticamente em 4 partições simultâneas, somadas a downloads paralelos de 8 fragmentos no `yt-dlp`.
- **📦 Empacotamento Híbrido por Camadas (v3) + Dicionário Web Zstd (`WebDictV1`):** agrupa arquivos de texto/código em blocos sólidos independentes por profundidade de pasta, aplica dicionário Web embutido estilo Brotli (-15% a -45% em HTML/CSS/JS) e isola falhas parciais.
- **🖼️ Calibração Real para Imagem Estática no Fundo:** as barras de Segurança, Tamanho e Velocidade detectam automaticamente quando uma imagem parada é usada como fundo e ajustam os indicadores (com marcadores dinâmicos `5 Mín / 6 Rec` ou `6 Mín / 7 Rec`).
- **🔌 Estabilidade em Segundo Plano (SSE `/api/events`):** conexão SSE persistente substitui a dependência de `setInterval`, impedindo que o processo se encerre quando a janela fica minimizada no Edge/Chrome.
- **🛑 Cancelamento Real de Geração (`context.Context`):** cancelamentos encerram o FFmpeg imediatamente e removem arquivos parciais sem travar a aplicação.

---

## Como baixar e usar (Windows)

Você **não precisa instalar Go, Python nem dependências manualmente**:

1. Baixe **`BrightVault.exe`** diretamente da página de **[Releases](https://github.com/brightbringer/BrightVault/releases)** (ou clique em **Code → Download ZIP** e extraia a pasta).
2. Dê **duplo clique em `BrightVault.exe`**.
3. O BrightVault inicia um servidor local seguro exclusivamente em `127.0.0.1` (com token aleatório por sessão) e abre automaticamente a interface desktop.
4. Na primeira abertura, crie seu cofre com sua **senha mestra** (e opcionalmente um **Arquivo-Chave 2º Fator**).

> **FFmpeg e yt-dlp:** para manter o executável compacto, `ffmpeg` e `yt-dlp` não são embutidos no `.exe`. Caso ainda não estejam instalados no seu PC, o próprio BrightVault exibe um botão **"Instalar Automaticamente"** em 1 clique na tela inicial.

---

## Principais Recursos

- **Criptografia Pós-Quântica em Cascata (*Zero-Knowledge*):** `Argon2id` + `HKDF-SHA512` + **AES-256-GCM + XChaCha20-Poly1305** com chaves independentes por pasta e por vídeo.
- **Chave própria por pasta + Compartilhamento seletivo:** cada pasta possui sua própria chave criptográfica guardada no seu cofre (aberta automaticamente pela sua senha mestra). Gere códigos de compartilhamento de pasta (`BVFOLDER1`) contendo apenas links e chaves dos arquivos daquela pasta.
- **Envio de Vídeos Ópticos (Dados Brutos ou Vídeo Artístico/Furtivo):**
  - Perfis de dados de alta densidade (`y8-gray`, `y6-gray`, `y4-gray`, `y2-safe`) em até **4K**.
  - Máscaras furtivas (`stealth-stars`, `stealth-vhs`, `stealth-rain`, `stealth-film`) e padrões artísticos/fractais com prévia animada GIF em tempo real, vídeo/imagem base opcional e estúdio completo de trilha sonora generativa.
  - Escolha o **nome personalizado do vídeo** antes de enviar (e renomeie direto no YouTube pela aba **🎬 Meus Vídeos**).
- **Integração com YouTube Data API v3:**
  - Envio automático, verificação de processamento, edição de privacidade/título, backup atômico do catálogo (**Vídeo-Índice**) e botão rápido **🔄 Reconectar YouTube**.
  - Também funciona 100% em **Modo Manual Assistido** sem precisar configurar API.
- **Painel Flutuante de Desempenho:** controle máximo de **% de CPU**, **% de RAM**, modo segundo plano e perfis rápidos (*Máximo*, *Equilibrado*, *Mínimo*).
- **Organização, Prévia e Sincronização Incremental:**
  - Pastas hierárquicas reordenáveis, coleções, tags, linha do tempo, edição de metadados em lote, prévia instantânea antes do download (*Manifest-First*) e vídeos de atualização incremental (cadeias delta).
- **Temas e Idiomas:** tema padrão **✨ Noite Brilhante**, além de *Noite*, *Papel*, *Deserto*, *Light*, *Rainbow* e *Crepúsculo*, em **9 idiomas** (PT, EN, ES, RU, FR, IT, JA, AR, ZH).
- **Atualizações Automáticas:** verifica novas versões no GitHub ao abrir e 1× por dia, exibindo o indicador de versão (`v0.701`) no topo com atualização automática em 1 clique.

---

## Comandos CLI (Opcional)

O mesmo `BrightVault.exe` pode ser usado via terminal:

```powershell
# Codificar arquivo ou pasta inteira em vídeo óptico criptografado
.\BrightVault.exe encode -in "C:\Pasta\Fotos" -out "backup.mp4" -pass "SuaSenhaForte" -res 2k -profile y8-gray

# Inspecionar conteúdo do vídeo sem baixar/decodificar tudo (Manifest-First)
.\BrightVault.exe peek -in "backup.mp4" -pass "SuaSenhaForte"

# Decodificar e restaurar com verificação SHA-256
.\BrightVault.exe decode -in "backup.mp4" -out ".\restaurado" -pass "SuaSenhaForte"
```

---

## Como compilar a partir do código-fonte

Requisitos: **Go 1.24+** (Windows 64-bit).

```powershell
.\build.ps1
```

Ou manualmente:

```powershell
$env:CGO_ENABLED = "0"
go build -ldflags="-H=windowsgui" -o BrightVault.exe ./cmd/ytcloud
```

---

## Aviso Importante

Utilizar plataformas de vídeo como repositório genérico de arquivos pode violar os Termos de Serviço da plataforma. **Nunca mantenha sua única cópia de arquivos insubstituíveis em um único lugar** — siga sempre a regra de backup 3-2-1.

## Licença

Distribuído sob a licença MIT. Veja [LICENSE](LICENSE) para mais detalhes.
