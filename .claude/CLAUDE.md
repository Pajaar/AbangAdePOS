# POS Jamu — Backend

Sistem kasir (POS) dan pencatatan keuangan untuk warung jamu dan obat herbal milik client.
Repo ini adalah **backend** (Go). Frontend ada di repo terpisah (`pos-jamu-frontend`, Vue 3 + Vite + TypeScript, Pinia, Vue Router).

Dokumen ini merangkum semua keputusan desain yang sudah disepakati. Ikuti keputusan di sini; kalau ada yang tampak bertentangan atau belum jelas, tanyakan dulu, jangan menebak.

## Stack

- Go, router **Gin** (asumsi; kalau diganti, sesuaikan pencocokan pola route di middleware)
- MySQL 8.0.16+ (InnoDB, CHECK constraint aktif)
- `github.com/golang-jwt/jwt/v5`, `golang.org/x/crypto/bcrypt`, `github.com/shopspring/decimal`
- **sqlc** untuk query database (engine `mysql`): SQL ditulis di `internal/repository/queries/*.sql`, kode Go di-generate ke `internal/repository/sqlc/` (package `sqlcgen`). Skema dibaca dari `migrations/` (hanya file `*.up.sql`). Konfigurasi di `sqlc.yaml`. Override tipe kolom DECIMAL ke `decimal.Decimal` (juga versi nullable ke `decimal.NullDecimal`). Pola ini sudah dibandingkan dengan `@Query` gaya Spring Data JPA (SQL sebagai string di kode) dan dipilih sqlc.
- Migrasi (golang-migrate, tiap file punya pasangan `.down.sql`):
  - `000001_init_schema` = isi `pos_jamu_schema.sql`
  - `000002_seed_akses` = isi `pos_jamu_seed_akses.sql`
  - `000003_jamu_sachet_penyesuaian_stok`: tabel `jamu_sachet` (menggantikan `jamu.sachet_id`/`sachet_per_gelas`) dan `penyesuaian_stok`
  - `000004_seed_akun_cek_harga`: seed `akun` dan endpoint `POST /api/v1/kasir/cek-harga` (menu Kasir)
  - Perubahan skema berikutnya selalu lewat file migrasi **baru**; jangan ubah migrasi yang sudah ada karena mungkin sudah dijalankan di database. Setelah menambah migrasi, jalankan `sqlc generate`.

## Aturan wajib saat menulis kode

1. **Uang tidak boleh `float64`.** Pakai `decimal.Decimal`. Kolom database DECIMAL(14,2).
2. **Transaksi database dibuka di service** lewat `database.WithTx`, bukan di middleware. Semua langkah satu operasi bisnis (nota, detail, potong stok, piutang, jurnal) harus dalam satu transaksi.
3. **Data transaksi tidak pernah dihapus.** Gunakan kolom `aktif`/`status`, atau revoke.
4. **Snapshot harga dan HPP** disimpan di baris transaksi saat terjadi. Laporan tidak boleh menghitung ulang dari `produk.hpp_rata_rata` saat ini.
5. **Potong stok dengan row lock** (`SELECT ... FOR UPDATE`) di dalam transaksi supaya dua kasir tidak balapan.
6. Respons JSON seragam lewat `internal/util/response.go`. 401 = belum login / token habis, 403 = role tidak boleh.
7. **Tidak ada SQL berupa string di kode Go.** Semua query ditulis di `internal/repository/queries/` lalu `sqlc generate`. Jangan edit isi `internal/repository/sqlc/`.
   - Satu file `.sql` per domain (mis. `produk.sql`, `penjualan.sql`, `piutang.sql`, `jurnal.sql`). Tiap query diberi anotasi `-- name: NamaQuery :one|:many|:exec|:execresult`. Parameter pakai `?` atau `sqlc.arg(nama)` supaya nama field Go jelas.
   - Struct hasil sqlc hanya boleh dipakai di `service/repository` modul. Repository memetakannya ke tipe `tck/domain`; executor memetakan domain ke `dto` lewat `api/mapper`.
   - Transaksi: InService memanggil `database.WithTx`, lalu membuat repository dengan `*sql.Tx` itu (`repository.NewXRepository(tx)`) untuk semua langkah dalam operasi. Query `SELECT ... FOR UPDATE` ditulis di file `.sql` seperti query lain.
   - Kolom hasil agregat/window function (`SUM()`, `SUM() OVER`, dll.) dibungkus `CAST(... AS DECIMAL(14,2))` supaya terbaca sebagai `decimal.Decimal`, bukan `interface{}`/string.
   - Pengecualian hanya untuk query yang tidak didukung sqlc (filter dinamis banyak, atau fitur MySQL yang gagal di-generate): tulis manual di file `service/repository` modul terkait, beri komentar alasannya, dan tetap terima `sqlcgen.DBTX` supaya ikut transaksi yang sama.

