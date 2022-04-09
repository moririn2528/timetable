use timetable;

drop table deleted_normal_timetable;
drop table deleted_additional_timetable;
drop table normal_timetable;
drop table additional_timetable;
drop table holiday;
drop table class_struct_edge;
drop table classroom;
drop table teacher;
drop table duration;
drop table subject;
drop table frame;
drop table place;


create table classroom(
	id int not null primary key,
    name char(30) not null,
    available char(50) not null,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
create table class_struct_edge(
	from_id int not null,
    to_id int not null,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY fk_from(from_id) REFERENCES classroom(id),
    FOREIGN KEY fk_from(to_id) REFERENCES classroom(id),
    UNIQUE id_pair(from_id,to_id)
);
create table teacher(
	id int not null primary key,
    name char(30) not null,
    avoid char(50) not null,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
create table duration(
	id int not null primary key,
    name char(30) not null,
    start_date date not null,
    end_date date not null,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
create table subject(
	id int not null primary key,
    name char(30) not null,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
create table frame(
	id int not null primary key,
    day_week int not null,
    period int not null,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
create table place(
	id int not null primary key,
    name char(30) not null,
	count int not null,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
create table normal_timetable(
	id int not null auto_increment primary key,
	duration_id int not null,
	class_id int not null,
	teacher_id int not null,
    subject_id int not null,
    frame_id int not null,
	place_id int not null,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY fk_duration(duration_id) REFERENCES duration(id),
    FOREIGN KEY fk_class(class_id) REFERENCES classroom(id),
    FOREIGN KEY fk_teacher(teacher_id) REFERENCES teacher(id),
    FOREIGN KEY fk_subject(subject_id) REFERENCES subject(id),
    FOREIGN KEY fk_frame(frame_id) REFERENCES frame(id),
    FOREIGN KEY fk_place(place_id) REFERENCES place(id)
);
create table additional_timetable(
	id int not null auto_increment primary key,
	duration_id int not null,
	class_id int not null,
	teacher_id int not null,
    subject_id int not null,
    frame_id int not null,
	place_id int not null,
    day date not null,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY fk_duration(duration_id) REFERENCES duration(id),
    FOREIGN KEY fk_class(class_id) REFERENCES classroom(id),
    FOREIGN KEY fk_teacher(teacher_id) REFERENCES teacher(id),
    FOREIGN KEY fk_subject(subject_id) REFERENCES subject(id),
    FOREIGN KEY fk_frame(frame_id) REFERENCES frame(id),
    FOREIGN KEY fk_place(place_id) REFERENCES place(id)
);
create table deleted_normal_timetable(
	id int not null,
    day date not null,
	deleted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY fk_id(id) REFERENCES normal_timetable(id),
    unique id_day(id, day)
);
create table deleted_additional_timetable(
	id int not null auto_increment primary key,
	duration_id int not null,
	class_id int not null,
	teacher_id int not null,
    subject_id int not null,
    frame_id int not null,
	place_id int not null,
    day date not null,
	deleted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
create table holiday(
	day date not null unique
);

insert into classroom(id, name, available) values 
(0, "1組", "111100000000000000000000000000000000000"),
(1, "2組", "111100000000000000000000000000000000000"),
(2, "3組", "111000000000000000000000000000000000000"),
(3, "1年", "111000000000000000000000000000000000000"),
(4, "松", "111100000000000000000000000000000000000"),
(5, "男子英数", "111000000000000000000000000000000000000"),
(10, "高3",  "111111011111111111110111111011111111111"),
(11, "高3男子英数",  "111111011111111111110111111011111111111"),
(12, "高3女子英数",  "111111011111111111110111111011111111111"),
(13, "高3女子特進",  "111111011111111111110111111011111111111"),
(20, "3A",  "111111011111111111110111111011111111111"),
(21, "3B",  "111111011111111111110111111011111111111"),
(22, "3C",  "111111011111111111110111111011111111111"),
(23, "3D",  "111111011111111111110111111011111111111"),
(24, "3E",  "111111011111111111110111111011111111111"),
(25, "3F",  "111111011111111111110111111011111111111"),
(26, "3G",  "111111011111111111110111111011111111111"),
(27, "3H",  "111111011111111111110111111011111111111"),
(28, "3I",  "111111011111111111110111111011111111111"),
(29, "3J",  "111111011111111111110111111011111111111"),
(30, "3S",  "111111011111111111110111111011111111111"),
(31, "3S文",  "111111011111111111110111111011111111111"),
(32, "3S文物理",  "111111011111111111110111111011111111111"),
(33, "3S文生物",  "111111011111111111110111111011111111111"),
(34, "3S文世界史",  "111111011111111111110111111011111111111"),
(35, "3S文日本史",  "111111011111111111110111111011111111111"),
(40, "3S理",  "111111011111111111110111111011111111111"),
(41, "3Mab",  "111111011111111111110111111011111111111"),
(42, "3Ma",  "111111011111111111110111111011111111111"),
(43, "3Mb",  "111111011111111111110111111011111111111"),
(44, "3Mab物理",  "111111011111111111110111111011111111111"),
(45, "3Mab生物",  "111111011111111111110111111011111111111"),
(46, "3K類",  "111111011111111111110111111011111111111"),
(47, "3K類物理",  "111111011111111111110111111011111111111"),
(48, "3K類生物",  "111111011111111111110111111011111111111"),
(50, "3S理地理",  "111111011111111111110111111011111111111"),
(51, "3S理世界史",  "111111011111111111110111111011111111111"),
(52, "3S理日本史",  "111111011111111111110111111011111111111"),
(60, "3英数",  "111111011111111111110111111011111111111"),
(61, "3M12",  "111111011111111111110111111011111111111"),
(62, "3M1",  "111111011111111111110111111011111111111"),
(63, "3M2男",  "111111011111111111110111111011111111111"),
(64, "3M2女",  "111111011111111111110111111011111111111"),
(70, "3M12物理1",  "111111011111111111110111111011111111111"),
(71, "3M12物理2",  "111111011111111111110111111011111111111"),
(72, "3M12生物",  "111111011111111111110111111011111111111"),
(73, "3M12地理",  "111111011111111111110111111011111111111"),
(74, "3M12世界史",  "111111011111111111110111111011111111111"),
(75, "3M12日本史",  "111111011111111111110111111011111111111"),
(80, "3英数文",  "111111011111111111110111111011111111111"),
(81, "3英数文1",  "111111011111111111110111111011111111111"),
(82, "3英数文2",  "111111011111111111110111111011111111111"),
(83, "3英数文世界史",  "111111011111111111110111111011111111111"),
(84, "3英数文日本史",  "111111011111111111110111111011111111111"),
(85, "3英数文物理",  "111111011111111111110111111011111111111"),
(86, "3英数文化学",  "111111011111111111110111111011111111111"),
(87, "3英数文生物",  "111111011111111111110111111011111111111"),
(90, "3EF",  "111111011111111111110111111011111111111"),
(91, "3EF文",  "111111011111111111110111111011111111111"),
(100, "3EF理",  "111111011111111111110111111011111111111");
insert into teacher(id, name, avoid) values 
(0, "石川先生","000000000000000000000000000000000000000"),
(1, "西川先生","000000000000000000000000000000000000000"),
(2, "赤坂先生","000000000000000000000000000000000000000");
insert into duration(id, name, start_date, end_date) values (0, "令和3年1学期", "2021-4-1", "2021-7-31");
insert into subject(id, name) values (0, "英語"), (1, "数学"), (2, "情報");
insert into frame(id, day_week, period) values (0,1,0),(1,1,1),(2,1,2),(3,1,3),(4,1,4),(5,1,5),(6,1,6),(7,2,0),(8,2,1),(9,2,2),(10,2,3),(11,2,4),(12,2,5),(13,2,6),(14,3,0),(15,3,1),(16,3,2),(17,3,3),(18,3,4),(19,3,5),(20,3,6),(21,4,0),(22,4,1),(23,4,2),(24,4,3),(25,4,4),(26,4,5),(27,4,6),(28,5,0),(29,5,1),(30,5,2),(31,5,3),(32,5,4),(33,5,5),(34,5,6),(35,6,0),(36,6,1),(37,6,2),(38,6,3);
insert into place(id, name, count) values (0,"自教室",100),(1,"コンピューター室",1);
insert into normal_timetable(duration_id,class_id,teacher_id,subject_id,frame_id,place_id) values
(0,3,0,0,0,0),(0,4,0,0,1,0),(0,5,2,1,1,0),(0,4,1,2,2,1),(0,5,0,0,2,0),(0,0,1,2,3,0),(0,1,0,0,3,0);
insert into deleted_normal_timetable(id,day) values
(1,"2021-4-5");
insert into additional_timetable(duration_id,class_id,teacher_id,subject_id,frame_id,place_id,day) values
(0,3,2,1,0,0,"2021-4-5");
insert into class_struct_edge(from_id, to_id) values 
(3,4), (3,5), (5,0), (5,1), (5,2),
(10,11),(10,12),(10,13),(10,30),(10,60),(10,90),
(11,20),(11,21),(11,22),(11,23),(12,26),(12,27),(12,28),(12,29),(13,24),(13,25),
(30,20),(30,21),(30,28),(30,29),(30,31),(30,40),
(31,32),(31,33),(31,34),(31,35),
(40,41),(40,46),(40,50),(40,51),(40,52),
(41,42),(41,43),(41,44),(41,45),
(46,47),(46,48),
(60,22),(60,23),(60,26),(60,27),(60,61),(60,80),
(61,62),(61,63),(61,64),(61,70),(61,71),(61,72),(61,73),(61,74),(61,75),
(80,81),(80,82),(80,83),(80,84),(80,85),(80,86),(80,87),
(90,24),(90,25),(90,91),(90,100);
insert into holiday(day) values ("2021-1-1"),("2021-1-11"),
("2021-2-11"),("2021-2-23"),("2021-3-20"),("2021-4-29"),("2021-5-3"),
("2021-5-4"),("2021-5-5"),("2021-7-22"),("2021-7-23"),("2021-8-8"),("2021-9-20"),
("2021-9-23"),("2021-11-3"),("2021-11-23");
