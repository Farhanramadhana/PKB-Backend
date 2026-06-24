-- =============================================================
-- SEED: Data dummy PKB — 10 wajib pajak, 10 kendaraan, dst.
-- Semua password: Pretest@2025 (bcrypt cost 12)
-- Idempotent: ON CONFLICT DO NOTHING
-- =============================================================

-- -------------------------
-- 1. WAJIB PAJAK (10 rows)
--    7 INDIVIDU + 3 BADAN_USAHA
-- -------------------------
INSERT INTO wajib_pajak
    (id, nama, jenis, nik, npwp, nib, alamat, no_telp, email, is_active)
VALUES
    -- INDIVIDU
    ('a0000000-0000-0000-0000-000000000001',
     'Budi Santoso', 'INDIVIDU', '3171050415850001', NULL, NULL,
     'Jl. Kebon Sirih No. 12, Menteng, Jakarta Pusat',
     '081234567801', 'budi.santoso@email.com', TRUE),

    ('a0000000-0000-0000-0000-000000000002',
     'Siti Rahayu', 'INDIVIDU', '3273106203900002', NULL, NULL,
     'Jl. Pasteur No. 45, Sukajadi, Bandung',
     '082345678902', 'siti.rahayu@email.com', TRUE),

    ('a0000000-0000-0000-0000-000000000003',
     'Ahmad Fauzi', 'INDIVIDU', '3578011501880003', NULL, NULL,
     'Jl. Raya Darmo No. 88, Wonokromo, Surabaya',
     '083456789003', 'ahmad.fauzi@email.com', TRUE),

    ('a0000000-0000-0000-0000-000000000004',
     'Dewi Permatasari', 'INDIVIDU', '3404126208920004', NULL, NULL,
     'Jl. Malioboro No. 15, Gedongtengen, Yogyakarta',
     '084567890104', 'dewi.permata@email.com', TRUE),

    ('a0000000-0000-0000-0000-000000000005',
     'Rizky Pratama', 'INDIVIDU', '3603010295950005', NULL, NULL,
     'Jl. Sudirman No. 5, Tangerang Kota',
     '085678901205', 'rizky.pratama@email.com', TRUE),

    ('a0000000-0000-0000-0000-000000000006',
     'Fitria Handayani', 'INDIVIDU', '1271085708880006', NULL, NULL,
     'Jl. Gatot Subroto No. 77, Medan Barat, Medan',
     '086789012306', 'fitria.handayani@email.com', TRUE),

    ('a0000000-0000-0000-0000-000000000007',
     'Hendra Kusuma', 'INDIVIDU', '7371230185870007', NULL, NULL,
     'Jl. Panakkukang No. 32, Makassar',
     '087890123407', 'hendra.kusuma@email.com', TRUE),

    -- BADAN USAHA
    ('a0000000-0000-0000-0000-000000000008',
     'PT Maju Bersama Tbk', 'BADAN_USAHA', NULL, '012345678901230', 'NIB20240800001',
     'Jl. Gatot Subroto Kav. 51, Jakarta Selatan',
     '02112345678', 'info@majubersama.co.id', TRUE),

    ('a0000000-0000-0000-0000-000000000009',
     'CV Sumber Rezeki', 'BADAN_USAHA', NULL, '023456789012340', 'NIB20240900002',
     'Jl. Raya Gubeng No. 18, Gubeng, Surabaya',
     '03112345679', 'admin@sumberrezeki.com', TRUE),

    ('a0000000-0000-0000-0000-000000000010',
     'PT Trans Nusantara Logistik', 'BADAN_USAHA', NULL, '034567890123450', 'NIB20241000003',
     'Jl. Soekarno-Hatta No. 100, Bandung',
     '02212345680', 'ops@transnusantara.id', TRUE)

ON CONFLICT (id) DO NOTHING;

-- -------------------------
-- 2. USERS (10 rows)
--    2 PETUGAS + 8 WAJIB_PAJAK
--    Password semua: Pretest@2025
-- -------------------------
INSERT INTO users
    (id, username, password_hash, role, wajib_pajak_id)
