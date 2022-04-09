use timetable;

SELECT c.id, c.name, duration.id, duration.name, 
frame.id, frame.day_week, frame.period, subject.id, subject.name, 
teacher.id, teacher.name FROM 
(SELECT * FROM normal_timetable WHERE duration_id = 0 AND class_id = 0) AS tb
LEFT JOIN classroom AS c ON tb.class_id = c.id
LEFT JOIN duration ON tb.duration_id = duration.id
LEFT JOIN frame ON tb.frame_id = frame.id
LEFT JOIN subject ON tb.subject_id = subject.id
LEFT JOIN teacher ON tb.teacher_id = teacher.id;

SELECT * FROM duration WHERE start_date = "2021-04-01";
SELECT * FROM duration WHERE start_date <= "2021-04-01" AND "2021-04-01" <= end_date;

select * from normal_timetable;

show tables;