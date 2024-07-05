declare
    group_id varchar2(32) := 'group_1';
    bucket_count int := 42;
begin
for x in (select * from test_tab where rownum < 100001)
    loop
        org$outbox_api.putUpdateEvent(
            p_key_n        => x.id,
            p_group_id     => group_id,
            p_bucket_count => bucket_count
        );
end loop;
end;
/

commit
/

declare
group_id varchar2(32) := 'group_1';
    bucket_count int := 42;
begin
for x in (select * from test_tab where rownum < 100001)
    loop
        org$outbox_api.putDeleteEvent(
            p_key_n        => x.id,
            p_group_id     => group_id,
            p_bucket_count => bucket_count
        );
end loop;
end;
/

commit
/

declare
    group_id varchar2(32) := 'group_2';
    bucket_count int := 42;
begin
for x in (select * from test_tab where rownum < 100001)
    loop
        org$outbox_api.putUpdateEvent(
            p_key_n        => x.id,
            p_group_id     => group_id,
            p_bucket_count => bucket_count
        );
end loop;
end;
/

commit
/

declare
    group_id varchar2(32) := 'group_3';
    bucket_count int := 42;
begin
for x in (select * from test_tab where rownum < 100001)
    loop
        org$outbox_api.putUpdateEvent(
            p_key_n        => x.id,
            p_group_id     => group_id,
            p_bucket_count => bucket_count
        );
end loop;
end;
/

commit
/