Rank	Fitur	Metode	Cocok?
🥇	Urutan PC (3)	Hungarian Algorithm	Paling cocok. Bukan ML supervised — optimization deterministik, nol data training. Langsung jalan hari pertama, hasil 100% valid selamat. Data kecil tak masalah.
🥈	Pergantian barang (2)	Weighted scoring + (opsional) Logistic Regression	Cocok. Formula bobot tak butuh data sama sekali. Kalau mau "ML", logistic regression butuh cuma ~30-50 label feedback admin — jauh lebih hemat sampel daripada XGBoost.
🥉	Penempatan barang (1)	k-NN / content-based similarity + rule fallback	Semi. Data label = "barang X ditaruh di lab Y" — bisa <50 row. Tree model pasti gagal. k-NN (cari barang serupa sebelumnya) atau rule-based lebih jujur di skala ini.

# Machine Learning — Sistem Inventaris Gudang Labkom

## Ringkasan Arsitektur

| Fitur | Metode | Cold-Start |
|---|---|---|
| Penempatan Barang Baru | LightGBM Classifier (multi-class) | Rule-based scoring |
| Pergantian Barang | Weighted Scoring → XGBRanker | Scoring formula (tanpa training) |
| Urutan PC di Lab | Hungarian Algorithm (linear assignment) | Langsung jalan (optimization) |

**Kenapa bukan Deep Learning?**
Dataset labkom universitas = ratusan–ribuan row. Tree-based methods (LightGBM/XGBoost) dominan di tabular data kecil, interpretable, cepat, tanpa GPU. Neural network butuh puluhan ribu row minimum.

**Library Python:**
```
pandas
numpy
scikit-learn
lightgbm        # fitur 1
xgboost          # fitur 2 (phase 2)
scipy            # fitur 3
```

---

## Fitur 1: Rekomendasi Penempatan Barang Baru

### Algoritma: LightGBM Classifier (Multi-Class)

Ketika barang baru masuk, sistem merekomendasikan lokasi/lab terbaik berdasarkan histori penempatan, kapasitas, kategori, dan kondisi barang existing.

### Input Features

| Feature | Source | Type |
|---|---|---|
| `item_category` | `item.category_id` | categorical |
| `item_condition` | `asset.asset_condition` | categorical |
| `lab_current_count` | `COUNT(asset) WHERE location_id = lab` | numeric |
| `lab_capacity` | `labkom.capacity` | numeric |
| `lab_fill_ratio` | `lab_current_count / lab_capacity` | numeric |
| `lab_same_category_count` | `COUNT(asset) WHERE location_id = lab AND category = Y` | numeric |
| `lab_avg_condition_score` | avg kondisi numerik per lab | numeric |
| `recent_movement_count` | `COUNT(asset_movement)` 6 bulan terakhir per lab | numeric |
| `lab_damage_rate` | `COUNT(damage_report) / total_asset` per lab | numeric |

### Target

`labkom_id` — lab mana barang ditempatkan (dari histori `asset.location_id`)

### Pipeline

```
SQL query join asset + item + labkom + damage_report
→ pandas DataFrame
→ LabelEncoder untuk categorical
→ train_test_split (stratified by labkom_id)
→ LightGBM multiclass
→ predict_proba → top-3 lab ranked by probability
```

### Metrik Evaluasi

- **Accuracy@3** — top-3 akurasi
- **Macro F1-score**

### Cold-Start (data < 100 row)

```python
def rule_based_placement(item_category, labs):
    # 1. Filter lab yang belum penuh
    candidates = [l for l in labs if l.fill_ratio < 0.9]
    # 2. Prioritas lab dengan kategori serupa terbanyak
    candidates.sort(key=lambda l: l.same_category_count, reverse=True)
    # 3. Tiebreak: lab dengan damage_rate terendah
    candidates.sort(key=lambda l: l.damage_rate)
    return candidates[:3]
```

Transisi: kumpulkan data, kalau row > 100, train model, A/B test lawan rule-based.

---

## Fitur 2: Rekomendasi Pergantian Barang Kondisi Buruk

### Algoritma: Weighted Scoring → XGBRanker

Goal bukan prediksi "kapan rusak" tapi **ranking prioritas** mana yang harus diganti duluan. Fase awal pakai formula scoring, lalu upgrade ke ML setelah ada feedback data.

### Input Features

| Feature | Source | Type |
|---|---|---|
| `asset_age_days` | `NOW() - asset.created_at` | numeric |
| `kondisi_score` | `asset.asset_condition` → numerik (BAIK=3, KURANG_BAIK=1, RUSAK=0) | numeric |
| `total_damage_count` | `COUNT(damage_report) WHERE asset_id = X` | numeric |
| `days_since_last_damage` | `NOW() - MAX(damage_report.created_at)` | numeric |
| `damage_frequency` | `total_damage_count / asset_age_days` | numeric |
| `total_peminjaman` | `COUNT(peminjaman_detail) WHERE asset_id = X` | numeric |
| `usage_intensity` | `total_peminjaman / asset_age_days` | numeric |
| `movement_count` | `COUNT(asset_movement) WHERE asset_id = X` | numeric |
| `is_critical_category` | 1 jika PC/monitor/server, 0 lainnya | binary |
| `repair_count` | jumlah kali status berubah dari RUSAK → BAIK | numeric |

