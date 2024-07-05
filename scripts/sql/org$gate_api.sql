create or replace package orgon.org$gate_api is

-- Author  : VODYANKO
-- Created : 17/04/2024 18:43:22
-- Purpose : The package implements the consumer to the transactional outbox and the data enrichment

/*
OVERVIEW

The package allows you to access the outbox via the API and receive the next events 
for the requested table and a specific part (bucket number).

If the events exist, the data is retrieved from the corresponding table (view or query join).
The data is encoded in XML format using the high-performance Oracle dbms_xmlgen core package (written in C).
For efficient transmission over the network, data is also compressed using the gzip algorithm.

*/

/*
TEST

-- Get the next (not processed yet) events as XML text (clob)
declare
  upd_cnt int;
  upd_xml clob;
  del_cnt int;
  del_xml clob;
begin
  org$gate_api.getNextEvents(p_group_id       => 'test_tab',
                             p_part_id        => 1,
                             p_rows           => 100,
                             p_qry_columns    => q'[*]',
                             p_qry_from       => 'select t.* from test_tab t',
                             p_qry_pk_column  => 'id',
                             r_upd_rows_dump  => upd_xml,
                             r_upd_rows_count => upd_cnt,
                             r_del_rows_dump  => del_xml,
                             r_del_rows_count => del_cnt       
                             );

  dbms_output.put_line('upd_xml=' || substr(upd_xml, 1, 1000));
  dbms_output.put_line('upd_cnt=' || upd_cnt);
  
  dbms_output.put_line('del_xml=' || substr(del_xml, 1, 1000));
  dbms_output.put_line('del_cnt=' || del_cnt);
end;

-- Get the next (not processed yet) events as compressed XML (blob)
declare
  upd_cnt int;
  upd_xml blob;
  del_cnt int;
  del_xml blob;
begin
  org$gate_api.getNextEvents(p_group_id       => 'test_tab',
                             p_part_id        => 2,
                             p_rows           => 1000,
                             p_qry_columns    => '*',
                             p_qry_from       => 'select t.* from test_tab t',
                             p_qry_pk_column  => 'id',
                             r_upd_rows_dump  => upd_xml,
                             r_upd_rows_count => upd_cnt,
                             r_del_rows_dump  => del_xml,
                             r_del_rows_count => del_cnt       
                             );

  dbms_output.put_line('upd_xml_len=' || length(upd_xml));
  dbms_output.put_line('upd_cnt=' || upd_cnt);
  
  dbms_output.put_line('del_xml_len=' || length(del_xml));
  dbms_output.put_line('del_cnt=' || del_cnt);
end;

*/

-- Get the next new events serialized in XML: symbolic representation
procedure getNextEvents(
  p_group_id in varchar2
, p_part_id in number
, p_rows in number

, p_qry_columns in varchar2
, p_qry_from in varchar2
, p_qry_pk_column in varchar2

, r_upd_rows_dump out nocopy clob
, r_upd_rows_count out number
, r_del_rows_dump out nocopy clob
, r_del_rows_count out number
);

-- Get the next new events serialized in XML: binary gzip representation UTF8 of XML-text
procedure getNextEvents(
  p_group_id in varchar2
, p_part_id in number
, p_rows in number

, p_qry_columns in varchar2
, p_qry_from in varchar2
, p_qry_pk_column in varchar2

, r_upd_rows_dump out nocopy blob
, r_upd_rows_count out number
, r_del_rows_dump out nocopy blob
, r_del_rows_count out number
);

end org$gate_api;
/

create or replace package body orgon.org$gate_api is

function toUnixTimestamp(
  p_ts in timestamp
, p_tz in varchar2 default DBTIMEZONE
) return number deterministic
is
  tsUTC timestamp;
  epoche number;
  intval INTERVAL DAY(9) TO SECOND;
begin
  tsUTC := FROM_TZ(p_ts, p_tz) AT TIME ZONE '00:00';
  intval := TO_DSINTERVAL(tsUTC - TIMESTAMP '1970-01-01 00:00:00');
  epoche := 
    EXTRACT(DAY FROM intval)*24*60*60 + 
    EXTRACT(HOUR FROM intval)*60*60 + 
    EXTRACT(MINUTE FROM intval)*60 + 
    EXTRACT(SECOND FROM intval);

  return round(epoche, 3) * 1000;
