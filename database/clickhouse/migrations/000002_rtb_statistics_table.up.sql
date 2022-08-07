CREATE TABLE rtb_statistics
(
    `StatDate`    Date,
    `FeedId`      UInt16,
    `Country`     LowCardinality(String),
    `Clicks`      UInt64,
    `Impressions` UInt64,
    `Cost`        UInt64,
    `Sign`        Int8
)
    ENGINE = CollapsingMergeTree(Sign)
        PARTITION BY toYYYYMM(StatDate)
        ORDER BY (StatDate, FeedId, Country)