create table if not exists Ivent (
            ID_Ivent UUID constraint PK_Ivent primary key default uuid_generate_v4(),
            Ivent_Title varchar (50) not null,
            Ivent_Description varchar (300) not null,
            Ivent_Date date not null,
            Ivent_Time time not null,
        );