end; /* toUnixTimestamp */

function makeSqlPkInListQuery(
  p_cols in varchar2
, p_from in varchar2
, p_pk_col in varchar2
, p_binds in out nocopy org$xml_factory.TBindParams
) return varchar2
is
  v_select varchar2(4000);
  v_where varchar2(16000);
begin
  v_select := 'select ' || p_cols || ' from ' || p_from;
  
  v_where := ' where ' || p_pk_col || ' in (';
  for i in 1..p_binds.count() loop
    v_where := v_where ||':'|| p_binds(i).Name || ',';
  end loop;    
  v_where := rtrim(v_where, ',') || ')';
  
  return (v_select || v_where);
end; /* makeSqlPkInListQuery */

procedure copmactAndSplitEvents(
  p_events in out nocopy org$outbox_api.TEventArray
, r_upd_events out nocopy org$outbox_api.TEventArray
, r_del_events out nocopy org$outbox_api.TEventArray
)
is
  type TCompList is table of integer index by varchar2(32);
  copmList TCompList;
  pk number;
  uc pls_integer := 0;
  dc pls_integer := 0;
  skipped pls_integer := 0;
begin
  r_upd_events := new org$outbox_api.TEventArray();
  r_del_events := new org$outbox_api.TEventArray();

  -- Compaction: last event wins
  for i in reverse 1..p_events.count() loop
    pk := p_events(i).key;
    if copmList.exists(pk) then
        skipped := skipped + 1;
        continue;
    else
      copmList(pk) := i;
    end if;
  end loop;
  
    -- Split and copy
  for i in 1..p_events.count() loop
    pk := p_events(i).key;
    if copmList(pk) = i then
      if p_events(i).op in (org$outbox_api.ACTION_INSERT, org$outbox_api.ACTION_UPDATE) then
        uc := uc + 1;
        r_upd_events.extend(1);
        r_upd_events(uc) := p_events(i);
      elsif p_events(i).op = org$outbox_api.ACTION_DELETE then
        dc := dc + 1;
        r_del_events.extend(1);
        r_del_events(dc) := p_events(i);
      end if;  
    end if;
  end loop;
end;

procedure dumpUpdatedRows(
  p_qry_columns in varchar2
, p_qry_from in varchar2
, p_qry_pk_column in varchar2
, p_events in out nocopy org$outbox_api.TEventArray
, r_rows_dump out nocopy clob
, r_rows_count out number
)
is
  qry varchar2(32000);
  binds org$xml_factory.TBindParams;
  alias constant varchar2(1) := 'q';

  function metaColumns return varchar2 is
  begin
    return
      /* action type */
      '''' || org$outbox_api.ACTION_UPDATE || '''' || ' "__op"' || ', ' || 
      /* pk column name */
      '''' || p_qry_pk_column || '''' || ' "__pk_name"' || ', ' || 
      /* pk column value */
      p_qry_pk_column || ' "__pk_val"'   || ', ' || 
      /* row state timestamp */
      'systimestamp AT TIME ZONE ''00:00'' "__ts"' || ', ' || 
      /* row state unix ts at UTC TZ */
      org$gate_api.toUnixTimestamp(systimestamp, '00:00') || '"__ux_ts"';
  end;  

  function wrapQueryColumns return varchar2 is
  begin
    if trim(p_qry_columns) = '*' then 
      return metaColumns() || ', '|| alias || '.*';
    else
      return metaColumns() || ', ' || p_qry_columns;
    end if;  
  end;  

  function wrapQueryFrom return varchar2 is
  begin
    return '('|| p_qry_from || ') ' || alias;
  end;   
  
