export default {
  importantSettingParamNames: {
    BOT_NAME: String,
    SPIDER_MODULES: Array,
    NEWSPIDER_MODULE: String,
    USER_AGENT: String,
    ROBOTSTXT_OBEY: Boolean,
    CONCURRENT_REQUESTS: Number,
    DOWNLOAD_DELAY: Number,
    CONCURRENT_REQUESTS_PER_DOMAIN: Number,
    CONCURRENT_REQUESTS_PER_IP: Number,
    COOKIES_ENABLED: Boolean,
    TELNETCONSOLE_ENABLED: Boolean,
    DEFAULT_REQUEST_HEADERS: Object,
    SPIDER_MIDDLEWARES: Object,
    DOWNLOADER_MIDDLEWARES: Object,
    EXTENSIONS: Object,
    ITEM_PIPELINES: Object,
    AUTOTHROTTLE_ENABLED: Boolean,
    AUTOTHROTTLE_START_DELAY: Number,
    AUTOTHROTTLE_MAX_DELAY: Number,
    AUTOTHROTTLE_TARGET_CONCURRENCY: Number,
    AUTOTHROTTLE_DEBUG: Boolean,
    HTTPCACHE_ENABLED: Boolean,
    HTTPCACHE_EXPIRATION_SECS: Number,
    HTTPCACHE_DIR: String,
    HTTPCACHE_IGNORE_HTTP_CODES: Array,
    HTTPCACHE_STORAGE: String
  },
  settingParamNames: [
    'AWS_ACCESS_KEY_ID',
    'AWS_SECRET_ACCESS_KEY',
    'AWS_ENDPOINT_URL',
    'AWS_USE_SSL',
    'AWS_VERIFY',
    'AWS_REGION_NAME',
    'BOT_NAME',
    'CONCURRENT_ITEMS',
    'CONCURRENT_REQUESTS',
    'CONCURRENT_REQUESTS_PER_DOMAIN',
    'CONCURRENT_REQUESTS_PER_IP',
    'DEFAULT_ITEM_CLASS',
    'DEFAULT_REQUEST_HEADERS',
    'DEPTH_LIMIT',
    'DEPTH_PRIORITY',
    'DEPTH_STATS_VERBOSE',
    'DNSCACHE_ENABLED',
    'DNSCACHE_SIZE',
    'DNS_TIMEOUT',
    'DOWNLOADER',
    'DOWNLOADER_HTTPCLIENTFACTORY',
    'DOWNLOADER_CLIENTCONTEXTFACTORY',
    'DOWNLOADER_CLIENT_TLS_CIPHERS',
    'DOWNLOADER_CLIENT_TLS_METHOD',
    'DOWNLOADER_CLIENT_TLS_VERBOSE_LOGGING',
    'DOWNLOADER_MIDDLEWARES',
    'DOWNLOADER_MIDDLEWARES_BASE',
    'DOWNLOADER_STATS',
    'DOWNLOAD_DELAY',
    'DOWNLOAD_HANDLERS',
    'DOWNLOAD_HANDLERS_BASE',
    'DOWNLOAD_TIMEOUT',
    'DOWNLOAD_MAXSIZE',
    'DOWNLOAD_WARNSIZE',
    'DOWNLOAD_FAIL_ON_DATALOSS',
    'DUPEFILTER_CLASS',
    'DUPEFILTER_DEBUG',
    'EDITOR',
    'EXTENSIONS',
    'EXTENSIONS_BASE',
    'FEED_TEMPDIR',
    'FTP_PASSIVE_MODE',
    'FTP_PASSWORD',
    'FTP_USER',
    'ITEM_PIPELINES',
    'ITEM_PIPELINES_BASE',
    'LOG_ENABLED',
    'LOG_ENCODING',
    'LOG_FILE',
    'LOG_FORMAT',
    'LOG_DATEFORMAT',
    'LOG_FORMATTER',
    'LOG_LEVEL',
    'LOG_STDOUT',
    'LOG_SHORT_NAMES',
    'LOGSTATS_INTERVAL',
    'MEMDEBUG_ENABLED',
    'MEMDEBUG_NOTIFY',
    'MEMUSAGE_ENABLED',
    'MEMUSAGE_LIMIT_MB',
    'MEMUSAGE_CHECK_INTERVAL_SECONDS',
    'MEMUSAGE_NOTIFY_MAIL',
    'MEMUSAGE_WARNING_MB',
    'NEWSPIDER_MODULE',
    'RANDOMIZE_DOWNLOAD_DELAY',
    'REACTOR_THREADPOOL_MAXSIZE',
    'REDIRECT_MAX_TIMES',
    'REDIRECT_PRIORITY_ADJUST',
    'RETRY_PRIORITY_ADJUST',
    'ROBOTSTXT_OBEY',
    'ROBOTSTXT_PARSER',
    'SCHEDULER',
    'SCHEDULER_DEBUG',
    'SCHEDULER_DISK_QUEUE',
    'SCHEDULER_MEMORY_QUEUE',
    'SCHEDULER_PRIORITY_QUEUE',
    'SPIDER_CONTRACTS',
    'SPIDER_CONTRACTS_BASE',
    'SPIDER_LOADER_CLASS',
    'SPIDER_LOADER_WARN_ONLY',
    'SPIDER_MIDDLEWARES',
    'SPIDER_MIDDLEWARES_BASE',
    'SPIDER_MODULES',
    'STATS_CLASS',
    'STATS_DUMP',
    'STATSMAILER_RCPTS',
    'TELNETCONSOLE_ENABLED',
    'TELNETCONSOLE_PORT',
    'TEMPLATES_DIR',
    'URLLENGTH_LIMIT',
    'USER_AGENT',
    'AJAXCRAWL_ENABLED',
    'AUTOTHROTTLE_DEBUG',
    'AUTOTHROTTLE_ENABLED',
    'AUTOTHROTTLE_MAX_DELAY',
    'AUTOTHROTTLE_START_DELAY',
    'AUTOTHROTTLE_TARGET_CONCURRENCY',
    'CLOSESPIDER_ERRORCOUNT',
    'CLOSESPIDER_ITEMCOUNT',
    'CLOSESPIDER_PAGECOUNT',
    'CLOSESPIDER_TIMEOUT',
    'COMMANDS_MODULE',
    'COMPRESSION_ENABLED',
    'COOKIES_DEBUG',
    'COOKIES_ENABLED',
    'FEED_EXPORTERS',
    'FEED_EXPORTERS_BASE',
    'FEED_EXPORT_ENCODING',
    'FEED_EXPORT_FIELDS',
    'FEED_EXPORT_INDENT',
    'FEED_FORMAT',
    'FEED_STORAGES',
    'FEED_STORAGES_BASE',
    'FEED_STORAGE_FTP_ACTIVE',
    'FEED_STORAGE_S3_ACL',
    'FEED_STORE_EMPTY',
    'FEED_URI',
    'FILES_EXPIRES',
    'FILES_RESULT_FIELD',
    'FILES_STORE',
    'FILES_STORE_GCS_ACL',
    'FILES_STORE_S3_ACL',
    'FILES_URLS_FIELD',
    'GCS_PROJECT_ID',
    'HTTPCACHE_ALWAYS_STORE',
    'HTTPCACHE_DBM_MODULE',
    'HTTPCACHE_DIR',
    'HTTPCACHE_ENABLED',
    'HTTPCACHE_EXPIRATION_SECS',
    'HTTPCACHE_GZIP',
    'HTTPCACHE_IGNORE_HTTP_CODES',
    'HTTPCACHE_IGNORE_MISSING',
    'HTTPCACHE_IGNORE_RESPONSE_CACHE_CONTROLS',
    'HTTPCACHE_IGNORE_SCHEMES',
    'HTTPCACHE_POLICY',
    'HTTPCACHE_STORAGE',
    'HTTPERROR_ALLOWED_CODES',
    'HTTPERROR_ALLOW_ALL',
    'HTTPPROXY_AUTH_ENCODING',
    'HTTPPROXY_ENABLED',
    'IMAGES_EXPIRES',
    'IMAGES_MIN_HEIGHT',
    'IMAGES_MIN_WIDTH',
    'IMAGES_RESULT_FIELD',
    'IMAGES_STORE',
    'IMAGES_STORE_GCS_ACL',
    'IMAGES_STORE_S3_ACL',
    'IMAGES_THUMBS',
    'IMAGES_URLS_FIELD',
    'MAIL_FROM',
    'MAIL_HOST',
    'MAIL_PASS',
    'MAIL_PORT',
    'MAIL_SSL',
    'MAIL_TLS',
    'MAIL_USER',
    'MEDIA_ALLOW_REDIRECTS',
    'METAREFRESH_ENABLED',
    'METAREFRESH_IGNORE_TAGS',
    'METAREFRESH_MAXDELAY',
    'REDIRECT_ENABLED',
    'REDIRECT_MAX_TIMES',
    'REFERER_ENABLED',
    'REFERRER_POLICY',
    'RETRY_ENABLED',
    'RETRY_HTTP_CODES',
    'RETRY_TIMES',
    'TELNETCONSOLE_HOST',
    'TELNETCONSOLE_PASSWORD',
    'TELNETCONSOLE_PORT',
    'TELNETCONSOLE_USERNAME',
    'REDIS_ITEMS_KEY',
    'REDIS_ITEMS_SERIALIZER',
    'REDIS_HOST',
    'REDIS_PORT',
    'REDIS_URL',
    'REDIS_PARAMS',
    'REDIS_START_URLS_AS_SET',
    'REDIS_START_URLS_KEY',
    'REDIS_ENCODING'
  ]
};