## Struktur folder

Arsitektur meniru proyek primecash-x-training (Spring Boot): setiap fitur adalah **modul** dengan tiga lapisan `tck` → `service` → `api`. Modul `auth` adalah contoh acuan; modul baru (produk, jamu, penjualan, piutang, pembelian, laporan, ...) mengikuti bentuk yang sama.

```
cmd/api/main.go                 # rakit modul, muat izin, jalankan router
cmd/seed-admin/main.go          # buat user pertama lewat UserInService
internal/
  config/  database/ (mysql.go, tx.go = WithTx)  auth/ (jwt, password, token)
  middleware/ recovery, request_id, logger, cors, auth, authorize, rate_limit
  router/router.go              # urutan middleware + grup route, memanggil RegisterRoutes tiap modul
  authorize/store.go            # izin role → endpoint di memori
  common/                       # ≈ library com.product.primecash.common
    processor/                  # IExecutorProcessor, IValidatorProcessor, Validate (≈ ValidationDetailsReport)
    route/                      # Groups (Public, Authenticated, Protected) + Registrar
  dto/                          # ≈ dto-tck: request/response JSON, dipakai bersama
  module/<fitur>/
    tck/                        # KONTRAK, tidak bergantung ke lapisan lain
      domain/                   # tipe domain (≈ entity), bukan struct sqlc
      inservice/                # interface I<X>InService (≈ tck/internal)
    service/
      repository/               # I<X>Repository + implementasi di atas sqlcgen, petakan row → domain
      inservice/                # implementasi InService: logika domain + database.WithTx (≈ service/internal)
    api/
      controller/               # handler Gin: bind, isi data HTTP (cookie/context), panggil executor, kirim respons
      executor/                 # satu struct per use case, Process(ctx, req) (res, error)
      validator/                # satu struct per aturan bisnis (bila perlu), dipanggil executor lewat processor.Validate
      mapper/                   # domain → dto
      config/                   # NewModule (rakit repository → inservice → executor → controller) + RegisterRoutes
  repository/
    queries/                    # file .sql per domain (sumber query, ditulis manual)
    sqlc/                       # hasil sqlc generate (package sqlcgen), jangan diedit
  integration/ openfoodfacts.go, upcitemdb.go   # klien API barcode pihak ketiga
  util/  response, apperror, money, barcode
migrations/
sqlc.yaml                       # konfigurasi sqlc (schema: migrations/, queries: internal/repository/queries/)
```

Alur satu request: **controller → executor → (validator) → InService (interface tck) → repository → sqlcgen**.
- Controller tipis: tidak ada logika bisnis.
- Executor hanya bergantung pada interface `tck/inservice`, sehingga bisa diuji dengan implementasi palsu tanpa database (contoh: `module/auth/api/executor/executor_test.go`).
- Operasi yang harus atomik (nota + detail + potong stok + piutang + jurnal) adalah satu method InService yang membuka satu `database.WithTx` dan memakai beberapa repository di dalamnya.
- Folder bernama `inservice`, bukan `internal`, karena di Go folder `internal` membatasi siapa yang boleh meng-import.
- Untuk menambah modul: buat `module/<fitur>/...`, lalu daftarkan `NewModule` di `cmd/api/main.go` (`Modules: []route.Registrar{...}`).

## Urutan middleware

```
Recovery → RequestID → Logger → CORS (hanya development)
  ├── /api/v1/auth/login            + RateLimit         (publik)
  ├── /api/v1/auth/refresh, logout                      (publik, pakai cookie)
  ├── /api/v1/auth/me, /auth/menus  + Auth              (semua user login)
  └── /api/v1/* lainnya             + Auth + Authorize
```