VALUES
    ('b0000000-0000-0000-0000-000000000001',
     'petugas1',
     '$2a$12$ilFIA9G0SnjRt16fPPAaJ.BBg.SkHJ3JXJFQyrITH7d/jR9jbgEtu',
     'PETUGAS', NULL),

    ('b0000000-0000-0000-0000-000000000002',
     'petugas2',
     '$2a$12$ilFIA9G0SnjRt16fPPAaJ.BBg.SkHJ3JXJFQyrITH7d/jR9jbgEtu',
     'PETUGAS', NULL),

    ('b0000000-0000-0000-0000-000000000003',
     'wp_budi',
     '$2a$12$ilFIA9G0SnjRt16fPPAaJ.BBg.SkHJ3JXJFQyrITH7d/jR9jbgEtu',
     'WAJIB_PAJAK', 'a0000000-0000-0000-0000-000000000001'),

    ('b0000000-0000-0000-0000-000000000004',
     'wp_siti',
     '$2a$12$ilFIA9G0SnjRt16fPPAaJ.BBg.SkHJ3JXJFQyrITH7d/jR9jbgEtu',
     'WAJIB_PAJAK', 'a0000000-0000-0000-0000-000000000002'),

    ('b0000000-0000-0000-0000-000000000005',
     'wp_ahmad',
     '$2a$12$ilFIA9G0SnjRt16fPPAaJ.BBg.SkHJ3JXJFQyrITH7d/jR9jbgEtu',
     'WAJIB_PAJAK', 'a0000000-0000-0000-0000-000000000003'),

    ('b0000000-0000-0000-0000-000000000006',
     'wp_dewi',
     '$2a$12$ilFIA9G0SnjRt16fPPAaJ.BBg.SkHJ3JXJFQyrITH7d/jR9jbgEtu',
     'WAJIB_PAJAK', 'a0000000-0000-0000-0000-000000000004'),

    ('b0000000-0000-0000-0000-000000000007',
     'wp_rizky',
     '$2a$12$ilFIA9G0SnjRt16fPPAaJ.BBg.SkHJ3JXJFQyrITH7d/jR9jbgEtu',
     'WAJIB_PAJAK', 'a0000000-0000-0000-0000-000000000005'),

    ('b0000000-0000-0000-0000-000000000008',
     'wp_fitria',
     '$2a$12$ilFIA9G0SnjRt16fPPAaJ.BBg.SkHJ3JXJFQyrITH7d/jR9jbgEtu',
     'WAJIB_PAJAK', 'a0000000-0000-0000-0000-000000000006'),

    ('b0000000-0000-0000-0000-000000000009',
     'wp_hendra',
     '$2a$12$ilFIA9G0SnjRt16fPPAaJ.BBg.SkHJ3JXJFQyrITH7d/jR9jbgEtu',
     'WAJIB_PAJAK', 'a0000000-0000-0000-0000-000000000007'),

    ('b0000000-0000-0000-0000-000000000010',
     'wp_majubersama',
     '$2a$12$ilFIA9G0SnjRt16fPPAaJ.BBg.SkHJ3JXJFQyrITH7d/jR9jbgEtu',
     'WAJIB_PAJAK', 'a0000000-0000-0000-0000-000000000008')

ON CONFLICT (id) DO NOTHING;

-- -------------------------
-- 3. KENDARAAN (10 rows)
--    Kendaraan populer Indonesia
-- -------------------------
INSERT INTO kendaraan
    (id, wajib_pajak_id, nomor_polisi, merk, model, tahun,
     jenis_kendaraan, bpkb, stnk, nilai_jual)
