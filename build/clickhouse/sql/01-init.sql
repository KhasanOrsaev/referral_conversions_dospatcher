create database if not exists cs;

create table if not exists cs.events
(
    event_id UUID,
    user_id Int64,
    application_id UUID,
    application_id_native String,
    offer_id String,
    offer_source String,
    offer_status String,
    offer_payout String,
    offer_link String,
    process_id int,
    source String,
    type String,
    dt_create DateTime default now(),
    dt_event DateTime,
    session_logger_id UUID,
    data_json String
) engine = ReplacingMergeTree(dt_create)
      partition by (toYYYYMM(dt_create), source)
      order by (user_id, application_id, source, type, dt_event)
      sample by (user_id);