CORS hanya dipasang bila `CORS_ALLOWED_ORIGINS` di `.env` terisi (lihat Deployment). Bila dipasang, harus sebelum Auth (preflight `OPTIONS` tidak membawa token).

## Autentikasi

- **Access token**: JWT, dikirim sebagai `Authorization: Bearer <token>`. Umur ±15 menit. Claims: `sub` (user_id), `role` (kode role), `exp`, `iat`. Tidak disimpan di database. Secret dari `JWT_SECRET` di `.env`.
- **Refresh token**: string acak, dikirim di cookie `httpOnly`. Umur ±7 hari. Di database hanya disimpan SHA-256-nya (`refresh_tokens.token_hash`).
- **Rotasi**: setiap refresh, token lama di-revoke, `replaced_by_id` menunjuk token baru. Kalau token yang sudah di-revoke dipakai lagi → revoke semua token user itu (indikasi dicuri).
  - Cek dan rotasi dilakukan dalam satu transaksi dengan `SELECT ... FOR UPDATE` (`SessionInService.Rotate`), jadi dua refresh bersamaan dengan token yang sama tidak bisa sama-sama lolos. Revoke massal saat reuse harus ikut ter-commit (fn transaksi return nil, error dikembalikan sesudahnya).
  - Akibatnya, frontend wajib menjalankan **satu refresh pada satu waktu**: request lain yang kena 401 menunggu hasil refresh yang sedang berjalan. Kalau tidak, refresh kedua dianggap pencurian dan user ter-logout.
- Refresh ditolak kalau user `aktif = false`.
- Login gagal: pesan umum "Username atau password salah".

## Otorisasi (role → menu → endpoint)

- Tabel: `roles`, `users`, `menus`, `api_endpoints`, `menu_endpoints`, `role_menus`.
- `api_endpoints.path` berisi **pola route** (`/api/v1/produk/:id`). Middleware `Authorize` membandingkan `c.Request.Method` + `c.FullPath()`, bukan URL asli.
- Izin dimuat ke memori saat start (map role → set "METHOD path"), dimuat ulang saat hak akses diubah lewat `PUT /api/v1/roles/:id/akses`.
- Headadmin diberi semua menu secara eksplisit di data, **bukan** bypass di kode.
- `menus.tipe = 'aksi'` adalah izin tanpa halaman. Saat ini: `lihat_hpp`. Service wajib mengecek izin ini dan **menghapus field HPP, HPP rata-rata, dan margin dari respons** kalau role tidak punya izin.

| Role | Menu |
|---|---|
| headadmin | Semua (termasuk Laporan, Pengguna, lihat_hpp) |
| admin | Produk, Menu Jamu, Pembelian (fokus input barang, tanpa laporan/HPP/keuntungan) |
| cashier | Kasir, Piutang |

## Model bisnis

### Jenis barang
- **`produk`**: semua barang yang distok (obat, sachet jamu, minyak urut, dll). Dikelompokkan lewat `kategori` (murni pengelompokan, tidak memengaruhi logika). Stok **selalu** dipotong saat terjual.
- **`jamu`**: menu jamu **seduh per gelas** (diseduh langsung saat dipesan, tidak pakai termos). Tidak punya stok sendiri. Sachet yang dipakai ada di `jamu_sachet` (boleh lebih dari satu jenis per jamu). Saat terjual, stok setiap sachet dipotong sebanyak `jumlah × qty gelas`. Harga seduh (`jamu.harga_jual`) berbeda dari harga sachet utuh (`produk.harga_normal`).
- **`additional_items`**: tambahan untuk jamu seduh, harga fix. Kalau `produk_id` terisi, stok produk itu dipotong; kalau kosong (madu, telur), tidak dilacak dan HPP-nya 0.
- `penjualan_detail` berisi **produk ATAU jamu** (CHECK constraint), tidak keduanya.

### HPP: rata-rata bergerak (moving average)
Saat pembelian:
```
hpp_baru = (stok_lama * hpp_lama + qty_beli * harga_beli) / (stok_lama + qty_beli)
jika stok_lama <= 0: hpp_baru = harga_beli
```
Stok awal (`penyesuaian_stok` jenis `stok_awal`) dihitung dengan rumus yang sama.

HPP jamu seduh per gelas = jumlah dari `jamu_sachet.jumlah * hpp_rata_rata` setiap sachet.

