create or replace package orgon.org$util is

-- Author  : VODYANKO
-- Created : 08/06/2024 12:52:42
-- Purpose : Utils

-- Character data compression using the standard utl_compress package.
-- Character data is converted to binary with UTF8 encoding.
procedure gzipPackage(
  text in out nocopy clob
, gzip out nocopy blob
);

end org$util;
/

create or replace package body orgon.org$util is

function clob_to_blob (p_data in clob) return blob
as
  l_blob         blob;
  l_dest_offset  pls_integer := 1;
  l_src_offset   pls_integer := 1;
  l_lang_context pls_integer := DBMS_LOB.default_lang_ctx;
  l_warning      pls_integer := DBMS_LOB.warn_inconvertible_char;
  l_blob_csid    integer := nls_charset_id('UTF8');
begin
  DBMS_LOB.createtemporary(
    lob_loc => l_blob,
    cache   => true);

  DBMS_LOB.converttoblob(
   dest_lob      => l_blob,
   src_clob      => p_data,
   amount        => DBMS_LOB.lobmaxsize,
   dest_offset   => l_dest_offset,
   src_offset    => l_src_offset,
   blob_csid     => l_blob_csid,
   lang_context  => l_lang_context,
   warning       => l_warning);
  
  return l_blob;
end;

procedure free_temporary (p_clob in out nocopy clob)
is
begin
  if dbms_lob.istemporary(p_clob) = 1 then
    dbms_lob.freetemporary(p_clob);
  end if;
end;

procedure free_temporary (p_blob in out nocopy blob)
is
begin
  if dbms_lob.istemporary(p_blob) = 1 then
    dbms_lob.freetemporary(p_blob);
  end if;
end;

procedure gzip_compress(
  p_src in blob
, r_dst in out nocopy blob
)
is
begin
  utl_compress.lz_compress(p_src, r_dst);
end; /* gzip_compress */

procedure gzip_compress(
  p_src in clob
, r_dst in out nocopy blob
)
is
  v_src_bytes blob;
begin
  v_src_bytes := clob_to_blob(p_src);
  gzip_compress(v_src_bytes, r_dst);
  free_temporary(v_src_bytes);
exception
  when others then
    free_temporary(v_src_bytes);
  raise;
end; /* gzip_compress */

procedure gzipPackage(
  text in out nocopy clob
, gzip out nocopy blob
)
is
begin
  if text is not null then
    begin
      dbms_lob.createtemporary(gzip, true);
      gzip_compress(text, gzip);
      dbms_lob.freetemporary(text);
    exception
      when others then
        dbms_lob.freetemporary(gzip);
        raise;
    end;
  end if;
end; /* gzipPackage */

end org$util;
/

