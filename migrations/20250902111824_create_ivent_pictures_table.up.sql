create table if not exists Event_Pictures (
            ID_Event_Picture UUID constraint PK_Event_Picture primary key default uuid_generate_v4(),
            Event_ID UUID not null references Events (ID_Event),
            Picture_Path bytea not null
        ); 