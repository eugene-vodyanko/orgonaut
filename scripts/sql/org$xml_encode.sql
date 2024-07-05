create or replace package orgon.org$xml_encode is

-- Author  : E.VODYANKO
-- Created : 18/04/2023 11:04:29
-- Purpose : Trivial stateful XML encoder for Rows metadata

/*
TEST 

declare 
  c clob;
begin
  org$xml_encode.initContext();
  
  for x in 
  (
    select level id, 'str:'||level str, systimestamp ts, sysdate dt from dual connect by level < 10
  )
  loop
    org$xml_encode.beginRow();
    -- <ROW>
      org$xml_encode.addColumn(x.id, 'id');
      org$xml_encode.addColumn(x.str, 'str'); 
      org$xml_encode.addColumn(x.ts, 'ts'); 
      org$xml_encode.addColumn(x.dt, 'dt'); 
      org$xml_encode.addColumn(3.14, 'pi'); 
    -- </ROW>
    org$xml_encode.endRow();
  end loop;
  
  org$xml_encode.closeContext(c);
  
  dbms_output.put_line(to_char(c));
end;  

*/

procedure initContext;

procedure closeContext(
  r_dump out nocopy clob
);

procedure beginRow;

procedure endRow;

procedure addColumn(
  obj in varchar2
, key in varchar2
);

procedure addColumn(
  obj in number
, key in varchar2
);

procedure addColumn(
  obj in date
, key in varchar2
);

procedure addColumn(
  obj in timestamp with time zone
, key in varchar2
);

end org$xml_encode;
/

create or replace package body orgon.org$xml_encode is

g_buf_str varchar2(32767);
g_dump_clob clob;

procedure addToClob(
  buf_lob in out nocopy clob
, buf_str in out nocopy varchar2
, str varchar2
)
is
begin
  if (length(str) >= 32767 - length(buf_str)) then
    dbms_lob.writeappend(buf_lob, length(buf_str), buf_str);
    buf_str := str;
  else
      buf_str := buf_str || str;
  end if;
end; /* addToClob */

procedure initContext 
is
begin
  g_buf_str := null;
  dbms_lob.createtemporary(g_dump_clob, true);

  addToClob(
    buf_lob => g_dump_clob, 
    buf_str => g_buf_str, 
    str => '<?xml version="1.0"?><ROWSET>' || chr(10)
  );
end; /* initContext */

procedure flushBuffer(
  buf_lob in out nocopy clob
, buf_str in out nocopy varchar2
)
is
begin
  dbms_lob.writeappend(buf_lob, length(buf_str), buf_str);
  buf_str := null;
end; /* flushBuffer */

procedure closeContext(r_dump out nocopy clob)
is
begin
  addToClob(
    buf_lob => g_dump_clob, 
    buf_str => g_buf_str, 
    str => '</ROWSET>' || chr(10)
  );
  
  if g_buf_str is not null then
    flushBuffer(buf_lob => g_dump_clob, buf_str => g_buf_str);
  end if;

  g_buf_str := null;
  r_dump := g_dump_clob;
end; /* closeContext */

procedure intBeginRow(
  buf in out nocopy clob
, buf_str in out nocopy varchar2
)
is
begin
  addToClob(buf, buf_str, '<ROW>' || chr(10));
end; /* intBeginRow */

procedure intEndRow(
  buf in out nocopy clob
, buf_str in out nocopy varchar2
)
is
begin
  addToClob(buf, buf_str, '</ROW>' || chr(10));
end; /* intEndRow */

procedure putString(
  buf in out nocopy clob
, buf_str in out nocopy varchar2
, val in varchar2
, key in varchar2
)
is
begin
  addToClob(buf, buf_str, '<' || key || '>' || val || '</' || key ||'>' || chr(10));
end; /* putString */

procedure addColumn(
  obj in varchar2
, key in varchar2
)
is
begin
  putString(g_dump_clob, g_buf_str, obj, key);
end; /* addColumn */


procedure addColumn(
  obj in number
, key in varchar2
)
is
begin
  putString(g_dump_clob, g_buf_str, obj, key);
end; /* addColumn */

procedure addColumn(
  obj in date
, key in varchar2
)
is
begin
  putString(g_dump_clob, g_buf_str, obj, key);
end; /* addColumn */

procedure addColumn(
  obj in timestamp with time zone
, key in varchar2
)
is
begin
  putString(g_dump_clob, g_buf_str, obj, key);
end; /* addColumn */

procedure beginRow
is
begin
  intBeginRow(g_dump_clob, g_buf_str);
end; /* beginRow */

procedure endRow
is
begin
  intEndRow(g_dump_clob, g_buf_str);
end; /* endRow */

begin
  execute immediate q'[alter session set nls_date_format = 'YYYY-MM-DD HH24:MI:SS']'; 
  execute immediate q'[alter session set nls_timestamp_format='YYYY-MM-DD HH24:MI:SS.FF6']';
  execute immediate q'[alter session set nls_timestamp_tz_format='YYYY-MM-DD"T"HH24:MI:SS.FF6 TZR']';
  execute immediate q'[alter session set nls_numeric_characters='.,']';
end org$xml_encode;
/

