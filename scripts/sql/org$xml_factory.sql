create or replace package orgon.org$xml_factory is

-- Author  : E.VODYANKO
-- Created : 18/09/2019 17:21:36
-- Purpose : SQL Query to XML row-serializer based on std. package "dbms_xmlgen"

/*
TEST

declare
  qry varchar2(32000) := 'select * from test_tab where id in (:id_1, :id_2, :id_3)';
  rowset clob;
  rowcnt  int;
  binds org$xml_factory.TBindParams;
  i int := 0;
begin
  for x in (select * from test_tab where rownum < = 3) loop
    i := i + 1;
    binds(i) := org$xml_factory.newBindParam('id_' || i, x.id);
  end loop;

  org$xml_factory.dumpCursorAsXml(p_query      => qry,
                                  p_binds      => binds,
                                  r_rows_count => rowcnt,
                                  r_rows_dump  => rowset);

  dbms_output.put_line('rowset: ' || substr(rowset, 1, 4000));
  dbms_output.put_line('rowset len: ' || to_char(length(rowset)));
  dbms_output.put_line('rowcnt: ' || to_char(rowcnt));
end;

*/

type TBindParam is record (
  /* private */
  name varchar2(32)
, value_dt date
, value_num number
, value_str varchar2(4000)

, is_dt boolean
, is_num boolean
, is_str boolean
);

type TBindParams is table of TBindParam index by binary_integer;

function newBindParam(p_name in varchar2, p_val in number) return TBindParam;

-- Dump rows returned by a parameterized query to XML.
procedure dumpCursorAsXml(
  p_query in varchar2
, p_binds in out nocopy TBindParams
, r_rows_count out number
, r_rows_dump out nocopy clob
);

end org$xml_factory;
/

create or replace package body orgon.org$xml_factory is

function newBindParam(p_name in varchar2, p_val in number) return TBindParam
is
  b TBindParam;
begin
  b.is_num := true;
  b.name := p_name;
  b.value_num := p_val;

  return b;
end; /* newBindParam */

procedure dumpCursorAsXml(
  p_query in varchar2
, p_binds in out nocopy TBindParams
, r_rows_count out number
, r_rows_dump out nocopy clob
)
is
  ctx  number;
begin
  ctx := dbms_xmlgen.newContext(p_query);
  dbms_xmlgen.setCheckInvalidChars(ctx, true);
  dbms_xmlgen.setConvertSpecialChars(ctx, true);
  dbms_xmlgen.setNullHandling(ctx, dbms_xmlgen.DROP_NULLS);
  
  for i in 1..p_binds.count() loop
    dbms_xmlgen.setBindValue(
      ctx => ctx
    , bindName => p_binds(i).name
    , bindValue => p_binds(i).value_num
    );
  end loop;    
  
  dbms_xmlgen.restartQuery(ctx);
  r_rows_dump := dbms_xmlgen.getXML(ctx);  
  r_rows_count := dbms_xmlgen.getNumRowsProcessed(ctx);
  
  dbms_xmlgen.closecontext(ctx); 
end; /* dumpCursorAsXml */

begin
  execute immediate q'[alter session set nls_date_format = 'YYYY-MM-DD HH24:MI:SS']'; 
  execute immediate q'[alter session set nls_timestamp_format='YYYY-MM-DD HH24:MI:SS.FF6']';
  execute immediate q'[alter session set nls_timestamp_tz_format='YYYY-MM-DD"T"HH24:MI:SS.FF6 TZR']';
  execute immediate q'[alter session set nls_numeric_characters='.,']';
end;
/

