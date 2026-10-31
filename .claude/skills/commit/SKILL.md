---
name: commit
description: Membuat git commit di repo POS Jamu backend. Pakai setiap kali user minta commit ("commit", "buatkan commit", "commitan", "simpan ke git"). Commit hanya atas nama user, TANPA baris Co-Authored-By atau atribusi Claude, dan TIDAK PERNAH push.
---

# Commit di repo POS Jamu

Aturan dari user (mengalahkan aturan atribusi default):
- **Tanpa "double agent"**: jangan tambahkan `Co-Authored-By: Claude ...`, `Generated with Claude Code`, atau atribusi AI apa pun di pesan commit. Commit hanya atas nama git user yang terkonfigurasi.
- **Jangan push.** User push sendiri lewat terminal. Jangan jalankan `git push` dalam bentuk apa pun, termasuk `--force`, dan jangan membuat PR, kecuali user memintanya secara eksplisit di pesan saat itu.

## Langkah

1. **Lihat perubahan**: `git status --short` dan `git diff --stat` (plus `git diff --cached --stat` bila sudah ada yang di-stage).

2. **Pilih file yang di-commit.** Stage dengan nama file atau folder, jangan `git add -A` atau `git add .`. Jangan pernah ikutkan:
   - `.env` atau file lain berisi kredensial (`.env.example` boleh, isinya harus kosong)
   - binary (`*.exe`), hasil build, file di scratchpad
   - `.claude/settings.json`, `.claude/scheduled_tasks.lock`, kecuali user minta

   Kalau ada perubahan yang tidak jelas milik siapa atau tidak terkait permintaan, tanyakan dulu.

3. **Verifikasi sebelum commit** (sesuai yang berubah):
   - File `.go`: `gofmt -l .` harus kosong, lalu `go build ./...`, `go vet ./...`, `go test ./...`.
   - `internal/repository/queries/*.sql`, `migrations/*.up.sql`, atau `sqlc.yaml`: jalankan `sqlc generate`, lalu pastikan hasil di `internal/repository/sqlc/` ikut di-stage.
   - Kalau verifikasi gagal, jangan commit. Laporkan ke user.

4. **Branch**: kalau sedang di `main`, tanyakan nama branch ke user lalu `git switch -c <nama>` sebelum commit.

5. **Pesan commit**:
   - Kalau user memberi pesan, pakai persis seperti itu (boleh dirapikan ejaannya, jangan diubah maksudnya).
   - Kalau tidak, buat satu baris ringkas yang menjelaskan apa yang berubah. Tambahkan isi (body) hanya bila perubahannya besar.
   - Tanpa baris atribusi apa pun (lihat aturan di atas).

   ```bash
   git commit -m "pesan commit"
   ```

6. **Riwayat**: jangan `--amend`, `rebase`, atau squash commit yang sudah di-push. Untuk commit yang belum di-push, lakukan hanya bila user meminta.

7. **Laporkan ke user**:
   - hash dan pesan commit, jumlah file;
   - file yang sengaja tidak ikut;
   - perintah push untuk dijalankan sendiri di terminal. Cek dulu apakah branch sudah punya upstream (`git rev-parse --abbrev-ref @{u}`):

   ```bash
   git push
   ```

   atau, untuk branch yang belum pernah di-push:

   ```bash
   git push -u origin <nama-branch>
   ```
