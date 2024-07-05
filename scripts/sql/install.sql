set define off
prompt
prompt Creating table EVENT_LOG
prompt ========================
prompt
@@org_event_log.sql
prompt
prompt Creating package ORG$GATE_API
prompt =============================
prompt
@@org$gate_api.sql
prompt
prompt Creating package ORG$OUTBOX_API
prompt ===============================
prompt
@@org$outbox_api.sql
prompt
prompt Creating package ORG$UTIL
prompt =========================
prompt
@@org$util.sql
prompt
prompt Creating package ORG$XML_ENCODE
prompt ===============================
prompt
@@org$xml_encode.sql
prompt
prompt Creating package ORG$XML_FACTORY
prompt ================================
prompt
@@org$xml_factory.sql
prompt Done
set define on
