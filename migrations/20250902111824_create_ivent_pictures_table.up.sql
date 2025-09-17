create table if not exists Ivent_Pictures (
            ID_Ivent_Picture UUID constraint PK_Ivent_Picture primary key default uuid_generate_v4(),
            Ivent_ID UUID not null references Ivents (ID_Ivent),
            Picture_Path bytea not null
        ); 