### Target/Label

**Phase 1 — Tidak ada label eksplisit.** Pakai weighted scoring:

```python
def replacement_priority_score(row):
    score = (
        0.30 * (1 - row['kondisi_score'] / 3) +       # kondisi buruk = skor tinggi
        0.25 * min(row['damage_frequency'] * 365, 1) +  # frekuensi rusak per tahun
        0.20 * min(row['asset_age_days'] / 1825, 1) +   # umur (cap 5 tahun)
        0.15 * min(row['usage_intensity'] * 30, 1) +    # intensitas pemakaian
        0.10 * row['is_critical_category']               # kritikal atau bukan
    )
    return score  # 0–1, higher = lebih prioritas diganti
```

**Phase 2 — Setelah kumpul approval data** (admin accept/reject rekomendasi):
- Label dari feedback → train XGBClassifier
- Output probability → rank

### Pipeline (Phase 1 → 2)

```
Phase 1: Weighted scoring → tampilkan ranking → admin approve/reject
Phase 2: Kumpul approval data 3-6 bulan → jadikan label → train XGBClassifier
Phase 3: Replace scoring dengan model, keep scoring sebagai feature tambahan
```

### Metrik Evaluasi

- **NDCG@K** — ranking quality
- **Precision@5** — dari top-5 rekomendasi, berapa yang admin setuju

---

## Fitur 3: Rekomendasi Urutan Penempatan PC di Lab

### Algoritma: Hungarian Algorithm (Linear Assignment)

Ini **optimization problem** (assignment), bukan classification/regression. Buat cost matrix `[n_pc × n_position]`, minimize total cost.

### Cost Function

```python
def cost(pc, position, lab_context):
    c = 0.0
    # PC bermasalah jangan di posisi high-traffic (depan/tengah)
    c += pc.damage_frequency * position.accessibility_score * 0.3

    # PC spek tinggi di posisi yang butuh spek tinggi
    spec_mismatch = abs(pc.spec_score - position.required_spec_score)
    c += spec_mismatch * 0.3

    # PC tua jangan di posisi yang sering dipakai
    c += (pc.age_days / 1825) * position.usage_frequency * 0.2

    # Minimize perpindahan dari posisi sekarang (stability bonus)
    if pc.current_position == position.number:
        c -= 0.2
    return c
```

### Input Features — Per PC

| Feature | Source | Type |
|---|---|---|
| `spec_score` | composite dari `pc_component` (RAM + CPU + storage normalized) | numeric |
| `damage_frequency` | `COUNT(damage_report) / age` | numeric |
| `age_days` | `NOW() - pc.created_at` | numeric |
| `current_position` | `pc.pc_number` | numeric |
| `total_repairs` | dari `damage_report` | numeric |
| `component_health` | avg kondisi semua `pc_component` | numeric |

### Input Features — Per Position

| Feature | Source | Type |
|---|---|---|
| `position_number` | 1..N | numeric |
| `accessibility_score` | depan=1.0, tengah=0.7, belakang=0.4 (configurable) | numeric |
| `usage_frequency` | dari `activity_log` / `peminjaman` per posisi | numeric |
| `required_spec_score` | kebutuhan lab (config manual atau dari jadwal mata kuliah) | numeric |

### Output

Mapping `{pc_id: optimal_position_number}` — full reassignment untuk satu lab.

### Pipeline

```
Query semua PC + components di lab X
→ Hitung spec_score, damage_freq, dll per PC
→ Define position properties (config / hardcode)
→ Build cost matrix [n_pc × n_position]
→ scipy.optimize.linear_sum_assignment(cost_matrix)
→ Output: reassignment plan + delta dari posisi sekarang
```

### Metrik Evaluasi

- Total cost sebelum vs sesudah optimasi
- Jumlah PC yang perlu pindah (fewer = better)
- Post-hoc: damage rate per lab setelah optimasi vs sebelum

### Cold-Start

Tidak ada masalah cold-start. Semua feature dihitung dari data existing (spesifikasi, umur). Langsung jalan hari pertama. Tune bobot cost function manual berdasarkan feedback admin.

---

## Strategi Transisi Universal

```
Day 1    → Launch dengan rule-based / formula / optimization (tanpa training data)
3-6 bln  → Kumpulkan data + feedback admin
6+ bln   → Train ML model, A/B test lawan baseline
Win      → Deploy ML, keep rule-based sebagai fallback permanen
```
