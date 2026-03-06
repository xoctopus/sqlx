-- MySQL type descriptor demo: NeedWidth / NeedPrecision / NeedWidthPrecision
-- Run in your DB, then query information_schema.columns to see how MySQL describes each type.

CREATE TABLE IF NOT EXISTS type_descriptor_demo (
    c_char_10       CHAR(10),
    c_varchat       VARCHAR,
    c_varchar_255   VARCHAR(255),
    c_binary_8      BINARY(8),
    c_varbinary_16  VARBINARY(16),
    c_bit_8         BIT(8),
    c_bit           BIT,
    c_time_3        TIME(3),
    c_datetime_3    DATETIME(3),
    c_timestamp_3   TIMESTAMP(3),
    c_time          TIME,
    c_datetime      DATETIME,
    c_timestamp     TIMESTAMP,
    c_decimal_10_2  DECIMAL(10, 2),
    c_numeric_10_2  NUMERIC(10, 2),
    c_decimal       DECIMAL,
    c_numeric       NUMERIC
);

SELECT
    column_name,
    data_type,
    column_type,
    character_maximum_length AS char_max_len,
    character_octet_length   AS char_octet_len,
    numeric_precision         AS num_precision,
    numeric_scale             AS num_scale,
    datetime_precision        AS dt_precision
FROM information_schema.columns
WHERE table_schema = DATABASE()
  AND table_name = 'type_descriptor_demo'
ORDER BY ordinal_position;
