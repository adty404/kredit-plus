-- Migrations UP

-- Tabel consumers
CREATE TABLE IF NOT EXISTS consumers (
    id BIGSERIAL PRIMARY KEY,
    nik VARCHAR(16) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255),
    tempat_lahir VARCHAR(100),
    tanggal_lahir DATE,
    gaji DECIMAL(15,2),
    overall_credit_limit DECIMAL(19,2) NOT NULL DEFAULT 0,
    foto_ktp VARCHAR(255),
    foto_selfie VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

-- Tabel consumer_credit_limits
CREATE TABLE IF NOT EXISTS consumer_credit_limits (
    id BIGSERIAL PRIMARY KEY,
    consumer_id BIGINT NOT NULL,
    tenor_months INT NOT NULL,
    credit_limit DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_consumer_limit FOREIGN KEY (consumer_id) REFERENCES consumers(id) ON DELETE CASCADE,
    UNIQUE (consumer_id, tenor_months)
    );

-- Tabel transactions
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    consumer_id BIGINT NOT NULL,
    consumer_credit_limit_id BIGINT NOT NULL,
    nomor_kontrak VARCHAR(50) UNIQUE NOT NULL,
    tanggal_kontrak DATE NOT NULL,
    otr DECIMAL(19,2) NOT NULL,
    uang_muka DECIMAL(19,2) DEFAULT 0,
    admin_fee DECIMAL(19,2) DEFAULT 0,
    pokok_pembiayaan_awal DECIMAL(19,2) NOT NULL,
    nilai_cicilan_per_periode DECIMAL(19,2) NOT NULL,
    tenor_bulan INT NOT NULL,
    total_bunga DECIMAL(19,2) NOT NULL,
    total_kewajiban_pembayaran DECIMAL(19,2) NOT NULL,
    nama_asset VARCHAR(255),
    jenis_asset VARCHAR(50),
    status_kontrak VARCHAR(30) NOT NULL,
    catatan TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_transaction_consumer FOREIGN KEY (consumer_id) REFERENCES consumers(id) ON DELETE RESTRICT,
    CONSTRAINT fk_transaction_consumer_limit FOREIGN KEY (consumer_credit_limit_id) REFERENCES consumer_credit_limits(id) ON DELETE RESTRICT
    );
