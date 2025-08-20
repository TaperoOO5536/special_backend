create table if not exists Ivent (
            ID_Ivent serial not null constraint PK_Ivent primary key,
            Ivent_Title varchar (50) not null,
            Ivent_Description varchar (300) not null,
            Ivent_Date date not null,
            Ivent_Time time not null,
        );