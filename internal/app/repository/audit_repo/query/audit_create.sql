INSERT INTO audit
(id, obj_id, obj_entity, obj_name, actor_id, action, data, message, timestamp)
VALUES (:id, :obj_id, :obj_entity, :obj_name, :actor_id, :action, :data, :message, :timestamp)