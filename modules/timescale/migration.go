package timescale

func migrateTimescale() {
	log.Info("Migrating timescale db")

	// #nosec G104
	ConnectionPool.Exec(Context, `CREATE TABLE IF NOT EXISTS investment_exporter (
		time TIMESTAMPTZ NOT NULL,
		total_btc_on_kraken DOUBLE PRECISION,
		total_cache_to_kraken DOUBLE PRECISION,
		total_kraken_fees_lost DOUBLE PRECISION,
		eur_on_kraken DOUBLE PRECISION,
		btc_price_eur DOUBLE PRECISION,
		btc_price_usd DOUBLE PRECISION,
		btc_in_wallet DOUBLE PRECISION,
		eur_in_wallet DOUBLE PRECISION,
		pending_fiat DOUBLE PRECISION,
		total_scrape_time DOUBLE PRECISION,
		next_dca_order_time TIMESTAMPTZ
	);`)

	// #nosec G104
	ConnectionPool.Exec(Context, "SELECT create_hypertable ('investment_exporter', by_range ('time', INTERVAL '1 day'), if_not_exists => TRUE);")

	// #nosec G104
	ConnectionPool.Exec(Context, "ALTER TABLE investment_exporter SET (timescaledb.compress, timescaledb.compress_orderby = 'time DESC');")

	// #nosec G104
	ConnectionPool.Exec(Context, "SELECT add_compression_policy('investment_exporter', compress_after => INTERVAL '60d');")

	// #nosec G104
	ConnectionPool.Exec(Context, `CREATE TABLE IF NOT EXISTS purchases (
		refid TEXT PRIMARY KEY,
		time TIMESTAMPTZ NOT NULL,
		amount DOUBLE PRECISION,
		fee DOUBLE PRECISION
	);`)

	// #nosec G104
	ConnectionPool.Exec(Context, "SELECT create_hypertable ('purchases', by_range ('time', INTERVAL '1 month'), if_not_exists => TRUE);")

	// #nosec G104
	ConnectionPool.Exec(Context, `CREATE TABLE IF NOT EXISTS utxo_balances (
		address     TEXT                NOT NULL,
		btc         DOUBLE PRECISION    NOT NULL,
		PRIMARY KEY (address),
		UNIQUE (address)
	);`)

	log.Info("Timescale migration done")
}
