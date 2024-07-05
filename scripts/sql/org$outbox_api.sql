create or replace package orgon.org$outbox_api is

-- Author  : VODYANKO
-- Created : 17/04/2024 18:43:22
-- Purpose : The package implements API-methods for working with the "Transactional outbox" pattern.

/*
OVERVIEW

The package only provides an API for inserting events into outbox "queue"-table and fetching.
The implementation of the data capture mechanism itself remains outside the scope of the package and 
can be implemented by embedding API calls directly in the application code or in triggers (after insertion/update/deletion) 
on the corresponding tracked tables.

Insertion into the underlying outbox "queue"-table can become a bottleneck due to the features of monotonous 
sequential indexes (based on the b-tree) in highly parallel access conditions.

Therefore, the concept of sharding is used to reduce hotspots:
- Inserts are distributed in buckets.
- Deterministic hash function applied to the passed key is used (the same keys always end up in the same package).

The user must select a unique code (string) for a specific load.
The user must also determine the number of buckets:
- Too high a value is not desirable due to reduced caching efficiency.
- Too small a value may cause synchronization conflicts during recording into DBMS.

The likely optimal value may be in the range of 10-100 (the greater the number of CPU cores and the number of 
parallel sessions, the greater the value).

*/

/*
TEST

-- Emulating update events
declare
  group_id varchar2(32) := 'test_tab';
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

-- Emulating deletion events
declare
    group_id varchar2(32) := 'test_tab';
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

-- Getting the next events in the queue
declare
  eve org$outbox_api.TEventArray;
begin
  org$outbox_api.getNewEvents(p_part_id   => 1,
                              p_row_count => 100,
                              p_group_id  => 'test_tab',
                              r_events    => eve
                              );

  for i in 1 .. eve.count() loop
    dbms_output.put_line(eve(i).rid || ' ' || eve(i).key || ' ' || eve(i).op);
  end loop;
end;

*/

ACTION_INSERT constant varchar2(1) := 'c';
ACTION_UPDATE constant varchar2(1) := 'u';
ACTION_DELETE constant varchar2(1) := 'd';

STATE_NEW constant varchar2(1) := 'n';
STATE_PROCESSED constant varchar2(1) := 'p';


-- TEvent describes a change data capture event.
-- Currently, only events for tables with a primary numeric key are supported.
type TEvent is record (
  rid rowid,
  key number,
  op varchar2(1), 
  ts timestamp
);

-- TEventArray is used to increase throughput and reduce context switching.
type TEventArray is table of TEvent;
  
-- Putting insert event in the "outbox-queue".
-- @p_group_id - unique payload code, selected by the user.
-- @p_bucket_count - number of buckets for sharding (see overview), selected by the user.
-- @p_key_n - primary key of the modified record (only numeric single column key are supported).
procedure putInsertEvent(
  p_group_id in varchar2
, p_key_n in number 
, p_bucket_count in number
);

-- Putting update event in the "outbox-queue".
procedure putUpdateEvent(
  p_group_id in varchar2
, p_key_n in number
, p_bucket_count in number
);

-- Putting delete event in the "outbox-queue".
procedure putDeleteEvent(
  p_group_id in varchar2
, p_key_n in number
, p_bucket_count in number
);

-- Mark the event as processed.
-- The method opens a transaction.
procedure markEventsAsProcessed(
  p_events in out nocopy TEventArray
);

-- Receiving the following events in the queue that have not yet been processed
-- The method is read-only, it is not thread-safe.
-- That is, a call in two sessions may return the same data.
procedure getNewEvents(
  p_part_id in number
, p_group_id in varchar2
, p_row_count in number
, r_events out nocopy TEventArray
);

end org$outbox_api;
/

create or replace package body orgon.org$outbox_api is

procedure putNewEvent(
  p_group_id in varchar2
, p_key_n in number
, p_action in varchar2
, p_bucket_count in number
) 
is
begin
  insert into EVENT_LOG(group_id, part_id, state, ts, key_n, action) 
    values(p_group_id, ora_hash(p_key_n, p_bucket_count - 1), STATE_NEW, systimestamp, p_key_n, p_action);
end; /* putNewEvent */

procedure putInsertEvent(
  p_group_id in varchar2
, p_key_n in number
, p_bucket_count in number
) 
is
begin
  putNewEvent(p_group_id, p_key_n, ACTION_INSERT, p_bucket_count);
end; /* putInsertEvent */

procedure putUpdateEvent(
  p_group_id in varchar2
, p_key_n in number
, p_bucket_count in number
) 
is
begin
  putNewEvent(p_group_id, p_key_n, ACTION_UPDATE, p_bucket_count);
end; /* putUpdateEvent */

procedure putDeleteEvent(
  p_group_id in varchar2
, p_key_n in number
, p_bucket_count in number
) 
is
begin
  putNewEvent(p_group_id, p_key_n, ACTION_DELETE, p_bucket_count);
end; /* putDeleteEvent */

procedure markEventsAsProcessed(
  p_events in out nocopy TEventArray
) 
is
begin
  forall i in 1.. p_events.count()
    update EVENT_LOG set state = STATE_PROCESSED
      where rowid = p_events(i).rid;
end; /* markEventsAsProcessed */

procedure getNewEvents(
  p_part_id in number
, p_group_id in varchar2
, p_row_count in number
, r_events out nocopy TEventArray
)
is
begin
  select rowid, key_n, action, ts bulk collect into r_events from 
  (
    select /*+ FIRST_ROWS(1) DYNAMIC_SAMPLING(0) */ key_n, action, ts from EVENT_LOG
    where group_id = p_group_id
      and part_id = p_part_id
      and state = STATE_NEW
    order by ts
  )
  where rownum <= p_row_count;
end; /* getNewEvents */

begin 
  execute immediate 'alter session set nls_sort = BINARY';
end org$outbox_api;
/

