# WAF Golang Migration & UI/UX Refinement Changelog
Dokumen ini mendokumentasikan evolusi kode dari integrasi Golang hingga perbaikan detail visual (Contrast, Hover, & Layout) dengan mengikuti workflow `@dokumentasi-perubahan.md`.

---

## 1. Global Theme & Border Fix
**1. Informasi Umum**
*   **Judul Perubahan:** Global Dark Mode Activation
*   **Tanggal:** 2026-02-01
*   **Requester:** User (Visual Polish)
*   **Agent / Modul Terkait:** UI Layout
*   **File Tujuan Dokumentasi:** [layout.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/layout.tsx)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Munculnya border putih tebal dan padding terang pada komponen Shadcn karena sistem menganggap aplikasi dalam mode terang secara default.
*   **Dampak Jika Tidak Diubah:** Estetika Dark Mode tidak konsisten, beberapa komponen terlihat pecah secara visual.

**3. Ruang Lingkup Perubahan**
*   Modifikasi pada file `layout.tsx` untuk memaksa class `dark`.

**4. Komparasi Kode**

**Sebelum Diubah**
```tsx
<html lang="en">
```

**Sesudah Diubah**
```tsx
<html lang="en" className="dark">
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Menambahkan `className="dark"` pada tag html utama.
*   **Alasan Pemilihan Solusi:** Cara tercepat untuk memicu variabel CSS dark mode.
*   **Lokasi Kode:** [layout.tsx:L26](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/layout.tsx#L26)

**6. Risiko & Dampak**
*   **Potensi Risiko:** Tidak ada.

**7. Testing & Validasi**
*   **Hasil Testing:** Seluruh komponen Shadcn sekarang menggunakan skema warna gelap yang halus.

**8. Status**
*   **Status Perubahan:** Selesai.

**9. Catatan Tambahan**
*   PIC: Antigravity

---

## 2. Tabel Log (Hover & Focus Polish)
**1. Informasi Umum**
*   **Judul Perubahan:** Log Table Row Interaction Polish
*   **Tanggal:** 2026-02-01
*   **Requester:** User (Visual Polish)
*   **Agent / Modul Terkait:** Dashboard UI
*   **File Tujuan Dokumentasi:** [page.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/page.tsx)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Hover baris tabel terlalu terang dan muncul outline biru saat diklik.
*   **Dampak Jika Tidak Diubah:** UI terasa kurang premium.

**4. Komparasi Kode**

**Sebelum Diubah**
```tsx
<TableRow 
    className="border-slate-800 hover:bg-slate-800/50 cursor-pointer transition-colors"
    onClick={() => toggleExpand(log.id)}
>
```

**Sesudah Diubah**
```tsx
<TableRow
    className="border-slate-800 hover:bg-slate-900 cursor-pointer transition-colors focus:outline-none focus:ring-0"
    onClick={() => toggleExpand(log.id)}