VALUES
    -- k01: Toyota Avanza milik Budi (wp01)
    ('c0000000-0000-0000-0000-000000000001',
     'a0000000-0000-0000-0000-000000000001',
     'B 1234 ABC', 'Toyota', 'Avanza 1.3 G', 2020,
     'Mobil Penumpang', 'BPKB-B-00001-2020', 'STNK-B-00001-2020',
     165000000.00),

    -- k02: Honda Beat milik Siti (wp02)
    ('c0000000-0000-0000-0000-000000000002',
     'a0000000-0000-0000-0000-000000000002',
     'D 5678 DEF', 'Honda', 'Beat 110 CBS ISS', 2021,
     'Sepeda Motor', 'BPKB-D-00002-2021', 'STNK-D-00002-2021',
     18500000.00),

    -- k03: Mitsubishi Xpander milik Ahmad (wp03)
    ('c0000000-0000-0000-0000-000000000003',
     'a0000000-0000-0000-0000-000000000003',
     'L 2345 GHI', 'Mitsubishi', 'Xpander Ultimate AT', 2019,
     'Mobil Penumpang', 'BPKB-L-00003-2019', 'STNK-L-00003-2019',
     215000000.00),

    -- k04: Yamaha NMAX milik Dewi (wp04)
    ('c0000000-0000-0000-0000-000000000004',
     'a0000000-0000-0000-0000-000000000004',
     'AB 6789 JKL', 'Yamaha', 'NMAX 155 Connected', 2022,
     'Sepeda Motor', 'BPKB-AB-00004-2022', 'STNK-AB-00004-2022',
     30500000.00),

    -- k05: Suzuki Ertiga milik Rizky (wp05)
    ('c0000000-0000-0000-0000-000000000005',
     'a0000000-0000-0000-0000-000000000005',
     'A 3456 MNO', 'Suzuki', 'Ertiga GX MT', 2018,
     'Mobil Penumpang', 'BPKB-A-00005-2018', 'STNK-A-00005-2018',
     155000000.00),

    -- k06: Honda Vario milik Fitria (wp06)
    ('c0000000-0000-0000-0000-000000000006',
     'a0000000-0000-0000-0000-000000000006',
     'BK 7890 PQR', 'Honda', 'Vario 160 CBS ISS', 2023,
     'Sepeda Motor', 'BPKB-BK-00006-2023', 'STNK-BK-00006-2023',
     24000000.00),

    -- k07: Toyota Rush milik Hendra (wp07)
    ('c0000000-0000-0000-0000-000000000007',
     'a0000000-0000-0000-0000-000000000007',
     'DD 4567 STU', 'Toyota', 'Rush TRD Sportivo AT', 2020,
     'Mobil Penumpang', 'BPKB-DD-00007-2020', 'STNK-DD-00007-2020',
     200000000.00),

    -- k08: Daihatsu Xenia milik PT Maju Bersama (wp08)
    ('c0000000-0000-0000-0000-000000000008',
     'a0000000-0000-0000-0000-000000000008',
     'B 8901 VWX', 'Daihatsu', 'Xenia 1.3 R Deluxe AT', 2021,
     'Mobil Penumpang', 'BPKB-B-00008-2021', 'STNK-B-00008-2021',
     175000000.00),

    -- k09: Isuzu Elf milik CV Sumber Rezeki (wp09) — armada angkutan
    ('c0000000-0000-0000-0000-000000000009',
     'a0000000-0000-0000-0000-000000000009',
     'L 5678 YZA', 'Isuzu', 'Elf NLR 55 BL Microbus', 2019,
     'Mobil Penumpang', 'BPKB-L-00009-2019', 'STNK-L-00009-2019',
     290000000.00),

    -- k10: Hino Ranger milik PT Trans Nusantara (wp10) — truk logistik
    ('c0000000-0000-0000-0000-000000000010',
     'a0000000-0000-0000-0000-000000000010',
     'D 9012 BCD', 'Hino', 'Ranger FM 260 JD', 2018,
     'Mobil Barang', 'BPKB-D-00010-2018', 'STNK-D-00010-2018',
     380000000.00)

ON CONFLICT (id) DO NOTHING;

-- -------------------------
-- 4. KEWAJIBAN PAJAK (10 rows)
--    Tahun pajak 2025
--    Pokok pajak ≈ 1,5–2% NJKB per tahun
-- -------------------------
INSERT INTO kewajiban_pajak
    (id, kendaraan_id, wajib_pajak_id, tahun_pajak,
     periode_awal, periode_final, pokok_pajak, status, total_dibayar)
