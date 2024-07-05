-- ALTER TABLE EVENT_LOG MOVE INITRANS 10;
-- ALTER INDEX EVENT_LOG_IDX REBUILD INITRANS 20;

create table EVENT_LOG
(
  ts       TIMESTAMP(3),
  group_id VARCHAR2(64),
  part_id  NUMBER,
  state    VARCHAR2(1),
  action   VARCHAR2(1),
  key_n    NUMBER
);

create index EVENT_LOG_IDX on ORGON.EVENT_LOG (GROUP_ID, PART_ID, STATE, TS, KEY_N, ACTION);