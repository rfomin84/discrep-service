CREATE TABLE detailed_feed_statistics
(
    `StatDate`    DateTime,
    `FeedId`      UInt16,
    `BillingType` LowCardinality(String),
    `Country`     LowCardinality(String),
    `Clicks`      UInt64,
    `Impressions` UInt64,
    `Cost`        UInt64,
    `Sign`        Int8
)
    ENGINE = CollapsingMergeTree(Sign)
        PARTITION BY toYYYYMM(StatDate)
        ORDER BY (StatDate, FeedId, BillingType, Country)
