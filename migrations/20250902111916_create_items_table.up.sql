create table if not exists Items (
            ID_Item UUID constraint PK_Item primary key default uuid_generate_v4(),
            Item_Title varchar (50) not null,
            Item_Description varchar (300) not null,
            Item_Price int not null,
            Little_Picture bytea not null
        ); 