### Harga di kasir dan peringatan di bawah modal
- **Role cashier tidak boleh menerima angka HPP, modal, nett, tengah, atau margin dari API dalam bentuk apa pun.** Termasuk tidak dikirim tersembunyi di respons lalu disembunyikan di frontend.
- Produk dengan `bisa_ditawar = true` (biasanya obat): kasir mengetik harga jadi per baris. Produk lain dan jamu memakai harga tetap.
- Diskon **hanya per nota** (`penjualan.diskon`), input nominal manual.
- Indikator **merah** (peringatan saja, **tidak diblokir**) muncul bila harga jadi suatu baris di bawah HPP baris itu, atau total nota setelah diskon di bawah total HPP nota.
- Pengecekan dilakukan di server lewat `POST /api/v1/kasir/cek-harga`: menerima isi keranjang dan diskon, membalas hanya boolean per baris dan per nota (`di_bawah_modal`), tanpa angka. Frontend memanggilnya dengan debounce saat harga atau diskon berubah. Saat transaksi disimpan, server menghitung ulang dan mengembalikan flag yang sama.
- Bila HPP produk belum diketahui (`hpp_rata_rata = 0`), balas `modal_diketahui: false` untuk baris itu. Frontend tidak memerahkan, cukup tanda netral.
- Harga **nett** = `hpp_rata_rata * (1 + margin_min_persen/100)` dan **tengah** = titik tengah harga normal dan nett, keduanya dibulatkan ke atas ke kelipatan Rp500. Hanya ditampilkan untuk role dengan izin `lihat_hpp` (misalnya halaman Produk), tidak disimpan. Bila nett ≥ harga normal, tandai produk "harga perlu dinaikkan".

### Stok
- Penjualan **boleh** membuat stok 0 atau minus, karena barang lama belum semua diinput. Tampilkan peringatan di kasir, jangan blokir.
- Stok awal barang lama dan koreksi stock opname dicatat di `penyesuaian_stok` (jenis `stok_awal`, `opname`, `rusak`).
- Selama HPP sebuah produk masih 0, penjualannya tercatat dengan HPP 0, sehingga laba terlihat lebih besar dari aslinya dan indikator merah tidak berfungsi. Stok awal sebaiknya diinput sebelum sistem dipakai.

### Struk
Cetak per nota: nama barang, qty, harga, subtotal, diskon, total, dibayar, kembalian atau sisa utang. Untuk produk yang punya `aturan_pakai`, cetak teksnya di bawah nama barang. Teks panjang dibungkus ke baris berikutnya, jangan dipotong.

### Diskon vs piutang (dua hal berbeda)
- `penjualan.diskon`: potongan nominal per nota, input manual, untuk pelanggan langganan. Mengurangi tagihan.
- Piutang: tagihan tidak berkurang, hanya belum dibayar. `sisa_piutang = total - dibayar`.
- Kalau `dibayar < total`, transaksi hanya boleh disimpan bila ditandai utang **dan** `pelanggan_id` terisi (dijaga juga oleh CHECK constraint).
- `pelanggan` hanya untuk orang yang ngutang. Pembeli tunai: `pelanggan_id = NULL`. Nama pelanggan tidak dihapus saat lunas.
- Pembayaran piutang: insert `pembayaran_piutang` + kurangi `penjualan.sisa_piutang` dalam satu transaksi.

### Barcode
- `produk_barcode`: satu produk boleh banyak barcode. `isi_per_kemasan` = jumlah stok yang dipotong per scan (barcode kotak isi 10 → 10).
- Normalisasi sebelum simpan/cari: UPC-A 12 digit → EAN-13 (tambah `0` di depan). Validasi checksum.
- **Di kasir**: hanya cari di database lokal (`GET /api/v1/barcode/:kode`). Tidak ketemu → 404, status `draft` → 409.
- **Di tambah produk / pembelian**: boleh lookup API publik (`POST /api/v1/barcode/lookup`): cek lokal → cek `barcode_lookup_cache` → Open Food Facts → UPCitemdb, timeout 3 detik per API. Simpan hasil (termasuk "tidak ditemukan") ke cache. **Error/timeout jangan di-cache.** Nama dari API hanya saran, user wajib cek.
- Produk hasil lookup yang belum diisi harga/kategori: `status = 'draft'`, tidak tampil di kasir.

