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
    site_name       TEXT NOT NULL,
    show_name       TEXT,
    `order`         INTEGER NOT NULL,
    host            TEXT,
    cookie          TEXT,
    api_token       TEXT,
    passkey         TEXT,
    rss_key         TEXT,
    user_agent      TEXT,
    custom_header   TEXT,
    proxy           BOOLEAN NOT NULL,
    is_active       BOOLEAN NOT NULL,
    timeout         INTEGER NOT NULL,
    create_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- 创建索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_seed_sync_site_name ON seed_sync_site(site_name);
CREATE INDEX IF NOT EXISTS idx_seed_sync_site_order ON seed_sync_site(`order`);


-- 站点流控表
CREATE TABLE IF NOT EXISTS seed_sync_site_flow_control
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    site_name       TEXT NOT NULL,
    max_per_min     INTEGER NOT NULL,
    max_per_hour    INTEGER NOT NULL,
    max_per_day     INTEGER NOT NULL,
    create_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- 创建索引
CREATE INDEX IF NOT EXISTS idx_seed_sync_site_flow_control_site_name ON seed_sync_site_flow_control(site_name);


-- 定时任务表
CREATE TABLE IF NOT EXISTS seed_sync_schedule_task
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    task_name       TEXT NOT NULL,
    cron            TEXT NOT NULL,
    execute_content TEXT NOT NULL,
    execute_status  TEXT NOT NULL,
    last_execute_time TIMESTAMP,
    next_execute_time TIMESTAMP,
    last_execute_result TEXT,
    active            BOOLEAN NOT NULL DEFAULT TRUE,
    create_user       TEXT,
    create_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- 创建索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_seed_sync_schedule_task_task_name ON seed_sync_schedule_task(task_name);

-- 下载器表
CREATE TABLE IF NOT EXISTS seed_sync_downloader
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL,
    url             TEXT NOT NULL,
    username        TEXT,
    password        TEXT,
    type            TEXT NOT NULL,
    create_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- 创建索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_seed_sync_downloader_name ON seed_sync_downloader(name);

-- 辅种任务表
CREATE TABLE IF NOT EXISTS seed_sync_seed_task
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    task_name       TEXT NOT NULL,
    site_list       TEXT,
    downloader_id   INTEGER NOT NULL,
    exclude_path    TEXT,
    exclude_tag     TEXT,
    min_size        INTEGER,
    add_tag         TEXT,
    status          TEXT NOT NULL,
    create_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)