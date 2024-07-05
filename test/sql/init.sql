create table TEST_TAB (
    id NUMBER,
    col_date DATE,
    col_timestamp  TIMESTAMP(6),
    col_integer INTEGER,
    col_float FLOAT,
    col_varchar VARCHAR2(100),
    col_clob CLOB,
    col_raw  RAW(100),
    col_blob CLOB
);

alter table TEST_TAB add constraint TEST_TAB_PK primary key (ID);

insert into test_tab (id, col_date, col_timestamp, col_integer, col_float, col_varchar, col_clob, col_raw, col_blob)
select
    level id,
    sysdate - level/24 dt,
    systimestamp ts,
    dbms_random.random,
    dbms_random.random/7,
    'string:'||dbms_random.value,
    to_clob('clob:' ||dbms_random.value),
    utl_raw.cast_to_raw('raw:' ||dbms_random.value),
    utl_raw.cast_to_raw('blob:' ||dbms_random.value)
from dual
    connect by level < 100001;