## Jurnal otomatis

Setiap transaksi membuat `jurnal` + `jurnal_detail` di transaksi database yang sama, dengan `sumber_tipe`/`sumber_id`. Total debit harus sama dengan total kredit.

| Transaksi | Debit | Kredit |
|---|---|---|
| Penjualan | Kas (dibayar), Piutang Usaha (sisa), Potongan Penjualan (diskon nota + selisih tawar) | Penjualan (total harga normal) |
| Penjualan (biaya) | HPP | Persediaan (total snapshot HPP) |
| Bayar piutang | Kas | Piutang Usaha |
| Pembelian | Persediaan | Kas (dibayar), Utang Usaha (sisa) |
| Bayar utang supplier | Utang Usaha | Kas |
| Pengeluaran | Beban Operasional | Kas |
| Stok awal | Persediaan | Modal Pemilik |
| Opname/rusak (stok berkurang) | HPP | Persediaan |
| Opname (stok bertambah) | Persediaan | HPP |

**Selisih tawar** = jumlah per baris `(harga normal − harga jadi) × qty`, hanya untuk produk `bisa_ditawar` yang harga jadinya di bawah harga normal. Contoh: harga normal 20.000, harga jadi 18.000, qty 2, diskon nota 1.000, dibayar tunai → Kas 35.000 + Potongan Penjualan 5.000 = Penjualan 40.000.

Buku besar = query `jurnal_detail` per akun. Buku pembantu piutang/utang = query `penjualan`/`pembayaran_piutang` per pelanggan dan `pembelian`/`pembayaran_utang` per supplier.

## Laporan (hanya headadmin)
Laba rugi (Penjualan, Potongan Penjualan, Penjualan Bersih, HPP, Laba Kotor, Beban, Laba Bersih), margin per produk dan per jamu, total diskon, buku besar.

## Deployment
Frontend dan backend disajikan dari **satu origin** lewat reverse proxy (Caddy atau Nginx): frontend di `/`, API di `/api/`. Repo tetap terpisah.
- Cookie refresh token: `HttpOnly; Secure; SameSite=Strict; Path=/api/v1/auth`.
- CORS hanya diaktifkan saat development (Vite dev server), lewat `CORS_ALLOWED_ORIGINS` di `.env`. Di produksi dikosongkan.

## Daftar akun (usulan default, belum dikonfirmasi pemilik)
Dipakai sebagai seed awal `akun`. Kode dan nama boleh diubah setelah pemilik melihat datanya.

| Kode | Nama | Tipe |
|---|---|---|
| 1-101 | Kas | aset |
| 1-102 | Piutang Usaha | aset |
| 1-103 | Persediaan | aset |
| 2-101 | Utang Usaha | kewajiban |
| 3-101 | Modal Pemilik | modal |
| 4-101 | Penjualan | pendapatan |
| 4-102 | Potongan Penjualan | pendapatan (saldo normal debit, pengurang penjualan) |
| 5-101 | HPP | beban |
| 6-101 | Beban Operasional | beban |

Sengaja dibuat minimal. Semua pengeluaran masuk Beban Operasional, jenisnya dibedakan lewat `pengeluaran.keterangan`. Selisih stock opname dan barang rusak dimasukkan ke HPP. Akun baru bisa ditambah nanti tanpa mengubah kode.

## Belum diputuskan (tanyakan sebelum implementasi bagian terkait)
- Konfirmasi daftar akun di atas oleh pemilik.
- Apakah ada pembayaran selain tunai (QRIS, transfer)? Jika ada, perlu akun Bank dan kolom metode bayar di penjualan.
- Server di cloud atau di komputer toko (jaringan lokal)? Menentukan apakah kasir tetap bisa dipakai saat internet mati.
- Ukuran printer struk (thermal 58 mm atau 80 mm), untuk lebar baris struk.
- Endpoint dan menu untuk penyesuaian stok (stok awal, opname, rusak): role mana yang boleh, dan masuk menu Produk atau menu sendiri. Tabelnya sudah ada, endpoint belum didaftarkan di seed.
- Bolehkah harga jadi di atas harga normal? Jika boleh, jurnalnya bagaimana (misalnya kredit Penjualan memakai harga jadi supaya selisih tawar tidak negatif).
