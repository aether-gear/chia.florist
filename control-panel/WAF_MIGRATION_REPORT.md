# Laporan Perubahan: Migrasi WAF Sentinel ke Ekosistem Chia Florist

Dokumen ini ditujukan untuk menyelaraskan pemahaman antara arsitektur `WAF-Golang` yang lama (berdiri sendiri) dengan arsitektur baru yang terintegrasi di dalam proyek `chia.florist/control-panel`.

Silakan berikan dokumen ini kepada AI Pembuat Laporan Anda agar ia dapat memperbarui konteks proyek terbaru.

---

## 1. Konsep Utama (The Core Concept)
**Tidak Ada Perubahan pada Filosofi.**
Pendekatan *Security-First* dan *Clean Architecture* tetap dipertahankan. WAF Golang tetap bertindak sebagai garda terdepan (Middleware) yang independen. 
**Perubahan Utama:** Manajemen antarmuka (Dashboard) kini tidak lagi ego-sentris (berdiri sendiri hanya untuk keamanan), melainkan **disatukan ke dalam satu Control Panel terpusat** yang juga memantau E-commerce, Penjualan, dan metrik AI.

## 2. Arsitektur WAF Sentinel (Golang Backend)
**Status: 100% SAMA (Dipertahankan)**
Kode sumber `main.go`, logika keamanan, dan arsitektur data dari WAF Golang dipindahkan secara utuh ke dalam folder `chia.florist/control-panel`.
*   **Bahasa:** Tetap menggunakan Golang.
*   **Posisi:** Tetap sebagai Middleware / Reverse Proxy (`:8080`).
*   **Penyimpanan (Persistence):** Tetap menggunakan Flat-File JSON (`waf-logs.json`, `waf-rules.json`, `waf-blocked.json`, `waf-filters.json`). Tidak ada transisi ke SQL.
*   **Konkurensi:** Tetap menggunakan `sync.RWMutex` untuk keamanan operasi file pada lalu lintas padat.
*   **Engine Deteksi:** Tetap menggunakan Regex standar Go (`regexp`).

## 3. Arsitektur Dashboard WAF (Frontend)
**Status: PERUBAHAN BESAR (RE-WRITE TERINTEGRASI)**
Inilah area di mana tim Anda melakukan perombakan besar untuk menyelaraskan UI/UX keseluruhan sistem Chia Florist.

*   **Tech Stack Lama:** Next.js (React App Router).
*   **Tech Stack Baru:** **Vite + React (Single Page Application) + TypeScript.**
*   **Komponen UI:** Kini menggunakan **Tailwind CSS v3** yang dikombinasikan dengan library komponen modern **shadcn/ui** (seperti yang digunakan di halaman E-Commerce utama).
*   **Struktur Halaman:** Dashboard keamanan WAF kini menjadi salah satu menu utama (berada di menu `/security`) di dalam `control-panel`, bersanding dengan menu manajemen toko lainnya (Shop, Products, Orders).
*   **Komunikasi Data:** 
    *   *Sebelumnya:* Next.js melakukan polling REST API ke Go (`/api/stats`).
    *   *Sekarang (Development Phase):* Karena Frontend (Vite) dan Backend (Go) digabung dalam satu folder kerja untuk memudahkan kolaborasi, frontend saat ini difasilitasi untuk **membaca file JSON secara langsung** (melalui `src/data/wafData.ts` yang meng-import `waf-logs.json` sebagai modul statis). Ke depannya, ketika dide-deploy ke production, ini akan dikembalikan ke metode *API Polling*.
*   **Visualisasi:** 
    *   Tetap menggunakan **Recharts** untuk membuat *Threat Landscape Chart* (Area Chart yang mengkomparasi *Allowed* vs *Blocked* request).
    *   *Live Attack Map* (Leaflet.js) untuk sementara ditiadakan/belum di-porting ke versi Vite, difokuskan pada ringkasan tabel ancaman (*Recent Logs*), *Anomaly Score*, dan manajemen aturan (*WAF Rules Status*).

## 4. Alur Kerja Sistem Terintegrasi (Request Lifecycle)
**Status: SAMA (Dipertahankan)**
Alur kerja yang disepakati di Laporan Form 3 masih berlaku seratus persen. 
1.  Request masuk ke WAF Golang.
2.  Pemindaian IP pada `waf-blocked.json` dan `waf-filters.json`.
3.  Pemindaian Payload melalui Regex dari `waf-rules.json`.
4.  Jika bahaya: Diblokir (403 Forbidden) & dicatat ke file.
5.  Jika aman: Diteruskan ke E-Commerce Backend (layanan inti Chia Florist).
6.  AI Analytics membaca log JSON tersebut untuk *Business Insights*.

---
### Kesimpulan untuk AI Laporan:
Instruksikan AI Laporan Anda dengan kalimat ini:
> *"Modul WAF Golang (Backend) tidak mengalami perubahan arsitektur sedikitpun, hanya lokasinya saja yang dipindahkan. Namun, **Dashboard Next.js yang lama telah ditinggalkan**. Sebagai gantinya, antarmuka pemantauan WAF dibangun ulang dari awal menggunakan **React + Vite + shadcn/ui** dan diintegrasikan ke dalam **Control Panel Utama (chia.florist/control-panel)** agar Admin toko dapat memantau keamanan siber, penjualan produk, dan analitik AI dalam satu aplikasi web yang sama."*
