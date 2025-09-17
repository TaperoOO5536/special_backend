create table if not exists Item_Pictures (
            ID_Item_Picture UUID constraint PK_Item_Picture primary key default uuid_generate_v4(),
            Item_ID UUID not null references Items (ID_Item),
            Picture_Path bytea not null
        ); 