begin
  for i in 1..p_events.count() loop
    binds(i) := org$xml_factory.newBindParam(i, p_events(i).key);
  end loop;
  
  qry := makeSqlPkInListQuery(
    p_cols    => wrapQueryColumns()
  , p_from    => wrapQueryFrom() -- 'test_tab m'
  , p_pk_col  => p_qry_pk_column -- 'id'
  , p_binds   => binds
  );
  
  org$xml_factory.dumpCursorAsXml(
    p_query      => qry
  , p_binds      => binds
  , r_rows_count => r_rows_count
  , r_rows_dump  => r_rows_dump
  );
end; /* dumpUpdatedRows */

procedure dumpDeletedRows(
  p_qry_pk_column in varchar2
, p_events in out nocopy org$outbox_api.TEventArray
, r_rows_dump out nocopy clob
, r_rows_count out number
)
is
begin
  org$xml_encode.initContext();

  for i in 1..p_events.count() loop
    -- <ROW>
    org$xml_encode.beginRow();
      org$xml_encode.addColumn(p_events(i).key, upper(p_qry_pk_column));
      org$xml_encode.addColumn(p_qry_pk_column, '__pk_name');
      org$xml_encode.addColumn(p_events(i).key, '__pk_val');
      org$xml_encode.addColumn(p_events(i).op, '__op'); 
      org$xml_encode.addColumn(toUnixTimestamp(p_events(i).ts), '__ux_ts'); 
      org$xml_encode.addColumn(FROM_TZ(p_events(i).ts, DBTIMEZONE) AT TIME ZONE '00:00', '__ts'); 
    org$xml_encode.endRow();      
    -- </ROW>  
  end loop;
  
  org$xml_encode.closeContext(r_rows_dump); 
  r_rows_count := p_events.count();
end; /* dumpDeletedRows */

procedure getNextEvents(
  p_group_id in varchar2
, p_part_id in number
, p_rows in number

, p_qry_columns in varchar2
, p_qry_from in varchar2
, p_qry_pk_column in varchar2

, r_upd_rows_dump out nocopy clob
, r_upd_rows_count out number
, r_del_rows_dump out nocopy clob
, r_del_rows_count out number
)
is
  all_events org$outbox_api.TEventArray;
  upd_events org$outbox_api.TEventArray;
  del_events org$outbox_api.TEventArray;
begin
  r_upd_rows_count := 0;
  
  org$outbox_api.getNewEvents(
    p_part_id   => p_part_id
  , p_group_id  => p_group_id
  , p_row_count => p_rows
  , r_events    => all_events
  );

  copmactAndSplitEvents(all_events, upd_events, del_events);
  
  if upd_events.count() > 0 then
    dumpUpdatedRows(p_qry_columns, p_qry_from, p_qry_pk_column, upd_events, r_upd_rows_dump, r_upd_rows_count);
  end if;    
  
  if del_events.count() > 0 then
    dumpDeletedRows(p_qry_pk_column, del_events, r_del_rows_dump, r_del_rows_count);
  end if;    

  org$outbox_api.markEventsAsProcessed(all_events);
end; /* getNextEvents */

procedure getNextEvents(
  p_group_id in varchar2
, p_part_id in number
, p_rows in number

, p_qry_columns in varchar2
, p_qry_from in varchar2
, p_qry_pk_column in varchar2

, r_upd_rows_dump out nocopy blob
, r_upd_rows_count out number
, r_del_rows_dump out nocopy blob
, r_del_rows_count out number
)
is
  v_upd_xml_text clob;
  v_del_xml_text clob;
begin
  getNextEvents(
    p_group_id       => p_group_id,
    p_part_id        => p_part_id,
    p_rows           => p_rows,
    p_qry_columns    => p_qry_columns,
    p_qry_from       => p_qry_from,
    p_qry_pk_column  => p_qry_pk_column,
    r_upd_rows_dump  => v_upd_xml_text,
    r_upd_rows_count => r_upd_rows_count,
    r_del_rows_dump  => v_del_xml_text,
    r_del_rows_count => r_del_rows_count
  );

  org$util.gzipPackage(v_upd_xml_text, r_upd_rows_dump);
  org$util.gzipPackage(v_del_xml_text, r_del_rows_dump);
end; /* getNextEvents */

end org$gate_api;
/