>
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Mengganti warna hover ke `bg-slate-900` dan mematikan focus ring.
*   **Lokasi Kode:** [page.tsx:L329](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/page.tsx#L329)

**8. Status**
*   **Status Perubahan:** Selesai.

---

## 3. SIEM-style Advanced Visualization (Bucket Aggregation)
**1. Informasi Umum**
*   **Judul Perubahan:** SIEM Analytics Dashboard (Zero-Filling)
*   **Tanggal:** 2026-02-02
*   **Requester:** User (SIEM Experience)
*   **Agent / Modul Terkait:** Performance & Analytics
*   **File Tujuan Dokumentasi:** [TrafficChart.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/components/TrafficChart.tsx)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Grafik traffic sebelumnya tidak linear (lompat jika data kosong).
*   **Dampak Jika Tidak Diubah:** Analisis tren waktu tidak akurat.

**4. Komparasi Kode**

**Sebelum Diubah**
```tsx
const chartData = data.map(log => ({ time: log.time, count: 1 }));
```

**Sesudah Diubah**
```tsx
while (current <= endTime) {
    buckets.push({
        timestamp: current.getTime(),
        label: formatLabel(current, timeRange),
        requests: 0 
    });
    current = new Date(current.getTime() + intervalMs);
}
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Implementasi loop untuk pembuatan bucket waktu kosong (Zero-Filling).
*   **Lokasi Kode:** [TrafficChart.tsx:L78-85](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/components/TrafficChart.tsx#L78-L85)

**8. Status**
*   **Status Perubahan:** Selesai.

---

## 4. Responsive & Anti-Lag Chart Optimization
**1. Informasi Umum**
*   **Judul Perubahan:** High-Performance Responsive Chart
*   **Tanggal:** 2026-02-02
*   **Requester:** User (Performance)
*   **Agent / Modul Terkait:** Traffic Analysis
*   **File Tujuan Dokumentasi:** [TrafficChart.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/components/TrafficChart.tsx)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Grafik macet saat range waktu panjang (data terlalu banyak).
*   **Dampak Jika Tidak Diubah:** Browser freeze pada range 30-hari.

**4. Komparasi Kode**

**Sebelum Diubah**
```tsx
const estimatedBuckets = rangeMs / intervalMs;
// Tidak ada pembatasan jumlah bucket
```

**Sesudah Diubah**
```tsx
const estimatedBuckets = rangeMs / intervalMs;
if (estimatedBuckets > 500) {
    intervalMs = Math.ceil(rangeMs / 500); 
}
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Menambahkan logic limitator (cap) 500 bucket.
*   **Lokasi Kode:** [TrafficChart.tsx:L55-62](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/components/TrafficChart.tsx#L55-L62)

**8. Status**
*   **Status Perubahan:** Selesai.

---

## 5. Security & Privacy: API Key Masking
**1. Informasi Umum**
*   **Judul Perubahan:** API Key UI Scrubbing
*   **Tanggal:** 2026-02-02
*   **Requester:** User (Security Concerns)
*   **Agent / Modul Terkait:** Security
*   **File Tujuan Dokumentasi:** [page.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/page.tsx)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** API Key terlihat oleh orang lain.

**4. Komparasi Kode**

**Sebelum Diubah**
```tsx
<div className="flex gap-4 items-center">
    <Input value={vtKey} /> // Key terlihat
```

**Sesudah Diubah**
```tsx
// Komponen Input API key dihapus sepenuhnya dari render
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Menghapus elemen JSX yang me-render API keys.
*   **Lokasi Kode:** [page.tsx:L186-206](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/page.tsx#L186-L206) (Header cleaning)

**8. Status**
*   **Status Perubahan:** Selesai.

---

## 6. Definitive Persistence Fix (Initialization Guard)
**1. Informasi Umum**
*   **Judul Perubahan:** Persistent Dashboard Preferences (Race-Condition Fix)
*   **Tanggal:** 2026-02-02
*   **Requester:** User (Lost Settings on Reload)
*   **Agent / Modul Terkait:** Frontend State
*   **File Tujuan Dokumentasi:** [page.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/page.tsx)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Pilihan user ter-reset ke default setiap refresh.

**4. Komparasi Kode**

**Sebelum Diubah**
```tsx
useEffect(() => {
    localStorage.setItem('waf_chart_range', chartTimeRange);
}, [chartTimeRange]);
```

**Sesudah Diubah**
```tsx
const [isReady, setIsReady] = useState(false);

useEffect(() => {
    // Restore logic...
    setIsReady(true);
}, []);

useEffect(() => {
    if (!isReady) return; 
    localStorage.setItem('waf_chart_range', chartTimeRange);
}, [chartTimeRange, isReady]);
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Menambahkan state `isReady` sebagai gerendel (guard).
*   **Lokasi Kode:** [page.tsx:L27-52](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/page.tsx#L27-L52)

**8. Status**
*   **Status Perubahan:** Selesai.

---

## 7. Smart IP Suggestions (Log Scanning)
**1. Informasi Umum**
*   **Judul Perubahan:** IP Action Suggestions from Logs
*   **Tanggal:** 2026-02-02
*   **Requester:** User (Typo avoidance)
*   **Agent / Modul Terkait:** Configuration UI
*   **File Tujuan Dokumentasi:** [config/page.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/config/page.tsx)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Admin kesulitan mem-ban IP baru karena harus ketik manual.

**4. Komparasi Kode**

**Sebelum Diubah**
```tsx
// Tidak ada pendeteksian IP otomatis
```

**Sesudah Diubah**
```tsx
// Scan stats.logs untuk IP baru yang belum terdaftar
const suggestionList = logsIPs.filter(ip => !existingIPs.has(ip));
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Implementasi filter log untuk ekstraksi IP unik yang belum dimanajemen.
*   **Lokasi Kode:** [config/page.tsx:L120-160](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/config/page.tsx#L120-L160)

**8. Status**
*   **Status Perubahan:** Selesai.

---

## 8. Backend IP Logic Refinement (IPv6 Support)
**1. Informasi Umum**
*   **Judul Perubahan:** Fix Missing IP (IPv6 Bracket Issue)
*   **Tanggal:** 2026-02-02
*   **Requester:** User (Bug Report)
*   **Agent / Modul Terkait:** Backend Go
*   **File Tujuan Dokumentasi:** [main.go](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/main.go)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** IP localhost terbaca sebagai `[` bukan `::1`.

**4. Komparasi Kode**

**Sebelum Diubah**
```go
return strings.Split(r.RemoteAddr, ":")[0]
```

**Sesudah Diubah**
```go
if strings.Contains(ip, "[") && strings.Contains(ip, "]") {
    parts := strings.Split(ip, "]")
    ip = strings.TrimPrefix(parts[0], "[")
}
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Penambahan parser string khusus untuk menangani format bracket IPv6 Go.
*   **Lokasi Kode:** [main.go:L215-243](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/main.go#L215-L243)

**8. Status**
*   **Status Perubahan:** Selesai.

---

## 9. IP Ban Reason / Notes Feature
**1. Informasi Umum**
*   **Judul Perubahan:** Add Note/Reason Feature for Banned IPs
*   **Tanggal:** 2026-02-24
*   **Requester:** User (Feature Request)
*   **Agent / Modul Terkait:** Configuration UI & WAF Backend
*   **File Tujuan Dokumentasi:** [main.go](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/main.go) dan [page.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/config/page.tsx)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Admin butuh memberikan catatan (reason) kenapa IP tertentu diban.

**4. Komparasi Kode**

**Sebelum Diubah (main.go)**
```go
type IPConfig struct {
	BlockedIPs  []string `json:"blocked_ips"`
```

**Sesudah Diubah (main.go)**
```go
type IPRecord struct {
	IP     string `json:"ip"`
	Reason string `json:"reason,omitempty"`
}
type IPConfig struct {
	BlockedIPs  []IPRecord `json:"blocked_ips"`
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Mengubah struktur data dari array string menjadi array objek yang menampung string IP dan reason. Menambahkan field "Reason" di UI untuk menampung input pengguna.
*   **Lokasi Kode:**
    *   [main.go](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/main.go)
    *   [config/page.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/config/page.tsx)

**8. Status**
*   **Status Perubahan:** Selesai.

---

## 10. Hydration Error Fix (Extension Conflict)
**1. Informasi Umum**
*   **Judul Perubahan:** Fix React Hydration Mismatch
*   **Tanggal:** 2026-02-24
*   **Requester:** User (Bug Report)
*   **Agent / Modul Terkait:** UI Layout
*   **File Tujuan Dokumentasi:** [layout.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/layout.tsx)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Muncul error Hydration Mismatch ("A tree hydrated but some attributes of the server rendered HTML didn't match the client properties") karena ada atribut `data-jetski-tab-id` yang diinjeksi oleh browser extension eksternal.
*   **Dampak Jika Tidak Diubah:** Console dipenuhi error merah dan React hydration gagal mencocokkan DOM tree.

**4. Komparasi Kode**

**Sebelum Diubah**
```tsx
<html lang="en" className="dark">
```

**Sesudah Diubah**
```tsx
<html lang="en" className="dark" suppressHydrationWarning>
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Menambahkan properti `suppressHydrationWarning` pada tag `<html>`.
*   **Alasan Pemilihan Solusi:** Cara standar dan direkomendasikan Next.js untuk mencegah peringatan/error saat browser extension mengubah struktur/atribut pada level `<html>` atau `<body>`.
*   **Lokasi Kode:** [layout.tsx:L26](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/layout.tsx#L26)

**8. Status**
*   **Status Perubahan:** Selesai.

## 11. Add SQLi Rule (OR Logic Prevention)
**1. Informasi Umum**
*   **Judul Perubahan:** SQL Injection Rule Addition
*   **Tanggal:** 2026-03-30
*   **Requester:** User (FR-02 Testing Validation)
*   **Agent / Modul Terkait:** WAF Rules Engine
*   **File Tujuan Dokumentasi:** [waf-rules.json](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/waf-rules.json)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Saat pengujian Capstone, WAF belum memiliki rule spesifik untuk mendeteksi *payload* *SQL Injection* bertipe logika `OR` (contoh: `' OR 1=1 --`).
*   **Dampak Jika Tidak Diubah:** WAF akan meloloskan serangan ini dan berpotensi meretas database *backend*.

**3. Ruang Lingkup Perubahan**
*   Penambahan rule baru dengan ID `1003` di file JSON aturan WAF.

**4. Komparasi Kode**

**Sebelum Diubah**
```json
// Hanya ada rule hingga ID 1002
```

**Sesudah Diubah**
```json
{
    "id": "1003",
    "description": "SQLi Detection (OR Logic)",
    "pattern": "(?i)('\\s*OR\\s+\\d+=\\d+\\s*--)",
    "tags": ["sqli", "injection"],
    "impact": "5"
}
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Menambahkan blok JSON berisi *Regular Expression* (Regex) yang melakukan pengecekan case-insensitive `(?i)` atas pola injeksi logika.
*   **Lokasi Kode:** [waf-rules.json](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/waf-rules.json)

**6. Risiko & Dampak**
*   **Potensi Risiko:** Ada risiko *crash* (Go Panic) jika karakter luput di-*escape* (contoh: backslash ganda di JSON). Ini sempat terjadi dan langsung disempurnakan.

**7. Testing & Validasi**
*   **Hasil Testing:** Serangan `/catalog?id=' OR 1=1 --` berhasil menghasilkan status 403 Forbidden dan diloging dengan warna merah.

**8. Status**
*   **Status Perubahan:** Selesai.

## 12. Fix VirusTotal/GeoIP Scan State Reset (React Re-rendering Bug)
**1. Informasi Umum**
*   **Judul Perubahan:** Prevent Child Component State Reset on Log Refresh
*   **Tanggal:** 2026-03-30
*   **Requester:** User (API Limit Protection & UX)
*   **Agent / Modul Terkait:** Frontend Dashboard UI
*   **File Tujuan Dokumentasi:** [page.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/page.tsx)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Saat *Auto-Refresh* (misal: 1s) berjalan, hasil *scan* VirusTotal dan GeoIP pada detail log yang sedang terbuka mendadak hilang/ter-reset jika ada log lalu lintas baru yang masuk. Admin harus menekan "Scan" ulang terus-menerus, membuang kuota/limit API.
*   **Dampak Jika Tidak Diubah:** Pemborosan kuota API VirusTotal/GeoIP secara drastis (Limit bisa habis dalam sekian menit) dan navigasi analisis (*Threat Intelligence*) yang sangat tidak nyaman sewaktu diserang telak.

**3. Ruang Lingkup Perubahan**
*   Modifikasi properti bawaan antarmuka bereaksi pada file *Frontend* pemonitor `page.tsx`.

**4. Komparasi Kode**

**Sebelum Diubah**
```tsx
<React.Fragment key={`${log.id}-${index}`}>
```

**Sesudah Diubah**
```tsx
<React.Fragment key={log.id}>
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Menghapus parameter `-${index}` dari `key` milik React element. Karena urutan index elemen otomatis bergeser turun saat log baru ditambahkan di paling atas, React akan membaca `index` yang lama tergeser dan menganggap barisan tabel itu dihancurkan (*unmount*) dan dibuat struktur baru. Ini memicu pengosongan komponen `<VirusTotalWidget>`, menghancurkan hasil *scanning*. Penggunaan pengenal *String ID* unik absolut memandu *React* murni cuma 'menggeser' visual DOM itu ke bawah tanpa menghapus keadaan objek (*Preserved State*).
*   **Lokasi Kode:** [page.tsx:L327](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/page.tsx#L327)

**8. Status**
*   **Status Perubahan:** Selesai.

## 13. Activate "Add Custom Rule" UI Feature
**1. Informasi Umum**
*   **Judul Perubahan:** Implementation of Custom Rule Generator UI
*   **Tanggal:** 2026-03-30
*   **Requester:** User (Feature Request)
*   **Agent / Modul Terkait:** Configuration UI & WAF Backend
*   **File Tujuan Dokumentasi:** [page.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/config/page.tsx)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Tombol "Add Custom Rule" di halaman konfigurasi sebelumnya hanyalah *mockup* pajangan. Pengguna harus membuka dan mengedit file `waf-rules.json` secara manual jika ingin menambahkan rule baru.
*   **Dampak Jika Tidak Diubah:** Alur kerja terasa tidak intuitif bagi Administrator WAF tanpa keahlian pemrograman *backend*.

**3. Ruang Lingkup Perubahan**
*   Pembuatan *form Input* pada halaman web.
*   Integrasi `fetch POST` ke endpoint Go `http://localhost:8080/api/rules`.

**4. Komparasi Kode**

**Sebelum Diubah**
```tsx
<Button size="sm" className="bg-blue-600 hover:bg-blue-500"><Plus className="mr-2 h-4 w-4" /> Add Custom Rule</Button>
// Tidak punya event onClick maupun Form penyerta
```

**Sesudah Diubah**
```tsx
const handleAddRule = async () => {
    // ... data construction
    await fetch('http://localhost:8080/api/rules', { method: 'POST', body: JSON.stringify(rulePayload) });
    // ... reset UI
};
// Menu toggle state ditambahkan bersama UI Form dinamis dengan kolom (Description, Regex, Tags, Impact)
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** UI ditambahkan fungsi *Toggle* `showRuleForm`. Saat terbuka, pengguna dapat mengisi parameter reguler expression *(Regex)*, deskripsi, skor bahaya *(Impact)*, dan *Tags*. Saat "*Save*", *React* mengirimkan objek JSON. *Backend* (yang sudah lebih dulu mendukung integrasi ini) akan otomatis memberikan ID baru (contoh `1004`) dan lalu menulisnya ke `waf-rules.json`.
*   **Lokasi Kode:** [config/page.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/config/page.tsx)

**8. Status**
*   **Status Perubahan:** Selesai.

## 14. Keyword Filter & URL Whitelist Feature
**1. Informasi Umum**
*   **Judul Perubahan:** Implementation of Keyword Blocking and URL Exceptions
*   **Tanggal:** 2026-03-30
*   **Requester:** User (Defense-in-Depth & Usability)
*   **Agent / Modul Terkait:** Go WAF Middleware & Configuration UI
*   **File Tujuan Dokumentasi:** [main.go](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/main.go), [page.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/config/page.tsx)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Penggunaan Regex untuk mendeteksi *payload* sangat presisi namun rumit. Selain itu, ada *request* aman (False Positives) yang tidak terhindarkan dari sistem blokir reguler.
*   **Dampak Jika Tidak Diubah:** WAF akan kurang fleksibel dalam menghadapi pola *threat* baru yang simpel, dan sulit untuk "mengizinkan" sebuah URL API masuk jika kebetulan mengandung pola Regex.

**3. Ruang Lingkup Perubahan**
*   Membuat struktur In-Memory baru `FilterConfig` pada Go.
*   Membuat endpoint `GET/POST /api/filters`.
*   Menambahkan Tab **Keywords & Exceptions** di Dasboard Frontend.
*   Menyelipkan logika *Middleware* untuk mengevaluasi *URL Whitelist* di puncak *(Bypass Filter)*, lalu melakukan *Keyword Checking* persis sebelum evaluasi logis Regex untuk mengamankan CPU Server.

**4. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** 
    1.  Tab baru dibuat di UI React.
    2.  `WAFMiddleware` sekarang memiliki dua lapisan baru:
        *   **Check 1:** `strings.HasPrefix(url, whitelistedUrl)` -> Jika *match*, loloskan!
        *   **Check 2:** `strings.Contains(payload, keyword)` -> Jika mengandung kata kunci *(substring)*, lansung Blokir `403 Forbidden` tanpa melalui regex WAF.
    3.  Tersimpan rapi di `/waf-filters.json`.
*   **Lokasi Kode:** [main.go](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/main.go) dan [ui/config/page.tsx](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/ui/app/config/page.tsx).

**8. Status**
*   **Status Perubahan:** Selesai.

## 15. Enforce Mutual Exclusivity for Keyword and URL Filters
**1. Informasi Umum**
*   **Judul Perubahan:** Mutually Exclusive Status for Keywords & Exceptions
*   **Tanggal:** 2026-03-30
*   **Requester:** User (Bug Report & UX Consistency)
*   **Agent / Modul Terkait:** Go WAF Backend API
*   **File Tujuan Dokumentasi:** [main.go](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/main.go)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** Sebelumnya, seseorang dapat menambahkan `MANGGA` ke daftar *Keyword Blocking* namun secara bersamaan bisa memasukkan kata yang sama ke dalam daftar *URL Whitelist*. Ini memicu kerancuan visual di UI dan kontradiksi logika karena satu obyek memiliki status ganda.
*   **Dampak Jika Tidak Diubah:** Potensi admin melakukan kesalahan konfigurasi secara fatal; admin mungkin mengira sebuah kata kunci berhasil di-*ban*, padahal sebelumnya benda yang mirip telanjur terdaftar di *whitelist*, menembus proteksi sistem.

**3. Ruang Lingkup Perubahan**
*   Modifikasi logika *REST API Endpoint* `POST /api/filters` di *Backend* Go.

**4. Komparasi Kode**

**Sebelum Diubah**
```go
if req.Type == "keyword" {
    if req.Action == "add" {
        filterConfig.Keywords = appendStringIfMissing(filterConfig.Keywords, req.Value)
    } 
// Sama sekali tidak ada evaluasi/pembersihan silang ke daftar lain
```

**Sesudah Diubah**
```go
if req.Type == "keyword" {
    if req.Action == "add" {
        filterConfig.Keywords = appendStringIfMissing(filterConfig.Keywords, req.Value)
        filterConfig.WhitelistedURLs = removeString(filterConfig.WhitelistedURLs, req.Value) // Prevent conflict
    }
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** Sistem kini dirancang *Mutually Exclusive* (Saling Menggugurkan). Jika sebuah string didaftarkan ke status/grup *Keyword* (Blokir), maka sistem secara aktif mencari elemen string itu di dalam grup *Whitelist* (Izinkan), lalu melenyapkannya jika ada. Hal yang sama berlaku sebaliknya.
*   **Lokasi Kode:** [main.go](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/main.go)

**8. Status**
*   **Status Perubahan:** Selesai.

## 16. Advanced WAF Evasion Mitigation (Double URL & Null Byte)
**1. Informasi Umum**
*   **Judul Perubahan:** Anti-Evasion Payload Normalizer
*   **Tanggal:** 2026-03-30
*   **Requester:** User (Advanced Threat Protection & Capstone Theory)
*   **Agent / Modul Terkait:** Go WAF Middleware (`main.go`)
*   **File Tujuan Dokumentasi:** [main.go](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/main.go)

**2. Deskripsi Masalah**
*   **Masalah / Kebutuhan:** *Hacker* yang cerdas dapat menggunakan teknik *Double URL Encoding* (seperti merubah `%` menjadi `%25` secara berulang) atau *Null Byte Injection* (`%00` atau `\x00`) agar *payload* berbahaya gagal terbaca oleh mesin pencegat WAF, namun tetap tereksekusi oleh server aplikasi dibelakangnya.
*   **Dampak Jika Tidak Diubah:** Sistem perlindungan WAF (baik regex maupun *keyword*) bisa di-*bypass* hanya dengan menggandakan *encode* *URL karakter*nya.

**3. Ruang Lingkup Perubahan**
*   Penerapan fungsi iteratif (*Loop*) terhadap *URL Unescape* pada `WAFMiddleware` hingga maksimal 3 kali.
*   Pembersihan karakter jebakan *Null-Byte* `\x00`.

**4. Komparasi Kode**

**Sebelum Diubah**
```go
decodedQuery, err := url.QueryUnescape(r.URL.RawQuery)
if err != nil {
    decodedQuery = r.URL.RawQuery
}
```

**Sesudah Diubah**
```go
decodedQuery := r.URL.RawQuery
// Mitigation: Double/Triple URL Encode Evasion
for i := 0; i < 3; i++ {
    if unescaped, err := url.QueryUnescape(decodedQuery); err == nil && unescaped != decodedQuery {
        decodedQuery = unescaped
    } else {
        break
    }
}
// Mitigation: Null Byte Injection (%00 / \x00)
decodedQuery = strings.ReplaceAll(decodedQuery, "\x00", "")
```

**5. Analisis Perubahan**
*   **Apa yang Berubah Secara Teknis:** *Middleware* kini membaca URL mentah *(Raw)* dan mencoba menguraikannya *(Decode)* berkali-kali (`max=3`) selama string-nya terus berubah. Hal ini akan menelanjangi *Double URL Encodings* hingga berbentuk *plaintext* utuh. Terakhir, *Null byte array* dihapus dari ingatan agar WAF memindai *payload* utuh 100%. Setelah itu *payload* baru diserahkan ke gerbang *Keyword Filter* & *Regex*.
*   **Lokasi Kode:** [main.go](file:///d:/Antigravity-Workspace/web-dummy/WAF-Golang/main.go)

**8. Status**
*   **Status Perubahan:** Selesai.

---
*Dokumentasi ini terus diperbarui seiring perkembangan request dari User.*
