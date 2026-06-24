-- ENUMs
CREATE TYPE jenis_wajib_pajak  AS ENUM ('INDIVIDU', 'BADAN_USAHA');
CREATE TYPE role_user           AS ENUM ('ADMIN', 'PETUGAS', 'WAJIB_PAJAK');
CREATE TYPE status_kewajiban    AS ENUM ('BELUM_BAYAR', 'LUNAS', 'KURANG_BAYAR', 'LEBIH_BAYAR');
CREATE TYPE status_pembayaran   AS ENUM ('LUNAS', 'KURANG_BAYAR', 'LEBIH_BAYAR');
CREATE TYPE jenis_denda         AS ENUM ('TELAT_BAYAR', 'KURANG_BAYAR');

-- WAJIB PAJAK (must be created before users due to FK)
CREATE TABLE wajib_pajak (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    nama       VARCHAR(255) NOT NULL,
    jenis      jenis_wajib_pajak NOT NULL,
    nik        CHAR(16)     UNIQUE,
    npwp       CHAR(15)     UNIQUE,
    nib        VARCHAR(30)  UNIQUE,
    alamat     TEXT         NOT NULL,
    no_telp    VARCHAR(20),
    email      VARCHAR(255),
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_individu_requires_nik
        CHECK (jenis <> 'INDIVIDU'    OR (nik IS NOT NULL AND npwp IS NULL AND nib IS NULL)),
    CONSTRAINT ck_badan_requires_npwp_nib
        CHECK (jenis <> 'BADAN_USAHA' OR (npwp IS NOT NULL AND nib IS NOT NULL AND nik IS NULL))
);

CREATE INDEX idx_wp_nama     ON wajib_pajak (nama);
CREATE INDEX idx_wp_jenis    ON wajib_pajak (jenis);
CREATE INDEX idx_wp_active   ON wajib_pajak (is_active);
CREATE INDEX idx_wp_nik      ON wajib_pajak (nik)  WHERE nik  IS NOT NULL;
CREATE INDEX idx_wp_npwp     ON wajib_pajak (npwp) WHERE npwp IS NOT NULL;

-- USERS
CREATE TABLE users (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    username        VARCHAR(100) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    role            role_user    NOT NULL,
    wajib_pajak_id  UUID         REFERENCES wajib_pajak(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_wp ON users (wajib_pajak_id) WHERE wajib_pajak_id IS NOT NULL;

-- KENDARAAN
CREATE TABLE kendaraan (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    wajib_pajak_id  UUID         NOT NULL REFERENCES wajib_pajak(id) ON DELETE RESTRICT,
    nomor_polisi    VARCHAR(20)  NOT NULL UNIQUE,
    merk            VARCHAR(100) NOT NULL,
    model           VARCHAR(100) NOT NULL,
    tahun           SMALLINT     NOT NULL CHECK (tahun >= 1900),
    jenis_kendaraan VARCHAR(50)  NOT NULL DEFAULT '',
    bpkb            VARCHAR(50)  NOT NULL DEFAULT '',
    stnk            VARCHAR(50)  NOT NULL DEFAULT '',
    nilai_jual      NUMERIC(15,2) NOT NULL DEFAULT 0 CHECK (nilai_jual >= 0),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_kendaraan_wp ON kendaraan (wajib_pajak_id);

-- KEWAJIBAN PAJAK
CREATE TABLE kewajiban_pajak (
    id              UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    kendaraan_id    UUID           NOT NULL REFERENCES kendaraan(id)    ON DELETE RESTRICT,
    wajib_pajak_id  UUID           NOT NULL REFERENCES wajib_pajak(id)  ON DELETE RESTRICT,
    tahun_pajak     SMALLINT       NOT NULL,
    periode_awal    DATE           NOT NULL,
    periode_final   DATE           NOT NULL,
    pokok_pajak     NUMERIC(15,2)  NOT NULL CHECK (pokok_pajak > 0),
    status          status_kewajiban NOT NULL DEFAULT 'BELUM_BAYAR',
    total_dibayar   NUMERIC(15,2)  NOT NULL DEFAULT 0 CHECK (total_dibayar >= 0),
    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_periode         CHECK (periode_final >= periode_awal),
    CONSTRAINT uq_kendaraan_tahun UNIQUE (kendaraan_id, tahun_pajak)
);

CREATE INDEX idx_kewajiban_wp     ON kewajiban_pajak (wajib_pajak_id);
CREATE INDEX idx_kewajiban_status ON kewajiban_pajak (status);

-- PEMBAYARAN
CREATE TABLE pembayaran (
    id              UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    kewajiban_id    UUID           NOT NULL REFERENCES kewajiban_pajak(id) ON DELETE RESTRICT,
    wajib_pajak_id  UUID           NOT NULL REFERENCES wajib_pajak(id)     ON DELETE RESTRICT,
    user_id         UUID           NOT NULL REFERENCES users(id)            ON DELETE RESTRICT,
    tanggal_bayar   DATE           NOT NULL,
    jumlah_bayar    NUMERIC(15,2)  NOT NULL CHECK (jumlah_bayar > 0),
    status          status_pembayaran NOT NULL,
    catatan_petugas TEXT,
    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pembayaran_kewajiban ON pembayaran (kewajiban_id);
CREATE INDEX idx_pembayaran_tanggal   ON pembayaran (tanggal_bayar);
CREATE INDEX idx_pembayaran_status    ON pembayaran (status);
CREATE INDEX idx_pembayaran_wp        ON pembayaran (wajib_pajak_id);

-- DENDA
CREATE TABLE denda (
    id              UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    kewajiban_id    UUID          NOT NULL REFERENCES kewajiban_pajak(id) ON DELETE RESTRICT,
    pembayaran_id   UUID          NOT NULL REFERENCES pembayaran(id)       ON DELETE RESTRICT,
    jenis           jenis_denda   NOT NULL,
    dasar           NUMERIC(15,2) NOT NULL CHECK (dasar > 0),
    tarif           NUMERIC(5,4)  NOT NULL,
    jumlah          NUMERIC(15,2) NOT NULL CHECK (jumlah > 0),
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_denda_kewajiban  ON denda (kewajiban_id);
CREATE INDEX idx_denda_pembayaran ON denda (pembayaran_id);