VALUES
    -- kw01: Avanza Budi → LUNAS (bayar tepat waktu)
    ('d0000000-0000-0000-0000-000000000001',
     'c0000000-0000-0000-0000-000000000001',
     'a0000000-0000-0000-0000-000000000001',
     2025, '2025-01-01', '2025-12-31',
     2500000.00, 'LUNAS', 2500000.00),

    -- kw02: Beat Siti → LUNAS (bayar terlambat, kena denda)
    ('d0000000-0000-0000-0000-000000000002',
     'c0000000-0000-0000-0000-000000000002',
     'a0000000-0000-0000-0000-000000000002',
     2025, '2025-01-01', '2025-12-31',
     500000.00, 'LUNAS', 500000.00),

    -- kw03: Xpander Ahmad → KURANG_BAYAR (bayar terlambat + kurang, denda gabungan)
    ('d0000000-0000-0000-0000-000000000003',
     'c0000000-0000-0000-0000-000000000003',
     'a0000000-0000-0000-0000-000000000003',
     2025, '2025-01-01', '2025-12-31',
     3000000.00, 'KURANG_BAYAR', 2000000.00),

    -- kw04: NMAX Dewi → LUNAS (bayar tepat waktu)
    ('d0000000-0000-0000-0000-000000000004',
     'c0000000-0000-0000-0000-000000000004',
     'a0000000-0000-0000-0000-000000000004',
     2025, '2025-01-01', '2025-12-31',
     750000.00, 'LUNAS', 750000.00),

    -- kw05: Ertiga Rizky → BELUM_BAYAR (belum ada pembayaran)
    ('d0000000-0000-0000-0000-000000000005',
     'c0000000-0000-0000-0000-000000000005',
     'a0000000-0000-0000-0000-000000000005',
     2025, '2025-01-01', '2025-12-31',
     2000000.00, 'BELUM_BAYAR', 0.00),

    -- kw06: Vario Fitria → LUNAS (bayar tepat waktu)
    ('d0000000-0000-0000-0000-000000000006',
     'c0000000-0000-0000-0000-000000000006',
     'a0000000-0000-0000-0000-000000000006',
     2025, '2025-01-01', '2025-12-31',
     600000.00, 'LUNAS', 600000.00),

    -- kw07: Rush Hendra → BELUM_BAYAR
    ('d0000000-0000-0000-0000-000000000007',
     'c0000000-0000-0000-0000-000000000007',
     'a0000000-0000-0000-0000-000000000007',
     2025, '2025-01-01', '2025-12-31',
     2750000.00, 'BELUM_BAYAR', 0.00),

    -- kw08: Xenia PT Maju → LUNAS (bayar terlambat, kena denda telat)
    ('d0000000-0000-0000-0000-000000000008',
     'c0000000-0000-0000-0000-000000000008',
     'a0000000-0000-0000-0000-000000000008',
     2025, '2025-01-01', '2025-12-31',
     2200000.00, 'LUNAS', 2200000.00),

    -- kw09: Isuzu Elf CV Sumber → BELUM_BAYAR
    ('d0000000-0000-0000-0000-000000000009',
     'c0000000-0000-0000-0000-000000000009',
     'a0000000-0000-0000-0000-000000000009',
     2025, '2025-01-01', '2025-12-31',
     4500000.00, 'BELUM_BAYAR', 0.00),

    -- kw10: Hino PT Trans → KURANG_BAYAR (bayar tepat waktu tapi kurang)
    ('d0000000-0000-0000-0000-000000000010',
     'c0000000-0000-0000-0000-000000000010',
     'a0000000-0000-0000-0000-000000000010',
     2025, '2025-01-01', '2025-12-31',
     5000000.00, 'KURANG_BAYAR', 3500000.00)

ON CONFLICT (id) DO NOTHING;

-- -------------------------
-- 5. PEMBAYARAN (7 rows)
--    Dicatat oleh petugas1 dan petugas2
--    3 kewajiban (kw05, kw07, kw09) masih BELUM_BAYAR
-- -------------------------
INSERT INTO pembayaran
    (id, kewajiban_id, wajib_pajak_id, user_id,
     tanggal_bayar, jumlah_bayar, status, catatan_petugas)
