-- 系统参数表
CREATE TABLE IF NOT EXISTS seed_sync_system_param
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    `key`       TEXT NOT NULL,
    value       TEXT NOT NULL,
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_seed_sync_system_param_key ON seed_sync_system_param(`key`);


-- 用户表
CREATE TABLE IF NOT EXISTS user
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    username        TEXT NOT NULL,
    password        TEXT NOT NULL,
    token           TEXT,
    status          TEXT,
    is_two_factor   BOOLEAN NOT NULL DEFAULT FALSE,
    two_factor_type TEXT,
    create_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- 站点表
CREATE TABLE IF NOT EXISTS seed_sync_site
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL,
    `order`         INTEGER NOT NULL,
    url             TEXT,
    cookie          TEXT,
    api_key         TEXT,
    token           TEXT,
    custom_header   TEXT,
    passkey         TEXT,
    rss             TEXT,
    domains         TEXT,
    download_url    TEXT,
    torrent_list_url TEXT,
    proxy           BOOLEAN NOT NULL,
    timeout         INTEGER,
    is_override     BOOLEAN NOT NULL,
    is_active       BOOLEAN NOT NULL,
    create_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_seed_sync_site_name ON seed_sync_site(name);
CREATE INDEX IF NOT EXISTS idx_seed_sync_site_order ON seed_sync_site(`order`);