INSERT INTO joblist (id,correlation_id,"name",created_at,created_by,modified_at,modified_by,status,"source",destination,"type",sub_type,"action",action_details,progress,history,extra_data,priority,"rank")
select substring(replace(to_char(clock_timestamp(),'yyyymmddhh24missus') || (to_char(random()*1e9,'000000000')),' ',''),1,27),
md5(RANDOM()::TEXT), 
md5(RANDOM()::TEXT),
NOW() - (random() * (interval '90 days')) + '30 days',
md5(RANDOM()::TEXT),
NOW() - (random() * (interval '90 days')) + '30 days',
md5(RANDOM()::TEXT),
md5(RANDOM()::TEXT),
md5(RANDOM()::TEXT),
md5(RANDOM()::TEXT),
md5(RANDOM()::TEXT),
md5(RANDOM()::TEXT),
md5(RANDOM()::TEXT),
md5(RANDOM()::TEXT),
0,
md5(RANDOM()::TEXT),
md5(RANDOM()::TEXT),
2,
0
end from pg_catalog.generate_series(1,10)