VALUES
    -- p01: kw01 Avanza Budi — LUNAS tepat waktu
    ('e0000000-0000-0000-0000-000000000001',
     'd0000000-0000-0000-0000-000000000001',
     'a0000000-0000-0000-0000-000000000001',
     'b0000000-0000-0000-0000-000000000001',
     '2025-03-15', 2500000.00, 'LUNAS',
     'Pembayaran PKB 2025 lunas tepat waktu'),

    -- p02: kw02 Beat Siti — LUNAS tapi terlambat (setelah 2025-12-31)
    ('e0000000-0000-0000-0000-000000000002',
     'd0000000-0000-0000-0000-000000000002',
     'a0000000-0000-0000-0000-000000000002',
     'b0000000-0000-0000-0000-000000000001',
     '2026-01-20', 500000.00, 'LUNAS',
     'Pembayaran terlambat, dikenakan denda telat bayar'),

    -- p03: kw03 Xpander Ahmad — KURANG_BAYAR dan terlambat (denda gabungan)
    ('e0000000-0000-0000-0000-000000000003',
     'd0000000-0000-0000-0000-000000000003',
     'a0000000-0000-0000-0000-000000000003',
     'b0000000-0000-0000-0000-000000000002',
     '2026-02-10', 2000000.00, 'KURANG_BAYAR',
     'Pembayaran sebagian, terlambat. Sisa pokok 1.000.000 belum dilunasi'),

    -- p04: kw04 NMAX Dewi — LUNAS tepat waktu
    ('e0000000-0000-0000-0000-000000000004',
     'd0000000-0000-0000-0000-000000000004',
     'a0000000-0000-0000-0000-000000000004',
     'b0000000-0000-0000-0000-000000000002',
     '2025-06-01', 750000.00, 'LUNAS',
     'PKB 2025 lunas'),

    -- p05: kw06 Vario Fitria — LUNAS tepat waktu
    ('e0000000-0000-0000-0000-000000000005',
     'd0000000-0000-0000-0000-000000000006',
     'a0000000-0000-0000-0000-000000000006',
     'b0000000-0000-0000-0000-000000000001',
     '2025-09-30', 600000.00, 'LUNAS',
     'Pembayaran PKB motor 2025'),

    -- p06: kw08 Xenia PT Maju — LUNAS tapi terlambat
    ('e0000000-0000-0000-0000-000000000006',
     'd0000000-0000-0000-0000-000000000008',
     'a0000000-0000-0000-0000-000000000008',
     'b0000000-0000-0000-0000-000000000002',
     '2026-01-05', 2200000.00, 'LUNAS',
     'Pembayaran terlambat oleh PT Maju Bersama Tbk'),

    -- p07: kw10 Hino PT Trans — KURANG_BAYAR tepat waktu
    ('e0000000-0000-0000-0000-000000000007',
     'd0000000-0000-0000-0000-000000000010',
     'a0000000-0000-0000-0000-000000000010',
     'b0000000-0000-0000-0000-000000000001',
     '2025-10-15', 3500000.00, 'KURANG_BAYAR',
     'Pembayaran sebagian. Sisa 1.500.000 akan dilunasi bulan depan')

ON CONFLICT (id) DO NOTHING;

-- -------------------------
-- 6. DENDA (5 rows)
--    p01,p04,p05 tidak kena denda (tepat waktu & lunas)
--    p02: TELAT_BAYAR saja    (2% × 500.000 = 10.000)
--    p03: GABUNGAN             (telat + kurang)
--         - TELAT_BAYAR: 2% × 3.000.000 = 60.000
--         - KURANG_BAYAR: 1% × 1.000.000 = 10.000
--    p06: TELAT_BAYAR saja    (2% × 2.200.000 = 44.000)
--    p07: KURANG_BAYAR saja   (1% × 1.500.000 = 15.000)
-- -------------------------
INSERT INTO denda
    (id, kewajiban_id, pembayaran_id, jenis, dasar, tarif, jumlah)
VALUES
    -- d01: p02 Beat Siti — denda telat bayar
    ('f0000000-0000-0000-0000-000000000001',
     'd0000000-0000-0000-0000-000000000002',
     'e0000000-0000-0000-0000-000000000002',
     'TELAT_BAYAR', 500000.00, 0.0200, 10000.00),

    -- d02: p03 Xpander Ahmad — denda telat bayar (bagian gabungan)
    ('f0000000-0000-0000-0000-000000000002',
     'd0000000-0000-0000-0000-000000000003',
     'e0000000-0000-0000-0000-000000000003',
     'TELAT_BAYAR', 3000000.00, 0.0200, 60000.00),

    -- d03: p03 Xpander Ahmad — denda kurang bayar (bagian gabungan)
    --      selisih = 3.000.000 - 2.000.000 = 1.000.000
    ('f0000000-0000-0000-0000-000000000003',
     'd0000000-0000-0000-0000-000000000003',
     'e0000000-0000-0000-0000-000000000003',
     'KURANG_BAYAR', 1000000.00, 0.0100, 10000.00),

    -- d04: p06 Xenia PT Maju — denda telat bayar
    ('f0000000-0000-0000-0000-000000000004',
     'd0000000-0000-0000-0000-000000000008',
     'e0000000-0000-0000-0000-000000000006',
     'TELAT_BAYAR', 2200000.00, 0.0200, 44000.00),

    -- d05: p07 Hino PT Trans — denda kurang bayar
    --      selisih = 5.000.000 - 3.500.000 = 1.500.000
    ('f0000000-0000-0000-0000-000000000005',
     'd0000000-0000-0000-0000-000000000010',
     'e0000000-0000-0000-0000-000000000007',
     'KURANG_BAYAR', 1500000.00, 0.0100, 15000.00)

ON CONFLICT (id) DO NOTHING;
